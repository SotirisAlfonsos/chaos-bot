package v1

import (
	"chaos-slave/chaoslogger"
	"chaos-slave/proto"
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
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
func TestServiceManager_Start(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	serviceName := "simple"
	hostname, _ := os.Hostname()
	expectedMessage := fmt.Sprintf("Slave %s started service %s", hostname, serviceName)

	sm := &ServiceManager{logger}
	resp, err := sm.Start(context.TODO(), &proto.ServiceRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service Start request. err=%s", err)
	}

	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

// === End to end testing ===
func TestServiceManager_Stop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	serviceName := "simple"
	hostname, _ := os.Hostname()
	expectedMessage := fmt.Sprintf("Slave %s stopped service %s", hostname, serviceName)

	sm := &ServiceManager{logger}
	resp, err := sm.Stop(context.TODO(), &proto.ServiceRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service Start request. err=%s", err)
	}

	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

// === End to end testing ===
func TestDockerManager_Start(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	serviceName := "dummy"
	hostname, _ := os.Hostname()
	expectedMessage := fmt.Sprintf("Slave %s started docker container %s", hostname, serviceName)

	dm := &DockerManager{logger}
	resp, err := dm.Start(context.TODO(), &proto.DockerRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service Start request. err=%s", err)
	}

	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

// === End to end testing ===
func TestDockerManager_Stop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	serviceName := "dummy"
	hostname, _ := os.Hostname()
	expectedMessage := fmt.Sprintf("Slave %s stopped docker container %s", hostname, serviceName)

	dm := &DockerManager{logger}
	resp, err := dm.Stop(context.TODO(), &proto.DockerRequest{Name: serviceName})
	if err != nil {
		t.Fatalf("Error in Service Start request. err=%s", err)
	}

	assert.Equal(t, proto.StatusResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedMessage, resp.Message)
}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}
	return chaoslogger.New(allowLevel)
}