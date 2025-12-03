package items

import (
	"context"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

type RepositoryMock struct {
	HandleGet                 func(ctx context.Context, itemID string) (models.Item, apierrors.ApiError)
	HandleGetByUserID         func(ctx context.Context, userID string) (models.Items, apierrors.ApiError)
	HandleGetByShopID         func(ctx context.Context, shopID string) (models.Items, apierrors.ApiError)
	HandleGetByShopCategoryID func(ctx context.Context, shopID string, itemID string) (models.Items, apierrors.ApiError)
	HandleGetByIDs            func(ctx context.Context, itemsIDs []string) (models.Items, apierrors.ApiError)
	HandleSave                func(ctx context.Context, item models.Item) (interface{}, apierrors.ApiError)
	HandleUpdate              func(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError)
	HandleDelete              func(ctx context.Context, itemID string) (int64, apierrors.ApiError)

	HandleUpdateItemsCategories func(ctx context.Context, category *models.Category) apierrors.ApiError
	HandleGetByCategoryID       func(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError)
}

func NewItemsRepositoryMock() RepositoryMock {
	return RepositoryMock{}
}

func (mock RepositoryMock) Get(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
	if mock.HandleGet != nil {
		return mock.HandleGet(ctx, itemID)
	}
	return models.Item{}, nil
}

func (mock RepositoryMock) GetByUserID(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
	if mock.HandleGetByUserID != nil {
		return mock.HandleGetByUserID(ctx, userID)
	}
	return models.Items{}, nil
}

func (mock RepositoryMock) GetByShopID(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
	if mock.HandleGetByShopID != nil {
		return mock.HandleGetByShopID(ctx, shopID)
	}
	return models.Items{}, nil
}

func (mock RepositoryMock) GetByShopCategoryID(ctx context.Context, shopID string, itemID string) (models.Items, apierrors.ApiError) {
	if mock.HandleGetByShopCategoryID != nil {
		return mock.HandleGetByShopCategoryID(ctx, shopID, itemID)
	}
	return models.Items{}, nil
}

func (mock RepositoryMock) GetByIDs(ctx context.Context, itemsIDs []string) (models.Items, apierrors.ApiError) {
	if mock.HandleGetByIDs != nil {
		return mock.HandleGetByIDs(ctx, itemsIDs)
	}
	return models.Items{}, nil
}

func (mock RepositoryMock) Save(ctx context.Context, item models.Item) (interface{}, apierrors.ApiError) {
	if mock.HandleSave != nil {
		return mock.HandleSave(ctx, item)
	}
	return nil, nil
}

func (mock RepositoryMock) Update(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError) {
	if mock.HandleUpdate != nil {
		return mock.HandleUpdate(ctx, itemID, updateItem)
	}
	return -1, nil
}

func (mock RepositoryMock) Delete(ctx context.Context, itemID string) (int64, apierrors.ApiError) {
	if mock.HandleDelete != nil {
		return mock.HandleDelete(ctx, itemID)
	}
	return -1, nil
}

func (mock RepositoryMock) UpdateItemsCategories(ctx context.Context, category *models.Category) apierrors.ApiError {
	if mock.HandleUpdateItemsCategories != nil {
		return mock.HandleUpdateItemsCategories(ctx, category)
	}
	return nil
}

func (mock RepositoryMock) GetByCategoryID(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError) {
	if mock.HandleGetByCategoryID != nil {
		return mock.HandleGetByCategoryID(ctx, categoryID)
	}
	return []models.Item{}, nil
}
