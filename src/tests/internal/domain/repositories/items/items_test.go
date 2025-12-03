package items

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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var depMock mockdeppkg.Dependencies
var repository repositories.ItemsRepository

func TestMain(m *testing.M) {
	depMock = setup.BeforeMemongoTestCase()
	m.Run()
	setup.AfterMemongoTestCase()
	depMock.ItemsRepository = nil
}

func TestRepository_NewRepository(t *testing.T) {
	assert.NotNil(t, depMock.ItemsRepository)
}

func TestRepository_Save_Success(t *testing.T) {

	for i := 0; i < len(mocks.ItemsMock.Items); i++ {
		id, err := depMock.ItemsRepository.Save(context.TODO(), mocks.ItemsMock.Items[i])

		assert.EqualValues(t, nil, err)
		assert.NotEqualValues(t, nil, id)
	}
}

func TestRepository_Save_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		var item = mocks.ItemsMock.Items[0]

		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		id, err := repository.Save(context.TODO(), item)

		assert.EqualValues(mock, nil, id)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_Get_Success_By_Id(t *testing.T) {

	arrangetItem := mocks.ItemMockOne
	arrangetItem.ID = ""
	idinterface, err := depMock.ItemsRepository.Save(context.Background(), arrangetItem)
	if err != nil {
		log.Fatal(err)
	}
	hexID := fmt.Sprint(idinterface)[10 : len(fmt.Sprint(idinterface))-2]
	arrangetItem.ID = hexID

	assertItem, err := depMock.ItemsRepository.Get(context.TODO(), arrangetItem.ID)

	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, assertItem.Name, mocks.ItemMockOne.Name)
	assert.EqualValues(t, assertItem.Description, mocks.ItemMockOne.Description)
	assert.EqualValues(t, assertItem.Status, mocks.ItemMockOne.Status)
	assert.EqualValues(t, assertItem.Images, mocks.ItemMockOne.Images)
	assert.EqualValues(t, assertItem.Attributes, mocks.ItemMockOne.Attributes)
	assert.EqualValues(t, assertItem.Eligible, mocks.ItemMockOne.Eligible)
}

