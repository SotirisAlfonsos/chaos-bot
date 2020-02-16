package docker

import (
	"chaos-slave/proto"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

type Docker struct {
	Logger     log.Logger
	Name string
}

func (s *Docker) Start() (*proto.StatusResponse, error) {
	return &proto.StatusResponse{
		Status: proto.StatusResponse_SUCCESS,
		Message: s.constructMessage("started"),
	}, nil
}

func (s *Docker) Stop() (*proto.StatusResponse, error) {
	return &proto.StatusResponse{
		Status: proto.StatusResponse_SUCCESS,
		Message: s.constructMessage("stopped"),
	}, nil}

func (s *Docker) constructMessage(action string) string {
	hostname, err := os.Hostname()
	if err != nil {
		level.Warn(s.Logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}

	return fmt.Sprintf("Slave %s %s docker container %s", hostname, action, s.Name)
}