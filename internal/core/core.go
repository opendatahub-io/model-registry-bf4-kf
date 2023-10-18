package core

import (
	"context"
	"fmt"
	"log"

	"github.com/opendatahub-io/model-registry/internal/core/mapper"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/registry"
	"google.golang.org/grpc"
)

var (
	registeredModelType = "odh.RegisteredModel"
	modelVersionType    = "odh.ModelVersion"
	modelArtifactType   = "odh.ModelArtifact"
)

type ModelRegistryApi interface {
	CreateRegisteredModel(registeredModel *registry.RegisteredModel) (*int64, error)
	GetRegisteredModels() ([]*registry.RegisteredModel, error)
	GetRegisteredModel(name string) (*registry.RegisteredModel, error)

	CreateModelVersion(modelVersion *registry.VersionedModel) (*int64, error)
	GetModelVersion(name string, version string) (*registry.VersionedModel, error)
}

// modelRegistryService is the core library of the model registry
type modelRegistryService struct {
	mlmdClient proto.MetadataStoreServiceClient
	mapper     *mapper.Mapper
}

// NewModelRegistryService create a fresh instance of ModelRegistryService, taking care of setting up needed MLMD Types
func NewModelRegistryService(cc grpc.ClientConnInterface) (*modelRegistryService, error) {

	client := proto.NewMetadataStoreServiceClient(cc)

	// Setup the needed Type instances if not existing already

	registeredModelReq := proto.PutContextTypeRequest{
		ContextType: &proto.ContextType{
			Name: &registeredModelType,
		},
	}

	modelVersionReq := proto.PutContextTypeRequest{
		ContextType: &proto.ContextType{
			Name: &modelVersionType,
			Properties: map[string]proto.PropertyType{
				"model_name": proto.PropertyType_STRING,
				"version":    proto.PropertyType_STRING,
				"author":     proto.PropertyType_STRING,
			},
		},
	}

	modelArtifactReq := proto.PutArtifactTypeRequest{
		ArtifactType: &proto.ArtifactType{
			Name: &modelArtifactType,
			Properties: map[string]proto.PropertyType{
				"model_format": proto.PropertyType_STRING,
			},
		},
	}

	registeredModelResp, err := client.PutContextType(context.Background(), &registeredModelReq)
	if err != nil {
		log.Fatalf("Error setting up context type %s: %v", registeredModelType, err)
	}

	modelVersionResp, err := client.PutContextType(context.Background(), &modelVersionReq)
	if err != nil {
		log.Fatalf("Error setting up context type %s: %v", modelVersionType, err)
	}
	modelArtifactResp, err := client.PutArtifactType(context.Background(), &modelArtifactReq)
	if err != nil {
		log.Fatalf("Error setting up artifact type %s: %v", modelArtifactType, err)
	}

	return &modelRegistryService{
		mlmdClient: client,
		mapper:     mapper.NewMapper(registeredModelResp.GetTypeId(), modelVersionResp.GetTypeId(), modelArtifactResp.GetTypeId()),
	}, nil
}

// Create a registered model
func (serv *modelRegistryService) CreateRegisteredModel(registeredModel *registry.RegisteredModel) (*int64, error) {
	log.Printf("Creating registered model for %s", *registeredModel.Name)

	modelCtx, err := serv.mapper.MapFromRegisteredModel(registeredModel)
	if err != nil {
		return nil, err
	}

	modelCtxResp, err := serv.mlmdClient.PutContexts(context.Background(), &proto.PutContextsRequest{
		Contexts: []*proto.Context{
			modelCtx,
		},
	})
	if err != nil {
		return nil, err
	}

	return &modelCtxResp.ContextIds[0], nil
}

// Get all registered models without all versions
func (serv *modelRegistryService) GetRegisteredModels() ([]*registry.RegisteredModel, error) {

	resp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &registeredModelType,
	})
	if err != nil {
		return nil, err
	}

	models := []*registry.RegisteredModel{}
	for _, ctx := range resp.Contexts {
		model, err := serv.mapper.MapToRegisteredModel(ctx)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

// Get registered model by name
func (serv *modelRegistryService) GetRegisteredModel(name string) (*registry.RegisteredModel, error) {
	log.Printf("Getting registered model %s", name)

	regModelResp, err := serv.mlmdClient.GetContextByTypeAndName(context.Background(), &proto.GetContextByTypeAndNameRequest{
		TypeName:    &registeredModelType,
		ContextName: &name,
	})
	if err != nil {
		return nil, err
	}

	modelVersionsResp, err := serv.mlmdClient.GetChildrenContextsByContext(context.Background(), &proto.GetChildrenContextsByContextRequest{
		ContextId: regModelResp.Context.Id,
	})
	if err != nil {
		return nil, err
	}

	versions := []registry.VersionedModel{}
	for _, mvc := range modelVersionsResp.Contexts {

		// get all artifacts associated to that model version
		artResp, err := serv.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
			ContextId: mvc.Id,
		})
		if err != nil {
			return nil, err
		}

		mv, err := serv.mapper.MapToModelVersion(mvc, artResp.Artifacts)
		if err != nil {
			return nil, err
		}

		versions = append(versions, *mv)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(regModelResp.Context)
	if err != nil {
		return nil, err
	}

	regModel.Versions = &versions

	return regModel, nil
}

