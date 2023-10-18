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

// modelRegistryService is the core library of the model registry
type modelRegistryService struct {
	mlmdClient proto.MetadataStoreServiceClient
	mapper     *mapper.Mapper
}

// NewModelRegistryService create a fresh instance of ModelRegistryService, taking care of setting up needed MLMD Types
func NewModelRegistryService(cc grpc.ClientConnInterface) (ModelRegistryApi, error) {

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

// REGISTERED MODELS

func (serv *modelRegistryService) UpsertRegisteredModel(registeredModel *registry.RegisteredModel) (*registry.RegisteredModel, error) {
	log.Printf("Creating or updating registered model for %s", *registeredModel.Name)

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

	modelId := &modelCtxResp.ContextIds[0]
	model, err := serv.GetRegisteredModelById((*BaseResourceId)(modelId))
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (serv *modelRegistryService) GetRegisteredModelById(id *BaseResourceId) (*registry.RegisteredModel, error) {
	log.Printf("Getting registered model %d", *id)

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{int64(*id)},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) != 1 {
		return nil, fmt.Errorf("multiple registered models found for id %d", *id)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return regModel, nil
}

func (serv *modelRegistryService) GetRegisteredModelByParams(name *string, externalId *string) (*registry.RegisteredModel, error) {
	log.Printf("Getting registered model by params name=%v, externalId=%v", name, externalId)

	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &registeredModelType,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) != 1 {
		return nil, fmt.Errorf("multiple registered models found for name=%v, externalId=%v", *name, *externalId)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return regModel, nil
}

func (serv *modelRegistryService) GetRegisteredModels(listOptions ListOptions) ([]*registry.RegisteredModel, ListResult, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, ListResult{}, err
	}
	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: &registeredModelType,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, ListResult{}, err
	}

	results := []*registry.RegisteredModel{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToRegisteredModel(c)
		if err != nil {
			return nil, ListResult{}, err
		}
		results = append(results, mapped)
	}

	listResult := NewListResult(contextsResp.Contexts, listOptions, contextsResp.NextPageToken)
	return results, listResult, nil
}

// MODEL VERSIONS

func (serv *modelRegistryService) UpsertModelVersion(modelVersion *registry.VersionedModel) (*registry.VersionedModel, error) {
	panic("Method not yet implemented")
}

func (serv *modelRegistryService) GetModelVersionById(id *BaseResourceId) (*registry.VersionedModel, error) {
	panic("Method not yet implemented")
}

// TODO: name not clear on OpenAPI, search by registeredModelName and versionName is missing - there is just unclear `name` param.
func (serv *modelRegistryService) GetModelVersionByParams(name *string, externalId *string) (*registry.VersionedModel, error) {
	panic("Method not yet implemented")
}

func (serv *modelRegistryService) GetModelVersions(listOptions ListOptions, registeredModelId *BaseResourceId) ([]*registry.VersionedModel, ListResult, error) {
	panic("Method not yet implemented")
}

// MODEL ARTIFACTS

func (serv *modelRegistryService) UpsertModelArtifact(modelArtifact *registry.Artifact) (*registry.Artifact, error) {
	artifact := serv.mapper.MapFromModelArtifact(*modelArtifact)

	artifactsResp, err := serv.mlmdClient.PutArtifacts(context.Background(), &proto.PutArtifactsRequest{
		Artifacts: []*proto.Artifact{artifact},
	})
	if err != nil {
		return nil, err
	}
	modelArtifact.Id = &artifactsResp.ArtifactIds[0]

	// add explicit association between artifacts and model version
	attributions := []*proto.Attribution{}
	for _, a := range artifactsResp.ArtifactIds {
		attributions = append(attributions, &proto.Attribution{
			ContextId:  modelArtifact.ModelVersionId,
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

	return modelArtifact, nil
}

func (serv *modelRegistryService) GetModelArtifactById(id *BaseResourceId) (*registry.Artifact, error) {
	artifactsResp, err := serv.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{int64(*id)},
	})
	if err != nil {
		return nil, err
	}

	result, err := serv.mapper.MapToModelArtifact(artifactsResp.Artifacts[0])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (serv *modelRegistryService) GetModelArtifactByParams(name *string, externalId *string) (*registry.Artifact, error) {
	var artifact0 *proto.Artifact

	filterQuery := ""
	if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else if name != nil { // TODO: see comment about `name` field in OpenAPI
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId")
	}

	artifactsResponse, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: &modelArtifactType,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(artifactsResponse.Artifacts) > 1 {
		return nil, fmt.Errorf("more than an artifact detected matching criteria: %v", artifactsResponse.Artifacts)
	}
	artifact0 = artifactsResponse.Artifacts[0]

	result, err := serv.mapper.MapToModelArtifact(artifact0)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (serv *modelRegistryService) GetModelArtifacts(listOptions ListOptions, modelVersionId *BaseResourceId) ([]*registry.Artifact, ListResult, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, ListResult{}, err
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if modelVersionId != nil {
		ctxId := int64(*modelVersionId)
		artifactsResp, err := serv.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
			ContextId: &ctxId,
			Options:   listOperationOptions,
		})
		if err != nil {
			return nil, ListResult{}, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	} else {
		artifactsResp, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
			TypeName: &modelArtifactType,
			Options:  listOperationOptions,
		})
		if err != nil {
			return nil, ListResult{}, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	}

	results := []*registry.Artifact{}
	for _, a := range artifacts {
		mapped, err := serv.mapper.MapToModelArtifact(a)
		if err != nil {
			return nil, ListResult{}, err
		}
		results = append(results, mapped)
	}

	listResult := NewListResult(artifacts, listOptions, nextPageToken)
	return results, listResult, nil
}
