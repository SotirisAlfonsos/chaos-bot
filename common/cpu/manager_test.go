package cpu

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/SotirisAlfonsos/chaos-master/chaoslogger"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

var (
	logger = getLogger()
)

type dataItems struct {
	message  string
	threads  int
	expected map[string]string
}

func TestCPU_Start_and_Stop(t *testing.T) {
	dataItems := []dataItems{
		{
			message: "Start and stop cpu experiment for single thread",
			threads: 1,
			expected: map[string]string{
				"start": fmt.Sprintf("Bot %s started cpu injection", getHostname(t)),
				"stop":  fmt.Sprintf("Bot %s stopped cpu injection", getHostname(t)),
			},
		},
	}

	for _, dataItem := range dataItems {
		t.Logf(dataItem.message)
		cpu := New(logger)

		goroutinesBeforeExperiment := runtime.NumGoroutine()

		startResponse := startExperiment(t, dataItem.threads, cpu)

		goroutinesDuringInjection := runtime.NumGoroutine()
		assert.Equal(t, goroutinesBeforeExperiment+dataItem.threads, goroutinesDuringInjection)
		assert.Equal(t, dataItem.expected["start"], startResponse)

		stopResponse := stopExperiment(t, cpu)

		goroutinesAfterInjection := runtime.NumGoroutine()
		assert.Equal(t, goroutinesAfterInjection, goroutinesDuringInjection-dataItem.threads)
		assert.Equal(t, dataItem.expected["stop"], stopResponse)
	}
}

func startExperiment(t *testing.T, percentage int, cpu *CPU) string {
	startResponse, startErr := cpu.Start(percentage)
	if startErr != nil {
		t.Fatalf("Failed to start cpu injection. err=%s", startErr.Error())
	}
	time.Sleep(3 * time.Second)
	return startResponse
}

func stopExperiment(t *testing.T, cpu *CPU) string {
	stopResponse, stopErr := cpu.Stop()
	if stopErr != nil {
		t.Fatalf("Failed to stop cpu injection")
	}
	time.Sleep(1 * time.Second)
	return stopResponse
}

func TestCPU_Start_with_status_already_started(t *testing.T) {
	cpu := New(logger)
	cpu.status = started

	message, err := cpu.Start(50)

	assert.Equal(t, "Could not inject cpu failure", message)
	assert.NotNil(t, err)
	assert.Equal(t, "CPU injection already running. Stop it before starting another", err.Error())
	assert.Equal(t, started, cpu.status)
}

func TestCPU_Start_and_Stop_with_invalid_thread_count(t *testing.T) {
	dataItems := []dataItems{
		{
			message: "Should not start experiment for percentage leading to 0 threads",
			threads: 0,
			expected: map[string]string{
				"message": "Could not inject cpu failure",
				"error":   "base on the percentage specified and your number of CPUs we can only block 0 cpu cores. ( CPU num * percentage / 100 )",
			},
		},
		{
			message: "Should not start experiment for percentage leading to -1 threads",
			threads: -1,
			expected: map[string]string{
				"message": "Could not inject cpu failure",
				"error":   "base on the percentage specified and your number of CPUs we can only block 0 cpu cores. ( CPU num * percentage / 100 )",
			},
		},
	}
	for _, dataItem := range dataItems {
		t.Logf(dataItem.message)
		cpu := New(logger)

		message, err := cpu.Start(dataItem.threads)

		assert.Equal(t, dataItem.expected["message"], message)
		assert.NotNil(t, err)
		assert.Equal(t, dataItem.expected["error"], err.Error())
	}
}

func getHostname(t *testing.T) string {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("Can not get hostname")
	}
	return hostname
}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}

	return chaoslogger.New(allowLevel)
}
