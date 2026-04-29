package dockerapi

import (
	"context"
	"strings"

	"github.com/moby/moby/client"
)

type Container struct {
	ID     string
	Name   string
	Image  string
	State  string
	Status string
}

func FetchContainers() []Container {
	var results []Container
	cli := GetClient()
	if cli == nil {
		return results
	}

	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return results
	}

	for _, ctr := range containers.Items {
		name := ""
		if len(ctr.Names) > 0 {
			name = strings.TrimPrefix(ctr.Names[0], "/")
		}

		results = append(results, Container{
			ID:     ctr.ID[:12],
			Name:   name,
			Image:  ctr.Image,
			State:  string(ctr.State),
			Status: ctr.Status,
		})
	}

	return results
}
