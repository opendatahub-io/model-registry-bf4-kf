package server

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type key int

var dbKey key

type txContext struct {
	db      *gorm.DB
	wrapped bool
}

// NewContext returns a new Context that carries value txContext.
func NewContext(ctx context.Context, db *gorm.DB) context.Context {
	ctxDb, ok := FromContext(ctx)
	if !ok {
		return context.WithValue(ctx, dbKey, &txContext{db.Begin(), false})
	}
	return context.WithValue(ctx, dbKey, &txContext{ctxDb, true})
}

// FromContext returns the gorm.DB value stored in ctx, if any.
func FromContext(ctx context.Context) (*gorm.DB, bool) {
	var result *gorm.DB
	tx, ok := ctx.Value(dbKey).(*txContext)
	if ok {
		result = tx.db
	}
	return result, ok
}

func Begin(ctx context.Context, db *gorm.DB) (newCtx context.Context, dbConn *gorm.DB) {
	newCtx = NewContext(ctx, db)
	dbConn, _ = FromContext(newCtx)
	return
}

func Commit(ctx context.Context) error {
	tx, ok := ctx.Value(dbKey).(*txContext)
	if !ok {
		return fmt.Errorf("missing db transaction in context %s", ctx)
	}
	if !tx.wrapped {
		return tx.db.Commit().Error
	}
	// commit in unwrapped parent context
	return nil
}

func Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(dbKey).(*txContext)
	if !ok {
		return fmt.Errorf("missing db transaction in context %s", ctx)
	}
	if !tx.wrapped {
		return tx.db.Rollback().Error
	}
	// rollback in unwrapped parent context
	return nil
}
