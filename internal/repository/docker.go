package repository

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func DockerLogin(ctx context.Context) (string, error) {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	resp, err := cli.RegistryLogin(ctx, types.AuthConfig{
		ServerAddress: "",
		Username:      "",
		Password:      "",
	})
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}
