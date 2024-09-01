package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
)

// initialize docker client
func GetDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error initializing Docker client: %v", err)
	}
	return cli
}

// list containers
func ListContainers(all bool) ([]types.Container, error) {
	cli := GetDockerClient()
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: all})
	if err != nil {
		return nil, err
	}

	return containers, nil

}

// start container
func StartContainer(containerID string) error {
	cli := GetDockerClient()
	return cli.ContainerStart(context.Background(), containerID, container.StartOptions{})

}
