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

//GrpcHandler is holding the cache and grpc configuration
type GrpcHandler struct {
	Port       string
	Logger     log.Logger
	grpcServer *grpc.Server
	cache      *cache.Cache
}

//NewGrpcHandler creates and returns an instance of GrpcHandler
func NewGrpcHandler(port string, logger log.Logger, cache *cache.Cache) *GrpcHandler {
	grpcServer := grpc.NewServer()
	return &GrpcHandler{port, logger, grpcServer, cache}
}

//Run starts the slave grpc server
func (h *GrpcHandler) Run() error {
	_ = level.Info(h.Logger).Log("msg", "starting web server on port "+h.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", h.Port))
	if err != nil {
		return err
	}

	h.registerServices()

	if err := h.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (h *GrpcHandler) registerServices() {
	proto.RegisterHealthServer(h.grpcServer, &api.HealthCheckService{})
	proto.RegisterServiceServer(h.grpcServer, &api.ServiceManager{
		Cache:  h.cache,
		Logger: h.Logger,
	})
	proto.RegisterDockerServer(h.grpcServer, &api.DockerManager{
		Cache:  h.cache,
		Logger: h.Logger,
	})
	proto.RegisterStrategyServer(h.grpcServer, &api.StrategyManager{
		Cache:  h.cache,
		Logger: h.Logger,
	})
}

//Stop stops the slave grpc server
func (h *GrpcHandler) Stop() {
	h.grpcServer.GracefulStop()
}
