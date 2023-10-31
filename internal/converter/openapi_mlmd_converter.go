package converter

import (
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
)

type OpenAPIModelWrapper[M openapi.RegisteredModel | openapi.ModelVersion | openapi.ModelArtifact] struct {
	TypeId           int64
	Model            *M
	ParentResourceId *string // optional parent id
	ModelName        *string // optional registered model name
}

// TODO: we don't know the TypeId for MLMD types at thsi step, need to find out how to retrieve it
// goverter:converter
// goverter:output:file ./generated/openapi_mlmd_converter.gen.go
// goverter:wrapErrors
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend Int64ToString
// goverter:extend StringToInt64
// goverter:extend MapOpenAPICustomProperties
type OpenAPIToMLMDConverter interface {
	// goverter:autoMap Model
	// goverter:map Model Type | MapRegisteredModelType
	// goverter:map Model Properties | MapRegisteredModelProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertRegisteredModel(source *OpenAPIModelWrapper[openapi.RegisteredModel]) (*proto.Context, error)

	// TODO: note that we don't know the registeredModel here, therefore name cannot be prefixed at this step
	// goverter:autoMap Model
	// goverter:map . Name | MapModelVersionName
	// goverter:map Model Type | MapModelVersionType
	// goverter:map . Properties | MapModelVersionProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertModelVersion(source *OpenAPIModelWrapper[openapi.ModelVersion]) (*proto.Context, error)

	// TODO: note that we don't know the modelVersion here, therefore name cannot be prefixed at this step
	// goverter:autoMap Model
	// goverter:map . Name | MapModelArtifactName
	// goverter:map Model Type | MapModelArtifactType
	// goverter:map Model Properties | MapModelArtifactProperties
	// goverter:map Model.State State | MapOpenAPIModelArtifactState
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertModelArtifact(source *OpenAPIModelWrapper[openapi.ModelArtifact]) (*proto.Artifact, error)
}
