package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID          string       `json:"id" bson:"_id,omitempty"`
	ShopID      string       `json:"shop_id" bson:"shop_id,$set,omitempty"`
	UserID      string       `json:"user_id" bson:"user_id,$set,omitempty" binding:"required"`
	Name        string       `json:"name" binding:"required" bson:"name,$set,omitempty"`
	Description string       `json:"description" bson:"description,$set"`
	Status      string       `json:"status" bson:"status,omitempty" default:"active"`
	Category    ItemCategory `json:"category" binding:"required" bson:"category,$set,omitempty"`
	Delivery    Delivery     `json:"delivery" binding:"required" bson:"delivery,$set"`
	Attributes  Attributes   `json:"attributes" bson:"attributes,$set"`
	SizeGuide   *SizeGuide   `json:"size_guide,omitempty" bson:"size_guide,$set,omitempty"`
	Variants    []Variant    `json:"variants" bson:"variants,$set"`
	Price       Price        `json:"price" bson:"-"`
	Source      *Source      `json:"source,omitempty" bson:"source,$set,omitempty"`
}

type Source struct {
	SourceType        string            `json:"source_type,omitempty" bson:"source_type,$set,omitempty"`
	ExternalID        string            `json:"external_id,omitempty" bson:"external_id,$set,omitempty"`
	ExternalSKU       string            `json:"external_sku,omitempty" bson:"external_sku,$set,omitempty"`
	BatchID           string            `json:"batch_id,omitempty" bson:"batch_id,$set,omitempty"`
	ImportedAt        time.Time         `json:"imported_at,omitempty" bson:"imported_at,$set,omitempty"`
	SyncedAt          *time.Time        `json:"synced_at,omitempty" bson:"synced_at,$set,omitempty"`
	EtlVersion        string            `json:"etl_version,omitempty" bson:"etl_version,$set,omitempty"`
	TransformMetadata map[string]string `json:"transform_metadata,omitempty" bson:"transform_metadata,$set,omitempty"`
}

type Dimensions struct {
	Weight int `json:"weight" binding:"required" bson:"weight,$set"`
	Length int `json:"length" binding:"required" bson:"length,$set"`
	Height int `json:"height" binding:"required" bson:"height,$set"`
	Width  int `json:"width" binding:"required" bson:"width,$set"`
}

type Delivery struct {
	Fragile    bool       `json:"fragile" bson:"fragile"`
	Dimensions Dimensions `json:"dimensions" bson:"dimensions,$set"`
}

type MeliAttribute struct {
	ID        string `json:"id,omitempty" bson:"id,omitempty"`
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	ValueID   string `json:"value_id,omitempty" bson:"value_id,omitempty"`
	ValueName string `json:"value_name,omitempty" bson:"value_name,omitempty"`
}

type Attributes struct {
	Condition           string          `json:"condition,omitempty" bson:"condition,$set,omitempty"`
	Gender              string          `json:"gender,omitempty" bson:"gender,$set,omitempty"`
	Fit                 string          `json:"fit,omitempty" bson:"fit,$set,omitempty"`
	Composition         string          `json:"composition,omitempty" bson:"composition,$set,omitempty"`
	Elasticity          string          `json:"elasticity,omitempty" bson:"elasticity,$set,omitempty"`
	Season              string          `json:"season,omitempty" bson:"season,$set,omitempty"`
	WashingInstructions string          `json:"washing_instructions,omitempty" bson:"washing_instructions,$set,omitempty"`
	MeliAttributes      []MeliAttribute `json:"meli_attributes,omitempty" bson:"meli_attributes,$set,omitempty"`
}

type Subcategory struct {
	ID   string `json:"id" bson:"_id,$set,omitempty"`
	Name string `json:"name" bson:"name,$set,omitempty"`
}

type ItemCategory struct {
	ID          string       `json:"id" bson:"_id,$set,omitempty"`
	Name        string       `json:"name" binding:"required" bson:"name,$set,required"`
	Subcategory *Subcategory `json:"subcategory,omitempty" bson:"subcategory,$set,omitempty"`
}

