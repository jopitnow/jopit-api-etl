package clients

import (
	"context"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
)

type ShopClientMock struct {
	HandleGetShopByUserID func(ctx context.Context) (models.Shop, apierrors.ApiError)
}

func NewShopClientMock() ShopClientMock {
	return ShopClientMock{}
}

func (mock ShopClientMock) GetShopByUserID(ctx context.Context) (models.Shop, apierrors.ApiError) {
	if mock.HandleGetShopByUserID != nil {
		return mock.HandleGetShopByUserID(ctx)
	}
	return models.Shop{}, nil
}
