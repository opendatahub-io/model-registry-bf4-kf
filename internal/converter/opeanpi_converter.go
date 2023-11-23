package converter

import "github.com/opendatahub-io/model-registry/pkg/openapi"

type OpenapiUpdateWrapper[
	M openapi.RegisteredModel |
		openapi.ModelVersion |
		openapi.ModelArtifact |
		openapi.ServingEnvironment |
		openapi.InferenceService |
		openapi.ServeModel,
] struct {
	Existing *M
	Update   *M
}

func NewOpenapiUpdateWrapper[
	M openapi.RegisteredModel |
		openapi.ModelVersion |
		openapi.ModelArtifact |
		openapi.ServingEnvironment |
		openapi.InferenceService |
		openapi.ServeModel,
](existing *M, update *M) OpenapiUpdateWrapper[M] {
	return OpenapiUpdateWrapper[M]{
		Existing: existing,
		Update:   update,
	}
}

// goverter:converter
// goverter:output:file ./generated/openapi_converter.gen.go
// goverter:wrapErrors
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
type OpenAPIConverter interface {
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertRegisteredModelCreate(source *openapi.RegisteredModelCreate) (*openapi.RegisteredModel, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	ConvertRegisteredModelUpdate(source *openapi.RegisteredModelUpdate) (*openapi.RegisteredModel, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertModelVersionCreate(source *openapi.ModelVersionCreate) (*openapi.ModelVersion, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	ConvertModelVersionUpdate(source *openapi.ModelVersionUpdate) (*openapi.ModelVersion, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType
	ConvertModelArtifactCreate(source *openapi.ModelArtifactCreate) (*openapi.ModelArtifact, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType Name
	ConvertModelArtifactUpdate(source *openapi.ModelArtifactUpdate) (*openapi.ModelArtifact, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertServingEnvironmentCreate(source *openapi.ServingEnvironmentCreate) (*openapi.ServingEnvironment, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	ConvertServingEnvironmentUpdate(source *openapi.ServingEnvironmentUpdate) (*openapi.ServingEnvironment, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertInferenceServiceCreate(source *openapi.InferenceServiceCreate) (*openapi.InferenceService, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name RegisteredModelId ServingEnvironmentId
	ConvertInferenceServiceUpdate(source *openapi.InferenceServiceUpdate) (*openapi.InferenceService, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertServeModelCreate(source *openapi.ServeModelCreate) (*openapi.ServeModel, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ModelVersionId
	ConvertServeModelUpdate(source *openapi.ServeModelUpdate) (*openapi.ServeModel, error)

	// Ignore all fields that ARE editable
	// goverter:default InitRegisteredModelWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalID CustomProperties
	OverrideNotEditableForRegisteredModel(source OpenapiUpdateWrapper[openapi.RegisteredModel]) (openapi.RegisteredModel, error)

	// Ignore all fields that ARE editable
	// goverter:default InitModelVersionWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalID CustomProperties
	OverrideNotEditableForModelVersion(source OpenapiUpdateWrapper[openapi.ModelVersion]) (openapi.ModelVersion, error)

	// Ignore all fields that ARE editable
	// goverter:default InitModelArtifactWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalID CustomProperties Uri State ServiceAccountName ModelFormatName ModelFormatVersion StorageKey StoragePath
	OverrideNotEditableForModelArtifact(source OpenapiUpdateWrapper[openapi.ModelArtifact]) (openapi.ModelArtifact, error)

	// Ignore all fields that ARE editable
	// goverter:default InitServingEnvironmentWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalID CustomProperties
	OverrideNotEditableForServingEnvironment(source OpenapiUpdateWrapper[openapi.ServingEnvironment]) (openapi.ServingEnvironment, error)

	// Ignore all fields that ARE editable
	// goverter:default InitInferenceServiceWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalID CustomProperties ModelVersionId Runtime State
	OverrideNotEditableForInferenceService(source OpenapiUpdateWrapper[openapi.InferenceService]) (openapi.InferenceService, error)

	// Ignore all fields that ARE editable
	// goverter:default InitServeModelWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalID CustomProperties LastKnownState
	OverrideNotEditableForServeModel(source OpenapiUpdateWrapper[openapi.ServeModel]) (openapi.ServeModel, error)
}

func InitRegisteredModelWithUpdate(source OpenapiUpdateWrapper[openapi.RegisteredModel]) openapi.RegisteredModel {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.RegisteredModel{}
}

func InitModelVersionWithUpdate(source OpenapiUpdateWrapper[openapi.ModelVersion]) openapi.ModelVersion {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.ModelVersion{}
}

func InitModelArtifactWithUpdate(source OpenapiUpdateWrapper[openapi.ModelArtifact]) openapi.ModelArtifact {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.ModelArtifact{}
}

func InitServingEnvironmentWithUpdate(source OpenapiUpdateWrapper[openapi.ServingEnvironment]) openapi.ServingEnvironment {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.ServingEnvironment{}
}

func InitInferenceServiceWithUpdate(source OpenapiUpdateWrapper[openapi.InferenceService]) openapi.InferenceService {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.InferenceService{}
}

func InitServeModelWithUpdate(source OpenapiUpdateWrapper[openapi.ServeModel]) openapi.ServeModel {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.ServeModel{}
}
