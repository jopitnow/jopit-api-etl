package categories

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/jopit-api-etl/src/main/api/dependencies"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/handlers"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/mocks"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/services/categories"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/services/items"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/setup"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_Get_Success(t *testing.T) {
	service := categories.NewServiceMock()
	service.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return mocks.CategoryOne, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	var result models.Category
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/category/"+primitive.NewObjectID().Hex(), nil, "")
	err := json.Unmarshal(response.Body.Bytes(), &result)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, mocks.CategoryOne, result)
}

func TestHandler_Get_Not_Found_Error(t *testing.T) {
	service := categories.NewServiceMock()
	service.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return models.Category{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/category/"+primitive.NewObjectID().Hex(), nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, http.StatusNotFound, apiError.ErrorStatus)
}

func TestHandler_Get_Server_Internal_Error(t *testing.T) {
	service := categories.NewServiceMock()
	service.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return models.Category{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/category/"+primitive.NewObjectID().Hex(), nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_Get_Server_Bad_Request_Invalid_Hex(t *testing.T) {
	service := categories.NewServiceMock()
	service.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return models.Category{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/category/a", nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_GetAllCategories_Success(t *testing.T) {
	service := categories.NewServiceMock()
	service.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return []models.Category{mocks.CategoryOne}, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	var result dto.CategoriesDTO
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/categories", nil, "")
	err := json.Unmarshal(response.Body.Bytes(), &result)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, []models.Category{mocks.CategoryOne}, result.CategoryDTO)
}

func TestHandler_GetAllCategories_Not_Found_Error(t *testing.T) {
	service := categories.NewServiceMock()
	service.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return []models.Category{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/categories", nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, http.StatusNotFound, apiError.ErrorStatus)
}

func TestHandler_GetAllCategories_Internal_Server_Error(t *testing.T) {
	service := categories.NewServiceMock()
	service.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return []models.Category{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/categories", nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_CreateCategory_Success(t *testing.T) {
	service := categories.NewServiceMock()
	service.HandleCreate = func(ctx context.Context, input models.Category) apierrors.ApiError {
		return nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items/category", nil, mocks.CategoryToJson(mocks.CategoryOne))

	assert.Equal(t, http.StatusCreated, response.Code)
}

func TestHandler_CreateCategory_Bad_Request_Parse_Error(t *testing.T) {
	service := categories.NewServiceMock()

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items/category", nil, "fake")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_CreateCategory_Bad_Request_Empty_Name_Error(t *testing.T) {
	var model = mocks.CategoryOne
	model.Name = ""

	service := categories.NewServiceMock()

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items/category", nil, mocks.CategoryToJson(model))

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_CreateCategory_Internal_Server_Error(t *testing.T) {
	service := categories.NewServiceMock()
	service.HandleCreate = func(ctx context.Context, input models.Category) apierrors.ApiError {
		return apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, nil)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items/category", nil, mocks.CategoryToJson(mocks.CategoryOne))

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_UpdateCategory_Success(t *testing.T) {

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleUpdateItemsCategories = func(ctx context.Context, category models.Category) apierrors.ApiError {
		return nil
	}

	service := categories.NewServiceMock()
	service.HandleUpdate = func(ctx context.Context, input models.Category) apierrors.ApiError {
		return nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/category", nil, mocks.CategoryToJson(mocks.CategoryOne))

	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestHandler_UpdateCategory_Bad_Request_Parse_Error(t *testing.T) {

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleUpdateItemsCategories = func(ctx context.Context, category models.Category) apierrors.ApiError {
		return nil
	}

	service := categories.NewServiceMock()
	service.HandleUpdate = func(ctx context.Context, input models.Category) apierrors.ApiError {
		return nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/category", nil, "fake")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_UpdateCategory_Bad_Request_Invalid_Hex_ID(t *testing.T) {
	var model = mocks.CategoryOne
	model.ID = "a"

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleUpdateItemsCategories = func(ctx context.Context, category models.Category) apierrors.ApiError {
		return nil
	}

	service := categories.NewServiceMock()

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/category", nil, mocks.CategoryToJson(model))

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_UpdateCategory_Bad_Request_Empty_Name_Error(t *testing.T) {

	var model = mocks.CategoryOne
	model.Name = ""

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleUpdateItemsCategories = func(ctx context.Context, category models.Category) apierrors.ApiError {
		return nil
	}

	service := categories.NewServiceMock()

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/category", nil, mocks.CategoryToJson(model))

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_UpdateCategory_Internal_Server_Error(t *testing.T) {

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleUpdateItemsCategories = func(ctx context.Context, category models.Category) apierrors.ApiError {
		return nil
	}

	service := categories.NewServiceMock()
	service.HandleUpdate = func(ctx context.Context, input models.Category) apierrors.ApiError {
		return apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/category", nil, mocks.CategoryToJson(mocks.CategoryOne))

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_UpdateCategory_Item_Internal_Service_Error(t *testing.T) {

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleUpdateItemsCategories = func(ctx context.Context, category models.Category) apierrors.ApiError {
		return apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := categories.NewServiceMock()
	service.HandleUpdate = func(ctx context.Context, input models.Category) apierrors.ApiError {
		return apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/category", nil, mocks.CategoryToJson(mocks.CategoryOne))

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_DeleteCategory_Success(t *testing.T) {

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleGetItemsByShopCategoryID = func(ctx context.Context, shopID, categoryID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := categories.NewServiceMock()
	service.HandleDelete = func(ctx context.Context, items []models.Item, categoryID string) apierrors.ApiError {
		return nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "DELETE", "/items/category/"+primitive.NewObjectID().Hex(), nil, "")

	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestHandler_DeleteCategory_Success_Conflict_Error(t *testing.T) {

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleGetItemsByShopCategoryID = func(ctx context.Context, shopID, categoryID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := categories.NewServiceMock()
	service.HandleDelete = func(ctx context.Context, items []models.Item, categoryID string) apierrors.ApiError {
		return apierrors.NewApiError("mock err", "mock err", http.StatusConflict, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "DELETE", "/items/category/"+primitive.NewObjectID().Hex(), nil, "")

	assert.Equal(t, http.StatusConflict, response.Code)
}

func TestHandler_DeleteCategory_InternalError(t *testing.T) {

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleGetItemsByShopCategoryID = func(ctx context.Context, shopID, categoryID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := categories.NewServiceMock()
	service.HandleDelete = func(ctx context.Context, items []models.Item, categoryID string) apierrors.ApiError {
		return apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "DELETE", "/items/category/"+primitive.NewObjectID().Hex(), nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_DeleteCategory_Bad_Request_Invalid_Hex_ID(t *testing.T) {
	var model = mocks.CategoryOne
	model.ID = "a"

	itemsServiceMock := items.NewItemsServiceMock()
	itemsServiceMock.HandleGetItemsByShopCategoryID = func(ctx context.Context, shopID, categoryID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := categories.NewServiceMock()

	var depend dependencies.HandlersStruct
	handler := handlers.NewCategoriesHandler(service, itemsServiceMock)
	depend.Categories = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "DELETE", "/items/category/a", nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}
