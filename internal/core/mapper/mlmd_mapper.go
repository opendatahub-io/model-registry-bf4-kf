package mapper

import (
	"fmt"
	"log"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/registry"
)

type Mapper struct {
	RegisteredModelTypeId int64
	ModelVersionTypeId    int64
	ModelArtifactTypeId   int64
}

func NewMapper(registeredModelTypeId int64, modelVersionTypeId int64, modelArtifactTypeId int64) *Mapper {
	return &Mapper{
		RegisteredModelTypeId: registeredModelTypeId,
		ModelVersionTypeId:    modelVersionTypeId,
		ModelArtifactTypeId:   modelArtifactTypeId,
	}
}

// Internal Model --> MLMD

// TODO: implement
// Map generic map into MLMD [custom] properties object
func (m *Mapper) MapToProperties(data map[string]any) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)

	for key, v := range data {
		value := proto.Value{}

		switch typedValue := v.(type) {
		case int:
			value.Value = &proto.Value_IntValue{IntValue: int64(typedValue)}
		case string:
			value.Value = &proto.Value_StringValue{StringValue: typedValue}
		case float64:
			value.Value = &proto.Value_DoubleValue{DoubleValue: float64(typedValue)}
		default:
			log.Printf("Type mapping not found for %s:%v", key, v)
			continue
		}

		props[key] = &value
	}

	return props, nil
}

func (m *Mapper) MapFromRegisteredModel(registeredModel *registry.RegisteredModel) (*proto.Context, error) {

	return &proto.Context{
		Name:   registeredModel.Name,
		TypeId: &m.RegisteredModelTypeId,
		Id:     registeredModel.Id,
	}, nil
}

func (m *Mapper) MapFromModelVersion(modelVersion *registry.VersionedModel, registeredModelId int64) (*proto.Context, error) {
	fullName := fmt.Sprintf("%d:%s", registeredModelId, *modelVersion.Version)
	customProps := make(map[string]*proto.Value)
	if modelVersion.Metadata != nil {
		customProps, _ = m.MapToProperties(*modelVersion.Metadata)
	}
	ctx := &proto.Context{
		Name:   &fullName,
		TypeId: &m.ModelVersionTypeId,
		Properties: map[string]*proto.Value{
			"model_name": {
				Value: &proto.Value_StringValue{
					StringValue: *modelVersion.ModelName,
				},
			},
		},
		CustomProperties: customProps,
	}
	if modelVersion.Version != nil {
		ctx.Properties["version"] = &proto.Value{
			Value: &proto.Value_StringValue{
				StringValue: *modelVersion.Version,
			},
		}
	}
	if modelVersion.Author != nil {
		ctx.Properties["author"] = &proto.Value{
			Value: &proto.Value_StringValue{
				StringValue: *modelVersion.Author,
			},
		}
	}

	return ctx, nil
}

func (m *Mapper) MapFromModelArtifact(modelArtifact registry.Artifact) *proto.Artifact {
	return &proto.Artifact{
		TypeId: &m.ModelArtifactTypeId,
		Name:   modelArtifact.Name,
		Uri:    modelArtifact.Uri,
	}
}

func (m *Mapper) MapFromModelArtifacts(modelArtifacts *[]registry.Artifact) ([]*proto.Artifact, error) {
	artifacts := []*proto.Artifact{}
	if modelArtifacts == nil {
		return artifacts, nil
	}
	for _, a := range *modelArtifacts {
		artifacts = append(artifacts, m.MapFromModelArtifact(a))
	}
	return artifacts, nil
}

//  MLMD --> Internal Model

// TODO implement
// Maps MLMD properties into a generic <string, any> map
func (m *Mapper) MapFromProperties(props map[string]*proto.Value) (map[string]any, error) {
	data := make(map[string]any)

	for key, v := range props {
		data[key] = v.Value
	}

	return data, nil
}

func (m *Mapper) MapToRegisteredModel(ctx *proto.Context) (*registry.RegisteredModel, error) {
	if ctx.GetTypeId() != m.RegisteredModelTypeId {
		return nil, fmt.Errorf("invalid TypeId, exptected %d but received %d", m.RegisteredModelTypeId, ctx.GetTypeId())
	}

	_, err := m.MapFromProperties(ctx.CustomProperties)
	if err != nil {
		return nil, err
	}

	model := &registry.RegisteredModel{
		Id:   ctx.Id,
		Name: ctx.Name,
	}

	return model, nil
}

func (m *Mapper) MapToModelVersion(ctx *proto.Context, artifacts []*proto.Artifact) (*registry.VersionedModel, error) {
	if ctx.GetTypeId() != m.ModelVersionTypeId {
		return nil, fmt.Errorf("invalid TypeId, exptected %d but received %d", m.ModelVersionTypeId, ctx.GetTypeId())
	}

	metadata, err := m.MapFromProperties(ctx.CustomProperties)
	if err != nil {
		return nil, err
	}

	modelName := ctx.GetProperties()["model_name"].GetStringValue()
	version := ctx.GetProperties()["version"].GetStringValue()
	author := ctx.GetProperties()["author"].GetStringValue()

	modelVersion := &registry.VersionedModel{
		ModelName: &modelName,
		Id:        ctx.Id,
		Version:   &version,
		Author:    &author,
		Metadata:  &metadata,
	}

	modelArtifacts := []registry.Artifact{}
	if artifacts != nil {
		for _, a := range artifacts {
			art, err := m.MapToModelArtifact(a)
			if err != nil {
				return nil, err
			}
			modelArtifacts = append(modelArtifacts, *art)
		}
		modelVersion.ModelUri = *modelArtifacts[0].Uri
	}

	modelVersion.Artifacts = &modelArtifacts

	return modelVersion, nil
}

func (m *Mapper) MapToModelArtifact(artifact *proto.Artifact) (*registry.Artifact, error) {
	if artifact.GetTypeId() != m.ModelArtifactTypeId {
		return nil, fmt.Errorf("invalid TypeId, exptected %d but received %d", m.ModelArtifactTypeId, artifact.GetTypeId())
	}

	_, err := m.MapFromProperties(artifact.CustomProperties)
	if err != nil {
		return nil, err
	}

	_, err = m.MapFromProperties(artifact.Properties)
	if err != nil {
		return nil, err
	}

	modelArtifact := &registry.Artifact{
		Uri:  artifact.Uri,
		Name: artifact.Name,
	}

	return modelArtifact, nil
}
