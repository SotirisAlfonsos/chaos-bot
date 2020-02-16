package web

import (
	"chaos-slave/proto"
	api "chaos-slave/web/api/v1"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"net"
)

type GrpcHandler struct {
	Port       string
	Logger     log.Logger
	grpcServer *grpc.Server
}

func NewGrpcHandler(port string, logger log.Logger) *GrpcHandler {
	grpcServer := grpc.NewServer()
	return &GrpcHandler{port, logger, grpcServer}
}

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
	proto.RegisterServiceServer(h.grpcServer, &api.ServiceManager{})
}

func (h *GrpcHandler) Stop() {
	h.grpcServer.GracefulStop()
}