type Size struct {
	SizeEquivalence      string `json:"size_equivalence" bson:"size_equivalence,$set"`
	ChestCircumference   int    `json:"chest_circumference,omitempty" bson:"chest_circumference,$set,omitempty"`
	WaistCircumference   int    `json:"waist_circumference,omitempty" bson:"waist_circumference,$set,omitempty"`
	Height               int    `json:"height,omitempty" bson:"height,$set,omitempty"`
	HipCircumference     int    `json:"hip_circumference,omitempty" bson:"hip_circumference,$set,omitempty"`
	FootLengthCm         int    `json:"foot_length_cm,omitempty" bson:"foot_length_cm,$set,omitempty"`
	SizeArgMen           int    `json:"size_arg_men,omitempty" bson:"size_arg_men,$set,omitempty"`
	SizeArgWomen         int    `json:"size_arg_women,omitempty" bson:"size_arg_women,$set,omitempty"`
	SizeUSA              int    `json:"size_usa,omitempty" bson:"size_usa,$set,omitempty"`
	SizeEUR              int    `json:"size_eur,omitempty" bson:"size_eur,$set,omitempty"`
	GarmentChestWidth    *int   `json:"garment_chest_width,omitempty" bson:"garment_chest_width,$set,omitempty"`
	GarmentLength        *int   `json:"garment_length,omitempty" bson:"garment_length,$set,omitempty"`
	GarmentShoulderWidth *int   `json:"garment_shoulder_width,omitempty" bson:"garment_shoulder_width,$set,omitempty"`
	GarmentHipWidth      *int   `json:"garment_hip_width,omitempty" bson:"garment_hip_width,$set,omitempty"`
	GarmentWaistWidth    *int   `json:"garment_waist_width,omitempty" bson:"garment_waist_width,$set,omitempty"`
}

type SizeGuide struct {
	Type               string `json:"type" bson:"type,$set"`
	BodyPart           string `json:"body_part" bson:"body_part,$set"`
	HasMeasurements    bool   `json:"has_measurements" bson:"has_measurements"`
	IsOneSize          bool   `json:"is_one_size" bson:"is_one_size"`
	MeasurementSource  string `json:"measurement_source,omitempty" bson:"measurement_source,$set,omitempty"`
	ExternalSizeGridID string `json:"external_size_grid_id,omitempty" bson:"external_size_grid_id,$set,omitempty"`
	Sizes              []Size `json:"sizes" bson:"sizes,$set"`
}

type SizeStock struct {
	SizeLabel string `json:"size_label" bson:"size_label,$set"`
	Stock     int    `json:"stock" bson:"stock,$set"`
	SKU       string `json:"sku,omitempty" bson:"sku,$set,omitempty"`
}

type Variant struct {
	ColorID   string      `json:"color_id" bson:"color_id,$set"`
	ColorName string      `json:"color_name" bson:"color_name,$set"`
	ColorHex  string      `json:"color_hex" bson:"color_hex,$set"`
	IsMain    bool        `json:"is_main" bson:"is_main"`
	Images    []Image     `json:"images" bson:"images,$set"`
	SizeStock []SizeStock `json:"size_stock" bson:"size_stock,$set"`
}

type Image string

type ItemsIds struct {
	Items []string `json:"items"`
}

func (i *Item) SetCategoryIDs() {
	i.Category.SetIDs()
}

func (c *ItemCategory) SetIDs() {
	if c.ID == "" {
		c.ID = primitive.NewObjectID().Hex()
	}
	if c.Subcategory != nil && c.Subcategory.ID == "" {
		c.Subcategory.ID = primitive.NewObjectID().Hex()
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
		if i.Items[index].Variants == nil {
			i.Items[index].Variants = []Variant{}
		}
		if i.Items[index].SizeGuide != nil && i.Items[index].SizeGuide.Sizes == nil {
			i.Items[index].SizeGuide.Sizes = []Size{}
		}
	}
}

func (i *Item) ValidateEmptySlices() {
	if i.Variants == nil {
		i.Variants = []Variant{}
	}
	if i.SizeGuide != nil && i.SizeGuide.Sizes == nil {
		i.SizeGuide.Sizes = []Size{}
	}
}
