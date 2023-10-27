package generated_test

import (
	"testing"

	"github.com/opendatahub-io/model-registry/internal/converter"
	"github.com/opendatahub-io/model-registry/internal/converter/generated"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (converter.OpenAPIConverter, *assert.Assertions) {
	return &generated.OpenAPIConverterImpl{}, assert.New(t)
}

var (
	name        = "entity-name"
	externalId  = "entity-ext-id"
	customProps = map[string]openapi.MetadataValue{}
)

func TestConvertRegisteredModelCreate(t *testing.T) {
	converter, assertion := setup(t)

	source := openapi.RegisteredModelCreate{
		Name:             &name,
		ExternalID:       &externalId,
		CustomProperties: &customProps,
	}

	target, err := converter.ConvertRegisteredModelCreate(&source)
	assertion.Nilf(err, "conversion should have worked: %v", err)

	assertion.Equal(*source.Name, *target.Name)
	assertion.Equal(*source.ExternalID, *target.ExternalID)
	assertion.Equal(*source.CustomProperties, *target.CustomProperties)
	assertion.Nil(target.Id)
	assertion.Nil(target.CreateTimeSinceEpoch)
	assertion.Nil(target.LastUpdateTimeSinceEpoch)
}

func TestConvertRegisteredModelUpdate(t *testing.T) {
	converter, assertion := setup(t)

	source := openapi.RegisteredModelUpdate{
		ExternalID:       &externalId,
		CustomProperties: &customProps,
	}

	target, err := converter.ConvertRegisteredModelUpdate(&source)
	assertion.Nilf(err, "conversion should have worked: %v", err)

	assertion.Equal(*source.ExternalID, *target.ExternalID)
	assertion.Equal(*source.CustomProperties, *target.CustomProperties)
	assertion.Nil(target.Id)
	assertion.Nil(target.CreateTimeSinceEpoch)
	assertion.Nil(target.LastUpdateTimeSinceEpoch)
	assertion.Nil(target.Name)
}

func TestConvertModelVersionCreate(t *testing.T) {
	converter, assertion := setup(t)

	source := openapi.ModelVersionCreate{
		Name:             &name,
		ExternalID:       &externalId,
		CustomProperties: &customProps,
	}

	target, err := converter.ConvertModelVersionCreate(&source)
	assertion.Nilf(err, "conversion should have worked: %v", err)

	assertion.Equal(*source.Name, *target.Name)
	assertion.Equal(*source.ExternalID, *target.ExternalID)
	assertion.Equal(*source.CustomProperties, *target.CustomProperties)
	assertion.Nil(target.Id)
	assertion.Nil(target.CreateTimeSinceEpoch)
	assertion.Nil(target.LastUpdateTimeSinceEpoch)
}

func TestConvertModelVersionUpdate(t *testing.T) {
	converter, assertion := setup(t)

	source := openapi.ModelVersionUpdate{
		ExternalID:       &externalId,
		CustomProperties: &customProps,
	}

	target, err := converter.ConvertModelVersionUpdate(&source)
	assertion.Nilf(err, "conversion should have worked: %v", err)

	assertion.Equal(*source.ExternalID, *target.ExternalID)
	assertion.Equal(*source.CustomProperties, *target.CustomProperties)
	assertion.Nil(target.Id)
	assertion.Nil(target.CreateTimeSinceEpoch)
	assertion.Nil(target.LastUpdateTimeSinceEpoch)
	assertion.Nil(target.Name)
}

func TestConvertModelArtifactCreate(t *testing.T) {
	converter, assertion := setup(t)

	source := openapi.ModelArtifactCreate{
		Name:             &name,
		ExternalID:       &externalId,
		CustomProperties: &customProps,
	}

	target, err := converter.ConvertModelArtifactCreate(&source)
	assertion.Nilf(err, "conversion should have worked: %v", err)

	assertion.Equal(*source.Name, *target.Name)
	assertion.Equal(*source.ExternalID, *target.ExternalID)
	assertion.Equal(*source.CustomProperties, *target.CustomProperties)
	assertion.Equal("", target.ArtifactType)
	assertion.Nil(target.Id)
	assertion.Nil(target.CreateTimeSinceEpoch)
	assertion.Nil(target.LastUpdateTimeSinceEpoch)
}

func TestConvertModelArtifactUpdate(t *testing.T) {
	converter, assertion := setup(t)

	source := openapi.ModelArtifactUpdate{
		ExternalID:       &externalId,
		CustomProperties: &customProps,
	}

	target, err := converter.ConvertModelArtifactUpdate(&source)
	assertion.Nilf(err, "conversion should have worked: %v", err)

	assertion.Equal(*source.ExternalID, *target.ExternalID)
	assertion.Equal(*source.CustomProperties, *target.CustomProperties)
	assertion.Equal("", target.ArtifactType)
	assertion.Nil(target.Id)
	assertion.Nil(target.CreateTimeSinceEpoch)
	assertion.Nil(target.LastUpdateTimeSinceEpoch)
	assertion.Nil(target.LastUpdateTimeSinceEpoch)
	assertion.Nil(target.Name)
}
