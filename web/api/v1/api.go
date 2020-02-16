package v1

import (
	. "chaos-slave/common/docker"
	. "chaos-slave/common/service"
	"chaos-slave/proto"
	"context"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthCheckService struct {
}

func (hcs *HealthCheckService) Check(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{Status: proto.HealthCheckResponse_SERVING}, nil
}

func (hcs *HealthCheckService) Watch(req *proto.HealthCheckRequest, srv proto.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

type ServiceManager struct {
	Logger     log.Logger
}

func (sm *ServiceManager) Start(ctx context.Context, req *proto.ServiceRequest) (*proto.StatusResponse, error) {
	s := &Service{Name: req.Name, Logger: sm.Logger}
	return s.Start()
}

func (sm *ServiceManager) Stop(ctx context.Context, req *proto.ServiceRequest) (*proto.StatusResponse, error) {
	s := &Service{Name: req.Name, Logger: sm.Logger}
	return s.Stop()
}

type DockerManager struct {
	Logger     log.Logger
}

func (sm *DockerManager) Start(ctx context.Context, req *proto.DockerRequest) (*proto.StatusResponse, error) {
	s := &Docker{Name: req.Name, Logger: sm.Logger}
	return s.Start()
}

func (sm *DockerManager) Stop(ctx context.Context, req *proto.DockerRequest) (*proto.StatusResponse, error) {
	s := &Docker{Name: req.Name, Logger: sm.Logger}
	return s.Stop()
}