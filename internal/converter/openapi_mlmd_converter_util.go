package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	RegisteredModelTypeName = "odh.RegisteredModel"
	ModelVersionTypeName    = "odh.ModelVersion"
	ModelArtifactTypeName   = "odh.ModelArtifact"
)

// StringToInt64 converts string-based id to int64 if numeric, otherwise return error
func StringToInt64(id *string) (*int64, error) {
	if id == nil {
		return nil, nil
	}

	idAsInt, err := strconv.Atoi(*id)
	if err != nil {
		return nil, err
	}

	idInt64 := int64(idAsInt)
	return &idInt64, nil
}

// Int64ToString converts numeric id to string-based one
func Int64ToString(id *int64) *string {
	if id == nil {
		return nil
	}

	idAsString := strconv.FormatInt(*id, 10)
	return &idAsString
}

// MapOpenAPICustomProperties maps OpenAPI custom properties model to MLMD one
func MapOpenAPICustomProperties(source *map[string]openapi.MetadataValue) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)

	if source != nil {
		for key, v := range *source {
			value := proto.Value{}

			switch {
			// bool value
			case v.MetadataBoolValue != nil:
				value.Value = &proto.Value_BoolValue{BoolValue: *v.MetadataBoolValue.BoolValue}
			// int value
			case v.MetadataIntValue != nil:
				intValue, err := StringToInt64(v.MetadataIntValue.IntValue)
				if err != nil {
					return nil, fmt.Errorf("unable to decode as int64 %w for key %s", err, key)
				}
				value.Value = &proto.Value_IntValue{IntValue: *intValue}
			// double value
			case v.MetadataDoubleValue != nil:
				value.Value = &proto.Value_DoubleValue{DoubleValue: *v.MetadataDoubleValue.DoubleValue}
			// string value
			case v.MetadataStringValue != nil:
				value.Value = &proto.Value_StringValue{StringValue: *v.MetadataStringValue.StringValue}
			// struct value
			case v.MetadataStructValue != nil:
				data, err := base64.StdEncoding.DecodeString(*v.MetadataStructValue.StructValue)
				if err != nil {
					return nil, fmt.Errorf("unable to decode %w for key %s", err, key)
				}
				var asMap map[string]interface{}
				err = json.Unmarshal(data, &asMap)
				if err != nil {
					return nil, fmt.Errorf("unable to decode %w for key %s", err, key)
				}
				asStruct, err := structpb.NewStruct(asMap)
				if err != nil {
					return nil, fmt.Errorf("unable to decode %w for key %s", err, key)
				}
				value.Value = &proto.Value_StructValue{
					StructValue: asStruct,
				}
			default:
				return nil, fmt.Errorf("type mapping not found for %s:%v", key, v)
			}

			props[key] = &value
		}
	}

	return props, nil
}

// PrefixWhenOwned compose the mlmd fullname by using ownerId as prefix
// For owned entity such as ModelVersion
// for potentially owned entity such as ModelArtifact
func PrefixWhenOwned(ownerId *int64, entityName string) string {
	if ownerId != nil {
		return fmt.Sprintf("%d:%s", *ownerId, entityName)
	}
	uuidPrefix := uuid.New().String()
	return fmt.Sprintf("%s:%s", uuidPrefix, entityName)
}

// REGISTERED MODEL

// MapRegisteredModelProperties maps RegisteredModel fields to specific MLMD properties
func MapRegisteredModelProperties(source *openapi.RegisteredModel) (map[string]*proto.Value, error) {
	// NOTE: at the moment there is no property specified for odh.RegisteredModel context type
	props := make(map[string]*proto.Value)
	return props, nil
}

// MapRegisteredModelType returnd RegisteredModel corresponding MLMD context type
func MapRegisteredModelType(_ *openapi.RegisteredModel) *string {
	return Of(RegisteredModelTypeName)
}

// MODEL VERSION

// MapModelVersionProperties maps ModelVersion fields to specific MLMD properties
func MapModelVersionProperties(source *openapi.ModelVersion) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)
	if source != nil {
		// TODO: model_name is known only from the RegisteredModel, not available here
		props["model_name"] = &proto.Value{
			Value: &proto.Value_StringValue{
				StringValue: "",
			},
		}
		props["version"] = &proto.Value{
			Value: &proto.Value_StringValue{
				StringValue: *source.Name,
			},
		}
		// TODO: not available for now
		props["author"] = &proto.Value{
			Value: &proto.Value_StringValue{
				StringValue: "",
			},
		}
	}
	return props, nil
}

// MapModelVersionType returnd ModelVersion corresponding MLMD context type
func MapModelVersionType(_ *openapi.ModelVersion) *string {
	return Of(ModelVersionTypeName)
}

// MODEL ARTIFACT

// MapModelArtifactProperties maps ModelArtifact fields to specific MLMD properties
func MapModelArtifactProperties(source *openapi.ModelArtifact) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)
	if source != nil {
		if source.ModelFormatName != nil {
			props["model_format"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.ModelFormatName,
				},
			}
		}
	}
	return props, nil
}

// MapModelArtifactType returnd ModelArtifact corresponding MLMD context type
func MapModelArtifactType(_ *openapi.ModelArtifact) *string {
	return Of(ModelArtifactTypeName)
}

func MapOpenAPIModelArtifactState(source *openapi.ArtifactState) *proto.Artifact_State {
	if source == nil {
		return nil
	}

	state := (proto.Artifact_State)(proto.Artifact_State_value[string(*source)])
	return &state
}

// Of returns a pointer to the provided literal/const input
func Of[E any](e E) *E {
	return &e
}
