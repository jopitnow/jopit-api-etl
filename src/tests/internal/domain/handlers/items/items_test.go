package items

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/jopit-api-etl/src/main/api/dependencies"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/handlers"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/mocks"
	categories "github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/services/categories"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/services/items"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/setup"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_Get_Success(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGet = func(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
		return mocks.ItemMockOne, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	var result models.Item
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/"+primitive.NewObjectID().Hex(), nil, "")
	err := json.Unmarshal(response.Body.Bytes(), &result)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Nil(t, err)
	assert.Equal(t, mocks.ItemMockOne, result)
}

func TestHandler_Get_NotFound(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGet = func(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
		return models.Item{}, apierrors.NewApiError("mock err", "mock error", http.StatusNotFound, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	var result models.Item
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/"+primitive.NewObjectID().Hex(), nil, "")
	err := json.Unmarshal(response.Body.Bytes(), &result)

	var apiError mocks.ApiError
	err = json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, http.StatusNotFound, apiError.ErrorStatus)
}

func TestHandler_Get_Internal_Server_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGet = func(ctx context.Context, shopID string) (models.Item, apierrors.ApiError) {
		return models.Item{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/"+primitive.NewObjectID().Hex(), nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_Get_Bad_Request_Invalid_Hex(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGet = func(ctx context.Context, shopID string) (models.Item, apierrors.ApiError) {
		return models.Item{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/a", nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_GetItemsByUserID_Success(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByUserID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	var result models.Items
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items", nil, "")
	err := json.Unmarshal(response.Body.Bytes(), &result)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Nil(t, err)
	assert.Equal(t, mocks.ItemsMock, result)
}

func TestHandler_GetItemsByUserID_Internal_Server_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByUserID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items", nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_GetItemsByShopID_Success(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	var result models.Items
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/shop/"+primitive.NewObjectID().Hex(), nil, "")
	err := json.Unmarshal(response.Body.Bytes(), &result)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Nil(t, err)
	assert.Equal(t, mocks.ItemsMock, result)
}

func TestHandler_GetItemsByShopID_Internal_Server_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/shop/"+primitive.NewObjectID().Hex(), nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_GetItemsByShopID_Not_Found_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "mock error", http.StatusNotFound, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/shop/"+primitive.NewObjectID().Hex(), nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, http.StatusNotFound, apiError.ErrorStatus)
}

func TestHandler_GetItemsByShopID_Bad_Request_Invalid_Hex(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", "/items/shop/a", nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_GetItemsByShopCategoryID_Success(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopCategoryID = func(ctx context.Context, shopID, categoryID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	var result models.Items
	endpoint := fmt.Sprintf("/items/shop/%s/category/%s", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex())
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", endpoint, nil, "")
	err := json.Unmarshal(response.Body.Bytes(), &result)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Nil(t, err)
	assert.Equal(t, mocks.ItemsMock, result)
}

func TestHandler_GetItemsByShopCategoryID_Internal_Server_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopCategoryID = func(ctx context.Context, shopID, categoryID string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	endpoint := fmt.Sprintf("/items/shop/%s/category/%s", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex())
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", endpoint, nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_GetItemsByShopCategoryID_Not_Found_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopCategoryID = func(ctx context.Context, shopID, categoryID string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock err", "mock eror", http.StatusNotFound, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	endpoint := fmt.Sprintf("/items/shop/%s/category/%s", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex())
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", endpoint, nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, http.StatusNotFound, apiError.ErrorStatus)
}

func TestHandler_GetItemsByShopCategoryID_Bad_Request_Invalid_Hex(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	endpoint := fmt.Sprintf("/items/shop/%s/category/%s", "a", primitive.NewObjectID().Hex())
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "GET", endpoint, nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_GetItemsByIDs_Success(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByIDs = func(ctx context.Context, items models.ItemsIds) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	var result models.Items
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items/list", nil, mocks.ItemsIdsToJson())
	err := json.Unmarshal(response.Body.Bytes(), &result)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Nil(t, err)
	assert.Equal(t, mocks.ItemsMock, result)
}

func TestHandler_GetItemsByIDs_Bad_Request_Empty_IDs_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByIDs = func(ctx context.Context, items models.ItemsIds) (models.Items, apierrors.ApiError) {
		return models.Items{}, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items/list", nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_GetItemsByIDs_Internal_Server_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByIDs = func(ctx context.Context, items models.ItemsIds) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items/list", nil, mocks.ItemsIdsToJson())

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_GetItemsByIDs_Not_Integral_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByIDs = func(ctx context.Context, items models.ItemsIds) (models.Items, apierrors.ApiError) {
		return models.Items{}, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items/list", nil, mocks.ItemsIdsToJson())

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "false", response.Header().Get("integrity"))
}

func TestHandler_GetItemsByIDs_Bad_Request_Invalid_Hex(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items/list", nil, "asd")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_CreateItem_Success(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleCreateItem = func(ctx context.Context, itemRequest dto.ItemDTO) (interface{}, apierrors.ApiError) {
		return mocks.ItemMockOne, nil
	}

	categoriesService := categories.NewServiceMock()
	categoriesService.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return mocks.ItemMockOne.Category, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, categoriesService)
	depend.Items = handler

	var result models.Item
	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items", nil, mocks.ItemToJson())
	err := json.Unmarshal(response.Body.Bytes(), &result)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Nil(t, err)
	assert.Equal(t, mocks.ItemMockOne, result)
}

