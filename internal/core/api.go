package core

import "github.com/opendatahub-io/model-registry/internal/model/openapi"

type BaseResourceId int64

type ListOptions struct {
	PageSize      *int32
	OrderBy       *string
	SortOrder     *string
	NextPageToken *string
}

type ListResult struct {
	PageSize      *int32
	Size          *int32
	NextPageToken *string
}

// type ListRes[T any] struct {
// 	Items         []*T
// 	PageSize      *int32
// 	Size          *int32
// 	NextPageToken *string
// }

type ModelRegistryApi interface {
	// REGISTERED MODEL

	// UpsertRegisteredModel create or update a registered model, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	UpsertRegisteredModel(registeredModel *openapi.RegisteredModel) (*openapi.RegisteredModel, error)

	GetRegisteredModelById(id *BaseResourceId) (*openapi.RegisteredModel, error)
	GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error)
	GetRegisteredModels(listOptions ListOptions) ([]*openapi.RegisteredModel, ListResult, error)

	// MODEL VERSION

	// Create or update a Model Version associated to a specific RegisteredModel
	// identified by ModelVersion.RegisteredModelId
	UpsertModelVersion(modelVersion *openapi.ModelVersion) (*openapi.ModelVersion, error)

	GetModelVersionById(id *BaseResourceId) (*openapi.ModelVersion, error)
	// TODO: name not clear on OpenAPI, search by registeredModelName and versionName is missing - there is just unclear `name` param.
	GetModelVersionByParams(name *string, externalId *string) (*openapi.ModelVersion, error)
	GetModelVersions(listOptions ListOptions, registeredModelId *BaseResourceId) ([]*openapi.ModelVersion, ListResult, error)

	// MODEL ARTIFACT

	// Create or update a Model Artifact associated to a specific ModelVersion
	// identified by ModelArtifact.ModelVersionId
	UpsertModelArtifact(modelArtifact *openapi.ModelArtifact) (*openapi.ModelArtifact, error)

	GetModelArtifactById(id *BaseResourceId) (*openapi.ModelArtifact, error)
	// TODO: what is this name?
	GetModelArtifactByParams(name *string, externalId *string) (*openapi.ModelArtifact, error)
	GetModelArtifacts(listOptions ListOptions, modelVersionId *BaseResourceId) ([]*openapi.ModelArtifact, ListResult, error)
}
