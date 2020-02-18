package docker

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Docker struct {
	Logger log.Logger
}

func (s *Docker) Start(name string) (string, error) {
	return constructMessage(s.Logger, "started", name), nil
}

func (s *Docker) Stop(name string) (string, error) {
	return constructMessage(s.Logger, "stopped", name), nil
}

func constructMessage(logger log.Logger, action string, name string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}

	return fmt.Sprintf("Slave %s %s docker container %s", hostname, action, name)
}
