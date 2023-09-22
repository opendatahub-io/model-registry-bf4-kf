package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func migrateDatabase(dbConn *gorm.DB) error {
	// using only needed RDBMS type for the scope under test
	err := dbConn.AutoMigrate(
		db.Type{},
	)
	if err != nil {
		return fmt.Errorf("db migration failed: %w", err)
	}
	return nil
}

func setup(tmpFile *os.File) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = migrateDatabase(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Bare minimal test of PutArtifactType with a given Name
func TestPutArtifactType(t *testing.T) {
	f, err := os.CreateTemp("", "model-registry-db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())
	db, err := setup(f)
	if err != nil {
		t.Errorf("Should expect DB connection: %v", err)
	}
	grpcServer := NewGrpcServer(db)
	artifactTypeName := "test"
	artifactType := proto.ArtifactType{
		Name: &artifactTypeName,
	}
	req := proto.PutArtifactTypeRequest{
		ArtifactType: &artifactType,
	}

	response, err := grpcServer.PutArtifactType(context.Background(), &req)
	if err != nil {
		t.Errorf("Should PutArtifactType: %v", err)
	}
	t.Logf("Should PutArtifactType: %v", response)

	if *response.TypeId > 0 {
		t.Logf("Should PutArtifactTypeResponse.TypeId an ID > 0: %v", response)
	} else {
		t.Errorf("Should PutArtifactTypeResponse.TypeId an ID > 0: %v", response)
	}
}

// Attemping gRPC of PutArtifactType without a Name should error
func TestPutArtifactTypeNoNameFailure(t *testing.T) {
	f, err := os.CreateTemp("", "model-registry-db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())
	db, err := setup(f)
	if err != nil {
		t.Errorf("Should expect DB connection: %v", err)
	}
	grpcServer := NewGrpcServer(db)
	versionName := "test"
	artifactType := proto.ArtifactType{
		Name:    nil,
		Version: &versionName,
	}
	req := proto.PutArtifactTypeRequest{
		ArtifactType: &artifactType,
	}

	response, err := grpcServer.PutArtifactType(context.Background(), &req)
	if err != nil {
		t.Logf("Should error on PutArtifactType with Request missing Name: %v", err)
	} else {
		t.Errorf("Should error on PutArtifactType with Request missing Name: %v, %v", response, err)
	}
}
