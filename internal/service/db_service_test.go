package service

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/model/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func migrateDatabase(dbConn *gorm.DB) error {
	// using only needed RDBMS type for the scope under test
	err := dbConn.AutoMigrate(
		db.Type{},
		db.TypeProperty{},
		// TODO: add as needed.
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

// Bare minimal test of PutArtifactType with a given Name, and Get.
func TestPutArtifactTypeThenGet(t *testing.T) {
	f, err := os.CreateTemp("", "model-registry-db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())
	db, err := setup(f)
	if err != nil {
		t.Errorf("Should expect DB connection: %v", err)
	}
	dal := Handle{
		db: db,
	}
	artifactName := "John Doe"
	at, err := dal.CreateArtifactType(artifactName, nil)
	if err != nil {
		t.Errorf("Should create ArtifactType: %v", err)
	}
	if at.ID < 0 {
		t.Errorf("Should have ID for ArtifactType: %v", at.ID)
	}
	if at.Name != artifactName {
		t.Errorf("Should have Name for ArtifactType per constant: %v", at.Name)
	}

	ats, err2 := dal.GetArtifactTypes(nil, &artifactName, nil)
	if err2 != nil {
		t.Errorf("Should get ArtifactType: %v", err2)
	}
	if len(ats) != 1 { // TODO if temp file is okay, this is superfluos
		t.Errorf("The test is running under different assumption")
	}
	at0 := ats[0]
	t.Logf("at0: %v", at0)
	if at0.ID != at.ID {
		t.Errorf("Should have same ID")
	}
	if at0.Name != at.Name {
		t.Errorf("Should have same Name")
	}

}

func TestGetArtifactTypesByCommonCriteria(t *testing.T) {
	f, err := os.CreateTemp("", "model-registry-db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())
	db, err := setup(f)
	if err != nil {
		t.Errorf("Should expect DB connection: %v", err)
	}
	dal := Handle{
		db: db,
	}

	fixVersion := "version"
	if _, err := dal.CreateArtifactType("at0", &fixVersion); err != nil {
		t.Errorf("Should create ArtifactType: %v", err)
	}
	if _, err := dal.CreateArtifactType("at1", &fixVersion); err != nil {
		t.Errorf("Should create ArtifactType: %v", err)
	}

	// TODO here only demonstrating criteria using "version", but likely more meaningful to use property as criteria
	results, err := dal.GetArtifactTypes(nil, nil, &fixVersion)
	t.Logf("results: %v", results)
	if err != nil {
		t.Errorf("Should get ArtifactTypes: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Should have retrieved 2 artifactTypes")
	}
}

func TestPutArtifactTypeSameNameDiffVersion(t *testing.T) {
	f, err := os.CreateTemp("", "model-registry-db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())
	db, err := setup(f)
	if err != nil {
		t.Errorf("Should expect DB connection: %v", err)
	}
	dal := Handle{
		db: db,
	}

	artifactName := "John Doe"
	v0 := "v0"
	v1 := "v1"
	at0, err := dal.CreateArtifactType(artifactName, &v0)
	if err != nil {
		t.Errorf("Should create ArtifactType: %v", err)
	}
	at1, err := dal.CreateArtifactType(artifactName, &v1)
	if err != nil {
		t.Errorf("Should create ArtifactType: %v", err)
	}
	if at0.ID > at1.ID {
		t.Errorf("ID invariant does not hold")
	}

	// TODO implement validation logic or RDBMS constraint/key
	// if _, err := dal.CreateArtifactType(artifactName, &v1); err == nil {
	// 	t.Errorf("Created multiple artifact with the same version")
	// }
}
