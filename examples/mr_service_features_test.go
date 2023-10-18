package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/opendatahub-io/model-registry/internal/core"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

type wdCtxKey struct{}
type testContainerCtxKey struct{}
type svcLayerCtxKey struct{}
type connCtxKey struct{}

func iHaveAConnectionToMR(ctx context.Context) (context.Context, error) {
	mlmdAddr := fmt.Sprintf("%s:%d", mlmdHostname, mlmdPort)
	conn, err := grpc.DialContext(
		context.Background(),
		mlmdAddr,
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Error dialing connection to mlmd server %s: %v", mlmdAddr, err)
		return nil, err
	}
	ctx = context.WithValue(ctx, connCtxKey{}, conn)
	service, err := core.NewModelRegistryService(conn)
	if err != nil {
		log.Fatalf("Error creating core service: %v", err)
		return nil, err
	}
	return context.WithValue(ctx, svcLayerCtxKey{}, service), nil
}

func iStoreARegisteredModelWithNameAndAChildVersionedModelWithNameAndAChildArtifactWithUri(ctx context.Context, registedModelName, modelVersionName, artifactURI string) error {
	service, ok := ctx.Value(svcLayerCtxKey{}).(core.ModelRegistryApi)
	if !ok {
		return fmt.Errorf("not found service layer connection in godog context")
	}
	var registeredModel *openapi.RegisteredModel
	var err error
	registeredModel, err = service.UpsertRegisteredModel(&openapi.RegisteredModel{
		Name: &registedModelName,
	})
	if err != nil {
		return err
	}
	registeredModelId, err := idToInt64(*registeredModel.Id)
	if err != nil {
		return err
	}

	var versionedModel *openapi.ModelVersion
	if versionedModel, err = service.UpsertModelVersion(&openapi.ModelVersion{Name: &modelVersionName}, (*core.BaseResourceId)(registeredModelId)); err != nil {
		return err
	}
	versionedModelId, err := idToInt64(*versionedModel.Id)
	if err != nil {
		return err
	}

	if _, err = service.UpsertModelArtifact(&openapi.ModelArtifact{Uri: &artifactURI}, (*core.BaseResourceId)(versionedModelId)); err != nil {
		return err
	}

	return nil
}

func idToInt64(idString string) (*int64, error) {
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	idInt64 := int64(idInt)

	return &idInt64, nil
}

func thereShouldBeAMlmdContextOfTypeNamed(ctx context.Context, arg1, arg2 string) error {
	conn := ctx.Value(connCtxKey{}).(*grpc.ClientConn)
	client := proto.NewMetadataStoreServiceClient(conn)
	query := fmt.Sprintf("type = \"%s\" and name = \"%s\"", arg1, arg2)
	fmt.Println("query: ", query)
	resp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{
		Options: &proto.ListOperationOptions{
			FilterQuery: &query,
		},
	})
	if err != nil {
		return err
	}
	if len(resp.Contexts) != 1 {
		return fmt.Errorf("Unexpected mlmd Context result size (%d), %v", len(resp.Contexts), resp.Contexts)
	}
	return nil
}

func thereShouldBeAMlmdContextOfTypeHavingPropertyNamedValorisedWithStringValue(ctx context.Context, arg1, arg2, arg3 string) error {
	conn := ctx.Value(connCtxKey{}).(*grpc.ClientConn)
	client := proto.NewMetadataStoreServiceClient(conn)
	query := fmt.Sprintf("type = \"%s\" and properties.%s.string_value = \"%s\"", arg1, arg2, arg3)
	fmt.Println("query: ", query)
	resp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{
		Options: &proto.ListOperationOptions{
			FilterQuery: &query,
		},
	})
	if err != nil {
		return err
	}
	if len(resp.Contexts) != 1 {
		return fmt.Errorf("Unexpected mlmd Context result size (%d), %v", len(resp.Contexts), resp.Contexts)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		wd, mlmdgrpc, err := setupTestContainer(ctx)
		if err != nil {
			return ctx, err
		}
		ctx = context.WithValue(ctx, wdCtxKey{}, wd)
		ctx = context.WithValue(ctx, testContainerCtxKey{}, mlmdgrpc)
		mappedHost, err := mlmdgrpc.Host(ctx)
		if err != nil {
			return ctx, err
		}
		mappedPort, err := mlmdgrpc.MappedPort(ctx, "8080")
		if err != nil {
			return ctx, err
		}
		// TODO: these are effectively global in main and could be worthy to revisit
		mlmdHostname = mappedHost
		mlmdPort = mappedPort.Int()
		return ctx, nil
	})
	ctx.Step(`^I have a connection to MR$`, iHaveAConnectionToMR)
	ctx.Step(`^I store a RegisteredModel with name "([^"]*)" and a child VersionedModel with name "([^"]*)" and a child Artifact with uri "([^"]*)"$`, iStoreARegisteredModelWithNameAndAChildVersionedModelWithNameAndAChildArtifactWithUri)
	ctx.Step(`^there should be a mlmd Context of type "([^"]*)" named "([^"]*)"$`, thereShouldBeAMlmdContextOfTypeNamed)
	ctx.Step(`^there should be a mlmd Context of type "([^"]*)" having property named "([^"]*)" valorised with string value "([^"]*)"$`, thereShouldBeAMlmdContextOfTypeHavingPropertyNamedValorisedWithStringValue)
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		conn := ctx.Value(connCtxKey{}).(*grpc.ClientConn)
		client := proto.NewMetadataStoreServiceClient(conn)
		resp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
		if err != nil {
			return nil, err
		}
		fmt.Printf("%v", resp)
		conn.Close()
		mlmdgrpc := ctx.Value(testContainerCtxKey{}).(testcontainers.Container)
		if err := mlmdgrpc.Terminate(ctx); err != nil {
			return ctx, err
		}
		wd := ctx.Value(wdCtxKey{}).(string)
		clearMetadataSqliteDB(wd)
		return ctx, nil
	})
}
