package server

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Server interface {
	StopUnix() (string, error)
}

type DefaultServer struct {
	logger log.Logger
}

// New will create a new Server instance
func New(logger log.Logger) Server {
	return &DefaultServer{
		logger: logger,
	}
}

// StopUnix will stop a Unix server after one minute
func (ds *DefaultServer) StopUnix() (string, error) {
	go func() {
		_ = level.Info(ds.logger).Log("msg", "server will be shutdown in one minute")
		time.Sleep(1 * time.Minute)
		err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
		if err != nil {
			_ = level.Error(ds.logger).Log("msg", "could not stop server", "err", err)
		}
	}()

	return constructMessage(ds.logger, "will stop"), nil
}

func constructMessage(logger log.Logger, action string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}
	return fmt.Sprintf("Bot %s %s server in 1 minute", hostname, action)
}
