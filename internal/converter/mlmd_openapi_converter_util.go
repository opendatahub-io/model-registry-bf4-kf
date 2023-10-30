package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
)

// MapMLMDCustomProperties maps MLMD custom properties model to OpenAPI one
func MapMLMDCustomProperties(source map[string]*proto.Value) (map[string]openapi.MetadataValue, error) {
	data := make(map[string]openapi.MetadataValue)

	for key, v := range source {
		// data[key] = v.Value
		customValue := openapi.MetadataValue{}

		switch typedValue := v.Value.(type) {
		case *proto.Value_BoolValue:
			customValue.MetadataBoolValue = &openapi.MetadataBoolValue{
				BoolValue: &typedValue.BoolValue,
			}
		case *proto.Value_IntValue:
			customValue.MetadataIntValue = &openapi.MetadataIntValue{
				IntValue: Int64ToString(&typedValue.IntValue),
			}
		case *proto.Value_DoubleValue:
			customValue.MetadataDoubleValue = &openapi.MetadataDoubleValue{
				DoubleValue: &typedValue.DoubleValue,
			}
		case *proto.Value_StringValue:
			customValue.MetadataStringValue = &openapi.MetadataStringValue{
				StringValue: &typedValue.StringValue,
			}
		case *proto.Value_StructValue:
			sv := typedValue.StructValue
			asMap := sv.AsMap()
			asJSON, err := json.Marshal(asMap)
			if err != nil {
				return nil, err
			}
			b64 := base64.StdEncoding.EncodeToString(asJSON)
			customValue.MetadataStructValue = &openapi.MetadataStructValue{
				StructValue: &b64,
			}
		default:
			return nil, fmt.Errorf("type mapping not found for %s:%v", key, v)
		}

		data[key] = customValue
	}

	return data, nil
}

// MapNameFromOwned derive the entity name from the mlmd fullname
// for owned entity such as ModelVersion
// for potentially owned entity such as ModelArtifact
func MapNameFromOwned(source *string) *string {
	if source == nil {
		return nil
	}

	exploded := strings.Split(*source, ":")
	if len(exploded) == 1 {
		return source
	}
	return &exploded[1]
}

// REGISTERED MODEL

// MODEL VERSION

// MODEL ARTIFACT

func MapArtifactType(source *proto.Artifact) (string, error) {
	if *source.Type == ModelArtifactTypeName {
		return "model-artifact", nil
	}
	return "", fmt.Errorf("invalid artifact type found")
}

func MapModelFormatName(source *proto.Artifact) *string {
	var value string
	property := source.Properties["model_format"]
	if property != nil {
		value = property.GetStringValue()
	}

	if value == "" {
		return nil
	}
	return &value
}

func MapMLMDModelArtifactState(source *proto.Artifact_State) *openapi.ArtifactState {
	if source == nil {
		return nil
	}

	state := source.String()
	return (*openapi.ArtifactState)(&state)
}

// MapModelArtifactProperties maps ModelArtifact fields to specific MLMD properties

// MapModelArtifactType returnd ModelArtifact corresponding MLMD context type
