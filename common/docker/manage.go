package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Docker is the interface implementation that manages chaos on Docker
type Docker struct {
	Name   string
	Logger log.Logger
}

// Start will perform a docker start on the container specified
func (d *Docker) Start() (string, error) {
	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		_ = level.Error(d.Logger).Log("msg", "Could not instantiate docker client", "err", err)
		return "Could not instantiate docker client", err
	}

	containerID, getIDErr := getContainerID(dockerClient, d.Name)
	if getIDErr != nil {
		_ = level.Error(d.Logger).Log("msg", "Could not get list of containers", "err", getIDErr)
	}

	errStart := dockerClient.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if errStart != nil {
		_ = level.Error(d.Logger).Log("msg", fmt.Sprintf("Could not start docker container %s", containerID), "err", errStart)
		return fmt.Sprintf("Could not start docker container %s", containerID), errStart
	}

	_ = level.Info(d.Logger).Log("msg", fmt.Sprintf("Started container with id %s", containerID))

	return constructMessage(d.Logger, "started", containerID), nil
}

// Stop will perform a docker stop on the container specified
func (d *Docker) Stop() (string, error) {
	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		_ = level.Error(d.Logger).Log("msg", "Could not instantiate docker client", "err", err)
		return "Could not instantiate docker client", err
	}

	containerID, getIDErr := getContainerID(dockerClient, d.Name)
	if getIDErr != nil {
		_ = level.Error(d.Logger).Log("msg", "Could not get list of containers", "err", getIDErr)
	}

	errStop := dockerClient.ContainerStop(context.Background(), containerID, nil)
	if errStop != nil {
		_ = level.Error(d.Logger).Log("msg", fmt.Sprintf("Could not stop docker container %s", containerID), "err", errStop)
		return fmt.Sprintf("Could not stop docker container %s", containerID), errStop
	}

	_ = level.Info(d.Logger).Log("msg", fmt.Sprintf("Stopped container with id %s", containerID))

	return constructMessage(d.Logger, "stopped", d.Name), nil
}

func constructMessage(logger log.Logger, action string, name string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}

	return fmt.Sprintf("Bot %s %s docker container %s", hostname, action, name)
}

func getContainerID(cli *client.Client, name string) (string, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return name, err
	}

	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == name {
				return container.ID, nil
			}
		}
	}

	return name, nil
}
