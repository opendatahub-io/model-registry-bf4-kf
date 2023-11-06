// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.

package generated

import openapi "github.com/opendatahub-io/model-registry/internal/model/openapi"

type OpenAPIConverterImpl struct{}

func (c *OpenAPIConverterImpl) ConvertModelArtifactCreate(source *openapi.ModelArtifactCreate) (*openapi.ModelArtifact, error) {
	var pOpenapiModelArtifact *openapi.ModelArtifact
	if source != nil {
		var openapiModelArtifact openapi.ModelArtifact
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).CustomProperties != nil {
			mapStringOpenapiMetadataValue := make(map[string]openapi.MetadataValue, len((*(*source).CustomProperties)))
			for key, value := range *(*source).CustomProperties {
				mapStringOpenapiMetadataValue[key] = c.openapiMetadataValueToOpenapiMetadataValue(value)
			}
			pMapStringOpenapiMetadataValue = &mapStringOpenapiMetadataValue
		}
		openapiModelArtifact.CustomProperties = pMapStringOpenapiMetadataValue
		var pString *string
		if (*source).Description != nil {
			xstring := *(*source).Description
			pString = &xstring
		}
		openapiModelArtifact.Description = pString
		var pString2 *string
		if (*source).ExternalID != nil {
			xstring2 := *(*source).ExternalID
			pString2 = &xstring2
		}
		openapiModelArtifact.ExternalID = pString2
		var pString3 *string
		if (*source).Uri != nil {
			xstring3 := *(*source).Uri
			pString3 = &xstring3
		}
		openapiModelArtifact.Uri = pString3
		var pOpenapiArtifactState *openapi.ArtifactState
		if (*source).State != nil {
			openapiArtifactState := openapi.ArtifactState(*(*source).State)
			pOpenapiArtifactState = &openapiArtifactState
		}
		openapiModelArtifact.State = pOpenapiArtifactState
		var pString4 *string
		if (*source).Name != nil {
			xstring4 := *(*source).Name
			pString4 = &xstring4
		}
		openapiModelArtifact.Name = pString4
		var pString5 *string
		if (*source).ModelFormatName != nil {
			xstring5 := *(*source).ModelFormatName
			pString5 = &xstring5
		}
		openapiModelArtifact.ModelFormatName = pString5
		var pString6 *string
		if (*source).Runtime != nil {
			xstring6 := *(*source).Runtime
			pString6 = &xstring6
		}
		openapiModelArtifact.Runtime = pString6
		var pString7 *string
		if (*source).StorageKey != nil {
			xstring7 := *(*source).StorageKey
			pString7 = &xstring7
		}
		openapiModelArtifact.StorageKey = pString7
		var pString8 *string
		if (*source).StoragePath != nil {
			xstring8 := *(*source).StoragePath
			pString8 = &xstring8
		}
		openapiModelArtifact.StoragePath = pString8
		var pString9 *string
		if (*source).ModelFormatVersion != nil {
			xstring9 := *(*source).ModelFormatVersion
			pString9 = &xstring9
		}
		openapiModelArtifact.ModelFormatVersion = pString9
		var pString10 *string
		if (*source).ServiceAccountName != nil {
			xstring10 := *(*source).ServiceAccountName
			pString10 = &xstring10
		}
		openapiModelArtifact.ServiceAccountName = pString10
		pOpenapiModelArtifact = &openapiModelArtifact
	}
	return pOpenapiModelArtifact, nil
}
func (c *OpenAPIConverterImpl) ConvertModelArtifactUpdate(source *openapi.ModelArtifactUpdate) (*openapi.ModelArtifact, error) {
	var pOpenapiModelArtifact *openapi.ModelArtifact
	if source != nil {
		var openapiModelArtifact openapi.ModelArtifact
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).CustomProperties != nil {
			mapStringOpenapiMetadataValue := make(map[string]openapi.MetadataValue, len((*(*source).CustomProperties)))
			for key, value := range *(*source).CustomProperties {
				mapStringOpenapiMetadataValue[key] = c.openapiMetadataValueToOpenapiMetadataValue(value)
			}
			pMapStringOpenapiMetadataValue = &mapStringOpenapiMetadataValue
		}
		openapiModelArtifact.CustomProperties = pMapStringOpenapiMetadataValue
		var pString *string
		if (*source).Description != nil {
			xstring := *(*source).Description
			pString = &xstring
		}
		openapiModelArtifact.Description = pString
		var pString2 *string
		if (*source).ExternalID != nil {
			xstring2 := *(*source).ExternalID
			pString2 = &xstring2
		}
		openapiModelArtifact.ExternalID = pString2
		var pString3 *string
		if (*source).Uri != nil {
			xstring3 := *(*source).Uri
			pString3 = &xstring3
		}
		openapiModelArtifact.Uri = pString3
		var pOpenapiArtifactState *openapi.ArtifactState
		if (*source).State != nil {
			openapiArtifactState := openapi.ArtifactState(*(*source).State)
			pOpenapiArtifactState = &openapiArtifactState
		}
		openapiModelArtifact.State = pOpenapiArtifactState
		var pString4 *string
		if (*source).ModelFormatName != nil {
			xstring4 := *(*source).ModelFormatName
			pString4 = &xstring4
		}
		openapiModelArtifact.ModelFormatName = pString4
		var pString5 *string
		if (*source).Runtime != nil {
			xstring5 := *(*source).Runtime
			pString5 = &xstring5
		}
		openapiModelArtifact.Runtime = pString5
		var pString6 *string
		if (*source).StorageKey != nil {
			xstring6 := *(*source).StorageKey
			pString6 = &xstring6
		}
		openapiModelArtifact.StorageKey = pString6
		var pString7 *string
		if (*source).StoragePath != nil {
			xstring7 := *(*source).StoragePath
			pString7 = &xstring7
		}
		openapiModelArtifact.StoragePath = pString7
		var pString8 *string
		if (*source).ModelFormatVersion != nil {
			xstring8 := *(*source).ModelFormatVersion
			pString8 = &xstring8
		}
		openapiModelArtifact.ModelFormatVersion = pString8
		var pString9 *string
		if (*source).ServiceAccountName != nil {
			xstring9 := *(*source).ServiceAccountName
			pString9 = &xstring9
		}
		openapiModelArtifact.ServiceAccountName = pString9
		pOpenapiModelArtifact = &openapiModelArtifact
	}
	return pOpenapiModelArtifact, nil
}
func (c *OpenAPIConverterImpl) ConvertModelVersionCreate(source *openapi.ModelVersionCreate) (*openapi.ModelVersion, error) {
	var pOpenapiModelVersion *openapi.ModelVersion
	if source != nil {
		var openapiModelVersion openapi.ModelVersion
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).CustomProperties != nil {
			mapStringOpenapiMetadataValue := make(map[string]openapi.MetadataValue, len((*(*source).CustomProperties)))
			for key, value := range *(*source).CustomProperties {
				mapStringOpenapiMetadataValue[key] = c.openapiMetadataValueToOpenapiMetadataValue(value)
			}
			pMapStringOpenapiMetadataValue = &mapStringOpenapiMetadataValue
		}
		openapiModelVersion.CustomProperties = pMapStringOpenapiMetadataValue
		var pString *string
		if (*source).Description != nil {
			xstring := *(*source).Description
			pString = &xstring
		}
		openapiModelVersion.Description = pString
		var pString2 *string
		if (*source).ExternalID != nil {
			xstring2 := *(*source).ExternalID
			pString2 = &xstring2
		}
		openapiModelVersion.ExternalID = pString2
		var pString3 *string
		if (*source).Name != nil {
			xstring3 := *(*source).Name
			pString3 = &xstring3
		}
		openapiModelVersion.Name = pString3
		pOpenapiModelVersion = &openapiModelVersion
	}
	return pOpenapiModelVersion, nil
}
func (c *OpenAPIConverterImpl) ConvertModelVersionUpdate(source *openapi.ModelVersionUpdate) (*openapi.ModelVersion, error) {
	var pOpenapiModelVersion *openapi.ModelVersion
	if source != nil {
		var openapiModelVersion openapi.ModelVersion
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).CustomProperties != nil {
			mapStringOpenapiMetadataValue := make(map[string]openapi.MetadataValue, len((*(*source).CustomProperties)))
			for key, value := range *(*source).CustomProperties {
				mapStringOpenapiMetadataValue[key] = c.openapiMetadataValueToOpenapiMetadataValue(value)
			}
			pMapStringOpenapiMetadataValue = &mapStringOpenapiMetadataValue
		}
		openapiModelVersion.CustomProperties = pMapStringOpenapiMetadataValue
		var pString *string
		if (*source).Description != nil {
			xstring := *(*source).Description
			pString = &xstring
		}
		openapiModelVersion.Description = pString
		var pString2 *string
		if (*source).ExternalID != nil {
			xstring2 := *(*source).ExternalID
			pString2 = &xstring2
		}
		openapiModelVersion.ExternalID = pString2
		pOpenapiModelVersion = &openapiModelVersion
	}
	return pOpenapiModelVersion, nil
}
func (c *OpenAPIConverterImpl) ConvertRegisteredModelCreate(source *openapi.RegisteredModelCreate) (*openapi.RegisteredModel, error) {
	var pOpenapiRegisteredModel *openapi.RegisteredModel
	if source != nil {
		var openapiRegisteredModel openapi.RegisteredModel
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).CustomProperties != nil {
			mapStringOpenapiMetadataValue := make(map[string]openapi.MetadataValue, len((*(*source).CustomProperties)))
			for key, value := range *(*source).CustomProperties {
				mapStringOpenapiMetadataValue[key] = c.openapiMetadataValueToOpenapiMetadataValue(value)
			}
			pMapStringOpenapiMetadataValue = &mapStringOpenapiMetadataValue
		}
		openapiRegisteredModel.CustomProperties = pMapStringOpenapiMetadataValue
		var pString *string
		if (*source).Description != nil {
			xstring := *(*source).Description
			pString = &xstring
		}
		openapiRegisteredModel.Description = pString
		var pString2 *string
		if (*source).ExternalID != nil {
			xstring2 := *(*source).ExternalID
			pString2 = &xstring2
		}
		openapiRegisteredModel.ExternalID = pString2
		var pString3 *string
		if (*source).Name != nil {
			xstring3 := *(*source).Name
			pString3 = &xstring3
		}
		openapiRegisteredModel.Name = pString3
		pOpenapiRegisteredModel = &openapiRegisteredModel
	}
	return pOpenapiRegisteredModel, nil
}
func (c *OpenAPIConverterImpl) ConvertRegisteredModelUpdate(source *openapi.RegisteredModelUpdate) (*openapi.RegisteredModel, error) {
	var pOpenapiRegisteredModel *openapi.RegisteredModel
	if source != nil {
		var openapiRegisteredModel openapi.RegisteredModel
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).CustomProperties != nil {
			mapStringOpenapiMetadataValue := make(map[string]openapi.MetadataValue, len((*(*source).CustomProperties)))
			for key, value := range *(*source).CustomProperties {
				mapStringOpenapiMetadataValue[key] = c.openapiMetadataValueToOpenapiMetadataValue(value)
			}
			pMapStringOpenapiMetadataValue = &mapStringOpenapiMetadataValue
		}
		openapiRegisteredModel.CustomProperties = pMapStringOpenapiMetadataValue
		var pString *string
		if (*source).Description != nil {
			xstring := *(*source).Description
			pString = &xstring
		}
		openapiRegisteredModel.Description = pString
		var pString2 *string
		if (*source).ExternalID != nil {
			xstring2 := *(*source).ExternalID
			pString2 = &xstring2
		}
		openapiRegisteredModel.ExternalID = pString2
		pOpenapiRegisteredModel = &openapiRegisteredModel
	}
	return pOpenapiRegisteredModel, nil
}
func (c *OpenAPIConverterImpl) openapiMetadataValueToOpenapiMetadataValue(source openapi.MetadataValue) openapi.MetadataValue {
	var openapiMetadataValue openapi.MetadataValue
	openapiMetadataValue.MetadataBoolValue = c.pOpenapiMetadataBoolValueToPOpenapiMetadataBoolValue(source.MetadataBoolValue)
	openapiMetadataValue.MetadataDoubleValue = c.pOpenapiMetadataDoubleValueToPOpenapiMetadataDoubleValue(source.MetadataDoubleValue)
	openapiMetadataValue.MetadataIntValue = c.pOpenapiMetadataIntValueToPOpenapiMetadataIntValue(source.MetadataIntValue)
	openapiMetadataValue.MetadataProtoValue = c.pOpenapiMetadataProtoValueToPOpenapiMetadataProtoValue(source.MetadataProtoValue)
	openapiMetadataValue.MetadataStringValue = c.pOpenapiMetadataStringValueToPOpenapiMetadataStringValue(source.MetadataStringValue)
	openapiMetadataValue.MetadataStructValue = c.pOpenapiMetadataStructValueToPOpenapiMetadataStructValue(source.MetadataStructValue)
	return openapiMetadataValue
}
func (c *OpenAPIConverterImpl) pOpenapiMetadataBoolValueToPOpenapiMetadataBoolValue(source *openapi.MetadataBoolValue) *openapi.MetadataBoolValue {
	var pOpenapiMetadataBoolValue *openapi.MetadataBoolValue
	if source != nil {
		var openapiMetadataBoolValue openapi.MetadataBoolValue
		var pBool *bool
		if (*source).BoolValue != nil {
			xbool := *(*source).BoolValue
			pBool = &xbool
		}
		openapiMetadataBoolValue.BoolValue = pBool
		pOpenapiMetadataBoolValue = &openapiMetadataBoolValue
	}
	return pOpenapiMetadataBoolValue
}
func (c *OpenAPIConverterImpl) pOpenapiMetadataDoubleValueToPOpenapiMetadataDoubleValue(source *openapi.MetadataDoubleValue) *openapi.MetadataDoubleValue {
	var pOpenapiMetadataDoubleValue *openapi.MetadataDoubleValue
	if source != nil {
		var openapiMetadataDoubleValue openapi.MetadataDoubleValue
		var pFloat64 *float64
		if (*source).DoubleValue != nil {
			xfloat64 := *(*source).DoubleValue
			pFloat64 = &xfloat64
		}
		openapiMetadataDoubleValue.DoubleValue = pFloat64
		pOpenapiMetadataDoubleValue = &openapiMetadataDoubleValue
	}
	return pOpenapiMetadataDoubleValue
}
func (c *OpenAPIConverterImpl) pOpenapiMetadataIntValueToPOpenapiMetadataIntValue(source *openapi.MetadataIntValue) *openapi.MetadataIntValue {
	var pOpenapiMetadataIntValue *openapi.MetadataIntValue
	if source != nil {
		var openapiMetadataIntValue openapi.MetadataIntValue
		var pString *string
		if (*source).IntValue != nil {
			xstring := *(*source).IntValue
			pString = &xstring
		}
		openapiMetadataIntValue.IntValue = pString
		pOpenapiMetadataIntValue = &openapiMetadataIntValue
	}
	return pOpenapiMetadataIntValue
}
func (c *OpenAPIConverterImpl) pOpenapiMetadataProtoValueToPOpenapiMetadataProtoValue(source *openapi.MetadataProtoValue) *openapi.MetadataProtoValue {
	var pOpenapiMetadataProtoValue *openapi.MetadataProtoValue
	if source != nil {
		var openapiMetadataProtoValue openapi.MetadataProtoValue
		var pString *string
		if (*source).Type != nil {
			xstring := *(*source).Type
			pString = &xstring
		}
		openapiMetadataProtoValue.Type = pString
		var pString2 *string
		if (*source).ProtoValue != nil {
			xstring2 := *(*source).ProtoValue
			pString2 = &xstring2
		}
		openapiMetadataProtoValue.ProtoValue = pString2
		pOpenapiMetadataProtoValue = &openapiMetadataProtoValue
	}
	return pOpenapiMetadataProtoValue
}
func (c *OpenAPIConverterImpl) pOpenapiMetadataStringValueToPOpenapiMetadataStringValue(source *openapi.MetadataStringValue) *openapi.MetadataStringValue {
	var pOpenapiMetadataStringValue *openapi.MetadataStringValue
	if source != nil {
		var openapiMetadataStringValue openapi.MetadataStringValue
		var pString *string
		if (*source).StringValue != nil {
			xstring := *(*source).StringValue
			pString = &xstring
		}
		openapiMetadataStringValue.StringValue = pString
		pOpenapiMetadataStringValue = &openapiMetadataStringValue
	}
	return pOpenapiMetadataStringValue
}
func (c *OpenAPIConverterImpl) pOpenapiMetadataStructValueToPOpenapiMetadataStructValue(source *openapi.MetadataStructValue) *openapi.MetadataStructValue {
	var pOpenapiMetadataStructValue *openapi.MetadataStructValue
	if source != nil {
		var openapiMetadataStructValue openapi.MetadataStructValue
		var pString *string
		if (*source).StructValue != nil {
			xstring := *(*source).StructValue
			pString = &xstring
		}
		openapiMetadataStructValue.StructValue = pString
		pOpenapiMetadataStructValue = &openapiMetadataStructValue
	}
	return pOpenapiMetadataStructValue
}
