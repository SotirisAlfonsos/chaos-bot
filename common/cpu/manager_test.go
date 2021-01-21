package cpu

import (
	"fmt"
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
	message    string
	percentage int
	expected   map[string]string
}

func TestCPU_Start_and_Stop(t *testing.T) {
	dataItems := []dataItems{
		{
			message:    "Start and stop cpu experiment for cpu 50 percent",
			percentage: 50,
			expected: map[string]string{
				"start": "Bot ubuntu started cpu injection",
				"stop":  "Bot ubuntu stopped cpu injection",
			},
		},
		{
			message:    "Start and stop cpu experiment for cpu 100 percent",
			percentage: 100,
			expected: map[string]string{
				"start": "Bot ubuntu started cpu injection",
				"stop":  "Bot ubuntu stopped cpu injection",
			},
		},
		{
			message:    "Start and stop cpu experiment for cpu 1 percent",
			percentage: 1,
			expected: map[string]string{
				"start": "Bot ubuntu started cpu injection",
				"stop":  "Bot ubuntu stopped cpu injection",
			},
		},
		{
			message:    "Start and stop cpu experiment for cpu 0 percent",
			percentage: 0,
			expected: map[string]string{
				"start": "Bot ubuntu started cpu injection",
				"stop":  "Bot ubuntu stopped cpu injection",
			},
		},
	}

	for _, dataItem := range dataItems {
		t.Logf(dataItem.message)
		cpu, err := New(dataItem.percentage, logger)
		if err != nil {
			t.Fatalf(err.Error())
		}
		goroutinesBeforeExperiment := runtime.NumGoroutine()

		startResponse := startExperiment(t, cpu)

		goroutinesDuringInjection := runtime.NumGoroutine()
		assert.Equal(t, goroutinesBeforeExperiment+cpu.threads, goroutinesDuringInjection)
		assert.Equal(t, dataItem.expected["start"], startResponse)

		stopResponse := stopExperiment(t, cpu)

		goroutinesAfterInjection := runtime.NumGoroutine()
		assert.Equal(t, goroutinesAfterInjection, goroutinesDuringInjection-cpu.threads)
		assert.Equal(t, dataItem.expected["stop"], stopResponse)
	}
}

func startExperiment(t *testing.T, cpu *CPU) string {
	startResponse, startErr := cpu.Start()
	if startErr != nil {
		t.Fatalf("Failed to start cpu injection")
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

func TestCPU_Start_and_Stop_with_invalid_percentage(t *testing.T) {
	dataItems := []dataItems{
		{
			message:    "Should not start experiment for cpu more than 100",
			percentage: 101,
			expected: map[string]string{
				"new cpu": "cpu injection percentage 101 is out of bounds. should be 0 to 100",
			},
		},
		{
			message:    "Should not start experiment for cpu less than 0",
			percentage: -1,
			expected: map[string]string{
				"new cpu": "cpu injection percentage -1 is out of bounds. should be 0 to 100",
			},
		},
	}

	for _, dataItem := range dataItems {
		t.Logf(dataItem.message)
		_, err := New(dataItem.percentage, logger)
		if err == nil {
			t.Fatalf("Should have error for invalid percentage %d", dataItem.percentage)
		}

		assert.Equal(t, dataItem.expected["new cpu"], err.Error())
	}
}

func getLogger() log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}

	return chaoslogger.New(allowLevel)
}
