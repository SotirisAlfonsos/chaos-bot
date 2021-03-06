package v1

import (
	"context"

	"github.com/SotirisAlfonsos/chaos-bot/common"
	"github.com/SotirisAlfonsos/chaos-bot/common/docker"
	"github.com/SotirisAlfonsos/chaos-bot/common/service"
	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HealthCheckService is the rpc
type HealthCheckService struct {
	*v1.UnimplementedHealthServer
}

// Check the health of the chaos bot
func (hcs *HealthCheckService) Check(context.Context, *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	return &v1.HealthCheckResponse{Status: v1.HealthCheckResponse_SERVING}, nil
}

// Watch is not used at the moment
func (hcs *HealthCheckService) Watch(*v1.HealthCheckRequest, v1.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

// ServiceManager is the rpc for services management
type ServiceManager struct {
	Logger log.Logger
	*v1.UnimplementedServiceServer
}

type response struct {
	message string
	err     error
}

// Start a service based on the name
func (sm *ServiceManager) Start(ctx context.Context, req *v1.ServiceRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.service.Start")
	defer span.End()

	resp := make(chan response, 1)

	serviceManage := &service.Service{Name: req.Name, Logger: sm.Logger}

	go func() {
		resp <- startTarget(serviceManage)
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(sm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

// Stop a service based on the name
func (sm *ServiceManager) Stop(ctx context.Context, req *v1.ServiceRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.service.Stop")
	defer span.End()

	resp := make(chan response, 1)

	serviceManage := &service.Service{Name: req.Name, Logger: sm.Logger}

	go func() {
		resp <- stopTarget(serviceManage)
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(sm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

// DockerManager is the rpc for docker management
type DockerManager struct {
	Logger log.Logger
	*v1.UnimplementedDockerServer
}

// Start a docker container based on the name
func (dm *DockerManager) Start(ctx context.Context, req *v1.DockerRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.docker.Start")
	defer span.End()

	resp := make(chan response, 1)
	dockerManage := &docker.Docker{Name: req.Name, Logger: dm.Logger}

	go func() {
		resp <- startTarget(dockerManage)
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(dm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

// Stop a docker container based on the name
func (dm *DockerManager) Stop(ctx context.Context, req *v1.DockerRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.docker.Stop")
	defer span.End()

	resp := make(chan response, 1)

	dockerManage := &docker.Docker{Name: req.Name, Logger: dm.Logger}

	go func() {
		resp <- stopTarget(dockerManage)
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(dm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

func startTarget(target common.Target) response {
	message, err := target.Start()

	return response{
		message: message,
		err:     err,
	}
}

func stopTarget(target common.Target) response {
	message, err := target.Stop()

	return response{
		message: message,
		err:     err,
	}
}

func prepareResponse(message string, err error) (*v1.StatusResponse, error) {
	if err != nil {
		return respFail(message, err)
	}

	return respSuccess(message)
}

func respSuccess(message string) (*v1.StatusResponse, error) {
	return &v1.StatusResponse{
		Status:  v1.StatusResponse_SUCCESS,
		Message: message,
	}, nil
}

func respFail(message string, err error) (*v1.StatusResponse, error) {
	return &v1.StatusResponse{
		Status:  v1.StatusResponse_FAIL,
		Message: message,
	}, err
}
