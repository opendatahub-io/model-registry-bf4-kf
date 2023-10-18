package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	useProvider = testcontainers.ProviderDefault // or explicit to testcontainers.ProviderPodman on local testing
)

func TestUsingContainer(t *testing.T) {
	ctx := context.Background()
	wd, mlmdgrpc, err := setupTestContainer(ctx)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := mlmdgrpc.Terminate(ctx); err != nil {
			t.Error(err)
		}
	}()
	mappedHost, err := mlmdgrpc.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	mappedPort, err := mlmdgrpc.MappedPort(ctx, "8080")
	if err != nil {
		t.Error(err)
	}
	uri := fmt.Sprintf("http://%s:%s", mappedHost, mappedPort.Port())
	t.Log("using uri: ", uri)

	mlmdHostname = mappedHost
	mlmdPort = mappedPort.Int()
	main()

	if err := clearMetadataSqliteDB(wd); err != nil {
		t.Error(err)
	}
}

func clearMetadataSqliteDB(wd string) error {
	if err := os.Remove(fmt.Sprintf("%s/%s", wd, "metadata.sqlite.db")); err != nil {
		return fmt.Errorf("Expected to clear sqlite file but didn't find: %v", err)
	}
	return nil
}

func setupTestContainer(ctx context.Context) (string, testcontainers.Container, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", nil, err
	}
	req := testcontainers.ContainerRequest{
		Image:        "gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0",
		ExposedPorts: []string{"8080/tcp"},
		Env: map[string]string{
			"METADATA_STORE_SERVER_CONFIG_FILE": "/tmp/shared/conn_config.pb",
		},
		Mounts: testcontainers.ContainerMounts{
			testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{
					HostPath: wd,
				},
				Target: "/tmp/shared",
			},
		},
		WaitingFor: wait.ForLog("Server listening on"),
	}
	mlmdgrpc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ProviderType:     useProvider,
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, err
	}
	return wd, mlmdgrpc, nil
}
