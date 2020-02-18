package v1

import (
	"chaos-slave/common"
	"chaos-slave/common/docker"
	"chaos-slave/common/service"
	"chaos-slave/proto"
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthCheckService struct {
}

//Check the health of the slave
func (hcs *HealthCheckService) Check(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{Status: proto.HealthCheckResponse_SERVING}, nil
}

//Watch is not used at the moment
func (hcs *HealthCheckService) Watch(req *proto.HealthCheckRequest, srv proto.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

type ServiceManager struct {
	Cache   *cache.Cache
	Service *service.Service
}

func (sm *ServiceManager) Start(ctx context.Context, req *proto.ServiceRequest) (*proto.StatusResponse, error) {
	message, err := sm.Service.Start(req.Name)
	if err == nil {
		sm.Cache.Delete(req.Name)
	}

	return prepareResponse(message, err)
}

func (sm *ServiceManager) Stop(ctx context.Context, req *proto.ServiceRequest) (*proto.StatusResponse, error) {
	message, err := sm.Service.Stop(req.Name)
	if err == nil {
		if cacheErr := sm.Cache.Add(req.Name, sm.Service, 0); cacheErr != nil {
			_ = level.Error(sm.Service.Logger).Log("msg", "Could not update cache after stopping service", "err", cacheErr)
		}
	}

	return prepareResponse(message, err)
}

type DockerManager struct {
	Cache  *cache.Cache
	Docker *docker.Docker
}

func (sm *DockerManager) Start(ctx context.Context, req *proto.DockerRequest) (*proto.StatusResponse, error) {
	return prepareResponse(sm.Docker.Start(req.Name))
}

func (sm *DockerManager) Stop(ctx context.Context, req *proto.DockerRequest) (*proto.StatusResponse, error) {
	return prepareResponse(sm.Docker.Stop(req.Name))
}

type StrategyManager struct {
	Logger log.Logger
	Cache  *cache.Cache
}

func (sm *StrategyManager) Recover(ctx context.Context, req *proto.RecoverRequest) (*proto.ResolveResponse, error) {
	var responses []*proto.StatusResponse

	var err error = nil

	for item := range sm.Cache.Items() {
		target, ok := sm.Cache.Get(item)
		if !ok {
			_ = level.Error(sm.Logger).Log("err", fmt.Sprintf("Could not find item %s in cache", item))
		}

		message, startErr := target.(common.Target).Start(item)
		if startErr == nil {
			sm.Cache.Delete(item)
			_ = level.Info(sm.Logger).Log("err", fmt.Sprintf("Started and removed item %s from cache", item))
		} else {
			err = errors.Wrap(err, startErr.Error())
		}

		resp, respErr := prepareResponse(message, err)
		if respErr != nil {
			err = errors.Wrap(err, respErr.Error())
		}

		responses = append(responses, resp)
	}

	return &proto.ResolveResponse{Response: responses}, err
}

func prepareResponse(message string, err error) (*proto.StatusResponse, error) {
	if err != nil {
		return respFail(message, err)
	}

	return respSuccess(message)
}

func respSuccess(message string) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{
		Status:  proto.StatusResponse_SUCCESS,
		Message: message,
	}, nil
}

func respFail(message string, err error) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{
		Status:  proto.StatusResponse_FAIL,
		Message: message,
	}, err
}
