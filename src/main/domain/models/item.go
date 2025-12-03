package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	ShopID      string     `json:"shop_id" bson:"shop_id,$set,omitempty"`
	UserID      string     `json:"user_id" bson:"user_id,$set,omitempty" binding:"required"`
	Name        string     `json:"name" binding:"required" bson:"name,$set,omitempty"`
	Description string     `json:"description" bson:"description,$set"`
	Status      string     `json:"status" bson:"status,omitempty" default:"active"`
	BatchID     string     `json:"batch_id" bson:"batch_id"`
	Fragile     bool       `json:"fragile" binding:"required" bson:"fragile"`
	Category    Category   `json:"category" binding:"required" bson:"category,$set,omitempty"`
	Dimensions  Dimensions `json:"dimensions" binding:"required" bson:"dimensions,$set"`
	Images      []Image    `json:"images" bson:"images,$set"`
	Attributes  Attributes `json:"attributes" bson:"attributes,$set"`
	Eligible    []Eligible `json:"eligible,omitempty" bson:"eligible,$set"`
	Price       Price      `json:"price" bson:"-"`
}

type Dimensions struct {
	Weight int `json:"weight" binding:"required"  bson:"weight,$set"`
	Length int `json:"length" binding:"required"  bson:"length,$set"`
	Height int `json:"height" binding:"required"  bson:"height,$set"`
	Width  int `json:"width" binding:"required"  bson:"width,$set"`
}

type Category struct {
	ID   string `json:"id" bson:"_id,$set,omitempty"`
	Name string `json:"name" binding:"required" bson:"name,$set,required"`
}

type Eligible struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Type       string   `json:"type"`
	IsRequired bool     `json:"is_required"`
	Options    []Option `json:"options"`
}

type Attributes map[string]string
type Option string
type Image string

type ItemsIds struct {
	Items []string `json:"items"`
}

func (i *Item) SetEligibleIDs() {
	for e := range i.Eligible {
		i.Eligible[e].ID = primitive.NewObjectID().Hex()
	}
}

func (i *Item) UpdateEligibleIDs() {
	for e := range i.Eligible {
		_, err := primitive.ObjectIDFromHex(i.Eligible[e].ID)

		if err != nil {
			i.Eligible[e].ID = primitive.NewObjectID().Hex()
		}
	}
}

func (i *Items) SetPriceToItems(response Prices) Items {

	iis := Items{}

	for _, pr := range response.Prices {

		for _, item := range i.Items {

			if pr.ItemID == item.ID {
				item.Price = pr
				iis.Items = append(iis.Items, item)
			}
		}
	}

	return iis
}

func (i *Items) GetItemsIds() []string {
	var ids []string

	for _, item := range i.Items {
		ids = append(ids, item.ID)
	}

	return ids
}

func (i *Items) ValidateEmptySlices() {
	for index := range i.Items {
		if i.Items[index].Eligible == nil {
			i.Items[index].Eligible = []Eligible{}
		}
		if i.Items[index].Attributes == nil {
			i.Items[index].Attributes = map[string]string{}
		}
		if i.Items[index].Images == nil {
			i.Items[index].Images = []Image{}
		}
	}
}

func (i *Item) ValidateEmptySlices() {
	if i.Eligible == nil {
		i.Eligible = []Eligible{}
	}
	if i.Attributes == nil {
		i.Attributes = map[string]string{}
	}
	if i.Images == nil {
		i.Images = []Image{}
	}
}
