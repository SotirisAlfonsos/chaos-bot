// +build integration

package server

import (
	"fmt"
	"os"
	"testing"

	"github.com/SotirisAlfonsos/chaos-master/chaoslogger"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

var (
	logger = getLogger()
)

func TestServerShutdown(t *testing.T) {
	hostname := getHostname()

	server := New(logger)
	message, err := server.StopUnix()

	assert.Equal(t, fmt.Sprintf("Bot %s will stop server in 1 minute", hostname), message)
	assert.Nil(t, err)
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown"
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
