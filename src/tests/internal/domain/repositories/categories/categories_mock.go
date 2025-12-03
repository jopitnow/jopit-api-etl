package categories

import (
	"context"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

type RepositoryMock struct {
	HandleGet              func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError)
	HandleGetAllCategories func(ctx context.Context) ([]models.Category, apierrors.ApiError)
	HandleCreate           func(ctx context.Context, input models.Category) (interface{}, apierrors.ApiError)
	HandleUpdate           func(ctx context.Context, input models.Category) (int64, apierrors.ApiError)
	HandleDelete           func(ctx context.Context, categoryID string) (int64, apierrors.ApiError)
}

func NewRepositoryMock() RepositoryMock {
	return RepositoryMock{}
}

func (mock RepositoryMock) Update(ctx context.Context, input models.Category) (int64, apierrors.ApiError) {
	if mock.HandleUpdate != nil {
		return mock.HandleUpdate(ctx, input)
	}

	return -1, nil

}

func (mock RepositoryMock) Delete(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
	if mock.HandleDelete != nil {
		return mock.HandleDelete(ctx, categoryID)
	}

	return -1, nil

}

func (mock RepositoryMock) Get(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {

	if mock.HandleGet != nil {
		return mock.HandleGet(ctx, categoryID)
	}
	return models.Category{}, nil

}

func (mock RepositoryMock) Create(ctx context.Context, input models.Category) (interface{}, apierrors.ApiError) {

	if mock.HandleCreate != nil {
		return mock.HandleCreate(ctx, input)
	}
	return "", nil

}

func (mock RepositoryMock) GetAllCategories(ctx context.Context) ([]models.Category, apierrors.ApiError) {
	if mock.HandleGetAllCategories != nil {
		return mock.HandleGetAllCategories(ctx)
	}
	return []models.Category{}, nil
}
