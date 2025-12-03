package categories

import (
	"context"
	"net/http"
	"testing"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/services"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/mocks"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/repositories/categories"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/stretchr/testify/assert"
)

func TestService_Get_Success(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return mocks.CategoryOne, nil
	}

	service := services.NewCategoriesService(repository)

	category, apiErr := service.Get(context.TODO(), "1")

	assert.Nil(t, apiErr)
	assert.Equal(t, mocks.CategoryOne, category)
}

func TestService_Get_Not_Found_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return models.Category{}, apierrors.NewApiError("mock error", "not_found", http.StatusNotFound, apierrors.CauseList{})
	}

	service := services.NewCategoriesService(repository)

	category, apiErr := service.Get(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, models.Category{}, category)
}

func TestService_Get_Repository_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGet = func(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
		return models.Category{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewCategoriesService(repository)

	category, apiErr := service.Get(context.TODO(), "1")

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status())
	assert.Equal(t, models.Category{}, category)
}

func TestService_GetAll_Success(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return mocks.Categories, nil
	}

	service := services.NewCategoriesService(repository)

	category, apiErr := service.GetAllCategories(context.TODO())

	assert.Nil(t, apiErr)
	assert.Equal(t, mocks.Categories, category)
}

func TestService_GetAll_Not_Found_Error_Internal_Server_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return []models.Category{}, nil
	}

	service := services.NewCategoriesService(repository)

	_, apiErr := service.GetAllCategories(context.TODO())

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status())
}

func TestService_GetAll_Repository_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return []models.Category{}, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewCategoriesService(repository)

	response, apiErr := service.GetAllCategories(context.TODO())

	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status())
	assert.Equal(t, []models.Category(nil), response)
}

func TestService_Create_Success(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleCreate = func(ctx context.Context, input models.Category) (interface{}, apierrors.ApiError) {
		return "1", nil
	}

	service := services.NewCategoriesService(repository)

	err := service.Create(context.TODO(), mocks.CategoryOne)

	assert.Nil(t, err)
}

func TestService_Create_Get_Categories_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return nil, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewCategoriesService(repository)

	err := service.Create(context.TODO(), mocks.CategoryOne)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Create_Already_Exist_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return mocks.Categories, nil
	}

	service := services.NewCategoriesService(repository)

	err := service.Create(context.TODO(), mocks.CategoryOne)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Status())
}

func TestService_Create_Repository_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleCreate = func(ctx context.Context, input models.Category) (interface{}, apierrors.ApiError) {
		return "", apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewCategoriesService(repository)

	err := service.Create(context.TODO(), mocks.CategoryOne)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Update_Success(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleUpdate = func(ctx context.Context, input models.Category) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewCategoriesService(repository)

	err := service.Update(context.TODO(), mocks.CategoryOne)

	assert.Nil(t, err)
}

func TestService_Update_Get_Categories_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return nil, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewCategoriesService(repository)

	err := service.Update(context.TODO(), mocks.CategoryOne)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Update_Already_Exist_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleGetAllCategories = func(ctx context.Context) ([]models.Category, apierrors.ApiError) {
		return mocks.Categories, nil
	}

	service := services.NewCategoriesService(repository)

	err := service.Update(context.TODO(), mocks.CategoryOne)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Status())
}

func TestService_Update_Repository_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleUpdate = func(ctx context.Context, input models.Category) (int64, apierrors.ApiError) {
		return -1, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewCategoriesService(repository)

	err := service.Update(context.TODO(), mocks.CategoryOne)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestService_Delete_Success(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleDelete = func(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewCategoriesService(repository)

	err := service.Delete(context.TODO(), []models.Item{}, "1")

	assert.Nil(t, err)
}

func TestService_Delete_Conflic_Deletion_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleDelete = func(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
		return 1, nil
	}

	service := services.NewCategoriesService(repository)

	err := service.Delete(context.TODO(), mocks.ItemsMock.Items, "1")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Status())
}

func TestService_Delete_Repository_Error(t *testing.T) {
	repository := categories.NewRepositoryMock()
	repository.HandleDelete = func(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
		return -1, apierrors.NewApiError("mock error", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	service := services.NewCategoriesService(repository)

	err := service.Delete(context.TODO(), []models.Item{}, "1")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}
