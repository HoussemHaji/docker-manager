package docker

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"

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

// stop container
func StopContainer(containerID string) error {
	cli := GetDockerClient()
	stopOptions := container.StopOptions{
		Timeout: nil, // You can set this to nil for default behavior, or a specific timeout value
	}
	return cli.ContainerStop(context.Background(), containerID, stopOptions)
}

// pause container
func PauseContainer(containerID string) error {
	cli := GetDockerClient()
	return cli.ContainerPause(context.Background(), containerID)
}

// unpause container
func UnpauseContainer(containerID string) error {
	cli := GetDockerClient()
	return cli.ContainerUnpause(context.Background(), containerID)
}

// delete container
func DeleteContainer(containerID string) error {
	cli := GetDockerClient()
	return cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true})
}

// execute command inside container
func ExecCommandInContainer(containerID string, cmd []string) error {
	cli := GetDockerClient()

	execConfig := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execIDResp, err := cli.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		return err
	}

	resp, err := cli.ContainerExecAttach(context.Background(), execIDResp.ID, container.ExecStartOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()

	_, err = io.Copy(os.Stdout, resp.Reader)
	return err
}

// retrieve container logs
func GetContainerLogs(containerID string) error {
	cli := GetDockerClient()

	logs, err := cli.ContainerLogs(context.Background(), containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}
	defer logs.Close()

	_, err = io.Copy(os.Stdout, logs)
	return err
}

// view container network and IP
func GetContainerNetworkInfo(containerID string) (map[string]string, error) {
	cli := GetDockerClient()
	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, err
	}

	networkInfo := make(map[string]string)
	for networkName, network := range containerJSON.NetworkSettings.Networks {
		networkInfo[networkName] = network.IPAddress
	}

	return networkInfo, nil
}

// filter containers by name
func FilterContainersByName(filterVal string) ([]types.Container, error) {
	cli := GetDockerClient()
	args := filters.NewArgs()
	args.Add("name", filterVal)
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true, Filters: args})

	if err != nil {
		return nil, err
	}

	return containers, nil
}

// filter containers by status
func FilterContainersByStatus(filterVal string) ([]types.Container, error) {
	cli := GetDockerClient()
	args := filters.NewArgs()
	args.Add("status", filterVal)
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true, Filters: args})

	if err != nil {
		return nil, err
	}

	return containers, nil
}

// restart container
func RestartContainer(containerID string) error {
	cli := GetDockerClient()
	restartOptions := container.StopOptions{
		Timeout: nil, // You can set this to nil for default behavior, or a specific timeout value
	}
	return cli.ContainerRestart(context.Background(), containerID, restartOptions)
}
