package web

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/SotirisAlfonsos/chaos-bot/config"
	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/SotirisAlfonsos/chaos-master/chaoslogger"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	logger = getLogger()
)

func TestHealthCheckGRPCCheckSuccess(t *testing.T) {
	const port = "8080"

	done := make(chan struct{})

	withTestGRPCServer(t, port, done)

	clientConn, err := withTestGRPCClientConn(port)
	if err != nil {
		t.Fatalf("Can not create client connection")
	}

	client := v1.NewHealthClient(clientConn)
	resp := performHChReqOnCheck(client)

	done <- struct{}{}

	err = clientConn.Close()
	if err != nil {
		t.Fatalf("Can not close client connection")
	}

	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, v1.HealthCheckResponse_SERVING, resp.Status)
}

// === End to end testing ===
func TestStartServiceGRPCSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	const port = "8080"

	done := make(chan struct{})

	withTestGRPCServer(t, port, done)

	clientConn, err := withTestGRPCClientConn(port)
	if err != nil {
		t.Fatalf("Can not create client connection")
	}

	client := v1.NewServiceClient(clientConn)
	resp := performStartServiceReq(client)

	done <- struct{}{}

	err = clientConn.Close()
	if err != nil {
		t.Fatalf("Can not close client connection")
	}

	assert.Nil(t, err)
	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.NotNil(t, resp.Message)
}

// === End to end testing ===
func TestStopServiceGRPCSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	const port = "8080"

	done := make(chan struct{})

	withTestGRPCServer(t, port, done)

	clientConn, err := withTestGRPCClientConn(port)
	if err != nil {
		t.Fatalf("Can not create client connection")
	}

	client := v1.NewServiceClient(clientConn)
	resp := performStopServiceReq(client)

	done <- struct{}{}

	err = clientConn.Close()
	if err != nil {
		t.Fatalf("Can not close client connection")
	}

	assert.Nil(t, err)
	assert.Equal(t, v1.StatusResponse_SUCCESS, resp.Status)
	assert.NotNil(t, resp.Message)
}

func withTestGRPCServer(t *testing.T, port string, done chan struct{}) {
	GRPCHandler, err := NewGRPCHandler(port, logger, &config.Config{})
	if err != nil {
		t.Fatalf("Could not create grpc handler")
	}

	go func() {
		GRPCHandler.Run()
	}()
	time.Sleep(2 * time.Second)

	go func(done chan struct{}) {
		<-done
		GRPCHandler.Stop()
	}(done)
}

func withTestGRPCClientConn(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%s", port), grpc.WithInsecure())
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed to dial server address", "err", err)
		return nil, err
	}

	return conn, nil
}

func performHChReqOnCheck(client v1.HealthClient) *v1.HealthCheckResponse {
	resp, err := client.Check(context.Background(), &v1.HealthCheckRequest{})
	if err != nil {
		_ = level.Error(logger).Log("msg", "Failed to call the Check method on the Health-check", "err", err)

		return nil
	}

	return resp
}

func performStartServiceReq(client v1.ServiceClient) *v1.StatusResponse {
	resp, err := client.Start(context.Background(), &v1.ServiceRequest{Name: "simple"})
	if err != nil {
		_ = level.Error(logger).Log("msg", "Failed to call the Start method on the Service", "err", err)
		return resp
	}

	return resp
}

func performStopServiceReq(client v1.ServiceClient) *v1.StatusResponse {
	resp, err := client.Stop(context.Background(), &v1.ServiceRequest{Name: "simple"})
	if err != nil {
		_ = level.Error(logger).Log("msg", "Failed to call the Start method on the Service", "err", err)
		return resp
	}

	return resp
}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}

	return chaoslogger.New(allowLevel)
}
