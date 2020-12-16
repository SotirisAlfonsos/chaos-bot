package service

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/takama/daemon"
)

// Service is the interface implementation that manages chaos on services
type Service struct {
	JobName string
	Name    string
	Logger  log.Logger
}

// Start will perform a service start on the service specified
func (s *Service) Start() (string, error) {
	dmn, err := daemon.New(s.Name, "", daemon.SystemDaemon)
	if err != nil {
		_ = level.Error(s.Logger).Log("msg", "Could not instantiate daemon", "err", err)
		return "Could not instantiate daemon", err
	}

	if _, startErr := dmn.Start(); startErr != nil {
		_ = level.Error(s.Logger).Log("msg", fmt.Sprintf("Could not start service %s", s.Name), "err", startErr)
		return fmt.Sprintf("Could not start service %s", s.Name), startErr
	}

	_ = level.Info(s.Logger).Log("msg", fmt.Sprintf("Started service with name %s", s.Name))

	return constructMessage(s.Logger, "started", s.Name), nil
}

// Stop will perform a service stop on the service specified
func (s *Service) Stop() (string, error) {
	dmn, err := daemon.New(s.Name, "", daemon.SystemDaemon)
	if err != nil {
		return "Could not instantiate daemon", err
	}

	_, err = dmn.Stop()
	if err != nil {
		return fmt.Sprintf("Could not stop service %s", s.Name), err
	}

	_ = level.Info(s.Logger).Log("msg", fmt.Sprintf("Stopped service with name %s", s.Name))

	return constructMessage(s.Logger, "stopped", s.Name), nil
}

func constructMessage(logger log.Logger, action string, name string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}

	return fmt.Sprintf("Slave %s %s service %s", hostname, action, name)
}
