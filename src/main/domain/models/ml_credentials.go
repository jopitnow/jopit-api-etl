package models

import "time"

const MeliType = "MELI"

type MercadoLibreURL struct {
	URL string `json:"url,omitempty"`
}

type MercadoLibreCredential struct {
	ID           string    `json:"id,omitempty" bson:"_id,omitempty"`
	ShopID       string    `json:"shop_id" binding:"required" bson:"shop_id,$set"`
	UserID       string    `json:"user_id" binding:"required" bson:"user_id,$set"`
	AccessToken  string    `json:"access_token" binding:"required" bson:"access_token,$set"`
	RefreshToken string    `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty" bson:"token_type,omitempty"`
	ExpiresIn    int       `json:"expires_in" bson:"expires_in,omitempty"`
	Scope        string    `json:"scope,omitempty" bson:"scope,omitempty"`
	UserIDMeli   int64     `json:"user_id_meli,omitempty" bson:"user_id_meli,omitempty"` // MercadoLibre's user ID
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at,omitempty"`
}
