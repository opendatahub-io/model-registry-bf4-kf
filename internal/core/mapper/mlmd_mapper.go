package mapper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/opendatahub-io/model-registry/internal/converter"
	"github.com/opendatahub-io/model-registry/internal/converter/generated"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
)

type Mapper struct {
	toMLMDConverter       converter.OpenAPIToMLMDConverter
	toOpenAPIConverter    converter.MLMDToOpenAPIConverter
	RegisteredModelTypeId int64
	ModelVersionTypeId    int64
	ModelArtifactTypeId   int64
}

func NewMapper(registeredModelTypeId int64, modelVersionTypeId int64, modelArtifactTypeId int64) *Mapper {
	return &Mapper{
		toMLMDConverter:       &generated.OpenAPIToMLMDConverterImpl{},
		toOpenAPIConverter:    &generated.MLMDToOpenAPIConverterImpl{},
		RegisteredModelTypeId: registeredModelTypeId,
		ModelVersionTypeId:    modelVersionTypeId,
		ModelArtifactTypeId:   modelArtifactTypeId,
	}
}

// Internal Model --> MLMD

func (m *Mapper) MapFromRegisteredModel(registeredModel *openapi.RegisteredModel) (*proto.Context, error) {
	ctx, err := m.toMLMDConverter.ConvertRegisteredModel(registeredModel)
	if err != nil {
		return nil, err
	}

	ctx.TypeId = &m.RegisteredModelTypeId
	return ctx, nil
}

func (m *Mapper) MapFromModelVersion(modelVersion *openapi.ModelVersion, registeredModelId int64, registeredModelName *string) (*proto.Context, error) {
	fullName := converter.PrefixWhenOwned(&registeredModelId, *modelVersion.Name)

	ctx, err := m.toMLMDConverter.ConvertModelVersion(modelVersion)
	if err != nil {
		return nil, err
	}

	ctx.TypeId = &m.ModelVersionTypeId
	ctx.Name = &fullName

	return ctx, nil
}

func (m *Mapper) MapFromModelArtifact(modelArtifact openapi.ModelArtifact, modelVersionId *int64) *proto.Artifact {
	// openapi.Artifact is defined with optional name, so build arbitrary name for this artifact if missing
	var artifactName string
	if modelArtifact.Name != nil {
		artifactName = *modelArtifact.Name
	} else {
		artifactName = uuid.New().String()
	}
	// build fullName for mlmd storage
	fullName := converter.PrefixWhenOwned(modelVersionId, artifactName)

	artifact, err := m.toMLMDConverter.ConvertModelArtifact(&modelArtifact)
	if err != nil {
		return nil
	}

	artifact.TypeId = &m.ModelArtifactTypeId
	artifact.Name = &fullName

	return artifact
}

func (m *Mapper) MapFromModelArtifacts(modelArtifacts *[]openapi.ModelArtifact, modelVersionId *int64) ([]*proto.Artifact, error) {
	artifacts := []*proto.Artifact{}
	if modelArtifacts == nil {
		return artifacts, nil
	}
	for _, a := range *modelArtifacts {
		artifacts = append(artifacts, m.MapFromModelArtifact(a, modelVersionId))
	}
	return artifacts, nil
}

//  MLMD --> Internal Model

func (m *Mapper) MapToRegisteredModel(ctx *proto.Context) (*openapi.RegisteredModel, error) {
	if ctx.GetTypeId() != m.RegisteredModelTypeId {
		return nil, fmt.Errorf("invalid TypeId, expected %d but received %d", m.RegisteredModelTypeId, ctx.GetTypeId())
	}

	return m.toOpenAPIConverter.ConvertRegisteredModel(ctx)
}

func (m *Mapper) MapToModelVersion(ctx *proto.Context) (*openapi.ModelVersion, error) {
	if ctx.GetTypeId() != m.ModelVersionTypeId {
		return nil, fmt.Errorf("invalid TypeId, expected %d but received %d", m.ModelVersionTypeId, ctx.GetTypeId())
	}

	return m.toOpenAPIConverter.ConvertModelVersion(ctx)
}

func (m *Mapper) MapToModelArtifact(artifact *proto.Artifact) (*openapi.ModelArtifact, error) {
	if artifact.GetTypeId() != m.ModelArtifactTypeId {
		return nil, fmt.Errorf("invalid TypeId, expected %d but received %d", m.ModelArtifactTypeId, artifact.GetTypeId())
	}

	return m.toOpenAPIConverter.ConvertModelArtifact(artifact)
}
