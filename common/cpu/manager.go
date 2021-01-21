package cpu

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// CPU is the interface implementation that manages cpu failure injections
type CPU struct {
	threads int
	stop    chan int
	logger  log.Logger
}

func New(percentage int, logger log.Logger) (*CPU, error) {
	if percentage > 100 || percentage < 0 {
		return nil, fmt.Errorf("cpu injection percentage %d is out of bounds. should be 0 to 100", percentage)
	}

	return &CPU{
		threads: runtime.NumCPU() * percentage / 100,
		stop:    make(chan int),
		logger:  logger,
	}, nil
}

// Start will perform a cpu failure injection by starting goroutines in for loops
func (cpu *CPU) Start() (string, error) {
	if err := cpu.injection(); err != nil {
		return "Could not inject cpu failure", err
	}
	return constructMessage(cpu.logger, "started"), nil
}

// Start will recover cpu failure by closing all channels
func (cpu *CPU) Stop() (string, error) {
	close(cpu.stop)
	return constructMessage(cpu.logger, "stopped"), nil
}

func (cpu *CPU) injection() error {
	for i := 0; i < cpu.threads; i++ {
		go func() {
			for {
				select {
				case <-cpu.stop:
					return
				default: //nolint:staticcheck
					// left empty on purpose for cpu failure injection
				}
			}
		}()
	}
	return nil
}

func constructMessage(logger log.Logger, action string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}
	return fmt.Sprintf("Bot %s %s cpu injection", hostname, action)
}
