package clients

import (
	"context"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

type PriceClientMock struct {
	HandleGetPriceByItemID func(ctx context.Context, itemID string) (models.Price, apierrors.ApiError)
	HandleCreatePrice      func(ctx context.Context, price *models.Price) apierrors.ApiError
	HandleModifyPrice      func(ctx context.Context, price *models.Price) apierrors.ApiError
	HandleGetItemsPrices   func(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError)
	HandleDeletePrice      func(ctx context.Context, priceID string) apierrors.ApiError
}

func NewPriceClientMock() PriceClientMock {
	return PriceClientMock{}
}

func (mock PriceClientMock) GetPriceByItemID(ctx context.Context, itemID string) (models.Price, apierrors.ApiError) {
	if mock.HandleGetPriceByItemID != nil {
		return mock.HandleGetPriceByItemID(ctx, itemID)
	}
	return models.Price{}, nil
}

func (mock PriceClientMock) CreatePrice(ctx context.Context, price *models.Price) apierrors.ApiError {
	if mock.HandleCreatePrice != nil {
		return mock.HandleCreatePrice(ctx, price)
	}
	return nil
}

func (mock PriceClientMock) UpdatePrice(ctx context.Context, price *models.Price) apierrors.ApiError {
	if mock.HandleModifyPrice != nil {
		return mock.HandleModifyPrice(ctx, price)
	}
	return nil
}

func (mock PriceClientMock) GetItemsPrices(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError) {
	if mock.HandleGetItemsPrices != nil {
		return mock.HandleGetItemsPrices(ctx, itemsIDs)
	}
	return models.Prices{}, nil
}

func (mock PriceClientMock) DeletePrice(ctx context.Context, priceID string) apierrors.ApiError {
	if mock.HandleDeletePrice != nil {
		return mock.HandleDeletePrice(ctx, priceID)
	}
	return nil
}
