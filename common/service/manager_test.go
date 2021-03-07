package service

import (
	"fmt"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/SotirisAlfonsos/chaos-master/chaoslogger"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
)

var (
	logger      = getLogger()
	serviceName = "simple"
)

type dataItem struct {
	message        string
	serviceRequest *v1.ServiceRequest
	expectedResult *expected
}

type expected struct {
	message string
	status  codes.Code
}

func Test_Service_Success_Start(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	dataItems := []dataItem{
		{
			message:        "Should start service with existing name",
			serviceRequest: &v1.ServiceRequest{Name: serviceName},
			expectedResult: &expected{
				message: fmt.Sprintf("Bot %s %s service %s", getHostname(t), "started", serviceName),
				status:  codes.OK,
			},
		},
	}

	for _, dataItem := range dataItems {
		t.Run(dataItem.message, func(t *testing.T) {
			serviceManager := &Service{Logger: logger}

			resp, err := serviceManager.Start(dataItem.serviceRequest.Name)
			defer cleanUp(t, serviceManager, dataItem.serviceRequest.Name)

			assert.Nil(t, err)
			assert.Equal(t, dataItem.expectedResult.message, resp)
		})
	}
}

func Test_Service_Failure_Start(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	dataItems := []dataItem{
		{
			message:        "Should fail to start service with non existing name",
			serviceRequest: &v1.ServiceRequest{Name: "non existing service name"},
			expectedResult: &expected{
				message: "Could not start service {non existing service name}",
				status:  codes.Internal,
			},
		},
		{
			message:        "Should fail to start service that has already started",
			serviceRequest: &v1.ServiceRequest{},
			expectedResult: &expected{
				message: "Could not start service {}",
				status:  codes.Internal,
			},
		},
	}

	for _, dataItem := range dataItems {
		runTest(t, dataItem, "start")
	}
}

func TestShouldSucceedWhenStartingServiceThatIsAlreadyStarted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	serviceManager := &Service{Logger: logger}

	_, err := serviceManager.Start(serviceName)
	if err != nil {
		t.Fatalf("should be able to start service %s. Got err=%s", serviceName, err.Error())
	}
	resp, err := serviceManager.Start(serviceName)
	defer cleanUp(t, serviceManager, serviceName)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("Bot %s started service %s", getHostname(t), serviceName), resp)
}

func Test_Service_Success_Stop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	dataItems := []dataItem{
		{
			message:        "Should stop service with existing name",
			serviceRequest: &v1.ServiceRequest{Name: serviceName},
			expectedResult: &expected{
				message: fmt.Sprintf("Bot %s %s service %s", getHostname(t), "stopped", serviceName),
				status:  codes.OK,
			},
		},
	}

	for _, dataItem := range dataItems {
		t.Run(dataItem.message, func(t *testing.T) {
			serviceManager := &Service{Logger: logger}

			setUp(t, serviceManager, dataItem.serviceRequest.Name)

			resp, err := serviceManager.Stop(dataItem.serviceRequest.Name)

			assert.Nil(t, err)
			assert.Equal(t, dataItem.expectedResult.message, resp)
		})
	}
}

func Test_Service_Failure_Stop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	dataItems := []dataItem{
		{
			message:        "Should fail to stop service with non existing  name",
			serviceRequest: &v1.ServiceRequest{Name: "non existing container name"},
			expectedResult: &expected{
				message: "Could not stop service {non existing container name}",
				status:  codes.Internal,
			},
		},
		{
			message:        "Should fail to stop service that has already started",
			serviceRequest: &v1.ServiceRequest{},
			expectedResult: &expected{
				message: "Could not stop service {}",
				status:  codes.Internal,
			},
		},
	}

	for _, dataItem := range dataItems {
		runTest(t, dataItem, "stop")
	}
}

func runTest(t *testing.T, dataItem dataItem, action string) {
	t.Run(dataItem.message, func(t *testing.T) {
		serviceManager := &Service{Logger: logger}

		var resp string
		var err error

		switch {
		case action == "start":
			resp, err = serviceManager.Start(dataItem.serviceRequest.Name)
		case action == "stop":
			resp, err = serviceManager.Stop(dataItem.serviceRequest.Name)
		default:
			t.Fatal("no valid action")
		}

		assert.NotNil(t, err)
		s, ok := status.FromError(err)
		if !ok {
			t.Fatal("should be able to get the status from the error")
		}
		assert.Equal(t, dataItem.expectedResult.status, s.Code())
		assert.Regexp(t, dataItem.expectedResult.message, resp)
	})
}

func setUp(t *testing.T, d *Service, name string) {
	_, err := d.Start(name)
	if err != nil {
		t.Fatalf("Could not start service %s", name)
	}

	_ = level.Info(logger).Log("msg", "setup operation. Start service")
}

func cleanUp(t *testing.T, d *Service, name string) {
	_, err := d.Stop(name)
	if err != nil {
		t.Fatalf("Could not stop service %s", name)
	}

	_ = level.Info(logger).Log("msg", "cleanup operation. Stop service")
}

func getHostname(t *testing.T) string {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("Can not get hostname")
	}
	return hostname
}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}

	return chaoslogger.New(allowLevel)
}
