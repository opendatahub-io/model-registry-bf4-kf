package core

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/opendatahub-io/model-registry/internal/converter"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"google.golang.org/grpc"
)

var (
	registeredModelTypeName    = of(converter.RegisteredModelTypeName)
	modelVersionTypeName       = of(converter.ModelVersionTypeName)
	modelArtifactTypeName      = of(converter.ModelArtifactTypeName)
	servingEnvironmentTypeName = of(converter.ServingEnvironmentTypeName)
	inferenceServiceTypeName   = of(converter.InferenceServiceTypeName)
	serveModelTypeName         = of(converter.ServeModelTypeName)
)

// modelRegistryService is the core library of the model registry
type modelRegistryService struct {
	mlmdClient proto.MetadataStoreServiceClient
	mapper     *Mapper
}

// NewModelRegistryService create a fresh instance of ModelRegistryService, taking care of setting up needed MLMD Types
func NewModelRegistryService(cc grpc.ClientConnInterface) (ModelRegistryApi, error) {

	client := proto.NewMetadataStoreServiceClient(cc)

	// Setup the needed Type instances if not existing already

	registeredModelReq := proto.PutContextTypeRequest{
		ContextType: &proto.ContextType{
			Name: registeredModelTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
			},
		},
	}

	modelVersionReq := proto.PutContextTypeRequest{
		ContextType: &proto.ContextType{
			Name: modelVersionTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
				"model_name":  proto.PropertyType_STRING,
				"version":     proto.PropertyType_STRING,
				"author":      proto.PropertyType_STRING,
			},
		},
	}

	modelArtifactReq := proto.PutArtifactTypeRequest{
		ArtifactType: &proto.ArtifactType{
			Name: modelArtifactTypeName,
			Properties: map[string]proto.PropertyType{
				"description":          proto.PropertyType_STRING,
				"runtime":              proto.PropertyType_STRING,
				"model_format_name":    proto.PropertyType_STRING,
				"model_format_version": proto.PropertyType_STRING,
				"storage_key":          proto.PropertyType_STRING,
				"storage_path":         proto.PropertyType_STRING,
				"service_account_name": proto.PropertyType_STRING,
			},
		},
	}

	servingEnvironmentReq := proto.PutContextTypeRequest{
		ContextType: &proto.ContextType{
			Name: servingEnvironmentTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
			},
		},
	}

	inferenceServiceReq := proto.PutContextTypeRequest{
		ContextType: &proto.ContextType{
			Name: inferenceServiceTypeName,
			Properties: map[string]proto.PropertyType{
				"description":      proto.PropertyType_STRING,
				"model_version_id": proto.PropertyType_INT,
				// we could remove this as we will use ParentContext to keep track of this association
				"registered_model_id":    proto.PropertyType_INT,
				"serving_environment_id": proto.PropertyType_INT,
			},
		},
	}

	serveModelReq := proto.PutExecutionTypeRequest{
		ExecutionType: &proto.ExecutionType{
			Name: serveModelTypeName,
			Properties: map[string]proto.PropertyType{
				"description": proto.PropertyType_STRING,
				// we could remove this as we will use ParentContext to keep track of this association
				"model_version_id": proto.PropertyType_INT,
			},
		},
	}

	registeredModelResp, err := client.PutContextType(context.Background(), &registeredModelReq)
	if err != nil {
		glog.Fatalf("Error setting up context type %s: %v", *registeredModelTypeName, err)
	}

	modelVersionResp, err := client.PutContextType(context.Background(), &modelVersionReq)
	if err != nil {
		glog.Fatalf("Error setting up context type %s: %v", *modelVersionTypeName, err)
	}

	modelArtifactResp, err := client.PutArtifactType(context.Background(), &modelArtifactReq)
	if err != nil {
		glog.Fatalf("Error setting up artifact type %s: %v", *modelArtifactTypeName, err)
	}

	servingEnvironmentResp, err := client.PutContextType(context.Background(), &servingEnvironmentReq)
	if err != nil {
		glog.Fatalf("Error setting up context type %s: %v", *servingEnvironmentTypeName, err)
	}

	inferenceServiceResp, err := client.PutContextType(context.Background(), &inferenceServiceReq)
	if err != nil {
		glog.Fatalf("Error setting up context type %s: %v", *inferenceServiceTypeName, err)
	}

	serveModelResp, err := client.PutExecutionType(context.Background(), &serveModelReq)
	if err != nil {
		glog.Fatalf("Error setting up execution type %s: %v", *serveModelTypeName, err)
	}

	return &modelRegistryService{
		mlmdClient: client,
		mapper: NewMapper(
			registeredModelResp.GetTypeId(),
			modelVersionResp.GetTypeId(),
			modelArtifactResp.GetTypeId(),
			servingEnvironmentResp.GetTypeId(),
			inferenceServiceResp.GetTypeId(),
			serveModelResp.GetTypeId(),
		),
	}, nil
}

