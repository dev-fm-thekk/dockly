package dockerapi

import (
	"context"

	"github.com/moby/moby/client"
)

type NetworkConfig struct {
	Container      Container
	Host           string
	Port           string
	Container_port string
}

func GetNetworkConfig(ctn Container) NetworkConfig {
	cli := GetClient()
	ctx := context.Background()

	contJson, _ := cli.ContainerInspect(ctx, ctn.ID, client.ContainerInspectOptions{})

	result := NetworkConfig{
		Container: ctn,
	}
	for port, bindings := range contJson.Container.NetworkSettings.Ports {
		result.Container_port = port.String()
		for _, b := range bindings {
			result.Host = b.HostIP.String()
			result.Port = b.HostPort
		}
	}

	return result
}

func FetchNetConfigAll() []NetworkConfig {
	ctns := FetchContainers()

	result := []NetworkConfig{}
	for ctn := range ctns {
		nconf := GetNetworkConfig(ctns[ctn])
		result = append(result, nconf)
	}

	return result
}
