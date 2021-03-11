package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/SotirisAlfonsos/chaos-bot/common/cpu"
	"github.com/SotirisAlfonsos/chaos-bot/common/server"
	"github.com/SotirisAlfonsos/chaos-bot/config"
	protoV1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	apiV1 "github.com/SotirisAlfonsos/chaos-bot/web/api/v1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCHandler is holding the GRPC configuration
type GRPCHandler struct {
	Port       string
	Logger     log.Logger
	GRPCServer *grpc.Server
}

type AuthInterceptor struct {
	peerToken string
}

// NewGRPCHandler creates and returns an instance of GRPCHandler
func NewGRPCHandler(
	port string,
	logger log.Logger,
	conf *config.Config,
) (*GRPCHandler, error) {
	opts := make([]grpc.ServerOption, 0)
	if conf.CertFile != "" && conf.KeyFile != "" {
		tlsCredentials, err := loadTLSCredentials(conf.CertFile, conf.KeyFile)
		if err != nil {
			return nil, err
		}

		authInterceptor := &AuthInterceptor{
			peerToken: conf.PeerToken,
		}

		opts = append(opts, grpc.UnaryInterceptor(authInterceptor.validateToken))
		opts = append(opts, grpc.Creds(tlsCredentials))
	} else {
		_ = level.Warn(logger).Log("msg", "Insecure... Starting bot without tls certificates")
	}

	GRPCServer := grpc.NewServer(opts...)

	return &GRPCHandler{port, logger, GRPCServer}, nil
}

func loadTLSCredentials(certFile string, keyFile string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		err = errors.Wrap(err, "failed to load key pair")
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	return credentials.NewTLS(tlsConfig), nil
}

func (auth *AuthInterceptor) validateToken(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if !valid(md["authorization"], auth.peerToken) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}

func valid(authorization []string, peerToken string) bool {
	if len(authorization) < 1 {
		return false
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	return token == peerToken
}

// Run starts the bot GRPC server
func (h *GRPCHandler) Run() {
	_ = level.Info(h.Logger).Log("msg", "starting web server on port "+h.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", h.Port))
	if err != nil {
		_ = level.Error(h.Logger).Log("msg", "could not listen on port "+h.Port, "err", err)
		os.Exit(1)
	}

	h.registerServices()

	c := make(chan os.Signal, 1)
	e := make(chan error)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := h.GRPCServer.Serve(lis); err != nil {
			e <- err
		}
	}()

	select {
	case err := <-e:
		_ = level.Error(h.Logger).Log("msg", "server error", "err", err)
		os.Exit(1)
	case <-c:
		_ = level.Info(h.Logger).Log("msg", "Gracefully shutting down GRPC server")
		h.GRPCServer.Stop()
		os.Exit(0)
	}
}

func (h *GRPCHandler) registerServices() {
	protoV1.RegisterHealthServer(h.GRPCServer, &apiV1.HealthCheckService{})
	protoV1.RegisterServiceServer(h.GRPCServer, apiV1.NewServiceHandler(h.Logger))
	protoV1.RegisterDockerServer(h.GRPCServer, &apiV1.DockerHandler{
		Logger: h.Logger,
	})
	protoV1.RegisterCPUServer(h.GRPCServer, &apiV1.CPUManager{
		CPU:    cpu.New(h.Logger),
		Logger: h.Logger,
	})
	protoV1.RegisterServerServer(h.GRPCServer, &apiV1.ServerManager{
		Server: server.New(h.Logger),
		Logger: h.Logger,
	})
	protoV1.RegisterNetworkServer(h.GRPCServer, apiV1.NewNetworkManager(h.Logger))
}

// Stop stops the bot GRPC server
func (h *GRPCHandler) Stop() {
	h.GRPCServer.GracefulStop()
}
