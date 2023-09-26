package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/db"
	"github.com/opendatahub-io/model-registry/internal/server"
)

type ArtifactTypeHandler interface {
	CreateArtifactType(name string, version *string) (result *db.Type, err error)
	CreateArtifactTypeAll(name string,
		version *string,
		description *string,
		inputType *string,
		outputType *string,
		externalID *string,
		properties map[string]proto.PropertyType) (result *db.Type, err error)
	GetArtifactTypes(id *int64, name *string, version *string) ([]db.Type, error)
	GetArtifactType(id *int64, name *string, version *string) (*db.Type, error)

	DeleteArtifactTypesByName(name string) ([]db.Type, error)
	DeleteArtifactTypeByNameVersion(name string, version string) (*db.Type, error)
	UpdateArtifactType(artifactType db.Type) (*db.Type, error)
}

type ArtifactHandler interface {
	CreateArtifact(name string, version string, artifactType string) (result *db.Artifact, err error)
	CreateArtifactAll(name string, // TODO: make this optional and auto uuid the name?
		version *string, // TODO: has no mapping on RDBMS?
		artifactType string,
		uri *string,
		state *int64, // TODO: relation with grpc?
		externalID *string,
		properties map[string]*proto.Value,
		customProperties map[string]*proto.Value,
	) (result *db.Artifact, err error)
	GetArtifacts(id *int64, name *string, artifactType *string, version *string) ([]db.Artifact, error)
	GetArtifact(id *int64, name *string, artifactType *string, version *string) (*db.Artifact, error)

	// ...Delete, Update, etc.
}

// ...entityHandler(s) to be repeated for other Entities

type TypeKind int32

// TODO: move from grpc to just here
const (
	EXECUTION_TYPE TypeKind = iota
	ARTIFACT_TYPE
	CONTEXT_TYPE
)

type Handle struct {
	db *gorm.DB
}

// create if not existing (or update (by name) if already existing)
func (h *Handle) CreateArtifactType(name string, version *string) (result *db.Type, err error) {
	return h.CreateArtifactTypeAll(name, version, nil, nil, nil, nil, nil)
}

func (h *Handle) CreateArtifactTypeAll(name string,
	version *string,
	description *string,
	inputType *string,
	outputType *string,
	externalID *string,
	properties map[string]proto.PropertyType) (result *db.Type, err error) {
	ctx, _ := server.Begin(context.Background(), h.db)
	defer handleTransaction(ctx, &err)

	value := &db.Type{
		Name:     name,
		Version:  version,
		TypeKind: int32(ARTIFACT_TYPE),
	}
	if err := h.db.Where("name = ?", value.Name).Assign(value).FirstOrCreate(value).Error; err != nil {
		err = fmt.Errorf("error creating type %s: %v", value.Name, err)
		return nil, err
	}
	// TODO handle remaining attributes/properties
	return value, nil
}

func (h *Handle) GetArtifactTypes(id *int64, name *string, version *string) ([]db.Type, error) {
	by := db.Type{TypeKind: int32(ARTIFACT_TYPE), Version: version}
	if id != nil {
		by.ID = *id
	}
	if name != nil {
		by.Name = *name
	}
	var results []db.Type
	rx := h.db.Find(&results, by)
	if rx.Error != nil {
		return nil, rx.Error
	}
	return results, nil
}

func (h *Handle) GetArtifactType(id *int64, name *string, version *string) (*db.Type, error) {
	return nil, nil
}

func (h *Handle) DeleteArtifactType() (*int64, error) {
	return nil, nil
}

func (h *Handle) UpdateArtifactType() (*db.Type, error) {
	return nil, nil
}

func (h *Handle) CreateArtifactAll(name string,
	version *string,
	artifactType string,
	uri *string,
	state *int64,
	externalID *string,
	properties map[string]*proto.Value,
	customProperties map[string]*proto.Value,
) (result *db.Artifact, err error) {
	// CreateTimeSinceEpoch     int64
	// LastUpdateTimeSinceEpoch int64
	return nil, nil
}

func handleTransaction(ctx context.Context, err *error) {
	// handle panic
	if perr := recover(); perr != nil {
		_ = server.Rollback(ctx)
		*err = status.Errorf(codes.Internal, "server panic: %v", perr)
		return
	}
	if err == nil || *err == nil {
		*err = server.Commit(ctx)
	} else {
		_ = server.Rollback(ctx)
		if _, ok := status.FromError(*err); !ok {
			*err = status.Errorf(codes.Internal, "internal error: %v", *err)
		}
	}
}
