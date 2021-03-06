package cpu

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

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
		logger: logger,
	}
}

// Start will perform a cpu failure injection by starting goroutines in for loops
func (cpu *CPU) Start(percentage int) (string, error) {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()

	if cpu.status == started {
		return "Could not inject cpu failure", errors.New("CPU injection already running. Stop it before starting another")
	}

	cpu.stop = make(chan int)

	if err := cpu.injection(percentage); err != nil {
		return "Could not inject cpu failure", err
	}

	cpu.status = started
	return constructStartMessage(cpu.logger, percentage), nil
}

// Start will recover cpu failure by closing all channels
func (cpu *CPU) Stop() (string, error) {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()

	if cpu.status != started {
		return "Could not stop cpu failure", errors.New("CPU injection is not running. Start it before trying to stop")
	}

	close(cpu.stop)
	cpu.status = stopped

	return constructMessage(cpu.logger, "stopped"), nil
}

func (cpu *CPU) injection(percent int) error {
	if percent < 0 || percent > 100 {
		return fmt.Errorf("cpu injection percentage %d is out of bounds. should be 0 to 100", percent)
	}

	sleepBaseOnPercentage := time.Duration(1000 * (100 - percent) / 100)

	for i := 0; i < runtime.NumCPU(); i++ {
		time.Sleep(time.Duration(1000/runtime.NumCPU()) * time.Millisecond)
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			for {
				select {
				case <-cpu.stop:
					return
				case <-ticker.C:
					time.Sleep(sleepBaseOnPercentage * time.Millisecond)
				default: //nolint:staticcheck
				}
			}
		}()
	}

	_ = level.Info(cpu.logger).Log("msg", fmt.Sprintf("Starting cpu injection for %d%%", percent))

	return nil
}

func constructStartMessage(logger log.Logger, percentage int) string {
	message := constructMessage(logger, "started")
	return fmt.Sprintf("%s at %d%%", message, percentage)
}

func constructMessage(logger log.Logger, action string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}
	return fmt.Sprintf("Bot %s %s cpu injection", hostname, action)
}
