package v1

import (
	"context"
	"fmt"
	"testing"

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

	statusResponse, err := serverHandler.Kill(context.TODO(), &v1.ServerRequest{})
	assert.Nil(t, err)
	assert.Equal(t, v1.StatusResponse_SUCCESS, statusResponse.Status)
}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}

	return chaoslogger.New(allowLevel)
}
