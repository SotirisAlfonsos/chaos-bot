package web

import (
	"chaos-slave/chaoslogger"
	"chaos-slave/proto"
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
	"time"
)

var (
	logger = getLogger()
)

func TestHealthCheckGrpcSuccess(t *testing.T) {
	done := make(chan struct{})
	errOnGrpcHandlerRun := make(chan error)
	port := "8080"
	withTestGrpcServer(port, done, errOnGrpcHandlerRun)
	resp, err := withTestGrpcClient(port, done)

	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Nil(t, <-errOnGrpcHandlerRun)
	assert.Equal(t, proto.HealthCheckResponse_SERVING, resp.Status)
}

func TestHealthCheckGrpcInvalidPort(t *testing.T) {
	done := make(chan struct{})
	errOnGrpcHandlerRun := make(chan error)
	port := "-1"
	withTestGrpcServer(port, done, errOnGrpcHandlerRun)
	resp, err := withTestGrpcClient(port, done)

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.EqualError(t, <-errOnGrpcHandlerRun, fmt.Sprintf("listen tcp: address %s: invalid port", port))
}

func withTestGrpcServer(port string, done chan struct{}, errOnGrpcHandlerRun chan error) {
	grpcHandler := NewGrpcHandler(port, logger)
	go func(errOnGrpcHandlerRun chan error) {
		if err := grpcHandler.Run(); err != nil {
			errOnGrpcHandlerRun <- err
		}
		errOnGrpcHandlerRun <- nil
	}(errOnGrpcHandlerRun)
	time.Sleep(2 * time.Second)

	go func(done chan struct{}) {
		select {
		case <-done:
			grpcHandler.Stop()
		}
	}(done)
}

func withTestGrpcClient(port string, done chan struct{}) (*proto.HealthCheckResponse, error) {
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%s", port), grpc.WithInsecure())
	if err != nil {
		level.Error(logger).Log("msg", "failed to dial server address", "err", err)
		return nil, err
	}
	defer conn.Close()

	client := proto.NewHealthClient(conn)
	resp, err := client.Check(context.Background(), &proto.HealthCheckRequest{})
	if err != nil {
		level.Error(logger).Log("msg", "Failed to call the Check method on the Health-check", "err", err)
		return nil, err
	}
	done <- struct{}{}
	return resp, err
}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}
	return chaoslogger.New(allowLevel)
}
