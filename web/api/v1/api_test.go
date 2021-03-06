package v1

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/SotirisAlfonsos/chaos-bot/common/cpu"

	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/SotirisAlfonsos/chaos-master/chaoslogger"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

var (
	logger = getLogger()
)

func TestHealthCheckService_Check(t *testing.T) {
	hcs := &HealthCheckService{}
	resp, err := hcs.Check(context.TODO(), &v1.HealthCheckRequest{})

	if err != nil {
		t.Fatalf("Error on Health Check request. err=%s", err)
	}

	assert.Equal(t, v1.HealthCheckResponse_SERVING, resp.Status)
}

type TestServer struct {
}

func (ts *TestServer) StopUnix() (string, error) {
	return "success", nil
}

func TestStopServer(t *testing.T) {
	testServer := &TestServer{}
	serverHandler := &ServerManager{
		Server: testServer,
		Logger: logger,
	}

	statusResponse, err := serverHandler.Stop(context.TODO(), &v1.ServerRequest{})
	assert.Nil(t, err)
	assert.Equal(t, v1.StatusResponse_SUCCESS, statusResponse.Status)
}

// === End to end testing ===
func TestServiceManager_e2e(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	serviceName := "simple"
	hostname, _ := os.Hostname()

	sm := &ServiceManager{Logger: logger}

	startService(sm, serviceName, t, hostname)
	startServiceFail(sm, serviceName, t)
	stopService(sm, serviceName, t, hostname)
}

func startService(sm *ServiceManager, serviceName string, t *testing.T, hostname string) {
	resp, err := sm.Start(context.TODO(), &v1.ServiceRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service Start request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s started service %s", hostname, serviceName)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func startServiceFail(sm *ServiceManager, serviceName string, t *testing.T) {
	resp, err := sm.Start(context.TODO(), &v1.ServiceRequest{Name: serviceName})

	expectedMessage := fmt.Sprintf("Could not start service %s", serviceName)

	assert.Error(t, err)
	assert.Equal(t, v1.StatusResponse_FAIL, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func stopService(sm *ServiceManager, serviceName string, t *testing.T, hostname string) {
	resp, err := sm.Stop(context.TODO(), &v1.ServiceRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service Stop request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s stopped service %s", hostname, serviceName)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

// === End to end testing ===
func TestDockerManager_e2e(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	dockerName := "zookeeper"
	hostname, _ := os.Hostname()

	dm := &DockerManager{Logger: logger}

	startDocker(dm, dockerName, t, hostname)
	stopDocker(dm, dockerName, t, hostname)
}

func startDocker(dm *DockerManager, dockerName string, t *testing.T, hostname string) {
	resp, err := dm.Start(context.TODO(), &v1.DockerRequest{Name: dockerName})
	if err != nil {
		t.Fatalf("Error in Docker Start request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s started docker container %s", hostname, dockerName)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func stopDocker(dm *DockerManager, dockerName string, t *testing.T, hostname string) {
	resp, err := dm.Stop(context.TODO(), &v1.DockerRequest{Name: dockerName})
	if err != nil {
		t.Fatalf("Error in Docker Stop request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s stopped docker container %s", hostname, dockerName)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func TestCPUManager_Start_Stop(t *testing.T) {
	hostname, _ := os.Hostname()

	cm := &CPUManager{CPU: cpu.New(logger), Logger: logger}

	startCPU(cm, int32(100), t, hostname)
	startCPUFailure(cm, int32(100), t)
	stopCPU(cm, t, hostname)
}

func startCPU(cm *CPUManager, percentage int32, t *testing.T, hostname string) {
	resp, err := cm.Start(context.TODO(), &v1.CPURequest{Percentage: percentage})
	if err != nil {
		t.Fatalf("Error in CPU Start request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s started cpu injection at %d%%", hostname, percentage)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func startCPUFailure(cm *CPUManager, percentage int32, t *testing.T) {
	resp, err := cm.Start(context.TODO(), &v1.CPURequest{Percentage: percentage})

	expectedMessage := "Could not inject cpu failure"

	assert.NotNil(t, err)
	assert.Equal(t, "CPU injection already running. Stop it before starting another", err.Error())
	assert.Equal(t, v1.StatusResponse_FAIL, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func stopCPU(cm *CPUManager, t *testing.T, hostname string) {
	resp, err := cm.Stop(context.TODO(), &v1.CPURequest{})
	if err != nil {
		t.Fatalf("Error in CPU Stop request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s stopped cpu injection", hostname)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

// === End to end testing ===
func TestNetworkManager_e2e(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	device := "lo"
	hostname, _ := os.Hostname()

	nm := NewNetworkManager(logger)

	startNetwork(t, nm, device, hostname)
	stopNetwork(t, nm, device, hostname)
}

func startNetwork(t *testing.T, nm *NetworkManager, device string, hostname string) {
	networkRequest := &v1.NetworkRequest{
		Device:  device,
		Latency: 10,
	}

	resp, err := nm.Start(context.TODO(), networkRequest)
	if err != nil {
		t.Fatalf("Error in Network Start request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s started network injection with details {Latency: 156250, Limit: 1000, Loss: 0, Gap: 0, Duplicate: 0, Jitter: 0}", hostname)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func stopNetwork(t *testing.T, nm *NetworkManager, device string, hostname string) {
	networkRequest := &v1.NetworkRequest{
		Device: device,
	}

	resp, err := nm.Stop(context.TODO(), networkRequest)
	if err != nil {
		t.Fatalf("Error in Network Stop request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s stopped network injection", hostname)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}

	return chaoslogger.New(allowLevel)
}