// REGISTERED MODELS

func (serv *modelRegistryService) UpsertRegisteredModel(registeredModel *openapi.RegisteredModel) (*openapi.RegisteredModel, error) {
	var err error
	var existing *openapi.RegisteredModel

	if registeredModel.Id == nil {
		glog.Info("Creating new registered model")
	} else {
		glog.Info("Updating registered model %s", *registeredModel.Id)
		existing, err = serv.GetRegisteredModelById(*registeredModel.Id)
		if err != nil {
			return nil, err
		}
	}

	// if already existing assure the name is the same
	if existing != nil && registeredModel.Name == nil {
		// user did not provide it
		// need to set it to avoid mlmd error "context name should not be empty"
		registeredModel.Name = existing.Name
	}

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

	idAsString := converter.Int64ToString(&modelCtxResp.ContextIds[0])
	model, err := serv.GetRegisteredModelById(*idAsString)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (serv *modelRegistryService) GetRegisteredModelById(id string) (*openapi.RegisteredModel, error) {
	glog.Info("Getting registered model %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple registered models found for id %s", id)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no registered model found for id %s", id)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return regModel, nil
}

func (serv *modelRegistryService) GetRegisteredModelByInferenceService(inferenceServiceId string) (*openapi.RegisteredModel, error) {
	panic("method not yet implemented")
}

func (serv *modelRegistryService) getRegisteredModelByVersionId(id string) (*openapi.RegisteredModel, error) {
	glog.Info("Getting registered model for model version %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getParentResp, err := serv.mlmdClient.GetParentContextsByContext(context.Background(), &proto.GetParentContextsByContextRequest{
		ContextId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple registered models found for model version %s", id)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no registered model found for model version %s", id)
	}

	regModel, err := serv.mapper.MapToRegisteredModel(getParentResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return regModel, nil
}

func (serv *modelRegistryService) GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error) {
	glog.Info("Getting registered model by params name=%v, externalId=%v", name, externalId)

	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: registeredModelTypeName,
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

func (serv *modelRegistryService) GetRegisteredModels(listOptions ListOptions) (*openapi.RegisteredModelList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}
	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: registeredModelTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.RegisteredModel{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToRegisteredModel(c)
		if err != nil {
			return nil, err
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.RegisteredModelList{
		NextPageToken: zeroIfNil(contextsResp.NextPageToken),
		PageSize:      zeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// MODEL VERSIONS

func (serv *modelRegistryService) UpsertModelVersion(modelVersion *openapi.ModelVersion, parentResourceId *string) (*openapi.ModelVersion, error) {
	var err error
	var existing *openapi.ModelVersion
	var registeredModel *openapi.RegisteredModel

	if modelVersion.Id == nil {
		// create
		glog.Info("Creating new model version")
		if parentResourceId == nil {
			return nil, fmt.Errorf("missing registered model id, cannot create model version without registered model")
		}
		registeredModel, err = serv.GetRegisteredModelById(*parentResourceId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Info("Updating model version %s", *modelVersion.Id)
		existing, err = serv.GetModelVersionById(*modelVersion.Id)
		if err != nil {
			return nil, err
		}
		registeredModel, err = serv.getRegisteredModelByVersionId(*modelVersion.Id)
		if err != nil {
			return nil, err
		}
	}

	// if already existing assure the name is the same
	if existing != nil && modelVersion.Name == nil {
		// user did not provide it
		// need to set it to avoid mlmd error "context name should not be empty"
		modelVersion.Name = existing.Name
	}

	modelCtx, err := serv.mapper.MapFromModelVersion(modelVersion, *registeredModel.Id, registeredModel.Name)
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
	if modelVersion.Id == nil {
		registeredModelId, err := converter.StringToInt64(registeredModel.Id)
		if err != nil {
			return nil, err
		}

		_, err = serv.mlmdClient.PutParentContexts(context.Background(), &proto.PutParentContextsRequest{
			ParentContexts: []*proto.ParentContext{{
				ChildId:  modelId,
				ParentId: registeredModelId}},
			TransactionOptions: &proto.TransactionOptions{},
		})
		if err != nil {
			return nil, err
		}
	}

	idAsString := converter.Int64ToString(modelId)
	model, err := serv.GetModelVersionById(*idAsString)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (serv *modelRegistryService) GetModelVersionById(id string) (*openapi.ModelVersion, error) {
	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple model versions found for id %s", id)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no model version found for id %s", id)
	}

	modelVer, err := serv.mapper.MapToModelVersion(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return modelVer, nil
}

func (serv *modelRegistryService) GetModelVersionByInferenceService(inferenceServiceId string) (*openapi.ModelVersion, error) {
	panic("method not yet implemented")
}

func (serv *modelRegistryService) getModelVersionByArtifactId(id string) (*openapi.ModelVersion, error) {
	glog.Info("Getting model version for model artifact %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getParentResp, err := serv.mlmdClient.GetContextsByArtifact(context.Background(), &proto.GetContextsByArtifactRequest{
		ArtifactId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple model versions found for model artifact %s", id)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no model version found for model artifact %s", id)
	}

	modelVersion, err := serv.mapper.MapToModelVersion(getParentResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return modelVersion, nil
}

func (serv *modelRegistryService) GetModelVersionByParams(versionName *string, parentResourceId *string, externalId *string) (*openapi.ModelVersion, error) {
	filterQuery := ""
	if versionName != nil && parentResourceId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(parentResourceId, *versionName))
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (versionName and parentResourceId), or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: modelVersionTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) != 1 {
		return nil, fmt.Errorf("multiple model versions found for versionName=%v, parentResourceId=%v, externalId=%v", zeroIfNil(versionName), zeroIfNil(parentResourceId), zeroIfNil(externalId))
	}

	modelVer, err := serv.mapper.MapToModelVersion(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return modelVer, nil
}

func (serv *modelRegistryService) GetModelVersions(listOptions ListOptions, parentResourceId *string) (*openapi.ModelVersionList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	if parentResourceId != nil {
		queryParentCtxId := fmt.Sprintf("parent_contexts_a.id = %s", *parentResourceId)
		listOperationOptions.FilterQuery = &queryParentCtxId
	}

	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: modelVersionTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.ModelVersion{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToModelVersion(c)
		if err != nil {
			return nil, err
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ModelVersionList{
		NextPageToken: zeroIfNil(contextsResp.NextPageToken),
		PageSize:      zeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// MODEL ARTIFACTS

func (serv *modelRegistryService) UpsertModelArtifact(modelArtifact *openapi.ModelArtifact, parentResourceId *string) (*openapi.ModelArtifact, error) {
	var err error
	var existing *openapi.ModelArtifact

	if modelArtifact.Id == nil {
		// create
		glog.Info("Creating new model artifact")
		if parentResourceId == nil {
			return nil, fmt.Errorf("missing model version id, cannot create model artifact without model version")
		}
		_, err = serv.GetModelVersionById(*parentResourceId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Info("Updating model artifact %s", *modelArtifact.Id)
		existing, err = serv.GetModelArtifactById(*modelArtifact.Id)
		if err != nil {
			return nil, err
		}
		_, err = serv.getModelVersionByArtifactId(*modelArtifact.Id)
		if err != nil {
			return nil, err
		}
	}

	// if already existing assure the name is the same
	if existing != nil {
		if modelArtifact.Name == nil {
			// user did not provide it
			// need to set it to avoid mlmd error "artifact name should not be empty"
			modelArtifact.Name = existing.Name
		}
	}

	artifact, err := serv.mapper.MapFromModelArtifact(modelArtifact, parentResourceId)
	if err != nil {
		return nil, err
	}

	artifactsResp, err := serv.mlmdClient.PutArtifacts(context.Background(), &proto.PutArtifactsRequest{
		Artifacts: []*proto.Artifact{artifact},
	})
	if err != nil {
		return nil, err
	}

	// add explicit association between artifacts and model version
	if parentResourceId != nil && modelArtifact.Id == nil {
		modelVersionId, err := converter.StringToInt64(parentResourceId)
		if err != nil {
			return nil, err
		}
		attributions := []*proto.Attribution{}
		for _, a := range artifactsResp.ArtifactIds {
			attributions = append(attributions, &proto.Attribution{
				ContextId:  modelVersionId,
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
	}

	idAsString := converter.Int64ToString(&artifactsResp.ArtifactIds[0])
	mapped, err := serv.GetModelArtifactById(*idAsString)
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

func (serv *modelRegistryService) GetModelArtifactById(id string) (*openapi.ModelArtifact, error) {
	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	artifactsResp, err := serv.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}

	if len(artifactsResp.Artifacts) > 1 {
		return nil, fmt.Errorf("multiple model artifacts found for id %s", id)
	}

	if len(artifactsResp.Artifacts) == 0 {
		return nil, fmt.Errorf("no model artifact found for id %s", id)
	}

	result, err := serv.mapper.MapToModelArtifact(artifactsResp.Artifacts[0])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (serv *modelRegistryService) GetModelArtifactByParams(artifactName *string, parentResourceId *string, externalId *string) (*openapi.ModelArtifact, error) {
	var artifact0 *proto.Artifact

	filterQuery := ""
	if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else if artifactName != nil && parentResourceId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(parentResourceId, *artifactName))
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and parentResourceId), or externalId")
	}

	artifactsResponse, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: modelArtifactTypeName,
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

func (serv *modelRegistryService) GetModelArtifacts(listOptions ListOptions, parentResourceId *string) (*openapi.ModelArtifactList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if parentResourceId != nil {
		ctxId, err := converter.StringToInt64(parentResourceId)
		if err != nil {
			return nil, err
		}
		artifactsResp, err := serv.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
			ContextId: ctxId,
			Options:   listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	} else {
		artifactsResp, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
			TypeName: modelArtifactTypeName,
			Options:  listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	}

	results := []openapi.ModelArtifact{}
	for _, a := range artifacts {
		mapped, err := serv.mapper.MapToModelArtifact(a)
		if err != nil {
			return nil, err
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ModelArtifactList{
		NextPageToken: zeroIfNil(nextPageToken),
		PageSize:      zeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// SERVING ENVIRONMENT

func (serv *modelRegistryService) UpsertServingEnvironment(servingEnvironment *openapi.ServingEnvironment) (*openapi.ServingEnvironment, error) {
	var err error
	var existing *openapi.ServingEnvironment

	if servingEnvironment.Id == nil {
		glog.Info("Creating new serving environment")
	} else {
		glog.Info("Updating serving environment %s", *servingEnvironment.Id)
		existing, err = serv.GetServingEnvironmentById(*servingEnvironment.Id)
		if err != nil {
			return nil, err
		}
	}

	// if already existing assure the name is the same
	if existing != nil && servingEnvironment.Name == nil {
		// user did not provide it
		// need to set it to avoid mlmd error "context name should not be empty"
		servingEnvironment.Name = existing.Name
	}

	protoCtx, err := serv.mapper.MapFromServingEnvironment(servingEnvironment)
	if err != nil {
		return nil, err
	}

	protoCtxResp, err := serv.mlmdClient.PutContexts(context.Background(), &proto.PutContextsRequest{
		Contexts: []*proto.Context{
			protoCtx,
		},
	})
	if err != nil {
		return nil, err
	}

	idAsString := converter.Int64ToString(&protoCtxResp.ContextIds[0])
	openapiModel, err := serv.GetServingEnvironmentById(*idAsString)
	if err != nil {
		return nil, err
	}

	return openapiModel, nil
}

func (serv *modelRegistryService) GetServingEnvironmentById(id string) (*openapi.ServingEnvironment, error) {
	glog.Info("Getting serving environment %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*idAsInt},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple serving environments found for id %s", id)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no serving environment found for id %s", id)
	}

	openapiModel, err := serv.mapper.MapToServingEnvironment(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return openapiModel, nil
}

func (serv *modelRegistryService) GetServingEnvironmentByParams(name *string, externalId *string) (*openapi.ServingEnvironment, error) {
	glog.Info("Getting serving environment by params name=%v, externalId=%v", name, externalId)

	filterQuery := ""
	if name != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", *name)
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either name or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: servingEnvironmentTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) != 1 {
		return nil, fmt.Errorf("could not find exactly one Context matching criteria: %v", getByParamsResp.Contexts)
	}

	openapiModel, err := serv.mapper.MapToServingEnvironment(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return openapiModel, nil
}

func (serv *modelRegistryService) GetServingEnvironments(listOptions ListOptions) (*openapi.ServingEnvironmentList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}
	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: servingEnvironmentTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.ServingEnvironment{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToServingEnvironment(c)
		if err != nil {
			return nil, err
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ServingEnvironmentList{
		NextPageToken: zeroIfNil(contextsResp.NextPageToken),
		PageSize:      zeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// INFERENCE SERVICE

func (serv *modelRegistryService) UpsertInferenceService(inferenceService *openapi.InferenceService) (*openapi.InferenceService, error) {
	var err error
	var existing *openapi.InferenceService
	var servingEnvironment *openapi.ServingEnvironment
	// for InferenceService, is part of model payload.
	parentResourceId := inferenceService.ServingEnvironmentId

	if inferenceService.Id == nil {
		// create
		glog.Info("Creating new InferenceService")
		servingEnvironment, err = serv.GetServingEnvironmentById(parentResourceId)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		glog.Info("Updating InferenceService %s", *inferenceService.Id)
		existing, err = serv.GetInferenceServiceById(*inferenceService.Id)
		if err != nil {
			return nil, err
		}
		servingEnvironment, err = serv.getServingEnvironmentByInferenceServiceId(*inferenceService.Id)
		if err != nil {
			return nil, err
		}
	}

	// validate RegisteredModelId is also valid
	if _, err := serv.GetRegisteredModelById(inferenceService.RegisteredModelId); err != nil {
		return nil, err
	}

	// if already existing assure the name is the same
	if existing != nil && inferenceService.Name == nil {
		// user did not provide it
		// need to set it to avoid mlmd error "context name should not be empty"
		inferenceService.Name = existing.Name
	}

	protoCtx, err := serv.mapper.MapFromInferenceService(inferenceService, *servingEnvironment.Id, servingEnvironment.Name)
	if err != nil {
		return nil, err
	}

	protoCtxResp, err := serv.mlmdClient.PutContexts(context.Background(), &proto.PutContextsRequest{
		Contexts: []*proto.Context{
			protoCtx,
		},
	})
	if err != nil {
		return nil, err
	}

	inferenceServiceId := &protoCtxResp.ContextIds[0]
	if inferenceService.Id == nil {
		servingEnvironmentId, err := converter.StringToInt64(servingEnvironment.Id)
		if err != nil {
			return nil, err
		}

		_, err = serv.mlmdClient.PutParentContexts(context.Background(), &proto.PutParentContextsRequest{
			ParentContexts: []*proto.ParentContext{{
				ChildId:  inferenceServiceId,
				ParentId: servingEnvironmentId}},
			TransactionOptions: &proto.TransactionOptions{},
		})
		if err != nil {
			return nil, err
		}
	}

	idAsString := converter.Int64ToString(inferenceServiceId)
	toReturn, err := serv.GetInferenceServiceById(*idAsString)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (serv *modelRegistryService) getServingEnvironmentByInferenceServiceId(id string) (*openapi.ServingEnvironment, error) {
	glog.Info("Getting ServingEnvironment for InferenceService %s", id)

	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getParentResp, err := serv.mlmdClient.GetParentContextsByContext(context.Background(), &proto.GetParentContextsByContextRequest{
		ContextId: idAsInt,
	})
	if err != nil {
		return nil, err
	}

	if len(getParentResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple ServingEnvironments found for InferenceService %s", id)
	}

	if len(getParentResp.Contexts) == 0 {
		return nil, fmt.Errorf("no ServingEnvironments found for InferenceService %s", id)
	}

	toReturn, err := serv.mapper.MapToServingEnvironment(getParentResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (serv *modelRegistryService) GetInferenceServiceById(id string) (*openapi.InferenceService, error) {
	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, err
	}

	getByIdResp, err := serv.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*idAsInt},
	})
	if err != nil {
		return nil, err
	}

	if len(getByIdResp.Contexts) > 1 {
		return nil, fmt.Errorf("multiple InferenceServices found for id %s", id)
	}

	if len(getByIdResp.Contexts) == 0 {
		return nil, fmt.Errorf("no InferenceService found for id %s", id)
	}

	toReturn, err := serv.mapper.MapToInferenceService(getByIdResp.Contexts[0])
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (serv *modelRegistryService) GetInferenceServiceByParams(name *string, parentResourceId *string, externalId *string) (*openapi.InferenceService, error) {
	filterQuery := ""
	if name != nil && parentResourceId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(parentResourceId, *name))
	} else if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (name and parentResourceId), or externalId")
	}

	getByParamsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: inferenceServiceTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(getByParamsResp.Contexts) != 1 {
		return nil, fmt.Errorf("multiple InferenceServices found for name=%v, parentResourceId=%v, externalId=%v", zeroIfNil(name), zeroIfNil(parentResourceId), zeroIfNil(externalId))
	}

	toReturn, err := serv.mapper.MapToInferenceService(getByParamsResp.Contexts[0])
	if err != nil {
		return nil, err
	}
	return toReturn, nil
}

func (serv *modelRegistryService) GetInferenceServices(listOptions ListOptions, parentResourceId *string) (*openapi.InferenceServiceList, error) {
	listOperationOptions, err := BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, err
	}

	if parentResourceId != nil {
		queryParentCtxId := fmt.Sprintf("parent_contexts_a.id = %s", *parentResourceId)
		listOperationOptions.FilterQuery = &queryParentCtxId
	}

	contextsResp, err := serv.mlmdClient.GetContextsByType(context.Background(), &proto.GetContextsByTypeRequest{
		TypeName: inferenceServiceTypeName,
		Options:  listOperationOptions,
	})
	if err != nil {
		return nil, err
	}

	results := []openapi.InferenceService{}
	for _, c := range contextsResp.Contexts {
		mapped, err := serv.mapper.MapToInferenceService(c)
		if err != nil {
			return nil, err
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.InferenceServiceList{
		NextPageToken: zeroIfNil(contextsResp.NextPageToken),
		PageSize:      zeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// SERVE MODEL

func (serv *modelRegistryService) UpsertServeModel(registeredModel *openapi.ServeModel, inferenceServiceId *string) (*openapi.ServeModel, error) {
	panic("method not yet implemented")
}

func (serv *modelRegistryService) GetServeModelById(id string) (*openapi.ServeModel, error) {
	panic("method not yet implemented")
}

func (serv *modelRegistryService) GetServeModels(listOptions ListOptions, inferenceServiceId *string) (*openapi.ServeModelList, error) {
	panic("method not yet implemented")
}

// of returns a pointer to the provided literal/const input
func of[E any](e E) *E {
	return &e
}
