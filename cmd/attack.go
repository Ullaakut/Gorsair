package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func gatherInformation(api *vulnerableDockerAPI) error {
	xyz := strings.Split(api.DockerVersion, ".")
	version := strings.Join(xyz[:2], ".")
	docker, err := client.NewClient(api.Endpoint, version, nil, nil)
	if err != nil {
		return err
	}

	info, err := docker.Info(context.Background())
	if err != nil {
		return err
	}

	api.Info.ContainersRunning = info.ContainersRunning
	api.Info.ContainersStopped = info.ContainersStopped
	api.Info.Images = info.Images
	api.Info.OS = info.OperatingSystem

	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		api.Containers = append(api.Containers, dockerContainer{
			Image:  container.Image,
			ID:     container.ID,
			Mounts: fmt.Sprintf("%v", container.Mounts),
		})
	}

	images, err := docker.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return err
	}

	for _, image := range images {
		if len(image.RepoTags) == 0 {
			continue
		}

		api.Images = append(api.Images, image.RepoTags[0])
	}

	return nil
}
