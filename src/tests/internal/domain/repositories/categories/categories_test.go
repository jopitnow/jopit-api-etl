package categories

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/jopitnow/jopit-api-etl/src/main/api/dependencies"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/repositories"
	mockdeppkg "github.com/jopitnow/jopit-api-etl/src/tests/internal/api/dependencies"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/mocks"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/setup"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var depMock mockdeppkg.Dependencies
var idTest *string

func TestMain(m *testing.M) {
	depMock = setup.BeforeMemongoTestCase()
	m.Run()
	setup.AfterMemongoTestCase()
	depMock.CategoriesRepository = nil
}

func TestRepository_NewRepository(t *testing.T) {
	assert.NotNil(t, depMock.CategoriesRepository)
}

func TestRepository_GetAll_Empty_Error(t *testing.T) {
	categories, err := depMock.CategoriesRepository.GetAllCategories(context.TODO())

	assert.EqualValues(t, []models.Category(nil), categories)
	assert.Nil(t, err)
}

func TestRepository_Save_Success(t *testing.T) {
	var categories = mocks.Categories

	for i := 0; i < len(categories); i++ {
		id, err := depMock.CategoriesRepository.Create(context.TODO(), categories[i])

		idTest = &categories[i].ID

		assert.EqualValues(t, nil, err)
		assert.NotEqualValues(t, nil, id)
	}
}

func TestRepository_Save_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		var category = mocks.CategoryOne

		repository := repositories.NewCategoriesRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		id, err := repository.Create(context.TODO(), category)

		assert.EqualValues(mock, nil, id)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetAll_Success(t *testing.T) {
	categories, err := depMock.CategoriesRepository.GetAllCategories(context.TODO())

	assert.True(t, len(categories) == 3)

	for i := 0; i < len(categories); i++ {
		assert.EqualValues(t, nil, err)
		assert.EqualValues(t, mocks.Categories[i].Name, categories[i].Name)
	}

}

func TestRepository_GetAll_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewCategoriesRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		categories, err := repository.GetAllCategories(context.TODO())

		assert.EqualValues(mock, []models.Category{}, categories)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetAll_Bad_Parse_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("BadParseError", func(mock *mtest.T) {
		repository := repositories.NewCategoriesRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"foo.bar",
				mtest.FirstBatch,
				bson.D{
					{"_id", "660a81b04f151e0d0f710192"},
					{"mock", "test"},
				},
			),
		)

		categories, err := repository.GetAllCategories(context.TODO())

		assert.EqualValues(mock, []models.Category{}, categories)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_Get_Success_By_Id(t *testing.T) {

	arrangeCat := models.Category{Name: "sape"}

	idinterface, err := depMock.CategoriesRepository.Create(context.TODO(), arrangeCat)
	if err != nil {
		log.Fatal(err)
	}
	hexID := fmt.Sprint(idinterface)[10 : len(fmt.Sprint(idinterface))-2]
	arrangeCat.ID = hexID

	assertCat, err := depMock.CategoriesRepository.Get(context.Background(), hexID)

	assert.EqualValues(t, nil, err)
	assert.Equal(t, arrangeCat, assertCat)
}

func TestRepository_Get_Error_With_No_Prices(t *testing.T) {
	category, err := depMock.CategoriesRepository.Get(context.TODO(), primitive.NewObjectID().Hex())

	assert.EqualValues(t, models.Category{}, category)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_Get_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewCategoriesRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		category, err := repository.Get(context.TODO(), "123")

		assert.EqualValues(t, models.Category{}, category)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_Get_Bad_Parse_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("BadParseError", func(mock *mtest.T) {
		repository := repositories.NewCategoriesRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"foo.bar",
				mtest.FirstBatch,
				bson.D{
					{"_id", 123},
					{"mock", "test"},
				},
			),
		)

		_, err := repository.Get(context.TODO(), primitive.NewObjectID().Hex())

		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_Update_Success(t *testing.T) {

	arrangeCat := models.Category{Name: "sape"}
	idinterface, err := depMock.CategoriesRepository.Create(context.TODO(), arrangeCat)
	if err != nil {
		log.Fatal(err)
	}
	hexID := fmt.Sprint(idinterface)[10 : len(fmt.Sprint(idinterface))-2]
	arrangeCat.ID = hexID
	//change name of cat
	arrangeCat.Name = "update"

	result, err := depMock.CategoriesRepository.Update(context.TODO(), arrangeCat)

	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, 1, result)
}

func TestRepository_Update_Not_Found_Error(t *testing.T) {
	category := models.Category{
		ID:   primitive.NewObjectID().Hex(),
		Name: "New Name",
	}

	result, err := depMock.CategoriesRepository.Update(context.TODO(), category)

	assert.EqualValues(t, -1, result)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_Update_Invalid_Id_Error(t *testing.T) {
	category := models.Category{
		ID:   "fake_id",
		Name: "New Name",
	}

	result, err := depMock.CategoriesRepository.Update(context.TODO(), category)

	assert.EqualValues(t, -1, result)
	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestRepository_Update_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		category := models.Category{
			ID:   primitive.NewObjectID().Hex(),
			Name: "New Name",
		}

		repository := repositories.NewCategoriesRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		result, err := repository.Update(context.TODO(), category)

		assert.EqualValues(mock, -1, result)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_Delete_Success(t *testing.T) {

	arrangeCat := models.Category{Name: "sape"}
	idinterface, err := depMock.CategoriesRepository.Create(context.TODO(), arrangeCat)
	if err != nil {
		log.Fatal(err)
	}
	hexID := fmt.Sprint(idinterface)[10 : len(fmt.Sprint(idinterface))-2]
	arrangeCat.ID = hexID

	result, err := depMock.CategoriesRepository.Delete(context.TODO(), hexID)

	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, 1, result)
}

func TestRepository_Delete_Not_Found_Error(t *testing.T) {
	result, err := depMock.CategoriesRepository.Delete(context.TODO(), primitive.NewObjectID().Hex())

	assert.EqualValues(t, -1, result)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_Delete_Invalid_Id_Error(t *testing.T) {
	result, err := depMock.CategoriesRepository.Delete(context.TODO(), "fake_id")

	assert.EqualValues(t, -1, result)
	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestRepository_Delete_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewCategoriesRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		result, err := repository.Delete(context.TODO(), mocks.CategoryOne.ID)

		assert.EqualValues(mock, -1, result)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}
