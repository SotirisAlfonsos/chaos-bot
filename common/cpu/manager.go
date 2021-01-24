package cpu

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// CPU is the interface implementation that manages cpu failure injections
type CPU struct {
	mu     sync.Mutex
	status status
	stop   chan int
	logger log.Logger
}

type status int

const (
	stopped status = iota
	started
)

// New will create a new CPU instance with the amount of threads to perform
// the injection and the channel that will be used to stop it
func New(logger log.Logger) *CPU {
	return &CPU{
		stop:   make(chan int),
		logger: logger,
	}
}

// Start will perform a cpu failure injection by starting goroutines in for loops
func (cpu *CPU) Start(threads int) (string, error) {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()

	if cpu.status == started {
		return "Could not inject cpu failure", errors.New("CPU injection already running. Stop it before starting another")
	}

	if err := cpu.injection(threads); err != nil {
		return "Could not inject cpu failure", err
	}

	cpu.status = started
	return constructMessage(cpu.logger, "started"), nil
}

// Start will recover cpu failure by closing all channels
func (cpu *CPU) Stop() (string, error) {
	close(cpu.stop)
	cpu.status = stopped
	return constructMessage(cpu.logger, "stopped"), nil
}

func (cpu *CPU) injection(threads int) error {
	if threads <= 0 {
		return errors.New("base on the percentage specified and your number of CPUs we can only block 0 cpu cores. ( CPU num * percentage / 100 )")
	}

	for i := 0; i < threads; i++ {
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
