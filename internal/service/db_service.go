package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/opendatahub-io/model-registry/internal/model/db"
	"github.com/opendatahub-io/model-registry/internal/server"
)

var _ DBService = dbServiceHandler{}
var _ DBService = (*dbServiceHandler)(nil)

func NewDBService(db *gorm.DB) DBService {
	return &dbServiceHandler{
		typeHandler:     &typeHandler{db: db},
		artifactHandler: &artifactHandler{db: db},
	}
}

type dbServiceHandler struct {
	*typeHandler
	*artifactHandler
}

type DBService interface {
	InsertType(db.Type) (*db.Type, error)
	UpsertType(db.Type) (*db.Type, error)
	ReadType(db.Type) (*db.Type, error)
	// Get functions to use a signature similar to the gorm `Where` func
	ReadAllType(query interface{}, args ...interface{}) ([]*db.Type, error)
	UpdateType(db.Type) (*db.Type, error)
	DeleteType(db.Type) (*db.Type, error)

	// InsertEEE(db.EEE) (*db.EEE, error)
	// UpsertEEE(db.EEE) (*db.EEE, error)
	// ReadEEE(db.EEE) (*db.EEE, error)
	// ReadAllEEE(query interface{}, args ...interface{}) ([]*db.EEE, error)
	// UpdateEEE(db.EEE) (*db.EEE, error)
	// DeleteEEE(db.EEE) (*db.EEE, error)
}

type typeHandler struct {
	db *gorm.DB
}

type artifactHandler struct {
	db *gorm.DB
}

func (h *typeHandler) InsertType(i db.Type) (r *db.Type, err error) {
	ctx, _ := server.Begin(context.Background(), h.db)
	defer handleTransaction(ctx, &err)

	result := h.db.Create(&i)
	if result.Error != nil {
		return nil, result.Error
	}
	return &i, nil
}

func (h *typeHandler) UpsertType(i db.Type) (r *db.Type, err error) {
	ctx, _ := server.Begin(context.Background(), h.db)
	defer handleTransaction(ctx, &err)

	if err := h.db.Where("name = ?", i.Name).Assign(i).FirstOrCreate(&i).Error; err != nil {
		err = fmt.Errorf("error creating type %s: %v", i.Name, err)
		return nil, err
	}
	return &i, nil
}

func (h *typeHandler) ReadType(i db.Type) (*db.Type, error) {
	var results []*db.Type
	rx := h.db.Find(&results)
	if rx.Error != nil {
		return nil, rx.Error
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("found more than one Type(s): %v", len(results))
	}
	return results[0], nil
}

func (h *typeHandler) ReadAllType(query interface{}, args ...interface{}) ([]*db.Type, error) {
	var results []*db.Type
	rx := h.db.Where(query, args).Find(&results)
	if rx.Error != nil {
		return nil, rx.Error
	}
	return results, nil
}

func (h *typeHandler) UpdateType(i db.Type) (result *db.Type, err error) {
	panic("unimplemented")
}

func (h *typeHandler) DeleteType(i db.Type) (r *db.Type, err error) {
	panic("unimplemented")
}

type TypeKind int32

// TODO: move from grpc to just here
const (
	EXECUTION_TYPE TypeKind = iota
	ARTIFACT_TYPE
	CONTEXT_TYPE
)

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
