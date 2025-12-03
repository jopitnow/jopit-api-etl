package items

import (
	"context"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

type ServiceMock struct {
	HandleGet                      func(ctx context.Context, itemID string) (models.Item, apierrors.ApiError)
	HandleDelete                   func(ctx context.Context, itemID string) apierrors.ApiError
	HandleCreateItem               func(ctx context.Context, itemRequest dto.ItemDTO) (interface{}, apierrors.ApiError)
	HandleGetAll                   func(ctx context.Context) ([]models.Item, apierrors.ApiError)
	HandleGetItemsByUserID         func(ctx context.Context, userID string) (models.Items, apierrors.ApiError)
	HandleGetItemsByShopID         func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError)
	HandleGetItemsByShopCategoryID func(ctx context.Context, shopID string, categoryID string) (models.Items, apierrors.ApiError)
	HandleGetItemsByIDs            func(ctx context.Context, items models.ItemsIds) (models.Items, apierrors.ApiError)
	HandleUpdate                   func(ctx context.Context, itemID string, itemRequest dto.ItemDTO) apierrors.ApiError

	HandleGetByCategoryID       func(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError)
	HandleUpdateItemsCategories func(ctx context.Context, category models.Category) apierrors.ApiError
}

func NewItemsServiceMock() ServiceMock {
	return ServiceMock{}
}

func (mock ServiceMock) Get(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
	if mock.HandleGet != nil {
		return mock.HandleGet(ctx, itemID)
	}
	return models.Item{}, nil
}

func (mock ServiceMock) Delete(ctx context.Context, itemID string) apierrors.ApiError {
	if mock.HandleDelete != nil {
		return mock.HandleDelete(ctx, itemID)
	}
	return nil
}

func (mock ServiceMock) CreateItem(ctx context.Context, itemRequest dto.ItemDTO) (interface{}, apierrors.ApiError) {
	if mock.HandleCreateItem != nil {
		return mock.HandleCreateItem(ctx, itemRequest)
	}
	return nil, nil
}

func (mock ServiceMock) GetAll(ctx context.Context) ([]models.Item, apierrors.ApiError) {
	if mock.HandleGetAll != nil {
		return mock.HandleGetAll(ctx)
	}
	return nil, nil
}

func (mock ServiceMock) GetItemsByUserID(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
	if mock.HandleGetItemsByUserID != nil {
		return mock.HandleGetItemsByUserID(ctx, userID)
	}
	return models.Items{}, nil
}

func (mock ServiceMock) GetItemsByShopID(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
	if mock.HandleGetItemsByShopID != nil {
		return mock.HandleGetItemsByShopID(ctx, shopID)
	}
	return models.Items{}, nil
}

func (mock ServiceMock) GetItemsByShopCategoryID(ctx context.Context, shopID string, categoryID string) (models.Items, apierrors.ApiError) {
	if mock.HandleGetItemsByShopCategoryID != nil {
		return mock.HandleGetItemsByShopCategoryID(ctx, shopID, categoryID)
	}
	return models.Items{}, nil
}

func (mock ServiceMock) GetItemsByIDs(ctx context.Context, items models.ItemsIds) (models.Items, apierrors.ApiError) {
	if mock.HandleGetItemsByIDs != nil {
		return mock.HandleGetItemsByIDs(ctx, items)
	}
	return models.Items{}, nil
}

func (mock ServiceMock) Update(ctx context.Context, itemID string, itemRequest dto.ItemDTO) apierrors.ApiError {
	if mock.HandleUpdate != nil {
		return mock.HandleUpdate(ctx, itemID, itemRequest)
	}
	return nil
}

func (mock ServiceMock) UpdateItemsCategories(ctx context.Context, category models.Category) apierrors.ApiError {
	if mock.HandleUpdateItemsCategories != nil {
		return mock.HandleUpdateItemsCategories(ctx, category)
	}
	return nil
}

func (mock ServiceMock) GetByCategoryID(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError) {
	if mock.HandleGetByCategoryID != nil {
		return mock.HandleGetByCategoryID(ctx, categoryID)
	}
	return []models.Item{}, nil
}