func TestRepository_Get_Error_With_No_Items(t *testing.T) {

	item, err := depMock.ItemsRepository.Get(context.TODO(), primitive.NewObjectID().Hex())

	assert.EqualValues(t, models.Item{}, item)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_Get_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		item, err := repository.Get(context.TODO(), "123")

		assert.EqualValues(t, models.Item{}, item)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_Get_Bad_Parse_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("BadParseError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

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

func TestRepository_GetByUserID_Success(t *testing.T) {

	itemTest, err := depMock.ItemsRepository.GetByUserID(context.TODO(), mocks.UserIdOne)

	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, mocks.ItemMockOne.ID, itemTest.Items[0].ID)
	assert.EqualValues(t, mocks.ItemMockOne.Name, itemTest.Items[0].Name)
	assert.EqualValues(t, mocks.ItemMockOne.UserID, itemTest.Items[0].UserID)

}

func TestRepository_GetByUserID_Error_With_No_Prices(t *testing.T) {
	item, err := depMock.ItemsRepository.GetByUserID(context.TODO(), "!")

	assert.EqualValues(t, models.Items{}, item)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_GetByUserID_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		item, err := repository.GetByUserID(context.TODO(), mocks.UserIdOne)

		assert.EqualValues(t, models.Items{}, item)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetByUserID_Bad_Parse_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("BadParseError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

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

		_, err := repository.GetByUserID(context.TODO(), mocks.UserIdOne)

		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetByShopID_Success(t *testing.T) {
	items, err := depMock.ItemsRepository.GetByShopID(context.TODO(), mocks.ShopIDOne)

	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, items.Items[0].Name, mocks.ItemMockOne.Name)
	assert.EqualValues(t, items.Items[0].Description, mocks.ItemMockOne.Description)
	assert.EqualValues(t, items.Items[0].Status, mocks.ItemMockOne.Status)
	assert.EqualValues(t, items.Items[0].Images, mocks.ItemMockOne.Images)
	assert.EqualValues(t, items.Items[0].Attributes, mocks.ItemMockOne.Attributes)
	assert.EqualValues(t, items.Items[0].Eligible, mocks.ItemMockOne.Eligible)
}

func TestRepository_GetByShopID_Error_With_No_Prices(t *testing.T) {
	item, err := depMock.ItemsRepository.GetByShopID(context.TODO(), "!")

	assert.EqualValues(t, models.Items{}, item)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_GetByShopID_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		item, err := repository.GetByShopID(context.TODO(), mocks.ShopIDOne)

		assert.EqualValues(t, models.Items{}, item)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetByShopID_Bad_Parse_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("BadParseError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

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

		_, err := repository.GetByShopID(context.TODO(), mocks.ShopIDOne)

		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetByShopCategoryID_Success(t *testing.T) {

	itemTest, err := depMock.ItemsRepository.GetByShopCategoryID(context.TODO(), mocks.ShopIDOne, mocks.CategoryIDOne)

	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, mocks.ItemMockOne.ID, itemTest.Items[0].ID)
	assert.EqualValues(t, mocks.ItemMockOne.Name, itemTest.Items[0].Name)
	assert.EqualValues(t, mocks.ItemMockOne.UserID, itemTest.Items[0].UserID)
}

func TestRepository_GetByShopCategoryID_Error_With_No_Prices(t *testing.T) {
	item, err := depMock.ItemsRepository.GetByShopCategoryID(context.TODO(), mocks.ShopIDOne, "!")

	assert.EqualValues(t, models.Items{}, item)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_GetByShopCategoryID_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		item, err := repository.GetByShopCategoryID(context.TODO(), mocks.CategoryIDOne, mocks.CategoryIDOne)

		assert.EqualValues(t, models.Items{}, item)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetByShopCategoryID_Bad_Parse_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("BadParseError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

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

		_, err := repository.GetByShopCategoryID(context.TODO(), mocks.CategoryIDOne, mocks.CategoryIDOne)

		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetByIDs_Success(t *testing.T) {

	arrangetItem := mocks.ItemMockOne
	arrangetItem.ID = ""
	idinterface, err := depMock.ItemsRepository.Save(context.Background(), arrangetItem)
	if err != nil {
		log.Fatal(err)
	}
	hexID := fmt.Sprint(idinterface)[10 : len(fmt.Sprint(idinterface))-2]
	arrangetItem.ID = hexID

	itemTest, err := depMock.ItemsRepository.GetByIDs(context.TODO(), []string{hexID})

	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, arrangetItem.ID, itemTest.Items[0].ID)
	assert.EqualValues(t, arrangetItem.Name, itemTest.Items[0].Name)
	assert.EqualValues(t, arrangetItem.UserID, itemTest.Items[0].UserID)
}

func TestRepository_GetByIDs_Error_With_No_Prices(t *testing.T) {
	item, err := depMock.ItemsRepository.GetByIDs(context.TODO(), []string{primitive.NewObjectID().Hex()})

	assert.EqualValues(t, models.Items{}, item)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_GetByIDs_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		item, err := repository.GetByIDs(context.TODO(), []string{primitive.NewObjectID().Hex()})

		assert.EqualValues(t, models.Items{}, item)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetByIDs_Bad_Parse_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("BadParseError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

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

		_, err := repository.GetByIDs(context.TODO(), []string{primitive.NewObjectID().Hex()})

		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_Update_Success(t *testing.T) {

	arrangetItem := mocks.ItemMockOne
	arrangetItem.ID = ""
	idinterface, err := depMock.ItemsRepository.Save(context.Background(), arrangetItem)
	if err != nil {
		log.Fatal(err)
	}
	hexID := fmt.Sprint(idinterface)[10 : len(fmt.Sprint(idinterface))-2]
	arrangetItem.ID = ""
	arrangetItem.Name = "update"

	result, update_err := depMock.ItemsRepository.Update(context.TODO(), hexID, &arrangetItem)
	arrangetItem.ID = hexID
	updated, _ := depMock.ItemsRepository.Get(context.TODO(), hexID)

	assert.EqualValues(t, arrangetItem.ID, updated.ID)
	assert.EqualValues(t, arrangetItem.Name, updated.Name)
	assert.Equal(t, arrangetItem.Name, updated.Name)
	assert.Nil(t, update_err)
	assert.EqualValues(t, 1, int(result))
}

func TestRepository_Update_Not_Found_Error(t *testing.T) {
	item := models.Item{
		Name: "New Name",
	}

	result, err := depMock.ItemsRepository.Update(context.TODO(), primitive.NewObjectID().Hex(), &item)

	assert.EqualValues(t, -1, result)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_Update_Invalid_Id_Error(t *testing.T) {
	item := models.Item{
		Name: "New Name",
	}

	result, err := depMock.ItemsRepository.Update(context.TODO(), "fake_id", &item)

	assert.EqualValues(t, -1, result)
	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestRepository_Update_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		item := models.Item{
			Name: "New Name",
		}

		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		result, err := repository.Update(context.TODO(), primitive.NewObjectID().Hex(), &item)

		assert.EqualValues(mock, -1, result)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_Delete_Success(t *testing.T) {

	arrangetItem := mocks.ItemMockOne
	arrangetItem.ID = ""
	idinterface, err := depMock.ItemsRepository.Save(context.Background(), arrangetItem)
	if err != nil {
		log.Fatal(err)
	}
	arrangetItem.ID = fmt.Sprint(idinterface)[10 : len(fmt.Sprint(idinterface))-2]

	result, delete_err := depMock.ItemsRepository.Delete(context.TODO(), arrangetItem.ID)
	_, get_err := depMock.ItemsRepository.Get(context.TODO(), mocks.ShopIDOne)

	assert.Nil(t, delete_err)
	assert.NotNil(t, get_err)
	assert.EqualValues(t, 1, result)
}

func TestRepository_Delete_Not_Found_Error(t *testing.T) {
	result, err := depMock.ItemsRepository.Delete(context.TODO(), primitive.NewObjectID().Hex())

	assert.EqualValues(t, -1, result)
	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_Delete_Invalid_Id_Error(t *testing.T) {
	result, err := depMock.ItemsRepository.Delete(context.TODO(), "fake_id")

	assert.EqualValues(t, -1, result)
	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestRepository_Delete_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		result, err := repository.Delete(context.TODO(), mocks.CategoryOne.ID)

		assert.EqualValues(mock, -1, result)
		assert.EqualValues(mock, "internal_server_error", err.Code())
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_UpdateItemsCategories_Success(t *testing.T) {

	arrangetItem := mocks.ItemMockOne
	arrangetItem.ID = ""
	_, err := depMock.ItemsRepository.Save(context.Background(), arrangetItem)
	if err != nil {
		log.Fatal(err)
	}

	arrangeCat := &mocks.CategoryOne
	arrangeCat.Name = "sape"
	arrangeCat.ID = mocks.CategoryIDOne

	err = depMock.ItemsRepository.UpdateItemsCategories(context.TODO(), arrangeCat)

	assert.Nil(t, err)
}

func TestRepository_UpdateItemsCategories_Not_Found_Error(t *testing.T) {

	catAssert := &mocks.CategoryOne
	catAssert.ID = primitive.NewObjectID().Hex()

	err := depMock.ItemsRepository.UpdateItemsCategories(context.TODO(), catAssert)

	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestRepository_UpdateItemsCategories_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {

		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		err := repository.UpdateItemsCategories(context.TODO(), &mocks.CategoryOne)

		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetByCategoryID_Success_By_Id(t *testing.T) {

	arrangetItem := mocks.ItemMockOne
	arrangetItem.ID = ""

	assertItems, err := depMock.ItemsRepository.GetByCategoryID(context.TODO(), mocks.CategoryIDOne)

	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, mocks.CategoryIDOne, assertItems[0].Category.ID)
}

func TestRepository_Get_With_No_Items_OK(t *testing.T) {

	items, err := depMock.ItemsRepository.GetByCategoryID(context.TODO(), primitive.NewObjectID().Hex())

	assert.EqualValues(t, []models.Item(nil), items)
	assert.Nil(t, err)
}

func TestRepository_GetByCategoryID_Internal_Server_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("InternalServerError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

		mock.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

		items, err := repository.GetByCategoryID(context.TODO(), primitive.NewObjectID().Hex())

		assert.EqualValues(t, []models.Item([]models.Item{}), items)
		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}

func TestRepository_GetByCategoryID_Bad_Parse_Error(t *testing.T) {
	mock := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mock.Cleanup(func() { mock.Client = nil })

	mock.Run("BadParseError", func(mock *mtest.T) {
		repository := repositories.NewItemsRepository(mock.DB.Collection(dependencies.KvsCategoriesCollection))

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

		_, err := repository.GetByCategoryID(context.TODO(), primitive.NewObjectID().Hex())

		assert.EqualValues(mock, http.StatusInternalServerError, err.Status())
	})
}
