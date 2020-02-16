package service

import (
	"chaos-slave/proto"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/takama/daemon"
	"os"
)

type Service struct {
	Logger     log.Logger
	Name string
}

func (s *Service) Start() (*proto.StatusResponse, error) {
	dmn, err := daemon.New(s.Name, "")
	if err != nil {
		return respFail("Could not instantiate daemon", err)
	}

	_, err = dmn.Start()
	if err != nil {
		return respFail(fmt.Sprintf("Could not start service %s",s.Name), err)
	}

	return respSuccess(s.constructMessage("started"))
}

func (s *Service) Stop() (*proto.StatusResponse, error) {
	dmn, err := daemon.New(s.Name, "")
	if err != nil {
		return respFail("Could not instantiate daemon", err)
	}

	_, err = dmn.Stop()
	if err != nil {
		return respFail(fmt.Sprintf("Could not stop service %s",s.Name), err)
	}

	return respSuccess(s.constructMessage("stopped"))
}

func (s *Service) constructMessage(action string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(s.Logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}

	return fmt.Sprintf("Slave %s %s service %s", hostname, action, s.Name)
}

func respFail(message string, err error) (*proto.StatusResponse, error){
	return &proto.StatusResponse{
		Status: proto.StatusResponse_FAIL,
		Message: message,
	}, err
}

func respSuccess(message string) (*proto.StatusResponse, error){
	return &proto.StatusResponse{
		Status: proto.StatusResponse_SUCCESS,
		Message: message,
	}, nil
}
