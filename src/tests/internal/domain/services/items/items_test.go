package items

import (
	"context"
	"net/http"
	"testing"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/services"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/clients"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/mocks"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/repositories/items"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/stretchr/testify/assert"
)

func TestService_Get_Success(t *testing.T) {
	var result = mocks.ItemMockOne
	result.Eligible = []models.Eligible{}
	result.Price = mocks.Price

	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetPriceByItemID = func(ctx context.Context, itemID string) (models.Price, apierrors.ApiError) {
		return mocks.Price, nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGet = func(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
		var item = mocks.ItemMockOne
		item.Eligible = []models.Eligible{}
		return item, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.Get(context.TODO(), "1")

	assert.Nil(t, apiErr)
	assert.Equal(t, result, item)
}

func TestService_Get_Repository_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	priceClient := clients.NewPriceClientMock()
	repository := items.NewItemsRepositoryMock()
	repository.HandleGet = func(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
		return models.Item{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.Get(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Item{}, item)
}

func TestService_Get_Client_NotFound(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetPriceByItemID = func(ctx context.Context, itemID string) (models.Price, apierrors.ApiError) {
		return models.Price{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGet = func(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
		return mocks.ItemMockOne, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.Get(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Item{}, item)
}

func TestService_Get_Client_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetPriceByItemID = func(ctx context.Context, itemID string) (models.Price, apierrors.ApiError) {
		return models.Price{}, apierrors.NewApiError("mock error", "mock err", http.StatusInternalServerError, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGet = func(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
		return mocks.ItemMockOne, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.Get(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status())
	assert.Equal(t, models.Item{}, item)
}

func TestService_GetItemsByUserID_Success(t *testing.T) {
	var result = mocks.ItemsMock

	for i, price := range mocks.Prices.Prices {
		if len(result.Items) > i {
			result.Items[i].ID = price.ItemID
		}
	}

	result.SetPriceToItems(mocks.Prices)

	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return mocks.Prices, nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByUserID = func(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	response, apiErr := service.GetItemsByUserID(context.TODO(), "1")

	assert.Nil(t, apiErr)
	assert.Equal(t, mocks.ItemsMock.Items[0].ID, response.Items[0].ID)
	assert.Equal(t, mocks.ItemsMock.Items[0].Name, response.Items[0].Name)
	assert.Equal(t, mocks.ItemsMock.Items[1].ID, response.Items[1].ID)
	assert.Equal(t, mocks.ItemsMock.Items[1].Name, response.Items[1].Name)
}

func TestService_GetItemsByUserID_Repository_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByUserID = func(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByUserID(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByUserID_Client_NotFound(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return models.Prices{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByUserID = func(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByUserID(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByUserID_Client_Err(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return models.Prices{}, apierrors.NewApiError("mock err", "mock err", http.StatusInternalServerError, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByUserID = func(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByUserID(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByShopID_Success(t *testing.T) {
	var result = mocks.ItemsMock
	result.SetPriceToItems(mocks.Prices)

	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return mocks.Prices, nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByShopID = func(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	response, apiErr := service.GetItemsByShopID(context.TODO(), "1")

	assert.Nil(t, apiErr)
	assert.Equal(t, mocks.ItemsMock.Items[0].ID, response.Items[0].ID)
	assert.Equal(t, mocks.ItemsMock.Items[0].Name, response.Items[0].Name)
	assert.Equal(t, mocks.ItemsMock.Items[1].ID, response.Items[1].ID)
	assert.Equal(t, mocks.ItemsMock.Items[1].Name, response.Items[1].Name)
}

func TestService_GetItemsByShopID_NotFound(t *testing.T) {
	var result = mocks.ItemsMock
	result.SetPriceToItems(mocks.Prices)

	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return mocks.Prices, nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByShopID = func(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "mock error", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByShopID(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByShopID_Repository_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByShopID = func(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByShopID(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByShopID_Client_NotFound(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return models.Prices{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByShopID = func(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByShopID(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByShopID_Client_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return models.Prices{}, apierrors.NewApiError("mock error", "mock err", http.StatusInternalServerError, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByShopID = func(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByShopID(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByShopCategoryID_Success(t *testing.T) {
	var result = mocks.ItemsMock
	result.SetPriceToItems(mocks.Prices)

	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return mocks.Prices, nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByShopCategoryID = func(ctx context.Context, shopID string, itemID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	response, apiErr := service.GetItemsByShopCategoryID(context.TODO(), "1", "2")

	assert.Nil(t, apiErr)
	assert.Equal(t, mocks.ItemsMock.Items[0].ID, response.Items[0].ID)
	assert.Equal(t, mocks.ItemsMock.Items[0].Name, response.Items[0].Name)
	assert.Equal(t, mocks.ItemsMock.Items[1].ID, response.Items[1].ID)
	assert.Equal(t, mocks.ItemsMock.Items[1].Name, response.Items[1].Name)
}

func TestService_GetItemsByShopCategoryID_Repository_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByShopCategoryID = func(ctx context.Context, shopID string, itemID string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByShopCategoryID(context.TODO(), "1", "2")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByShopCategoryID_Client_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return models.Prices{}, apierrors.NewApiError("mock err", "mock err", http.StatusInternalServerError, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByShopCategoryID = func(ctx context.Context, shopID string, itemID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByShopCategoryID(context.TODO(), "1", "2")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByShopCategoryID_Client_Price_NotFound(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return models.Prices{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByShopCategoryID = func(ctx context.Context, shopID string, itemID string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByShopCategoryID(context.TODO(), "1", "2")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}
func TestService_GetItemsByIDs_Success(t *testing.T) {
	var result = mocks.ItemsMock
	result.SetPriceToItems(mocks.Prices)

	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return mocks.Prices, nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByIDs = func(ctx context.Context, itemsIDs []string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	response, apiErr := service.GetItemsByIDs(context.TODO(), mocks.ItemIds)

	assert.Nil(t, apiErr)
	assert.Equal(t, mocks.ItemsMock.Items[0].ID, response.Items[0].ID)
	assert.Equal(t, mocks.ItemsMock.Items[0].Name, response.Items[0].Name)
	assert.Equal(t, mocks.ItemsMock.Items[1].ID, response.Items[1].ID)
	assert.Equal(t, mocks.ItemsMock.Items[1].Name, response.Items[1].Name)
}

func TestService_GetItemsByIDs_Repo_NotFound(t *testing.T) {
	var result = mocks.ItemsMock
	result.SetPriceToItems(mocks.Prices)

	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return mocks.Prices, nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByIDs = func(ctx context.Context, itemsIDs []string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "mock error", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	items, apiErr := service.GetItemsByIDs(context.TODO(), mocks.ItemIds)

	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{Items: []models.Item(nil)}, items)
}

func TestService_GetItemsByIDs_Repository_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByIDs = func(ctx context.Context, itemsIDs []string) (models.Items, apierrors.ApiError) {
		return models.Items{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByIDs(context.TODO(), models.ItemsIds{Items: mocks.ItemsMock.GetItemsIds()})

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByIDs_Client_Not_Found(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return models.Prices{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByIDs = func(ctx context.Context, itemsIDs []string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	item, apiErr := service.GetItemsByIDs(context.TODO(), models.ItemsIds{Items: mocks.ItemsMock.GetItemsIds()})

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Items{}, item)
}

func TestService_GetItemsByIDs_Client_Conn_Err(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetItemsPrices = func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
		return models.Prices{}, apierrors.NewApiError("mock err", "mock err", http.StatusInternalServerError, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByIDs = func(ctx context.Context, itemsIDs []string) (models.Items, apierrors.ApiError) {
		return mocks.ItemsMock, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	items, apiErr := service.GetItemsByIDs(context.TODO(), models.ItemsIds{Items: mocks.ItemsMock.GetItemsIds()})

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status())
	assert.Equal(t, models.Items{Items: []models.Item(nil)}, items)
}

func TestService_Create_Success(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	shopClient.HandleGetShopByUserID = func(ctx context.Context) (models.Shop, apierrors.ApiError) {
		return mocks.Shop, nil
	}

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleCreatePrice = func(ctx context.Context, price *models.Price) apierrors.ApiError {
		return nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleSave = func(ctx context.Context, item models.Item) (interface{}, apierrors.ApiError) {
		return primitive.NewObjectID(), nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	_, err := service.CreateItem(context.TODO(), mocks.ItemDTO)

	assert.Nil(t, err)
}

func TestService_Create_Empty_Price_Error(t *testing.T) {
	var mock = mocks.ItemDTO
	mock.Price = dto.PriceDTO{}

	shopClient := clients.NewShopClientMock()
	shopClient.HandleGetShopByUserID = func(ctx context.Context) (models.Shop, apierrors.ApiError) {
		return mocks.Shop, nil
	}

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleCreatePrice = func(ctx context.Context, price *models.Price) apierrors.ApiError {
		return nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleSave = func(ctx context.Context, item models.Item) (interface{}, apierrors.ApiError) {
		return primitive.NewObjectID(), nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	_, err := service.CreateItem(context.TODO(), mock)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status())
}

func TestService_Create_Shop_Client_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	shopClient.HandleGetShopByUserID = func(ctx context.Context) (models.Shop, apierrors.ApiError) {
		return models.Shop{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	priceClient := clients.NewPriceClientMock()

	repository := items.NewItemsRepositoryMock()
	repository.HandleSave = func(ctx context.Context, item models.Item) (interface{}, apierrors.ApiError) {
		return primitive.NewObjectID(), nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	_, err := service.CreateItem(context.TODO(), mocks.ItemDTO)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Create_Repository_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	shopClient.HandleGetShopByUserID = func(ctx context.Context) (models.Shop, apierrors.ApiError) {
		return mocks.Shop, nil
	}

	priceClient := clients.NewPriceClientMock()

	repository := items.NewItemsRepositoryMock()
	repository.HandleSave = func(ctx context.Context, item models.Item) (interface{}, apierrors.ApiError) {
		return "", apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	_, err := service.CreateItem(context.TODO(), mocks.ItemDTO)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Create_Price_Client_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	shopClient.HandleGetShopByUserID = func(ctx context.Context) (models.Shop, apierrors.ApiError) {
		return mocks.Shop, nil
	}

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleCreatePrice = func(ctx context.Context, price *models.Price) apierrors.ApiError {
		return apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleSave = func(ctx context.Context, item models.Item) (interface{}, apierrors.ApiError) {
		return primitive.NewObjectID(), nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	_, err := service.CreateItem(context.TODO(), mocks.ItemDTO)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Update_Success(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	priceClient := clients.NewPriceClientMock()
	repository := items.NewItemsRepositoryMock()
	repository.HandleUpdate = func(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Update(context.TODO(), "1", mocks.ItemDTO)

	assert.Nil(t, err)
}

func TestService_Update_Item_NotFound(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	priceClient := clients.NewPriceClientMock()
	repository := items.NewItemsRepositoryMock()
	repository.HandleUpdate = func(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError) {
		return -1, apierrors.NewApiError("mock err", "mock err", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Update(context.TODO(), "1", mocks.ItemDTO)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestService_Update_Empty_Price_Error(t *testing.T) {
	var mock = mocks.ItemDTO
	mock.Price = dto.PriceDTO{}

	shopClient := clients.NewShopClientMock()
	priceClient := clients.NewPriceClientMock()
	repository := items.NewItemsRepositoryMock()

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Update(context.TODO(), "1", mock)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status())
}

func TestService_Update_Price_Client_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetPriceByItemID = func(ctx context.Context, itemID string) (models.Price, apierrors.ApiError) {
		return models.Price{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleUpdate = func(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Update(context.TODO(), "1", mocks.ItemDTO)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Update_Price_NotFound_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetPriceByItemID = func(ctx context.Context, itemID string) (models.Price, apierrors.ApiError) {
		return models.Price{}, apierrors.NewApiError("mock err", "mock err", http.StatusNotFound, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleUpdate = func(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Update(context.TODO(), primitive.NewObjectID().Hex(), mocks.ItemDTO)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestService_Update_Repository_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetPriceByItemID = func(ctx context.Context, itemID string) (models.Price, apierrors.ApiError) {
		return mocks.Price, nil
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleUpdate = func(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError) {
		return -1, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Update(context.TODO(), "1", mocks.ItemDTO)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Update_Price_Update_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	priceClient := clients.NewPriceClientMock()
	priceClient.HandleGetPriceByItemID = func(ctx context.Context, itemID string) (models.Price, apierrors.ApiError) {
		return mocks.Price, nil
	}
	priceClient.HandleModifyPrice = func(ctx context.Context, price *models.Price) apierrors.ApiError {
		return nil
	}

	priceClient.HandleModifyPrice = func(ctx context.Context, price *models.Price) apierrors.ApiError {
		return apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	repository := items.NewItemsRepositoryMock()
	repository.HandleUpdate = func(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Update(context.TODO(), "1", mocks.ItemDTO)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Delete_Success(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	priceClient := clients.NewPriceClientMock()
	repository := items.NewItemsRepositoryMock()
	repository.HandleDelete = func(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Delete(context.TODO(), "1")

	assert.Nil(t, err)
}

func TestService_Delete_Repository_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	priceClient := clients.NewPriceClientMock()
	repository := items.NewItemsRepositoryMock()
	repository.HandleDelete = func(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
		return -1, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Delete(context.TODO(), "1")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Delete_Prices_Client_Not_Found(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	priceClient := clients.NewPriceClientMock()
	priceClient.HandleDeletePrice = func(ctx context.Context, priceID string) apierrors.ApiError {
		return apierrors.NewApiError("mock err", "mock err", http.StatusNotFound, apierrors.CauseList{})
	}
	repository := items.NewItemsRepositoryMock()
	repository.HandleDelete = func(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Delete(context.TODO(), "1")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestService_Delete_Prices_Client_Internal_Server_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	priceClient := clients.NewPriceClientMock()
	priceClient.HandleDeletePrice = func(ctx context.Context, priceID string) apierrors.ApiError {
		return apierrors.NewApiError("mock err", "mock err", http.StatusInternalServerError, apierrors.CauseList{})
	}
	repository := items.NewItemsRepositoryMock()
	repository.HandleDelete = func(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Delete(context.TODO(), "1")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Delete_Prices_Client_Invalid_Hex_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()
	priceClient := clients.NewPriceClientMock()
	priceClient.HandleDeletePrice = func(ctx context.Context, priceID string) apierrors.ApiError {
		return apierrors.NewApiError("mock err", "mock err", http.StatusBadRequest, apierrors.CauseList{})
	}
	repository := items.NewItemsRepositoryMock()
	repository.HandleDelete = func(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewItemsService(repository, priceClient, shopClient)

	err := service.Delete(context.TODO(), "1")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status())
}

func TestService_UpdateItemsCategories_Success(t *testing.T) {

	shopClient := clients.NewShopClientMock()
	repository := items.NewItemsRepositoryMock()
	repository.HandleUpdateItemsCategories = func(ctx context.Context, category *models.Category) apierrors.ApiError {
		return nil
	}
	service := services.NewItemsService(repository, nil, shopClient)

	err := service.UpdateItemsCategories(context.TODO(), mocks.CategoryOne)

	assert.Nil(t, err)
}

func TestService_UpdateItemsCategories_Items_NotFound_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	repository := items.NewItemsRepositoryMock()
	repository.HandleUpdateItemsCategories = func(ctx context.Context, category *models.Category) apierrors.ApiError {
		return apierrors.NewApiError("mock err", "mock err", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, nil, shopClient)

	err := service.UpdateItemsCategories(context.TODO(), models.Category{})

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestService_UpdateItemsCategories_Repository_Error(t *testing.T) {
	shopClient := clients.NewShopClientMock()

	repository := items.NewItemsRepositoryMock()
	repository.HandleUpdateItemsCategories = func(ctx context.Context, category *models.Category) apierrors.ApiError {
		return apierrors.NewApiError("mock err", "mock err", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, nil, shopClient)

	err := service.UpdateItemsCategories(context.TODO(), models.Category{})

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_GetByCategoryID_Success(t *testing.T) {

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByCategoryID = func(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError) {
		return mocks.ItemsMock.Items, nil
	}
	service := services.NewItemsService(repository, nil, nil)

	items, apiErr := service.GetByCategoryID(context.TODO(), "1")

	assert.Nil(t, apiErr)
	assert.Equal(t, mocks.CategoryIDOne, items[0].Category.ID)
}

func TestService_GetByCategoryID_NoItemsFound_OK(t *testing.T) {

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByCategoryID = func(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError) {
		return []models.Item{}, nil
	}
	service := services.NewItemsService(repository, nil, nil)

	items, apiErr := service.GetByCategoryID(context.TODO(), "1")

	assert.Nil(t, apiErr)
	assert.Equal(t, []models.Item([]models.Item{}), items)
}

func TestService_GetByCategoryID_Repository_Error(t *testing.T) {

	repository := items.NewItemsRepositoryMock()
	repository.HandleGetByCategoryID = func(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError) {
		return []models.Item{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewItemsService(repository, nil, nil)

	item, apiErr := service.GetByCategoryID(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, []models.Item([]models.Item{}), item)
}