// Versioned Models

// Create a model version
func (serv *modelRegistryService) CreateModelVersion(modelVersion *registry.VersionedModel) (*int64, error) {
	fullName := fmt.Sprintf("%s:%s", *modelVersion.ModelName, *modelVersion.Version)
	log.Printf("Creating model version for %s", fullName)

	// check if RegisteredModel is present, if not create it
	checkRegModelResp, err := serv.mlmdClient.GetContextByTypeAndName(context.Background(), &proto.GetContextByTypeAndNameRequest{
		TypeName:    &registeredModelType,
		ContextName: modelVersion.ModelName,
	})

	var registeredModelId *int64
	if err != nil || checkRegModelResp.Context == nil {
		// create new simple registered model
		regModel := &registry.RegisteredModel{
			Name: modelVersion.ModelName,
		}
		registeredModelId, err = serv.CreateRegisteredModel(regModel)
		if err != nil {
			return nil, err
		}
	} else {
		registeredModelId = checkRegModelResp.Context.Id
	}

	// create a model version
	versionCtx, err := serv.mapper.MapFromModelVersion(modelVersion, *registeredModelId)
	if err != nil {
		return nil, err
	}

	versionCtxResp, err := serv.mlmdClient.PutContexts(context.Background(), &proto.PutContextsRequest{
		Contexts: []*proto.Context{
			versionCtx,
		},
	})
	if err != nil {
		return nil, err
	}
	modelVersionId := versionCtxResp.GetContextIds()[0]

	// add explicit association between model version and registered model
	_, err = serv.mlmdClient.PutParentContexts(context.Background(), &proto.PutParentContextsRequest{
		ParentContexts: []*proto.ParentContext{
			{
				ParentId: registeredModelId,
				ChildId:  &modelVersionId,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// create artifacts associated to the newly created model version
	artifacts, err := serv.mapper.MapFromModelArtifacts(modelVersion.Artifacts)
	if err != nil {
		return nil, err
	}

	artifactsResp, err := serv.mlmdClient.PutArtifacts(context.Background(), &proto.PutArtifactsRequest{
		Artifacts: artifacts,
	})
	if err != nil {
		return nil, err
	}

	// add explicit association between artifacts and model version
	attributions := []*proto.Attribution{}
	for _, a := range artifactsResp.ArtifactIds {
		attributions = append(attributions, &proto.Attribution{
			ContextId:  &modelVersionId,
			ArtifactId: &a,
		})
	}

	_, err = serv.mlmdClient.PutAttributionsAndAssociations(context.Background(), &proto.PutAttributionsAndAssociationsRequest{
		Attributions: attributions,
		Associations: make([]*proto.Association, 0),
	})
	if err != nil {
		return nil, err
	}

	return &modelVersionId, nil
}

// Get a specific model version
func (serv *modelRegistryService) GetModelVersion(name string, version string) (*registry.VersionedModel, error) {
	log.Printf("Getting model version %s:%s", name, version)

	// get the model version as context instance
	query := fmt.Sprintf("type = \"%s\" and properties.model_name.string_value = \"%s\" and properties.version.string_value = \"%s\"", modelVersionType, name, version)
	resp, err := serv.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{
		Options: &proto.ListOperationOptions{
			FilterQuery: &query,
		},
	})
	if err != nil {
		return nil, err
	}

	ctx := resp.Contexts[0]

	// get all artifacts associated to that model version
	artResp, err := serv.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
		ContextId: ctx.Id,
	})
	if err != nil {
		return nil, err
	}

	modelVersion, err := serv.mapper.MapToModelVersion(ctx, artResp.Artifacts)
	if err != nil {
		return nil, err
	}

	return modelVersion, nil
}

// Additional Info

// Create a new deployment for a specific model version
func (serv *modelRegistryService) CreateDeployment(name string, version string, env *registry.Environment) error {
	// TODO implement me
	panic("implement me")
}
