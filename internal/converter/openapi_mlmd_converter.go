package converter

import (
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
)

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
	// goverter:map . Type | MapRegisteredModelType
	// goverter:map . Properties | MapRegisteredModelProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch TypeId Type
	ConvertRegisteredModel(source *openapi.RegisteredModel) (*proto.Context, error)

	// TODO: note that we don't know the registeredModel here, therefore name cannot be prefixed at this step
	// goverter:map . Type | MapModelVersionType
	// goverter:map . Properties | MapModelVersionProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch TypeId Type
	ConvertModelVersion(source *openapi.ModelVersion) (*proto.Context, error)

	// TODO: note that we don't know the modelVersion here, therefore name cannot be prefixed at this step
	// goverter:map . Type | MapModelArtifactType
	// goverter:map . Properties | MapModelArtifactProperties
	// goverter:map State | MapOpenAPIModelArtifactState
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch TypeId Type
	ConvertModelArtifact(source *openapi.ModelArtifact) (*proto.Artifact, error)
}
