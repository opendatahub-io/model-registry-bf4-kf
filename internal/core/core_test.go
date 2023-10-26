package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/core"
	"github.com/opendatahub-io/model-registry/internal/core/mapper"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	testutils "github.com/opendatahub-io/model-registry/test/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// common utility test variables
var (
	// generic
	ascOrderDirection  string
	descOrderDirection string
	// registered models
	modelName       string
	modelExternalId string
	owner           string
	// model version
	modelVersionName  string
	versionExternalId string
	author            string
	// model artifact
	artifactName  string
	artifactExtId string
	artifactState string
	artifactUri   string
)

func setup(t *testing.T) (*assert.Assertions, *grpc.ClientConn, proto.MetadataStoreServiceClient, func(t *testing.T)) {
	// initialize test variable before each test
	ascOrderDirection = "ASC"
	descOrderDirection = "DESC"
	modelName = "MyAwesomeModel"
	modelExternalId = "org.myawesomemodel"
	owner = "owner"
	modelVersionName = "v1"
	versionExternalId = "org.myawesomemodel@v1"
	author = "author1"
	artifactName = "Pickle model"
	artifactExtId = "org.myawesomemodel@v1:pickle"
	artifactState = "LIVE"
	artifactUri = "path/to/model/v1"

	conn, client, teardown := testutils.SetupMLMDTestContainer(t)
	return assert.New(t), conn, client, teardown
}

// initialize model registry service and assert no error is thrown
func initModelRegistryService(assertion *assert.Assertions, conn *grpc.ClientConn) core.ModelRegistryApi {
	service, err := core.NewModelRegistryService(conn)
	assertion.Nilf(err, "error creating core service: %v", err)
	return service
}

// utility function that register a new simple model and return its ID
func registerModel(assertion *assert.Assertions, service core.ModelRegistryApi, overrideModelName *string, overrideExternalId *string) string {
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	if overrideModelName != nil {
		registeredModel.Name = overrideModelName
	}

	if overrideExternalId != nil {
		registeredModel.ExternalID = overrideExternalId
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	return *createdModel.Id
}

// utility function that register a new simple model and return its ID
func registerModelVersion(
	assertion *assert.Assertions,
	service core.ModelRegistryApi,
	overrideModelName *string,
	overrideExternalId *string,
	overrideVersionName *string,
	overrideVersionExtId *string,
) string {
	registeredModelId := registerModel(assertion, service, overrideModelName, overrideExternalId)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	if overrideVersionName != nil {
		modelVersion.Name = overrideVersionName
	}

	if overrideVersionExtId != nil {
		modelVersion.ExternalID = overrideVersionExtId
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating model version: %v", err)

	return *createdVersion.Id
}

func TestModelRegistryTypes(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	_ = initModelRegistryService(assertion, conn)

	// assure the types have been correctly setup at startup
	ctx := context.Background()
	regModelResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: &core.RegisteredModelTypeName,
	})
	assertion.NotNilf(regModelResp.ContextType, "registered model type %s should exists", core.RegisteredModelTypeName)
	assertion.Equal(core.RegisteredModelTypeName, *regModelResp.ContextType.Name)

	modelVersionResp, _ := client.GetContextType(ctx, &proto.GetContextTypeRequest{
		TypeName: &core.ModelVersionTypeName,
	})
	assertion.NotNilf(modelVersionResp.ContextType, "model version type %s should exists", core.ModelVersionTypeName)
	assertion.Equal(core.ModelVersionTypeName, *modelVersionResp.ContextType.Name)

	modelArtifactResp, _ := client.GetArtifactType(ctx, &proto.GetArtifactTypeRequest{
		TypeName: &core.ModelArtifactTypeName,
	})
	assertion.NotNilf(modelArtifactResp.ArtifactType, "model version type %s should exists", core.ModelArtifactTypeName)
	assertion.Equal(core.ModelArtifactTypeName, *modelArtifactResp.ArtifactType.Name)
}

