package dto

type PriceDTO struct {
	ID       string      `json:"id"`
	Amount   float64     `json:"amount" binding:"required"`
	Currency CurrencyDTO `json:"currency" binding:"required"`
}

type CurrencyDTO struct {
	ID               string `json:"id"`
	Symbol           string `json:"symbol" binding:"required"`
	DecimalDivider   string `json:"decimal_divider" binding:"required"`
	ThousandsDivider string `json:"thousands_divider" binding:"required"`
}

type RequestItemsList struct {
	ItemsIDs []string `json:"items_ids" binding:"required"`
}
