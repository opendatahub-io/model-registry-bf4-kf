package converter

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func setup(t *testing.T) *assert.Assertions {
	return assert.New(t)
}

func TestStringToInt64(t *testing.T) {
	assertion := setup(t)

	valid := "12345"
	converted, err := StringToInt64(&valid)
	assertion.Nil(err)
	assertion.Equal(int64(12345), *converted)
	assertion.Nil(StringToInt64(nil))
}

func TestStringToInt64InvalidNumber(t *testing.T) {
	assertion := setup(t)

	invalid := "not-a-number"
	converted, err := StringToInt64(&invalid)
	assertion.NotNil(err)
	assertion.Nil(converted)
}

func TestInt64ToString(t *testing.T) {
	assertion := setup(t)

	valid := int64(54321)
	converted := Int64ToString(&valid)
	assertion.Equal("54321", *converted)
	assertion.Nil(Int64ToString(nil))
}

func TestMetadataValueBool(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my bool"
	mdValue := true
	data[key] = openapi.MetadataBoolValueAsMetadataValue(&openapi.MetadataBoolValue{BoolValue: &mdValue})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueInt(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my int"
	mdValue := "987"
	data[key] = openapi.MetadataIntValueAsMetadataValue(&openapi.MetadataIntValue{IntValue: &mdValue})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueIntFailure(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my int"
	mdValue := "not a number"
	data[key] = openapi.MetadataIntValueAsMetadataValue(&openapi.MetadataIntValue{IntValue: &mdValue})

	assertion := setup(t)
	asGRPC, err := MapOpenAPICustomProperties(&data)
	if err == nil {
		assertion.Fail("Did not expected a converted value but an error: %v", asGRPC)
	}
}

func TestMetadataValueDouble(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my double"
	mdValue := 3.1415
	data[key] = openapi.MetadataDoubleValueAsMetadataValue(&openapi.MetadataDoubleValue{DoubleValue: &mdValue})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueString(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my string"
	mdValue := "Hello, World!"
	data[key] = openapi.MetadataStringValueAsMetadataValue(&openapi.MetadataStringValue{StringValue: &mdValue})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueStruct(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my struct"

	myMap := make(map[string]interface{})
	myMap["name"] = "John Doe"
	myMap["age"] = 47
	asJSON, err := json.Marshal(myMap)
	if err != nil {
		t.Error(err)
	}
	b64 := base64.StdEncoding.EncodeToString(asJSON)
	data[key] = openapi.MetadataStructValueAsMetadataValue(&openapi.MetadataStructValue{StructValue: &b64})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueProtoUnsupported(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my proto"

	myMap := make(map[string]interface{})
	myMap["name"] = "John Doe"
	myMap["age"] = 47
	asJSON, err := json.Marshal(myMap)
	if err != nil {
		t.Error(err)
	}
	b64 := base64.StdEncoding.EncodeToString(asJSON)
	typeDef := "map[string]openapi.MetadataValue"
	data[key] = openapi.MetadataProtoValueAsMetadataValue(&openapi.MetadataProtoValue{
		Type:       &typeDef,
		ProtoValue: &b64,
	})

	assertion := setup(t)
	asGRPC, err := MapOpenAPICustomProperties(&data)
	if err == nil {
		assertion.Fail("Did not expected a converted value but an error: %v", asGRPC)
	}
}

func roundTripAndAssert(t *testing.T, data map[string]openapi.MetadataValue, key string) {
	assertion := setup(t)

	// first half
	asGRPC, err := MapOpenAPICustomProperties(&data)
	if err != nil {
		t.Error(err)
	}
	assertion.Contains(maps.Keys(asGRPC), key)

	// second half
	unmarshall, err := MapMLMDCustomProperties(asGRPC)
	if err != nil {
		t.Error(err)
	}
	assertion.Equal(data, unmarshall, "result of round-trip shall be equal to original data")
}

func TestPrefixWhenOwned(t *testing.T) {
	assertion := setup(t)

	owner := "owner"
	entity := "name"
	assertion.Equal("owner:name", PrefixWhenOwned(&owner, entity))
}

func TestPrefixWhenOwnedWithoutOwner(t *testing.T) {
	assertion := setup(t)

	entity := "name"
	prefixed := PrefixWhenOwned(nil, entity)
	assertion.Equal(2, len(strings.Split(prefixed, ":")))
	assertion.Equal("name", strings.Split(prefixed, ":")[1])
}

func TestMapRegisteredModelProperties(t *testing.T) {
	assertion := setup(t)

	props, err := MapRegisteredModelProperties(&openapi.RegisteredModel{})
	assertion.Nil(err)
	assertion.Equal(0, len(props))
}

func TestMapRegisteredModelType(t *testing.T) {
	assertion := setup(t)

	typeName := MapRegisteredModelType(&openapi.RegisteredModel{})
	assertion.NotNil(typeName)
	assertion.Equal(RegisteredModelTypeName, *typeName)
}

func TestMapModelVersionProperties(t *testing.T) {
	assertion := setup(t)

	props, err := MapModelVersionProperties(&OpenAPIModelWrapper[openapi.ModelVersion]{
		TypeId:           123,
		ParentResourceId: of("123"),
		ModelName:        of("MyModel"),
		Model: &openapi.ModelVersion{
			Name: of("v1"),
		},
	})
	assertion.Nil(err)
	assertion.Equal(3, len(props))
	// TODO check all 3 props
}

func TestMapModelVersionType(t *testing.T) {
	assertion := setup(t)

	typeName := MapModelVersionType(&openapi.ModelVersion{})
	assertion.NotNil(typeName)
	assertion.Equal(ModelVersionTypeName, *typeName)
}

func TestMapModelVersionName(t *testing.T) {
	assertion := setup(t)

	name := MapModelVersionName(&OpenAPIModelWrapper[openapi.ModelVersion]{
		TypeId:           123,
		ParentResourceId: of("123"),
		ModelName:        of("MyModel"),
		Model: &openapi.ModelVersion{
			Name: of("v1"),
		},
	})
	assertion.NotNil(name)
	assertion.Equal("123:v1", *name)
}

func TestMapModelArtifactProperties(t *testing.T) {
	assertion := setup(t)

	props, err := MapModelArtifactProperties(&openapi.ModelArtifact{
		Name:            of("v1"),
		ModelFormatName: of("sklearn"),
	})
	assertion.Nil(err)
	assertion.Equal(1, len(props))
	// TODO check value

	props, err = MapModelArtifactProperties(&openapi.ModelArtifact{
		Name: of("v1"),
	})
	assertion.Nil(err)
	assertion.Equal(0, len(props))
}

func TestMapModelArtifactType(t *testing.T) {
	assertion := setup(t)

	typeName := MapModelArtifactType(&openapi.ModelArtifact{})
	assertion.NotNil(typeName)
	assertion.Equal(ModelArtifactTypeName, *typeName)
}

func TestMapModelArtifactName(t *testing.T) {
	assertion := setup(t)

	name := MapModelArtifactName(&OpenAPIModelWrapper[openapi.ModelArtifact]{
		TypeId:           123,
		ParentResourceId: of("parent"),
		Model: &openapi.ModelArtifact{
			Name: of("v1"),
		},
	})
	assertion.NotNil(name)
	assertion.Equal("parent:v1", *name)

	name = MapModelArtifactName(&OpenAPIModelWrapper[openapi.ModelArtifact]{
		TypeId:           123,
		ParentResourceId: of("parent"),
		Model: &openapi.ModelArtifact{
			Name: nil,
		},
	})
	assertion.NotNil(name)
	assertion.Regexp("parent:.*", *name)

	name = MapModelArtifactName(&OpenAPIModelWrapper[openapi.ModelArtifact]{
		TypeId: 123,
		Model: &openapi.ModelArtifact{
			Name: of("v1"),
		},
	})
	assertion.NotNil(name)
	assertion.Regexp(".*:v1", *name)
}

func TestMapOpenAPIModelArtifactState(t *testing.T) {
	assertion := setup(t)

	state := MapOpenAPIModelArtifactState(of(openapi.ARTIFACTSTATE_LIVE))
	assertion.NotNil(state)
	assertion.Equal(string(openapi.ARTIFACTSTATE_LIVE), state.String())

	state = MapOpenAPIModelArtifactState(nil)
	assertion.Nil(state)
}
