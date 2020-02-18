package service

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/takama/daemon"
)

type Service struct {
	Logger log.Logger
}

func (s *Service) Start(name string) (string, error) {
	dmn, err := daemon.New(name, "")
	if err != nil {
		_ = level.Error(s.Logger).Log("msg", "Could not instantiate daemon", "err", err)
		return "Could not instantiate daemon", err
	}

	_, startErr := dmn.Start()
	if startErr != nil {
		_ = level.Error(s.Logger).Log("msg", fmt.Sprintf("Could not start service %s", name), "err", startErr)
		return fmt.Sprintf("Could not start service %s", name), startErr
	}

	return constructMessage(s.Logger, "started", name), nil
}

func (s *Service) Stop(name string) (string, error) {
	dmn, err := daemon.New(name, "")
	if err != nil {
		return "Could not instantiate daemon", err
	}

	_, err = dmn.Stop()
	if err != nil {
		return fmt.Sprintf("Could not stop service %s", name), err
	}

	return constructMessage(s.Logger, "stopped", name), nil
}

func constructMessage(logger log.Logger, action string, name string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}

	return fmt.Sprintf("Slave %s %s service %s", hostname, action, name)
}
