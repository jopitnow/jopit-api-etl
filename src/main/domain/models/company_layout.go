package models

// --- Structured request definition ---
type RequestConfig struct {
	Method      string            `bson:"method" json:"method"`             // e.g. "GET", "POST"
	Endpoint    string            `bson:"endpoint" json:"endpoint"`         // base URL
	Headers     map[string]string `bson:"headers" json:"headers"`           // custom headers
	QueryParams map[string]string `bson:"query_params" json:"query_params"` // query string params
	Body        map[string]string `bson:"body" json:"body,omitempty"`       // optional body for POST/PUT
}

type ResponseConfig struct {
	FieldMap map[string]string `bson:"field_map" json:"field_map"`
}

// --- Company config with flexible request definition ---
type CompanyLayout struct {
	ID          string            `json:"id" bson:"_id,omitempty"`
	UserID      string            `bson:"user_id" json:"user_id"`
	ShopID      string            `bson:"shop_id" json:"shop_id"`
	Name        string            `bson:"name" json:"name"`
	Request     RequestConfig     `bson:"request" json:"request"`
	ItemMap     map[string]string `bson:"item_map" json:"field_map"`
	CategoryMap map[string]string `bson:"category_map" json:"category_map"`
}

// ++++++++++++++++++

// Default internal mapping for Item schema
var InternalItemMap = map[string]string{
	"id":          "id",
	"shop_id":     "shop_id",
	"user_id":     "user_id",
	"name":        "name",
	"description": "description",
	"status":      "status",
	"fragile":     "fragile",

	"category.id":   "category.id",
	"category.name": "category.name",

	"dimensions.weight": "dimensions.weight",
	"dimensions.length": "dimensions.length",
	"dimensions.height": "dimensions.height",
	"dimensions.width":  "dimensions.width",

	"images":     "images",
	"attributes": "attributes",

	"eligible[].id":          "eligible.id",
	"eligible[].title":       "eligible.title",
	"eligible[].type":        "eligible.type",
	"eligible[].is_required": "eligible.is_required",
	"eligible[].options":     "eligible.options",

	"price.amount":   "price.amount",
	"price.currency": "price.currency",
}

func MergeMaps(defaults, overrides map[string]string) map[string]string {
	merged := make(map[string]string)

	// start with defaults
	for k, v := range defaults {
		merged[k] = v
	}

	// apply overrides
	for k, v := range overrides {
		merged[k] = v
	}

	return merged
}
