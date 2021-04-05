// +build integration

package service

import (
	"fmt"
	"os"
	"testing"

	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/SotirisAlfonsos/chaos-master/chaoslogger"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func Test_Service_Success_Recover(t *testing.T) {
	dataItems := []dataItem{
		{
			message:        "Should recover service with existing name",
			serviceRequest: &v1.ServiceRequest{Name: serviceName},
			expectedResult: &expected{
				message: fmt.Sprintf("Bot %s %s service %s", getHostname(t), "recovered", serviceName),
				status:  codes.OK,
			},
		},
	}

	for _, dataItem := range dataItems {
		t.Run(dataItem.message, func(t *testing.T) {
			serviceManager := &Service{Logger: logger}

			resp, err := serviceManager.Recover(dataItem.serviceRequest.Name)
			defer cleanUp(t, serviceManager, dataItem.serviceRequest.Name)

			assert.Nil(t, err)
			assert.Equal(t, dataItem.expectedResult.message, resp)
		})
	}
}

func Test_Service_Failure_Recover(t *testing.T) {
	dataItems := []dataItem{
		{
			message:        "Should fail to recover service with non existing name",
			serviceRequest: &v1.ServiceRequest{Name: "non existing service name"},
			expectedResult: &expected{
				message: "Could not recover service {non existing service name}",
				status:  codes.Internal,
			},
		},
		{
			message:        "Should fail to recover service that has already recovered",
			serviceRequest: &v1.ServiceRequest{},
			expectedResult: &expected{
				message: "Could not recover service {}",
				status:  codes.Internal,
			},
		},
	}

	for _, dataItem := range dataItems {
		runTest(t, dataItem, "recover")
	}
}

func TestShouldSucceedWhenRecoveringServiceThatIsAlreadyrecovered(t *testing.T) {
	serviceManager := &Service{Logger: logger}

	_, err := serviceManager.Recover(serviceName)
	if err != nil {
		t.Fatalf("should be able to recover service %s. Got err=%s", serviceName, err.Error())
	}
	resp, err := serviceManager.Recover(serviceName)
	defer cleanUp(t, serviceManager, serviceName)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("Bot %s recovered service %s", getHostname(t), serviceName), resp)
}

func Test_Service_Success_Kill(t *testing.T) {
	dataItems := []dataItem{
		{
			message:        "Should kill service with existing name",
			serviceRequest: &v1.ServiceRequest{Name: serviceName},
			expectedResult: &expected{
				message: fmt.Sprintf("Bot %s %s service %s", getHostname(t), "killed", serviceName),
				status:  codes.OK,
			},
		},
	}

	for _, dataItem := range dataItems {
		t.Run(dataItem.message, func(t *testing.T) {
			serviceManager := &Service{Logger: logger}

			setUp(t, serviceManager, dataItem.serviceRequest.Name)

			resp, err := serviceManager.Kill(dataItem.serviceRequest.Name)

			assert.Nil(t, err)
			assert.Equal(t, dataItem.expectedResult.message, resp)
		})
	}
}

func Test_Service_Failure_Kill(t *testing.T) {
	dataItems := []dataItem{
		{
			message:        "Should fail to kill service with non existing  name",
			serviceRequest: &v1.ServiceRequest{Name: "non existing container name"},
			expectedResult: &expected{
				message: "Could not kill service {non existing container name}",
				status:  codes.Internal,
			},
		},
		{
			message:        "Should fail to kill service that has already recovered",
			serviceRequest: &v1.ServiceRequest{},
			expectedResult: &expected{
				message: "Could not kill service {}",
				status:  codes.Internal,
			},
		},
	}

	for _, dataItem := range dataItems {
		runTest(t, dataItem, "kill")
	}
}

func runTest(t *testing.T, dataItem dataItem, action string) {
	t.Run(dataItem.message, func(t *testing.T) {
		serviceManager := &Service{Logger: logger}

		var resp string
		var err error

		switch {
		case action == "recover":
			resp, err = serviceManager.Recover(dataItem.serviceRequest.Name)
		case action == "kill":
			resp, err = serviceManager.Kill(dataItem.serviceRequest.Name)
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
	_, err := d.Recover(name)
	if err != nil {
		t.Fatalf("Could not recover service %s", name)
	}

	_ = level.Info(logger).Log("msg", "setup operation. Recover service")
}

func cleanUp(t *testing.T, d *Service, name string) {
	_, err := d.Kill(name)
	if err != nil {
		t.Fatalf("Could not kill service %s", name)
	}

	_ = level.Info(logger).Log("msg", "cleanup operation. Kill service")
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
