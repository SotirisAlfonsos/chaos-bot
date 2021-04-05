package docker

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Docker is the interface implementation that manages chaos on Docker
type Docker struct {
	Logger log.Logger
	client *client.Client
}

// Recover will perform a docker start on the container specified
func (d *Docker) Recover(container string) (string, error) {
	err := d.initClient()
	if err != nil {
		return "Could not instantiate docker client", err
	}

	containerID := getContainerID(d.client, container, d.Logger)

	errRecover := d.client.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if errRecover != nil {
		_ = level.Error(d.Logger).Log("msg", fmt.Sprintf("Could not recover docker container {%s}", containerID), "err", errRecover)
		return fmt.Sprintf("Could not recover docker container {%s}", containerID), status.Error(codes.Internal, errRecover.Error())
	}

	_ = level.Info(d.Logger).Log("msg", fmt.Sprintf("Recovered container with id {%s}", containerID))

	return constructMessage(d.Logger, "recovered", containerID), nil
}

// Kill will perform a docker stop on the container specified
func (d *Docker) Kill(container string) (string, error) {
	err := d.initClient()
	if err != nil {
		return "Could not instantiate docker client", err
	}

	containerID := getContainerID(d.client, container, d.Logger)

	errKill := d.client.ContainerStop(context.Background(), containerID, nil)
	if errKill != nil {
		_ = level.Error(d.Logger).Log("msg", fmt.Sprintf("Could not kill docker container {%s}", containerID), "err", errKill)
		return fmt.Sprintf("Could not kill docker container {%s}", containerID), status.Error(codes.Internal, errKill.Error())
	}

	_ = level.Info(d.Logger).Log("msg", fmt.Sprintf("Killed container with id %s", containerID))

	return constructMessage(d.Logger, "killed", containerID), nil
}

func (d *Docker) initClient() error {
	if d.client != nil {
		return nil
	}

	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		_ = level.Error(d.Logger).Log("msg", "Could not instantiate docker client", "err", err)
		return status.Error(codes.Internal, err.Error())
	}

	d.client = dockerClient
	return nil
}

func constructMessage(logger log.Logger, action string, name string) string {
	hostname, err := os.Hostname()
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get hostname", "err", err)
		hostname = "Unknown"
	}

	return fmt.Sprintf("Bot %s %s docker container %s", hostname, action, name)
}

func getContainerID(cli *client.Client, name string, logger log.Logger) string {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Could not get list of containers", "err", err)
		return name
	}

	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == name {
				return container.ID
			}
		}
	}

	return name
}
