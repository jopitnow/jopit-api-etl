package categories

import (
	"context"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

type ServiceMock struct {
	HandleGet              func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError)
	HandleGetAllCategories func(ctx context.Context) ([]models.Category, apierrors.ApiError)
	HandleCreate           func(ctx context.Context, input models.Category) apierrors.ApiError
	HandleUpdate           func(ctx context.Context, input models.Category) apierrors.ApiError
	HandleDelete           func(ctx context.Context, items []models.Item, categoryID string) apierrors.ApiError
}

func NewServiceMock() ServiceMock {
	return ServiceMock{}
}

func (mock ServiceMock) Update(ctx context.Context, input models.Category) apierrors.ApiError {
	if mock.HandleUpdate != nil {
		return mock.HandleUpdate(ctx, input)
	}
	return nil
}

func (mock ServiceMock) Delete(ctx context.Context, items []models.Item, categoryID string) apierrors.ApiError {
	if mock.HandleDelete != nil {
		return mock.HandleDelete(ctx, items, categoryID)
	}
	return nil
}

func (mock ServiceMock) Get(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
	if mock.HandleGet != nil {
		return mock.HandleGet(ctx, categoryID)
	}
	return models.Category{}, nil
}

func (mock ServiceMock) Create(ctx context.Context, input models.Category) apierrors.ApiError {
	if mock.HandleCreate != nil {
		return mock.HandleCreate(ctx, input)
	}
	return nil
}

func (mock ServiceMock) GetAllCategories(ctx context.Context) ([]models.Category, apierrors.ApiError) {
	if mock.HandleGetAllCategories != nil {
		return mock.HandleGetAllCategories(ctx)
	}
	return []models.Category{}, nil
}
