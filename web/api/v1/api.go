package v1

import (
	"context"
	"fmt"

	"github.com/SotirisAlfonsos/chaos-bot/common/server"

	"github.com/SotirisAlfonsos/chaos-bot/common"
	"github.com/SotirisAlfonsos/chaos-bot/common/cpu"
	"github.com/SotirisAlfonsos/chaos-bot/common/docker"
	"github.com/SotirisAlfonsos/chaos-bot/common/service"
	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/patrickmn/go-cache"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HealthCheckService is the rpc
type HealthCheckService struct {
	*v1.UnimplementedHealthServer
}

// Check the health of the chaos bot
func (hcs *HealthCheckService) Check(ctx context.Context,
	req *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	return &v1.HealthCheckResponse{Status: v1.HealthCheckResponse_SERVING}, nil
}

// Watch is not used at the moment
func (hcs *HealthCheckService) Watch(req *v1.HealthCheckRequest, srv v1.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

// ServiceManager is the rpc for services management
type ServiceManager struct {
	Cache  *cache.Cache
	Logger log.Logger
	*v1.UnimplementedServiceServer
}

type response struct {
	message string
	err     error
}

// Start a service based on the name. Delete the item from the cache if it had been cached previously
func (sm *ServiceManager) Start(ctx context.Context, req *v1.ServiceRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.service.Start")
	defer span.End()

	resp := make(chan response, 1)

	serviceManage := &service.Service{JobName: req.JobName, Name: req.Name, Logger: sm.Logger}

	go func() {
		resp <- startTarget(serviceManage, sm.Cache, req.Name)
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

// Stop a service based on the name. Cache it if the service is stopped successfully
func (sm *ServiceManager) Stop(ctx context.Context, req *v1.ServiceRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.service.Stop")
	defer span.End()

	resp := make(chan response, 1)

	serviceManage := &service.Service{JobName: req.JobName, Name: req.Name, Logger: sm.Logger}

	go func() {
		resp <- stopTarget(serviceManage, sm.Cache, req.Name, sm.Logger)
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
	Cache  *cache.Cache
	Logger log.Logger
	*v1.UnimplementedDockerServer
}

// Start a docker container based on the name. Delete the item from the cache if it had been cached previously
func (dm *DockerManager) Start(ctx context.Context, req *v1.DockerRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.docker.Start")
	defer span.End()

	resp := make(chan response, 1)
	dockerManage := &docker.Docker{JobName: req.JobName, Name: req.Name, Logger: dm.Logger}

	go func() {
		resp <- startTarget(dockerManage, dm.Cache, req.Name)
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

// Stop a docker container based on the name. Cache it if the docker container is stopped successfully
func (dm *DockerManager) Stop(ctx context.Context, req *v1.DockerRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.docker.Stop")
	defer span.End()

	resp := make(chan response, 1)

	dockerManage := &docker.Docker{JobName: req.JobName, Name: req.Name, Logger: dm.Logger}

	go func() {
		resp <- stopTarget(dockerManage, dm.Cache, req.Name, dm.Logger)
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

func startTarget(target common.Target, cache *cache.Cache, name string) response {
	message, err := target.Start()
	if err == nil {
		cache.Delete(name)
	}

	return response{
		message: message,
		err:     err,
	}
}

func stopTarget(target common.Target, cache *cache.Cache, name string, logger log.Logger) response {
	message, err := target.Stop()
	if err == nil {
		if cacheErr := cache.Add(name, target, 0); cacheErr != nil {
			_ = level.Error(logger).Log("msg",
				fmt.Sprintf("Could not update cache after stopping target %s", name), "err", cacheErr)
		}
	}

	return response{
		message: message,
		err:     err,
	}
}

type CPUManager struct {
	CPU    *cpu.CPU
	Logger log.Logger
	*v1.UnimplementedCPUServer
}

func (cm *CPUManager) Start(ctx context.Context, req *v1.CPURequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.cpu.Start")
	defer span.End()

	resp := make(chan response, 1)

	go func() {
		resp <- cm.startCPU(req)
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(cm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

func (cm *CPUManager) Stop(ctx context.Context, req *v1.CPURequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.cpu.Stop")
	defer span.End()

	resp := make(chan response, 1)

	go func() {
		resp <- cm.stopCPU()
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(cm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

func (cm *CPUManager) startCPU(req *v1.CPURequest) response {
	message, err := cm.CPU.Start(int(req.Percentage))

	return response{
		message: message,
		err:     err,
	}
}

func (cm *CPUManager) stopCPU() response {
	message, err := cm.CPU.Stop()

	return response{
		message: message,
		err:     err,
	}
}

type ServerManager struct {
	Server server.Server
	Logger log.Logger
	*v1.UnimplementedServerServer
}

func (sm *ServerManager) Stop(ctx context.Context, req *v1.ServerRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.server.Stop")
	defer span.End()

	resp := make(chan response, 1)

	go func() {
		resp <- sm.stop()
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

func (sm *ServerManager) stop() response {
	message, err := sm.Server.StopUnix()

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
