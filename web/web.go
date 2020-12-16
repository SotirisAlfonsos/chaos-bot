package web

import (
	"fmt"
	"net"

	"github.com/SotirisAlfonsos/chaos-slave/proto"
	api "github.com/SotirisAlfonsos/chaos-slave/web/api/v1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

// GRPCHandler is holding the cache and GRPC configuration
type GRPCHandler struct {
	Port       string
	Logger     log.Logger
	GRPCServer *grpc.Server
	cache      *cache.Cache
}

// NewGRPCHandler creates and returns an instance of GRPCHandler
func NewGRPCHandler(port string, logger log.Logger, cache *cache.Cache) *GRPCHandler {
	GRPCServer := grpc.NewServer()
	return &GRPCHandler{port, logger, GRPCServer, cache}
}

// Run starts the slave GRPC server
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

// Stop stops the slave GRPC server
func (h *GRPCHandler) Stop() {
	h.GRPCServer.GracefulStop()
}
