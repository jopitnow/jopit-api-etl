package models

type Price struct {
	ID       string   `json:"id"`
	ItemID   string   `json:"item_id"`
	ShopID   string   `json:"shop_id"`
	Amount   float64  `json:"amount"`
	Currency Currency `json:"currency"`
}

type Currency struct {
	ID               string `json:"id"`
	Symbol           string `json:"symbol"`
	DecimalDivider   string `json:"decimal_divider"`
	ThousandsDivider string `json:"thousands_divider"`
}

type Prices struct {
	Prices []Price `json:"prices"`
}
