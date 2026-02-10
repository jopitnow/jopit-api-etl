package dto

// MercadoLibreAuthRedirectDTO is the request body when MercadoLibre redirects with auth code
type MercadoLibreAuthRedirectDTO struct {
	Code string `json:"code" binding:"required"`
}

// MercadoLibreAuthRequestDTO is the request body for OAuth token exchange
type MercadoLibreAuthRequestDTO struct {
	ClientID     string  `json:"client_id"`
	ClientSecret string  `json:"client_secret"`
	GrantType    string  `json:"grant_type"`
	Code         *string `json:"code,omitempty"`
	CodeVerifier *string `json:"code_verifier,omitempty"`
	RedirectURI  *string `json:"redirect_uri,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}

// MercadoLibreAuthResponse represents OAuth response from MercadoLibre
type MercadoLibreAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UserID       int64  `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

// MercadoLibreUserResponse represents user info from MercadoLibre API
type MercadoLibreUserResponse struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	SiteID    string `json:"site_id"`
}

// MercadoLibreSearchFilters contains params for the search endpoint
type MercadoLibreSearchFilters struct {
	SiteID   string            `json:"site_id"`
	SellerID string            `json:"seller_id,omitempty"`
	Query    string            `json:"query,omitempty"`
	Params   map[string]string `json:"params,omitempty"`
}

// ============== Item DTOs ==============

// MeliPicture represents an item picture
type MeliPicture struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Size      string `json:"size"`
	MaxSize   string `json:"max_size"`
	Quality   string `json:"quality"`
}

// MeliAttributeValue represents a value within an attribute
type MeliAttributeValue struct {
	ID     *string     `json:"id"`
	Name   string      `json:"name"`
	Struct interface{} `json:"struct"`
}

// MeliAttribute represents item attributes
type MeliAttribute struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	ValueID     *string              `json:"value_id"`
	ValueName   string               `json:"value_name"`
	ValueStruct interface{}          `json:"value_struct"`
	Values      []MeliAttributeValue `json:"values"`
	ValueType   string               `json:"value_type"`
}

// MeliShipping represents shipping information
type MeliShipping struct {
	Mode         string      `json:"mode"`
	Methods      []string    `json:"methods"`
	Tags         []string    `json:"tags"`
	Dimensions   interface{} `json:"dimensions"`
	LocalPickUp  bool        `json:"local_pick_up"`
	FreeShipping bool        `json:"free_shipping"`
	LogisticType string      `json:"logistic_type"`
	StorePickUp  bool        `json:"store_pick_up"`
}

// MeliVariation represents an item variation
type MeliVariation struct {
	ID                    int64           `json:"id"`
	Price                 float64         `json:"price"`
	AttributeCombinations []MeliAttribute `json:"attribute_combinations"`
	AvailableQuantity     int             `json:"available_quantity"`
	SoldQuantity          int             `json:"sold_quantity"`
	SaleTerms             []interface{}   `json:"sale_terms"`
	PictureIDs            []string        `json:"picture_ids"`
	SellerCustomField     *string         `json:"seller_custom_field"`
	CatalogProductID      *string         `json:"catalog_product_id"`
	InventoryID           *string         `json:"inventory_id"`
	ItemRelations         []interface{}   `json:"item_relations"`
	UserProductID         string          `json:"user_product_id"`
}

// MeliItemResponse represents a full item from MercadoLibre
type MeliItemResponse struct {
	ID                 string          `json:"id"`
	SiteID             string          `json:"site_id"`
	Title              string          `json:"title"`
	SellerID           int64           `json:"seller_id"`
	CategoryID         string          `json:"category_id"`
	Price              float64         `json:"price"`
	BasePrice          float64         `json:"base_price"`
	OriginalPrice      *float64        `json:"original_price"`
	CurrencyID         string          `json:"currency_id"`
	InitialQuantity    int             `json:"initial_quantity"`
	AvailableQuantity  int             `json:"available_quantity"`
	SoldQuantity       int             `json:"sold_quantity"`
	SaleTerms          []interface{}   `json:"sale_terms"`
	BuyingMode         string          `json:"buying_mode"`
	ListingTypeID      string          `json:"listing_type_id"`
	StartTime          string          `json:"start_time"`
	StopTime           string          `json:"stop_time"`
	EndTime            string          `json:"end_time"`
	ExpirationTime     string          `json:"expiration_time"`
	Condition          string          `json:"condition"`
	Permalink          string          `json:"permalink"`
	ThumbnailID        string          `json:"thumbnail_id"`
	Thumbnail          string          `json:"thumbnail"`
	Pictures           []MeliPicture   `json:"pictures"`
	VideoID            *string         `json:"video_id"`
	Descriptions       []interface{}   `json:"descriptions"`
	AcceptsMercadoPago bool            `json:"accepts_mercadopago"`
	Shipping           MeliShipping    `json:"shipping"`
	Attributes         []MeliAttribute `json:"attributes"`
	Variations         []MeliVariation `json:"variations"`
	Status             string          `json:"status"`
	SubStatus          []string        `json:"sub_status"`
	Tags               []string        `json:"tags"`
	Warranty           string          `json:"warranty"`
	CatalogProductID   *string         `json:"catalog_product_id"`
	DomainID           string          `json:"domain_id"`
	DealIDs            []string        `json:"deal_ids"`
	AutomaticRelist    bool            `json:"automatic_relist"`
	DateCreated        string          `json:"date_created"`
	LastUpdated        string          `json:"last_updated"`
	CatalogListing     bool            `json:"catalog_listing"`
	Channels           []string        `json:"channels"`
}

// MeliUserItemsSearchResponse represents the search response with item IDs
type MeliUserItemsSearchResponse struct {
	Results []string `json:"results"`
	Paging  struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"paging"`
}

// MeliSearchResponse represents the public search response
type MeliSearchResponse struct {
	SiteID string `json:"site_id"`
	Query  string `json:"query"`
	Paging struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"paging"`
	Results []MeliItemResponse `json:"results"`
}