// REGISTERED MODELS

func TestCreateRegisteredModel(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	assertion.Nilf(err, "error creating registered model: %v", err)
	assertion.NotNilf(createdModel.Id, "created registered model should not have nil Id")

	byTypeAndNameResp, err := client.GetContextByTypeAndName(context.Background(), &proto.GetContextByTypeAndNameRequest{
		TypeName:    &core.RegisteredModelTypeName,
		ContextName: &modelName,
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctxId := mapper.IdToString(*byTypeAndNameResp.Context.Id)
	assertion.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	assertion.Equal(modelName, *byTypeAndNameResp.Context.Name, "saved model name should match the provided one")
	assertion.Equal(modelExternalId, *byTypeAndNameResp.Context.ExternalId, "saved external id should match the provided one")
	assertion.Equal(owner, byTypeAndNameResp.Context.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")
}

func TestUpdateRegisteredModel(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	assertion.Nilf(err, "error creating registered model: %v", err)
	assertion.NotNilf(createdModel.Id, "created registered model should not have nil Id")

	// checks created model matches original one except for Id
	assertion.Equal(*registeredModel.Name, *createdModel.Name, "returned model name should match the original one")
	assertion.Equal(*registeredModel.ExternalID, *createdModel.ExternalID, "returned model external id should match the original one")
	assertion.Equal(*registeredModel.CustomProperties, *createdModel.CustomProperties, "returned model custom props should match the original one")

	byTypeAndNameResp, err := client.GetContextByTypeAndName(context.Background(), &proto.GetContextByTypeAndNameRequest{
		TypeName:    &core.RegisteredModelTypeName,
		ContextName: &modelName,
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctxId := mapper.IdToString(*byTypeAndNameResp.Context.Id)
	assertion.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	assertion.Equal(modelName, *byTypeAndNameResp.Context.Name, "saved model name should match the provided one")
	assertion.Equal(modelExternalId, *byTypeAndNameResp.Context.ExternalId, "saved external id should match the provided one")
	assertion.Equal(owner, byTypeAndNameResp.Context.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	// update existing model
	newModelExternalId := "newExternalId"
	newOwner := "newOwner"

	createdModel.ExternalID = &newModelExternalId
	(*createdModel.CustomProperties)["owner"] = openapi.MetadataValue{
		MetadataStringValue: &openapi.MetadataStringValue{
			StringValue: &newOwner,
		},
	}

	// update the model
	createdModel, err = service.UpsertRegisteredModel(createdModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	// still one registered model
	getAllResp, err = client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	byTypeAndNameResp, err = client.GetContextByTypeAndName(context.Background(), &proto.GetContextByTypeAndNameRequest{
		TypeName:    &core.RegisteredModelTypeName,
		ContextName: &modelName,
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctxId = mapper.IdToString(*byTypeAndNameResp.Context.Id)
	assertion.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	assertion.Equal(modelName, *byTypeAndNameResp.Context.Name, "saved model name should match the provided one")
	assertion.Equal(newModelExternalId, *byTypeAndNameResp.Context.ExternalId, "saved external id should match the provided one")
	assertion.Equal(newOwner, byTypeAndNameResp.Context.CustomProperties["owner"].GetStringValue(), "saved owner custom property should match the provided one")
}

func TestGetRegisteredModelById(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"owner": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &owner,
				},
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	assertion.Nilf(err, "error creating registered model: %v", err)

	getModelById, err := service.GetRegisteredModelById(*createdModel.Id)
	assertion.Nilf(err, "error getting registered model by id %s: %v", *createdModel.Id, err)

	// checks created model matches original one except for Id
	assertion.Equal(*registeredModel.Name, *getModelById.Name, "saved model name should match the original one")
	assertion.Equal(*registeredModel.ExternalID, *getModelById.ExternalID, "saved model external id should match the original one")
	assertion.Equal(*registeredModel.CustomProperties, *getModelById.CustomProperties, "saved model custom props should match the original one")
}

func TestGetRegisteredModelByParamsName(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	byName, err := service.GetRegisteredModelByParams(&modelName, nil)
	assertion.Nilf(err, "error getting registered model by name: %v", err)

	assertion.Equalf(*createdModel.Id, *byName.Id, "the returned model id should match the retrieved by name")
}

func TestGetRegisteredModelByParamsExternalId(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	byName, err := service.GetRegisteredModelByParams(nil, &modelExternalId)
	assertion.Nilf(err, "error getting registered model by external id: %v", err)

	assertion.Equalf(*createdModel.Id, *byName.Id, "the returned model id should match the retrieved by name")
}

func TestGetRegisteredModelsOrderedById(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	orderBy := "ID"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	_, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	_, err = service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	_, err = service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	orderedById, err := service.GetRegisteredModels(core.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		assertion.Less(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}

	orderedById, err = service.GetRegisteredModels(core.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		assertion.Greater(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}
}

func TestGetRegisteredModelsOrderedByLastUpdate(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	orderBy := "LAST_UPDATE_TIME"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	firstModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	secondModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	thirdModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	// update second model
	secondModel.ExternalID = nil
	_, err = service.UpsertRegisteredModel(secondModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	orderedById, err := service.GetRegisteredModels(core.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	assertion.Equal(*firstModel.Id, *orderedById.Items[0].Id)
	assertion.Equal(*thirdModel.Id, *orderedById.Items[1].Id)
	assertion.Equal(*secondModel.Id, *orderedById.Items[2].Id)

	orderedById, err = service.GetRegisteredModels(core.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(3, int(orderedById.Size))
	assertion.Equal(*secondModel.Id, *orderedById.Items[0].Id)
	assertion.Equal(*thirdModel.Id, *orderedById.Items[1].Id)
	assertion.Equal(*firstModel.Id, *orderedById.Items[2].Id)
}

func TestGetRegisteredModelsWithPageSize(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	pageSize := int32(1)
	pageSize2 := int32(2)
	modelName := "PricingModel1"
	modelExternalId := "myExternalId1"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       &modelName,
		ExternalID: &modelExternalId,
	}

	firstModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	secondModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = &newModelName
	registeredModel.ExternalID = &newModelExternalId
	thirdModel, err := service.UpsertRegisteredModel(registeredModel)
	assertion.Nilf(err, "error creating registered model: %v", err)

	truncatedList, err := service.GetRegisteredModels(core.ListOptions{
		PageSize: &pageSize,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(1, int(truncatedList.Size))
	assertion.NotEqual("", truncatedList.NextPageToken, "next page token should not be empty")
	assertion.Equal(*firstModel.Id, *truncatedList.Items[0].Id)

	truncatedList, err = service.GetRegisteredModels(core.ListOptions{
		PageSize:      &pageSize2,
		NextPageToken: &truncatedList.NextPageToken,
	})
	assertion.Nilf(err, "error getting registered models: %v", err)

	assertion.Equal(2, int(truncatedList.Size))
	assertion.Equal("", truncatedList.NextPageToken, "next page token should be empty as list item returned")
	assertion.Equal(*secondModel.Id, *truncatedList.Items[0].Id)
	assertion.Equal(*thirdModel.Id, *truncatedList.Items[1].Id)
}

// MODEL VERSIONS

func TestCreateModelVersionWithInvalidRegisteredModelId(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	notExistingRegisteredModelId := "9999"

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &notExistingRegisteredModelId)
	assertion.NotNil(err, "model version should fail because registered model id does not exist")
	assertion.Equal(fmt.Sprintf("not a valid registered model id: %s", notExistingRegisteredModelId), err.Error())
	assertion.Nil(createdVersion)
}

func TestCreateModelVersion(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")

	createdVersionId, _ := mapper.IdToInt64(*createdVersion.Id)

	byId, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	assertion.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	assertion.Equal(*createdVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	assertion.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	assertion.Equal(versionExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	assertion.Equal(author, byId.Contexts[0].CustomProperties["author"].GetStringValue(), "saved author custom property should match the provided one")
	assertion.Equalf(core.ModelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", core.ModelVersionTypeName)

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	assertion.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd")
}

func TestUpdateModelVersion(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := mapper.IdToInt64(*createdVersion.Id)

	newExternalId := "org.my_awesome_model@v1"
	newScore := 0.95

	createdVersion.ExternalID = &newExternalId
	// keep the original name. TODO remove once https://github.com/opendatahub-io/model-registry/pull/79 got merged
	createdVersion.Name = &modelVersionName
	(*createdVersion.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: &openapi.MetadataDoubleValue{
			DoubleValue: &newScore,
		},
	}

	updatedVersion, err := service.UpsertModelVersion(createdVersion, &registeredModelId)
	assertion.Nilf(err, "error updating new model version for %d: %v", registeredModelId, err)

	updateVersionId, _ := mapper.IdToInt64(*updatedVersion.Id)
	assertion.Equal(*createdVersionId, *updateVersionId, "created and updated model version should have same id")

	byId, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	assertion.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	assertion.Equal(*updateVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	assertion.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	assertion.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	assertion.Equal(author, byId.Contexts[0].CustomProperties["author"].GetStringValue(), "saved author custom property should match the provided one")
	assertion.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	assertion.Equalf(core.ModelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", core.ModelVersionTypeName)

	getAllResp, err := client.GetContexts(context.Background(), &proto.GetContextsRequest{})
	assertion.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	fmt.Printf("%+v", getAllResp.Contexts)
	assertion.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd")
}

func TestGetModelVersionById(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := mapper.IdToInt64(*createdVersion.Id)

	getById, err := service.GetModelVersionById(*createdVersion.Id)
	assertion.Nilf(err, "error getting model version with id %d", *createdVersionId)

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	assertion.Equal(*getById.Id, *mapper.IdToString(*ctx.Id), "returned model version id should match the mlmd context one")
	// TODO uncomment once https://github.com/opendatahub-io/model-registry/pull/79 got merged, for now *getById.Name == *ctx.Name
	// assertion.Equal(fmt.Sprintf("%d:%s", registeredModelId, *getById.Name), *ctx.Name, "saved model name should match the provided one")
	assertion.Equal(*getById.ExternalID, *modelVersion.ExternalID, "saved external id should match the provided one")
	assertion.Equal(*(*getById.CustomProperties)["author"].MetadataStringValue.StringValue, author, "saved author custom property should match the provided one")
}

func TestGetModelVersionByParamsName(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := mapper.IdToInt64(*createdVersion.Id)

	// TODO use just modelVersionName once https://github.com/opendatahub-io/model-registry/pull/79 got merged
	ctxName := fmt.Sprintf("%s:%s", registeredModelId, modelVersionName)

	// TODO: fix name + parentId
	getByName, err := service.GetModelVersionByParams(&ctxName, nil, nil)
	assertion.Nilf(err, "error getting model version by name %d", *createdVersionId)

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	assertion.Equal(*mapper.IdToString(*ctx.Id), *getByName.Id, "returned model version id should match the mlmd context one")
	// TODO uncomment once https://github.com/opendatahub-io/model-registry/pull/79 got merged, for now *getById.Name == *ctx.Name
	// assertion.Equal(fmt.Sprintf("%d:%s", registeredModelId, *getById.Name), *ctx.Name, "saved model name should match the provided one")
	assertion.Equal(*ctx.ExternalId, *getByName.ExternalID, "saved external id should match the provided one")
	assertion.Equal(ctx.CustomProperties["author"].GetStringValue(), *(*getByName.CustomProperties)["author"].MetadataStringValue.StringValue, "saved author custom property should match the provided one")
}

func TestGetModelVersionByParamsExternalId(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := mapper.IdToInt64(*createdVersion.Id)

	getByExternalId, err := service.GetModelVersionByParams(nil, nil, modelVersion.ExternalID)
	assertion.Nilf(err, "error getting model version by external id %d", *modelVersion.ExternalID)

	ctxById, err := client.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	assertion.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	assertion.Equal(*mapper.IdToString(*ctx.Id), *getByExternalId.Id, "returned model version id should match the mlmd context one")
	// TODO uncomment once https://github.com/opendatahub-io/model-registry/pull/79 got merged, for now *getById.Name == *ctx.Name
	// assertion.Equal(fmt.Sprintf("%d:%s", registeredModelId, *getById.Name), *ctx.Name, "saved model name should match the provided one")
	assertion.Equal(*ctx.ExternalId, *getByExternalId.ExternalID, "saved external id should match the provided one")
	assertion.Equal(ctx.CustomProperties["author"].GetStringValue(), *(*getByExternalId.CustomProperties)["author"].MetadataStringValue.StringValue, "saved author custom property should match the provided one")
}

func TestGetModelVersions(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	registeredModelId := registerModel(assertion, service, nil, nil)

	modelVersion1 := &openapi.ModelVersion{
		Name:       &modelVersionName,
		ExternalID: &versionExternalId,
	}

	secondModelVersionName := "v2"
	secondModelVersionExtId := "org.myawesomemodel@v2"
	modelVersion2 := &openapi.ModelVersion{
		Name:       &secondModelVersionName,
		ExternalID: &secondModelVersionExtId,
	}

	thirdModelVersionName := "v3"
	thirdModelVersionExtId := "org.myawesomemodel@v3"
	modelVersion3 := &openapi.ModelVersion{
		Name:       &thirdModelVersionName,
		ExternalID: &thirdModelVersionExtId,
	}

	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	createdVersion3, err := service.UpsertModelVersion(modelVersion3, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	anotherRegModelName := "AnotherModel"
	anotherRegModelExtId := "org.another"
	anotherRegisteredModelId := registerModel(assertion, service, &anotherRegModelName, &anotherRegModelExtId)

	anotherModelVersionName := "v1.0"
	anotherModelVersionExtId := "org.another@v1.0"
	modelVersionAnother := &openapi.ModelVersion{
		Name:       &anotherModelVersionName,
		ExternalID: &anotherModelVersionExtId,
	}

	_, err = service.UpsertModelVersion(modelVersionAnother, &anotherRegisteredModelId)
	assertion.Nilf(err, "error creating new model version for %d", anotherRegisteredModelId)

	createdVersionId1, _ := mapper.IdToInt64(*createdVersion1.Id)
	createdVersionId2, _ := mapper.IdToInt64(*createdVersion2.Id)
	createdVersionId3, _ := mapper.IdToInt64(*createdVersion3.Id)

	getAll, err := service.GetModelVersions(core.ListOptions{}, nil)
	assertion.Nilf(err, "error getting all model versions")
	assertion.Equal(int32(4), getAll.Size, "expected four model versions across all registered models")

	getAllByRegModel, err := service.GetModelVersions(core.ListOptions{}, &registeredModelId)
	assertion.Nilf(err, "error getting all model versions")
	assertion.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	assertion.Equal(*mapper.IdToString(*createdVersionId1), *getAllByRegModel.Items[0].Id)
	assertion.Equal(*mapper.IdToString(*createdVersionId2), *getAllByRegModel.Items[1].Id)
	assertion.Equal(*mapper.IdToString(*createdVersionId3), *getAllByRegModel.Items[2].Id)

	// order by last update time, expecting last created as first
	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByRegModel, err = service.GetModelVersions(core.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &registeredModelId)
	assertion.Nilf(err, "error getting all model versions")
	assertion.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	assertion.Equal(*mapper.IdToString(*createdVersionId1), *getAllByRegModel.Items[2].Id)
	assertion.Equal(*mapper.IdToString(*createdVersionId2), *getAllByRegModel.Items[1].Id)
	assertion.Equal(*mapper.IdToString(*createdVersionId3), *getAllByRegModel.Items[0].Id)

	// update the second version
	newVersionExternalId := "updated.org:v2"
	createdVersion2.ExternalID = &newVersionExternalId
	createdVersion2, err = service.UpsertModelVersion(createdVersion2, &registeredModelId)
	assertion.Nilf(err, "error creating new model version for %d", registeredModelId)

	assertion.Equal(newVersionExternalId, *createdVersion2.ExternalID)

	getAllByRegModel, err = service.GetModelVersions(core.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &registeredModelId)
	assertion.Nilf(err, "error getting all model versions")
	assertion.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	assertion.Equal(*mapper.IdToString(*createdVersionId1), *getAllByRegModel.Items[2].Id)
	assertion.Equal(*mapper.IdToString(*createdVersionId2), *getAllByRegModel.Items[0].Id)
	assertion.Equal(*mapper.IdToString(*createdVersionId3), *getAllByRegModel.Items[1].Id)
}

// MODEL ARTIFACTS

func TestCreateModelArtifact(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	assertion.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	assertion.Equal(artifactName, *createdArtifact.Name)
	assertion.Equal(*state, *createdArtifact.State)
	assertion.Equal(artifactUri, *createdArtifact.Uri)
	assertion.Equal(author, *(*createdArtifact.CustomProperties)["author"].MetadataStringValue.StringValue)

	createdArtifactId, _ := mapper.IdToInt64(*createdArtifact.Id)
	getById, err := client.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	assertion.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	assertion.Equal(*createdArtifact.Name, *getById.Artifacts[0].Name)
	assertion.Equal(string(*createdArtifact.State), getById.Artifacts[0].State.String())
	assertion.Equal(*createdArtifact.Uri, *getById.Artifacts[0].Uri)
	assertion.Equal(*(*createdArtifact.CustomProperties)["author"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["author"].GetStringValue())

	modelVersionIdAsInt, _ := mapper.IdToInt64(modelVersionId)
	byCtx, _ := client.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
		ContextId: (*int64)(modelVersionIdAsInt),
	})
	assertion.Equal(1, len(byCtx.Artifacts))
	assertion.Equal(*createdArtifactId, *byCtx.Artifacts[0].Id)
}

func TestUpdateModelArtifact(t *testing.T) {
	assertion, conn, client, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	newState := "MARKED_FOR_DELETION"
	createdArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertModelArtifact(createdArtifact, &modelVersionId)
	assertion.Nilf(err, "error updating model artifact for %d: %v", modelVersionId, err)

	createdArtifactId, _ := mapper.IdToInt64(*createdArtifact.Id)
	updatedArtifactId, _ := mapper.IdToInt64(*updatedArtifact.Id)
	assertion.Equal(createdArtifactId, updatedArtifactId)

	getById, err := client.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	assertion.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	assertion.Equal(*createdArtifact.Name, *getById.Artifacts[0].Name)
	assertion.Equal(string(newState), getById.Artifacts[0].State.String())
	assertion.Equal(*createdArtifact.Uri, *getById.Artifacts[0].Uri)
	assertion.Equal(*(*createdArtifact.CustomProperties)["author"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["author"].GetStringValue())
}

func TestGetModelArtifactById(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId, _ := mapper.IdToInt64(*createdArtifact.Id)

	getById, err := service.GetModelArtifactById(*createdArtifact.Id)
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	assertion.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	assertion.Equal(artifactName, *getById.Name)
	assertion.Equal(*state, *getById.State)
	assertion.Equal(artifactUri, *getById.Uri)
	assertion.Equal(author, *(*getById.CustomProperties)["author"].MetadataStringValue.StringValue)

	assertion.Equal(*createdArtifact, *getById, "artifacts returned during creation and on get by id should be equal")
}

func TestGetModelArtifactByParams(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalID: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId, _ := mapper.IdToInt64(*createdArtifact.Id)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)

	// TODO: fix name + parentId
	getByName, err := service.GetModelArtifactByParams(&artifactName, nil, nil)
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	assertion.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	assertion.Equal(artifactName, *getByName.Name)
	assertion.Equal(artifactExtId, *getByName.ExternalID)
	assertion.Equal(*state, *getByName.State)
	assertion.Equal(artifactUri, *getByName.Uri)
	assertion.Equal(author, *(*getByName.CustomProperties)["author"].MetadataStringValue.StringValue)

	assertion.Equal(*createdArtifact, *getByName, "artifacts returned during creation and on get by name should be equal")

	getByExtId, err := service.GetModelArtifactByParams(nil, nil, &artifactExtId)
	assertion.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	assertion.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	assertion.Equal(artifactName, *getByExtId.Name)
	assertion.Equal(artifactExtId, *getByExtId.ExternalID)
	assertion.Equal(*state, *getByExtId.State)
	assertion.Equal(artifactUri, *getByExtId.Uri)
	assertion.Equal(author, *(*getByExtId.CustomProperties)["author"].MetadataStringValue.StringValue)

	assertion.Equal(*createdArtifact, *getByExtId, "artifacts returned during creation and on get by ext id should be equal")
}

func TestGetModelArtifacts(t *testing.T) {
	assertion, conn, _, teardown := setup(t)
	defer teardown(t)

	// create mode registry service
	service := initModelRegistryService(assertion, conn)

	modelVersionId := registerModelVersion(assertion, service, nil, nil, nil, nil)

	modelArtifact1 := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalID: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	secondArtifactName := "second-name"
	secondArtifactExtId := "second-ext-id"
	secondArtifactUri := "second-uri"
	modelArtifact2 := &openapi.ModelArtifact{
		Name:       &secondArtifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &secondArtifactUri,
		ExternalID: &secondArtifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	thirdArtifactName := "third-name"
	thirdArtifactExtId := "third-ext-id"
	thirdArtifactUri := "third-uri"
	modelArtifact3 := &openapi.ModelArtifact{
		Name:       &thirdArtifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &thirdArtifactUri,
		ExternalID: &thirdArtifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"author": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: &author,
				},
			},
		},
	}

	createdArtifact1, err := service.UpsertModelArtifact(modelArtifact1, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)
	createdArtifact2, err := service.UpsertModelArtifact(modelArtifact2, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)
	createdArtifact3, err := service.UpsertModelArtifact(modelArtifact3, &modelVersionId)
	assertion.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId1, _ := mapper.IdToInt64(*createdArtifact1.Id)
	createdArtifactId2, _ := mapper.IdToInt64(*createdArtifact2.Id)
	createdArtifactId3, _ := mapper.IdToInt64(*createdArtifact3.Id)

	getAll, err := service.GetModelArtifacts(core.ListOptions{}, nil)
	assertion.Nilf(err, "error getting all model artifacts")
	assertion.Equalf(int32(3), getAll.Size, "expected three model artifacts")

	assertion.Equal(*mapper.IdToString(*createdArtifactId1), *getAll.Items[0].Id)
	assertion.Equal(*mapper.IdToString(*createdArtifactId2), *getAll.Items[1].Id)
	assertion.Equal(*mapper.IdToString(*createdArtifactId3), *getAll.Items[2].Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByModelVersion, err := service.GetModelArtifacts(core.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &modelVersionId)
	assertion.Nilf(err, "error getting all model artifacts for %d", modelVersionId)
	assertion.Equalf(int32(3), getAllByModelVersion.Size, "expected three model artifacts for model version %d", modelVersionId)

	assertion.Equal(*mapper.IdToString(*createdArtifactId1), *getAllByModelVersion.Items[2].Id)
	assertion.Equal(*mapper.IdToString(*createdArtifactId2), *getAllByModelVersion.Items[1].Id)
	assertion.Equal(*mapper.IdToString(*createdArtifactId3), *getAllByModelVersion.Items[0].Id)
}
