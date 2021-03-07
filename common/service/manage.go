package service

import (
	"fmt"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/takama/daemon"
)

// Service is the interface implementation that manages chaos on services
type Service struct {
	Logger log.Logger
}

// Start will perform a service start on the service specified
func (s *Service) Start(serviceName string) (string, error) {
	dmn, err := daemon.New(serviceName, "", daemon.SystemDaemon)
	if err != nil {
		_ = level.Error(s.Logger).Log("msg", "Could not instantiate daemon", "err", err)
		return "Could not instantiate daemon", status.Error(codes.Internal, err.Error())
	}

	_, err = dmn.Start()
	if err != nil && daemon.ErrAlreadyRunning != err {
		_ = level.Error(s.Logger).Log("msg", fmt.Sprintf("Could not start service {%s}", serviceName), "err", err)
		return fmt.Sprintf("Could not start service {%s}", serviceName), status.Error(codes.Internal, err.Error())
	}

	_ = level.Info(s.Logger).Log("msg", fmt.Sprintf("Started service with name %s", serviceName))

	return constructMessage(s.Logger, "started", serviceName), nil
}

// Stop will perform a service stop on the service specified
func (s *Service) Stop(serviceName string) (string, error) {
	dmn, err := daemon.New(serviceName, "", daemon.SystemDaemon)
	if err != nil {
		_ = level.Error(s.Logger).Log("msg", "Could not instantiate daemon", "err", err)
		return "Could not instantiate daemon", status.Error(codes.Internal, err.Error())
	}

	_, err = dmn.Stop()
	if err != nil && daemon.ErrAlreadyStopped != err {
		_ = level.Error(s.Logger).Log("msg", fmt.Sprintf("Could not stop service {%s}", serviceName), "err", err)
		return fmt.Sprintf("Could not stop service {%s}", serviceName), status.Error(codes.Internal, err.Error())
	}

	_ = level.Info(s.Logger).Log("msg", fmt.Sprintf("Stopped service with name %s", serviceName))

	return constructMessage(s.Logger, "stopped", serviceName), nil
}

func constructMessage(logger log.Logger, action string, name string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}

	return fmt.Sprintf("Bot %s %s service %s", hostname, action, name)
}