func TestHandler_CreateItem_Cat_Not_Found(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleCreateItem = func(ctx context.Context, itemRequest dto.ItemDTO) (interface{}, apierrors.ApiError) {
		return mocks.ItemMockOne, nil
	}

	categoriesService := categories.NewServiceMock()
	categoriesService.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return models.Category{}, apierrors.NewApiError("mock error", "mock errr", http.StatusNotFound, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, categoriesService)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items", nil, mocks.ItemToJson())

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, http.StatusNotFound, apiError.ErrorStatus)
}
func TestHandler_CreateItem_Bad_Request_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleCreateItem = func(ctx context.Context, itemRequest dto.ItemDTO) (interface{}, apierrors.ApiError) {
		return "fake", nil
	}

	categoriesService := categories.NewServiceMock()
	categoriesService.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return mocks.ItemMockOne.Category, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, categoriesService)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items", nil, "fake")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_CreateItem_Internal_Server_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleCreateItem = func(ctx context.Context, itemRequest dto.ItemDTO) (interface{}, apierrors.ApiError) {
		return nil, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	categoriesService := categories.NewServiceMock()
	categoriesService.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return mocks.ItemMockOne.Category, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, categoriesService)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "POST", "/items", nil, mocks.ItemToJson())

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_UpdateItem_Success(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleUpdate = func(ctx context.Context, itemID string, itemRequest dto.ItemDTO) apierrors.ApiError {
		return nil
	}

	categoriesService := categories.NewServiceMock()
	categoriesService.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return mocks.ItemMockOne.Category, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, categoriesService)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/"+primitive.NewObjectID().Hex(), nil, mocks.ItemToJson())

	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestHandler_UpdateItem_Cat_Not_Found(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleUpdate = func(ctx context.Context, itemID string, itemRequest dto.ItemDTO) apierrors.ApiError {
		return nil
	}

	categoriesService := categories.NewServiceMock()
	categoriesService.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return models.Category{}, apierrors.NewApiError("category not found", "mock error", http.StatusNotFound, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, categoriesService)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/"+primitive.NewObjectID().Hex(), nil, mocks.ItemToJson())

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, http.StatusNotFound, apiError.ErrorStatus)
}

func TestHandler_UpdateItem_Bad_Request_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleUpdate = func(ctx context.Context, itemID string, itemRequest dto.ItemDTO) apierrors.ApiError {
		return nil
	}
	categoriesService := categories.NewServiceMock()
	categoriesService.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return mocks.ItemMockOne.Category, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, categoriesService)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/"+primitive.NewObjectID().Hex(), nil, "fake")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_UpdateItem_Internal_Server_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleUpdate = func(ctx context.Context, itemID string, itemRequest dto.ItemDTO) apierrors.ApiError {
		return apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	categoriesService := categories.NewServiceMock()
	categoriesService.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return mocks.ItemMockOne.Category, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, categoriesService)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/"+primitive.NewObjectID().Hex(), nil, mocks.ItemToJson())

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_UpdateItem_Bad_Request_Invalid_Hex(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	categoriesService := categories.NewServiceMock()
	categoriesService.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return mocks.ItemMockOne.Category, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, categoriesService)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "PUT", "/items/a", nil, mocks.ItemToJson())

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}

func TestHandler_DeleteItem_Success(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleDelete = func(ctx context.Context, itemID string) apierrors.ApiError {
		return nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "DELETE", "/items/"+primitive.NewObjectID().Hex(), nil, "")

	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestHandler_DeleteItem_Internal_Server_Error(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleDelete = func(ctx context.Context, itemID string) apierrors.ApiError {
		return apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "DELETE", "/items/"+primitive.NewObjectID().Hex(), nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, http.StatusInternalServerError, apiError.ErrorStatus)
}

func TestHandler_DeleteItem_Bad_Request_Invalid_Hex(t *testing.T) {

	service := items.NewItemsServiceMock()
	service.HandleGetItemsByShopID = func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	var depend dependencies.HandlersStruct
	handler := handlers.NewItemsHandler(service, nil)
	depend.Items = handler

	response := setup.ExecuteRequest(setup.BuildRouter(depend), "DELETE", "/items/a", nil, "")

	var apiError mocks.ApiError
	err := json.Unmarshal(response.Body.Bytes(), &apiError)
	if err != nil {
		panic("Cannot decode error response body.")
	}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, http.StatusBadRequest, apiError.ErrorStatus)
}
