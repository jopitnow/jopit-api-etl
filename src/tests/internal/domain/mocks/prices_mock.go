package mocks

import (
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PriceId                  = "676695ae1d8a4720e8c60de8"
	Amount                   = 20.0
	CurrencyId               = "ARS"
	CurrencySymbol           = "$"
	CurrencyDecimalDivider   = ","
	CurrencyThousandsDivider = "."
)

var Prices = models.Prices{Prices: []models.Price{Price, Price2}}

var Price = models.Price{
	ID:     PriceId,
	ItemID: ItemIdOne,
	Amount: Amount,
	Currency: models.Currency{
		ID:               CurrencyId,
		Symbol:           CurrencySymbol,
		DecimalDivider:   CurrencyDecimalDivider,
		ThousandsDivider: CurrencyThousandsDivider,
	},
}

var Price2 = models.Price{
	ID:     primitive.NewObjectID().Hex(),
	ItemID: ItemIdTwo,
	Amount: Amount,
	Currency: models.Currency{
		ID:               CurrencyId,
		Symbol:           CurrencySymbol,
		DecimalDivider:   CurrencyDecimalDivider,
		ThousandsDivider: CurrencyThousandsDivider,
	},
}

var MockItemsIDs = []string{ItemIdOne, ItemIdTwo}
