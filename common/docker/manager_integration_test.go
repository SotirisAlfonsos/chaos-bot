// +build integration

package docker

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
	logger        = getLogger()
	containerName = "zookeeper"
)

type dataItem struct {
	message        string
	dockerRequest  *v1.DockerRequest
	expectedResult *expected
}

type expected struct {
	message string
	status  codes.Code
}

func Test_Docker_Success_Recover(t *testing.T) {
	dataItems := []dataItem{
		{
			message:       "Should recover docker with existing name",
			dockerRequest: &v1.DockerRequest{Name: containerName},
			expectedResult: &expected{
				message: fmt.Sprintf("Bot %s %s docker container %s", getHostname(t), "recovered", containerName),
				status:  codes.OK,
			},
		},
	}

	for _, dataItem := range dataItems {
		t.Run(dataItem.message, func(t *testing.T) {
			dockerManage := &Docker{Logger: logger}

			resp, err := dockerManage.Recover(dataItem.dockerRequest.Name)
			defer cleanUp(t, dockerManage, dataItem.dockerRequest.Name)

			assert.Nil(t, err)
			assert.Equal(t, dataItem.expectedResult.message, resp)
		})
	}
}

func Test_Docker_Failure_Recover(t *testing.T) {
	dataItems := []dataItem{
		{
			message:       "Should fail to recover docker with non existing name",
			dockerRequest: &v1.DockerRequest{Name: "non existing container name"},
			expectedResult: &expected{
				message: "Could not recover docker container {non existing container name}",
				status:  codes.Internal,
			},
		},
		{
			message:       "Should fail to recover docker that has already recovered",
			dockerRequest: &v1.DockerRequest{},
			expectedResult: &expected{
				message: "Could not recover docker container {}",
				status:  codes.Internal,
			},
		},
	}

	for _, dataItem := range dataItems {
		runTest(t, dataItem, "recover")
	}
}

func TestShouldSucceedWhenRecoveringDockerThatIsAlreadyRecovered(t *testing.T) {
	dockerManage := &Docker{Logger: logger}

	_, err := dockerManage.Recover(containerName)
	if err != nil {
		t.Fatalf("should be able to recover container %s. Got err=%s", containerName, err.Error())
	}
	resp, err := dockerManage.Recover(containerName)
	defer cleanUp(t, dockerManage, containerName)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("Bot %s recovered docker container %s", getHostname(t), containerName), resp)
}

func Test_Docker_Success_Kill(t *testing.T) {
	dataItems := []dataItem{
		{
			message:       "Should kill docker with existing name",
			dockerRequest: &v1.DockerRequest{Name: containerName},
			expectedResult: &expected{
				message: fmt.Sprintf("Bot %s %s docker container %s", getHostname(t), "killed", containerName),
				status:  codes.OK,
			},
		},
	}

	for _, dataItem := range dataItems {
		t.Run(dataItem.message, func(t *testing.T) {
			dockerManage := &Docker{Logger: logger}

			setUp(t, dockerManage, dataItem.dockerRequest.Name)

			resp, err := dockerManage.Kill(dataItem.dockerRequest.Name)

			assert.Nil(t, err)
			assert.Equal(t, dataItem.expectedResult.message, resp)
		})
	}
}

func Test_Docker_Failure_Kill(t *testing.T) {
	dataItems := []dataItem{
		{
			message:       "Should fail to kill docker with non existing  name",
			dockerRequest: &v1.DockerRequest{Name: "non existing container name"},
			expectedResult: &expected{
				message: "Could not kill docker container {non existing container name}",
				status:  codes.Internal,
			},
		},
		{
			message:       "Should fail to kill docker that has already recovered",
			dockerRequest: &v1.DockerRequest{},
			expectedResult: &expected{
				message: "Could not kill docker container {}",
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
		dockerManage := &Docker{Logger: logger}

		var resp string
		var err error

		switch {
		case action == "recover":
			resp, err = dockerManage.Recover(dataItem.dockerRequest.Name)
		case action == "kill":
			resp, err = dockerManage.Kill(dataItem.dockerRequest.Name)
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

func setUp(t *testing.T, d *Docker, container string) {
	_, err := d.Recover(container)
	if err != nil {
		t.Fatalf("Could not recover docker %s", container)
	}

	_ = level.Info(logger).Log("msg", "setup operation. Recover docker")
}

func cleanUp(t *testing.T, d *Docker, container string) {
	_, err := d.Kill(container)
	if err != nil {
		t.Fatalf("Could not kill docker %s", container)
	}

	_ = level.Info(logger).Log("msg", "cleanup operation. Kill docker")
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
