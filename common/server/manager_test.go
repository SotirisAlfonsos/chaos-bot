package server

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerShutdown(t *testing.T) {
	hostname := getHostname()

	server := &DefaultServer{}
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
