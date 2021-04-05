// +build integration

package v1

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/SotirisAlfonsos/chaos-bot/common/cpu"

	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/stretchr/testify/assert"
)

// === End to end testing ===
func TestServiceManager_e2e(t *testing.T) {
	serviceName := "simple"
	hostname, _ := os.Hostname()

	sm := NewServiceHandler(logger)

	recoverService(sm, serviceName, t, hostname)
	recoverService(sm, serviceName, t, hostname)
	killService(sm, serviceName, t, hostname)
}

func recoverService(sm *ServiceHandler, serviceName string, t *testing.T, hostname string) {
	resp, err := sm.Recover(context.TODO(), &v1.ServiceRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service recover request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s recovered service %s", hostname, serviceName)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func killService(sm *ServiceHandler, serviceName string, t *testing.T, hostname string) {
	resp, err := sm.Kill(context.TODO(), &v1.ServiceRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service kill request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s killed service %s", hostname, serviceName)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

// === End to end testing ===
func TestDockerManager_e2e(t *testing.T) {
	dockerName := "zookeeper"
	hostname, _ := os.Hostname()

	dm := &DockerHandler{Logger: logger}

	recoverDocker(dm, dockerName, t, hostname)
	killDocker(dm, dockerName, t, hostname)
}

func recoverDocker(dm *DockerHandler, dockerName string, t *testing.T, hostname string) {
	resp, err := dm.Recover(context.TODO(), &v1.DockerRequest{Name: dockerName})
	if err != nil {
		t.Fatalf("Error in Docker recover request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s recovered docker container %s", hostname, dockerName)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func killDocker(dm *DockerHandler, dockerName string, t *testing.T, hostname string) {
	resp, err := dm.Kill(context.TODO(), &v1.DockerRequest{Name: dockerName})
	if err != nil {
		t.Fatalf("Error in Docker kill request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s killed docker container %s", hostname, dockerName)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

// === End to end testing ===
func TestCPUManager_Start_Recover(t *testing.T) {
	hostname, _ := os.Hostname()

	cm := &CPUManager{CPU: cpu.New(logger), Logger: logger}

	startCPU(cm, int32(100), t, hostname)
	startCPUFailure(cm, int32(100), t)
	recoverCPU(cm, t, hostname)
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
	assert.Equal(t, "CPU injection already running. Recover before starting another", err.Error())
	assert.Equal(t, v1.StatusResponse_FAIL, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func recoverCPU(cm *CPUManager, t *testing.T, hostname string) {
	resp, err := cm.Recover(context.TODO(), &v1.CPURequest{})
	if err != nil {
		t.Fatalf("Error in CPU recover request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s recovered cpu injection", hostname)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

// === End to end testing ===
func TestNetworkManager_e2e(t *testing.T) {
	device := "lo"
	hostname, _ := os.Hostname()

	nm := NewNetworkManager(logger)

	startNetwork(t, nm, device, hostname)
	recoverNetwork(t, nm, device, hostname)
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

func recoverNetwork(t *testing.T, nm *NetworkManager, device string, hostname string) {
	networkRequest := &v1.NetworkRequest{
		Device: device,
	}

	resp, err := nm.Recover(context.TODO(), networkRequest)
	if err != nil {
		t.Fatalf("Error in Network recover request. err=%s", err)
	}

	expectedMessage := fmt.Sprintf("Bot %s recovered network injection", hostname)

	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}
