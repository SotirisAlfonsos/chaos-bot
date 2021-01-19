package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/SotirisAlfonsos/chaos-bot/config"
	"github.com/SotirisAlfonsos/chaos-bot/proto"
	api "github.com/SotirisAlfonsos/chaos-bot/web/api/v1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCHandler is holding the cache and GRPC configuration
type GRPCHandler struct {
	Port       string
	Logger     log.Logger
	GRPCServer *grpc.Server
	cache      *cache.Cache
}

type AuthInterceptor struct {
	peerToken string
}

// NewGRPCHandler creates and returns an instance of GRPCHandler
func NewGRPCHandler(
	port string,
	logger log.Logger,
	cache *cache.Cache,
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

	return &GRPCHandler{port, logger, GRPCServer, cache}, nil
}

func loadTLSCredentials(certFile string, keyFile string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		err = errors.Wrap(err, "failed to load key pair")
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
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
func (h *GRPCHandler) Run() error {
	_ = level.Info(h.Logger).Log("msg", "starting web server on port "+h.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", h.Port))
	if err != nil {
		return err
	}

	h.registerServices()

	if err := h.GRPCServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (h *GRPCHandler) registerServices() {
	proto.RegisterHealthServer(h.GRPCServer, &api.HealthCheckService{})
	proto.RegisterServiceServer(h.GRPCServer, &api.ServiceManager{
		Cache:  h.cache,
		Logger: h.Logger,
	})
	proto.RegisterDockerServer(h.GRPCServer, &api.DockerManager{
		Cache:  h.cache,
		Logger: h.Logger,
	})
	proto.RegisterStrategyServer(h.GRPCServer, &api.StrategyManager{
		Cache:  h.cache,
		Logger: h.Logger,
	})
}

// Stop stops the bot GRPC server
func (h *GRPCHandler) Stop() {
	h.GRPCServer.GracefulStop()
}
