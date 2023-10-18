package core_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/core"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	useProvider      = testcontainers.ProviderDefault // or explicit to testcontainers.ProviderPodman if needed
	mlmdImage        = "gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0"
	sqliteFile       = "metadata.sqlite.db"
	testConfigFolder = "test/config/ml-metadata"
)

func TestUsingContainer(t *testing.T) {
	conn, teardown := setupTestContainer(t)
	defer teardown(t)

	// [TEST CASE]

	// create mode registry service
	service, err := core.NewModelRegistryService(conn)
	if err != nil {
		t.Errorf("error creating core service: %v", err)
	}

	nameModel1 := "PricingModel"
	nameModel2 := "ForecastingModel"
	version1 := "v1"
	version2 := "v2"
	author1 := "John"
	author2 := "Jim"
	uri1 := "/path/to/pricing/v1"
	uri2 := "/path/to/pricing/v2"
	modelFormat1 := "tensorflow"
	modelFormat2 := "sklearn"

	// register two new models
	registerNewModel(service, nameModel1, version1, author1, uri1, modelFormat1)
	registerNewModel(service, nameModel2, version2, author2, uri2, modelFormat2)

	// retrieve models and log content
	allRegModels, _, err := service.GetRegisteredModels(core.ListOptions{})
	if err != nil {
		log.Fatalf("Error getting all registered models: %v", err)
	}
	allRegModelJson, _ := json.MarshalIndent(allRegModels, "", "  ")
	fmt.Printf("Found n. %d registered models: %v\n", len(allRegModels), string(allRegModelJson))

	id, _ := idToInt(*allRegModels[0].Id)
	regModel, err := service.GetRegisteredModelById((*core.BaseResourceId)(id))
	if err != nil {
		log.Fatalf("Error getting registered model: %v", err)
	}
	regModelJson, _ := json.MarshalIndent(regModel, "", "  ")
	fmt.Printf("Getting registered model: %+v, \n", string(regModelJson))
}

func registerNewModel(service core.ModelRegistryApi, name string, version string, author string, uri string, format string) {
	registeredModel := &openapi.RegisteredModel{
		Name: &name,
	}

	_, err := service.UpsertRegisteredModel(registeredModel)
	if err != nil {
		log.Fatalf("Error creating registered model: %v", err)
	}
}

// #################
// ##### Utils #####
// #################

func idToInt(idString string) (*int64, error) {
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	idInt64 := int64(idInt)

	return &idInt64, nil
}

func clearMetadataSqliteDB(wd string) error {
	if err := os.Remove(fmt.Sprintf("%s/%s", wd, sqliteFile)); err != nil {
		return fmt.Errorf("expected to clear sqlite file but didn't find: %v", err)
	}
	return nil
}

// setupTestContainer create a MLMD gRPC test container returning the mlmd uri and a teardown function
func setupTestContainer(t *testing.T) (*grpc.ClientConn, func(t *testing.T)) {
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("error getting working directory: %v", err)
	}
	wd = fmt.Sprintf("%s/../../%s", wd, testConfigFolder)
	t.Logf("using working directory: %s", wd)

	req := testcontainers.ContainerRequest{
		Image:        mlmdImage,
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
		t.Errorf("error setting up mlmd grpc container: %v", err)
	}

	mappedHost, err := mlmdgrpc.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	mappedPort, err := mlmdgrpc.MappedPort(ctx, "8080")
	if err != nil {
		t.Error(err)
	}
	mlmdAddr := fmt.Sprintf("%s:%s", mappedHost, mappedPort.Port())
	t.Log("MLMD test container setup at: ", mlmdAddr)

	// setup grpc connection
	conn, err := grpc.DialContext(
		context.Background(),
		mlmdAddr,
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Errorf("error dialing connection to mlmd server %s: %v", mlmdAddr, err)
	}

	return conn, func(t *testing.T) {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
		if err := mlmdgrpc.Terminate(ctx); err != nil {
			t.Error(err)
		}
		if err := clearMetadataSqliteDB(wd); err != nil {
			t.Error(err)
		}
	}
}
