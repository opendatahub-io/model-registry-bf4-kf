// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.

package generated

import (
	"fmt"
	converter "github.com/opendatahub-io/model-registry/internal/converter"
	proto "github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	openapi "github.com/opendatahub-io/model-registry/internal/model/openapi"
)

type OpenAPIToMLMDConverterImpl struct{}

func (c *OpenAPIToMLMDConverterImpl) ConvertInferenceService(source *converter.OpenAPIModelWrapper[openapi.InferenceService]) (*proto.Context, error) {
	var pProtoContext *proto.Context
	if source != nil {
		var protoContext proto.Context
		var pString *string
		if (*source).Model != nil {
			pString = (*source).Model.Id
		}
		pInt64, err := converter.StringToInt64(pString)
		if err != nil {
			return nil, fmt.Errorf("error setting field Id: %w", err)
		}
		protoContext.Id = pInt64
		var pString2 *string
		if (*source).Model != nil {
			pString2 = (*source).Model.Name
		}
		var pString3 *string
		if pString2 != nil {
			xstring := *pString2
			pString3 = &xstring
		}
		protoContext.Name = pString3
		pInt642 := (*source).TypeId
		protoContext.TypeId = &pInt642
		protoContext.Type = converter.MapInferenceServiceType((*source).Model)
		var pString4 *string
		if (*source).Model != nil {
			pString4 = (*source).Model.ExternalID
		}
		var pString5 *string
		if pString4 != nil {
			xstring2 := *pString4
			pString5 = &xstring2
		}
		protoContext.ExternalId = pString5
		mapStringPProtoValue, err := converter.MapInferenceServiceProperties((*source).Model)
		if err != nil {
			return nil, fmt.Errorf("error setting field Properties: %w", err)
		}
		protoContext.Properties = mapStringPProtoValue
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).Model != nil {
			pMapStringOpenapiMetadataValue = (*source).Model.CustomProperties
		}
		mapStringPProtoValue2, err := converter.MapOpenAPICustomProperties(pMapStringOpenapiMetadataValue)
		if err != nil {
			return nil, fmt.Errorf("error setting field CustomProperties: %w", err)
		}
		protoContext.CustomProperties = mapStringPProtoValue2
		pProtoContext = &protoContext
	}
	return pProtoContext, nil
}
func (c *OpenAPIToMLMDConverterImpl) ConvertModelArtifact(source *converter.OpenAPIModelWrapper[openapi.ModelArtifact]) (*proto.Artifact, error) {
	var pProtoArtifact *proto.Artifact
	if source != nil {
		var protoArtifact proto.Artifact
		var pString *string
		if (*source).Model != nil {
			pString = (*source).Model.Id
		}
		pInt64, err := converter.StringToInt64(pString)
		if err != nil {
			return nil, fmt.Errorf("error setting field Id: %w", err)
		}
		protoArtifact.Id = pInt64
		protoArtifact.Name = converter.MapModelArtifactName(source)
		pInt642 := (*source).TypeId
		protoArtifact.TypeId = &pInt642
		protoArtifact.Type = converter.MapModelArtifactType((*source).Model)
		var pString2 *string
		if (*source).Model != nil {
			pString2 = (*source).Model.Uri
		}
		var pString3 *string
		if pString2 != nil {
			xstring := *pString2
			pString3 = &xstring
		}
		protoArtifact.Uri = pString3
		var pString4 *string
		if (*source).Model != nil {
			pString4 = (*source).Model.ExternalID
		}
		var pString5 *string
		if pString4 != nil {
			xstring2 := *pString4
			pString5 = &xstring2
		}
		protoArtifact.ExternalId = pString5
		mapStringPProtoValue, err := converter.MapModelArtifactProperties((*source).Model)
		if err != nil {
			return nil, fmt.Errorf("error setting field Properties: %w", err)
		}
		protoArtifact.Properties = mapStringPProtoValue
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).Model != nil {
			pMapStringOpenapiMetadataValue = (*source).Model.CustomProperties
		}
		mapStringPProtoValue2, err := converter.MapOpenAPICustomProperties(pMapStringOpenapiMetadataValue)
		if err != nil {
			return nil, fmt.Errorf("error setting field CustomProperties: %w", err)
		}
		protoArtifact.CustomProperties = mapStringPProtoValue2
		var pOpenapiArtifactState *openapi.ArtifactState
		if (*source).Model != nil {
			pOpenapiArtifactState = (*source).Model.State
		}
		pProtoArtifact_State, err := converter.MapOpenAPIModelArtifactState(pOpenapiArtifactState)
		if err != nil {
			return nil, fmt.Errorf("error setting field State: %w", err)
		}
		protoArtifact.State = pProtoArtifact_State
		pProtoArtifact = &protoArtifact
	}
	return pProtoArtifact, nil
}
func (c *OpenAPIToMLMDConverterImpl) ConvertModelVersion(source *converter.OpenAPIModelWrapper[openapi.ModelVersion]) (*proto.Context, error) {
	var pProtoContext *proto.Context
	if source != nil {
		var protoContext proto.Context
		var pString *string
		if (*source).Model != nil {
			pString = (*source).Model.Id
		}
		pInt64, err := converter.StringToInt64(pString)
		if err != nil {
			return nil, fmt.Errorf("error setting field Id: %w", err)
		}
		protoContext.Id = pInt64
		protoContext.Name = converter.MapModelVersionName(source)
		pInt642 := (*source).TypeId
		protoContext.TypeId = &pInt642
		protoContext.Type = converter.MapModelVersionType((*source).Model)
		var pString2 *string
		if (*source).Model != nil {
			pString2 = (*source).Model.ExternalID
		}
		var pString3 *string
		if pString2 != nil {
			xstring := *pString2
			pString3 = &xstring
		}
		protoContext.ExternalId = pString3
		mapStringPProtoValue, err := converter.MapModelVersionProperties(source)
		if err != nil {
			return nil, fmt.Errorf("error setting field Properties: %w", err)
		}
		protoContext.Properties = mapStringPProtoValue
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).Model != nil {
			pMapStringOpenapiMetadataValue = (*source).Model.CustomProperties
		}
		mapStringPProtoValue2, err := converter.MapOpenAPICustomProperties(pMapStringOpenapiMetadataValue)
		if err != nil {
			return nil, fmt.Errorf("error setting field CustomProperties: %w", err)
		}
		protoContext.CustomProperties = mapStringPProtoValue2
		pProtoContext = &protoContext
	}
	return pProtoContext, nil
}
func (c *OpenAPIToMLMDConverterImpl) ConvertRegisteredModel(source *converter.OpenAPIModelWrapper[openapi.RegisteredModel]) (*proto.Context, error) {
	var pProtoContext *proto.Context
	if source != nil {
		var protoContext proto.Context
		var pString *string
		if (*source).Model != nil {
			pString = (*source).Model.Id
		}
		pInt64, err := converter.StringToInt64(pString)
		if err != nil {
			return nil, fmt.Errorf("error setting field Id: %w", err)
		}
		protoContext.Id = pInt64
		var pString2 *string
		if (*source).Model != nil {
			pString2 = (*source).Model.Name
		}
		var pString3 *string
		if pString2 != nil {
			xstring := *pString2
			pString3 = &xstring
		}
		protoContext.Name = pString3
		pInt642 := (*source).TypeId
		protoContext.TypeId = &pInt642
		protoContext.Type = converter.MapRegisteredModelType((*source).Model)
		var pString4 *string
		if (*source).Model != nil {
			pString4 = (*source).Model.ExternalID
		}
		var pString5 *string
		if pString4 != nil {
			xstring2 := *pString4
			pString5 = &xstring2
		}
		protoContext.ExternalId = pString5
		mapStringPProtoValue, err := converter.MapRegisteredModelProperties((*source).Model)
		if err != nil {
			return nil, fmt.Errorf("error setting field Properties: %w", err)
		}
		protoContext.Properties = mapStringPProtoValue
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).Model != nil {
			pMapStringOpenapiMetadataValue = (*source).Model.CustomProperties
		}
		mapStringPProtoValue2, err := converter.MapOpenAPICustomProperties(pMapStringOpenapiMetadataValue)
		if err != nil {
			return nil, fmt.Errorf("error setting field CustomProperties: %w", err)
		}
		protoContext.CustomProperties = mapStringPProtoValue2
		pProtoContext = &protoContext
	}
	return pProtoContext, nil
}
func (c *OpenAPIToMLMDConverterImpl) ConvertServeModel(source *converter.OpenAPIModelWrapper[openapi.ServeModel]) (*proto.Execution, error) {
	var pProtoExecution *proto.Execution
	if source != nil {
		var protoExecution proto.Execution
		var pString *string
		if (*source).Model != nil {
			pString = (*source).Model.Id
		}
		pInt64, err := converter.StringToInt64(pString)
		if err != nil {
			return nil, fmt.Errorf("error setting field Id: %w", err)
		}
		protoExecution.Id = pInt64
		var pString2 *string
		if (*source).Model != nil {
			pString2 = (*source).Model.Name
		}
		var pString3 *string
		if pString2 != nil {
			xstring := *pString2
			pString3 = &xstring
		}
		protoExecution.Name = pString3
		pInt642 := (*source).TypeId
		protoExecution.TypeId = &pInt642
		protoExecution.Type = converter.MapServeModelType((*source).Model)
		var pString4 *string
		if (*source).Model != nil {
			pString4 = (*source).Model.ExternalID
		}
		var pString5 *string
		if pString4 != nil {
			xstring2 := *pString4
			pString5 = &xstring2
		}
		protoExecution.ExternalId = pString5
		var pOpenapiExecutionState *openapi.ExecutionState
		if (*source).Model != nil {
			pOpenapiExecutionState = (*source).Model.LastKnownState
		}
		pProtoExecution_State, err := converter.MapLastKnownState(pOpenapiExecutionState)
		if err != nil {
			return nil, fmt.Errorf("error setting field LastKnownState: %w", err)
		}
		protoExecution.LastKnownState = pProtoExecution_State
		mapStringPProtoValue, err := converter.MapServeModelProperties((*source).Model)
		if err != nil {
			return nil, fmt.Errorf("error setting field Properties: %w", err)
		}
		protoExecution.Properties = mapStringPProtoValue
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).Model != nil {
			pMapStringOpenapiMetadataValue = (*source).Model.CustomProperties
		}
		mapStringPProtoValue2, err := converter.MapOpenAPICustomProperties(pMapStringOpenapiMetadataValue)
		if err != nil {
			return nil, fmt.Errorf("error setting field CustomProperties: %w", err)
		}
		protoExecution.CustomProperties = mapStringPProtoValue2
		pProtoExecution = &protoExecution
	}
	return pProtoExecution, nil
}
func (c *OpenAPIToMLMDConverterImpl) ConvertServingEnvironment(source *converter.OpenAPIModelWrapper[openapi.ServingEnvironment]) (*proto.Context, error) {
	var pProtoContext *proto.Context
	if source != nil {
		var protoContext proto.Context
		var pString *string
		if (*source).Model != nil {
			pString = (*source).Model.Id
		}
		pInt64, err := converter.StringToInt64(pString)
		if err != nil {
			return nil, fmt.Errorf("error setting field Id: %w", err)
		}
		protoContext.Id = pInt64
		var pString2 *string
		if (*source).Model != nil {
			pString2 = (*source).Model.Name
		}
		var pString3 *string
		if pString2 != nil {
			xstring := *pString2
			pString3 = &xstring
		}
		protoContext.Name = pString3
		pInt642 := (*source).TypeId
		protoContext.TypeId = &pInt642
		protoContext.Type = converter.MapServingEnvironmentType((*source).Model)
		var pString4 *string
		if (*source).Model != nil {
			pString4 = (*source).Model.ExternalID
		}
		var pString5 *string
		if pString4 != nil {
			xstring2 := *pString4
			pString5 = &xstring2
		}
		protoContext.ExternalId = pString5
		mapStringPProtoValue, err := converter.MapServingEnvironmentProperties((*source).Model)
		if err != nil {
			return nil, fmt.Errorf("error setting field Properties: %w", err)
		}
		protoContext.Properties = mapStringPProtoValue
		var pMapStringOpenapiMetadataValue *map[string]openapi.MetadataValue
		if (*source).Model != nil {
			pMapStringOpenapiMetadataValue = (*source).Model.CustomProperties
		}
		mapStringPProtoValue2, err := converter.MapOpenAPICustomProperties(pMapStringOpenapiMetadataValue)
		if err != nil {
			return nil, fmt.Errorf("error setting field CustomProperties: %w", err)
		}
		protoContext.CustomProperties = mapStringPProtoValue2
		pProtoContext = &protoContext
	}
	return pProtoContext, nil
}
