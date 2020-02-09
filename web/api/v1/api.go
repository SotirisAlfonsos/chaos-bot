package v1

import (
	"chaos-slave/proto"
	"context"
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
