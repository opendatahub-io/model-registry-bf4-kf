package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/opendatahub-io/model-registry/internal/core"
	"github.com/opendatahub-io/model-registry/internal/model/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InterceptorLogger(l *log.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			msg = fmt.Sprintf("DEBUG :%v", msg)
		case logging.LevelInfo:
			msg = fmt.Sprintf("INFO :%v", msg)
		case logging.LevelWarn:
			msg = fmt.Sprintf("WARN :%v", msg)
		case logging.LevelError:
			msg = fmt.Sprintf("ERROR :%v", msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
		l.Println(append([]any{"msg", msg}, fields...))
	})
}

var mlmdHostname = "localhost"
var mlmdPort = 8080

// Prerequisites: having MLMD store server running on localhost:808
func main() {
	glog.Infof("Setting up connection with MLMD server at %s:%v", mlmdHostname, mlmdPort)

	// logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	// lopts := []logging.Option{
	// 	logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent, logging.StartCall, logging.FinishCall),
	// 	// Add any other option (check functions starting with logging.With).
	// }

	mlmdAddr := fmt.Sprintf("%s:%d", mlmdHostname, mlmdPort)
	conn, err := grpc.DialContext(
		context.Background(),
		mlmdAddr,
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithChainUnaryInterceptor(
		// 	logging.UnaryClientInterceptor(InterceptorLogger(logger), lopts...),
		// ),
	)
	if err != nil {
		log.Fatalf("Error dialing connection to mlmd server %s: %v", mlmdAddr, err)
	}
	defer conn.Close()

	service, err := core.NewModelRegistryService(conn)
	if err != nil {
		log.Fatalf("Error creating core service: %v", err)
	}

	// [START]: demo showcasing model registration with multiple versions
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

	registerModelVersion(service, nameModel1, version1, author1, uri1, modelFormat1)
	registerModelVersion(service, nameModel2, version2, author2, uri2, modelFormat2)

	allRegModels, _, err := service.GetRegisteredModels(core.ListOptions{})
	if err != nil {
		log.Fatalf("Error getting all registered models: %v", err)
	}
	allRegModelJson, _ := json.MarshalIndent(allRegModels, "", "  ")
	fmt.Printf("Found n. %d registered models: %v\n", len(allRegModels), string(allRegModelJson))

	regModel, err := service.GetRegisteredModelById((*core.BaseResourceId)(allRegModels[0].Id))
	if err != nil {
		log.Fatalf("Error getting registered model: %v", err)
	}
	regModelJson, _ := json.MarshalIndent(regModel, "", "  ")
	fmt.Printf("Getting registered model: %+v, \n", string(regModelJson))

	// v1, err := service.GetModelVersion(nameModel, version1)
	// if err != nil {
	// 	log.Fatalf("Error getting model version v1: %v", err)
	// }
	// v1Json, _ := json.MarshalIndent(v1, "", "  ")
	// fmt.Printf("Getting model version: %v\n", string(v1Json))

	// [END]: demo showcasing model registration with multiple versions

	glog.Info("shutdown!")
}

func registerModelVersion(service core.ModelRegistryApi, name string, version string, author string, uri string, format string) {
	modelVersion := &registry.VersionedModel{}

	modelArtifactName := fmt.Sprintf("%s/%s", name, format)
	modelArtifact := registry.Artifact{
		Name: &modelArtifactName,
		Uri:  &uri,
	}

	artifacts := &[]registry.Artifact{
		modelArtifact,
	}

	modelVersion.ModelName = &name
	modelVersion.ModelUri = uri
	modelVersion.Version = &version
	modelVersion.Artifacts = artifacts
	modelVersion.Author = &author
	modelVersion.Metadata = &map[string]interface{}{
		"accuracy": 0.89,
		"not_supported_key": map[string]interface{}{
			"custom_key": "custom_value",
		},
	}

	registeredModel := &registry.RegisteredModel{
		Name: &name,
	}

	_, err := service.UpsertRegisteredModel(registeredModel)
	if err != nil {
		log.Fatalf("Error creating registered model: %v", err)
	}
}
