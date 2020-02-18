package v1

import (
	"chaos-slave/chaoslogger"
	"chaos-slave/common/service"
	"chaos-slave/proto"
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	logger = getLogger()
)

func TestHealthCheckService_Check(t *testing.T) {
	hcs := &HealthCheckService{}
	resp, err := hcs.Check(context.TODO(), &proto.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Error on Health Check request. err=%s", err)
	}

	assert.Equal(t, proto.HealthCheckResponse_SERVING, resp.Status)
}

// === End to end testing ===
func TestServiceManager_e2e(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	myCache := cache.New(0, 0)
	serviceName := "simple"
	hostname, _ := os.Hostname()

	sm := &ServiceManager{myCache, &service.Service{Logger: logger}}
	stratM := &StrategyManager{logger, myCache}
	startService(sm, serviceName, t, hostname)
	startServiceFail(sm, serviceName, t, hostname)
	recoverServiceEmpty(stratM, serviceName, t)
	stopService(sm, serviceName, t, hostname)
	recoverService(stratM, serviceName, t, hostname)
	stopService(sm, serviceName, t, hostname)
}

func startService(sm *ServiceManager, serviceName string, t *testing.T, hostname string) {
	resp, err := sm.Start(context.TODO(), &proto.ServiceRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service Start request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Slave %s started service %s", hostname, serviceName)
	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
	_, ok := sm.Cache.Get(serviceName)
	assert.False(t, ok)
	assert.Equal(t, 0, sm.Cache.ItemCount())
}

func startServiceFail(sm *ServiceManager, serviceName string, t *testing.T, hostname string) {
	resp, err := sm.Start(context.TODO(), &proto.ServiceRequest{Name: serviceName})

	expectedMessage := fmt.Sprintf("Could not start service %s", serviceName)
	assert.Error(t, err)
	assert.Equal(t, proto.StatusResponse_FAIL, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
	_, ok := sm.Cache.Get(serviceName)
	assert.False(t, ok)
	assert.Equal(t, 0, sm.Cache.ItemCount())
}

func stopService(sm *ServiceManager, serviceName string, t *testing.T, hostname string) {
	resp, err := sm.Stop(context.TODO(), &proto.ServiceRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service Stop request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Slave %s stopped service %s", hostname, serviceName)
	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
	serviceObj, ok := sm.Cache.Get(serviceName)
	if !ok {
		t.Fatalf(fmt.Sprintf("Could not retrieve item %s from cache", serviceName))
	}
	assert.Equal(t, sm.Service, serviceObj)
}

func recoverService(sm *StrategyManager, serviceName string, t *testing.T, hostname string) {
	resp, err := sm.Recover(context.TODO(), &proto.RecoverRequest{})
	if err != nil {
		t.Fatalf("Error in Service Recover request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Slave %s started service %s", hostname, serviceName)
	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Response[0].Status)
	assert.Equal(t, expectedMessage, resp.Response[0].Message)
	_, ok := sm.Cache.Get(serviceName)
	assert.False(t, ok)
	assert.Equal(t, 0, sm.Cache.ItemCount())
}

func recoverServiceEmpty(sm *StrategyManager, serviceName string, t *testing.T) {
	resp, err := sm.Recover(context.TODO(), &proto.RecoverRequest{})
	if err != nil {
		t.Fatalf("Error in Service Recover request. err=%s", err)
	}

	assert.Equal(t, 0, len(resp.Response))
	_, ok := sm.Cache.Get(serviceName)
	assert.False(t, ok)
	assert.Equal(t, 0, sm.Cache.ItemCount())
}

// === End to end testing ===
//func TestDockerManager_e2e(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping testing in short mode")
//	}
//
//	serviceName := "dummy"
//	hostname, _ := os.Hostname()
//
//	dm := &DockerManager{logger}
//	startDocker(dm, serviceName, t, hostname)
//	stopDocker(dm, serviceName, t, hostname)
//	recoverDocker(dm, serviceName, t, hostname)
//	stopDocker(dm, serviceName, t, hostname)
//}
//
//func startDocker(dm *DockerManager, serviceName string, t *testing.T, hostname string) {
//	resp, err := dm.Start(context.TODO(), &proto.DockerRequest{Name: serviceName})
//	if err != nil {
//		t.Fatalf("Error in Docker Start request. err=%s", err)
//	}
//
//	expectedMessage := fmt.Sprintf("Slave %s started docker container %s", hostname, serviceName)
//	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Status)
//	assert.Equal(t, expectedMessage, resp.Message)
//}
//
//func stopDocker(dm *DockerManager, serviceName string, t *testing.T, hostname string) {
//	resp, err := dm.Stop(context.TODO(), &proto.DockerRequest{Name: serviceName})
//	if err != nil {
//		t.Fatalf("Error in Docker Stop request. err=%s", err)
//	}
//
//	expectedMessage := fmt.Sprintf("Slave %s stopped docker container %s", hostname, serviceName)
//	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Status)
//	assert.Equal(t, expectedMessage, resp.Message)
//}
//
//func recoverDocker(dm *DockerManager, serviceName string, t *testing.T, hostname string) {
//	resp, err := dm.Recover(context.TODO(), &proto.DockerRequest{Name: serviceName})
//	if err != nil {
//		t.Fatalf("Error in Service Start request. err=%s", err)
//	}
//
//	expectedMessage := fmt.Sprintf("Slave %s started docker container %s", hostname, serviceName)
//	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Status)
//	assert.Equal(t, expectedMessage, resp.Message)
//}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}
	return chaoslogger.New(allowLevel)
}
