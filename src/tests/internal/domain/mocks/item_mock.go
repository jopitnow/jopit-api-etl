package mocks

import (
	"encoding/json"
	"log"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
)

const (
	ItemIdOne     = "676695ae1d8a4720e8c60dea"
	ItemIdTwo     = "676695ae1d8a4720e8c60deb"
	ShopIDOne     = "6766992ece7f50da3bd47c25"
	ShopIDTwo     = "6766992ece7f50da3bd47d12"
	UserIdOne     = "6766992ece7f50da3bd47c26"
	UserIdTwo     = "6766992ece7f50da3bd47c27"
	CategoryIDOne = "65988dc8e37778a9afb2eb47"
	EligibleIdOne = "65988dc8e37778a9afb2cb47"
)

var ItemMockOne = models.Item{
	ID:          ItemIdOne,
	ShopID:      ShopIDOne,
	UserID:      UserIdOne,
	Name:        "Example Item",
	Description: "This is an example item for testing purposes.",
	Fragile:     true,
	Dimensions: models.Dimensions{
		Weight: 100,
		Width:  100,
		Height: 100,
		Length: 100,
	},
	Status: "active",
	Category: models.Category{
		ID:   CategoryIDOne,
		Name: "Mock",
	},
	Price: models.Price{
		Currency: models.Currency{
			ID:               "ARS",
			Symbol:           "$",
			DecimalDivider:   ",",
			ThousandsDivider: ".",
		},
		Amount: 19.99,
	},
	Images: []models.Image{
		"https://example.com/image1.jpg",
		"https://example.com/image2.jpg",
	},
	Attributes: models.Attributes{
		"Color": "Blue",
		"Size":  "M",
		"Brand": "Example Brand",
	},
	Eligible: []models.Eligible{
		{
			ID:         EligibleIdOne,
			Title:      "Eligible 2",
			Type:       "color",
			IsRequired: false,
			Options: []models.Option{
				"Option A",
				"Option B",
				"Option C",
			},
		},
	},
}

var ItemMockTwo = models.Item{
	ID:          ItemIdTwo,
	UserID:      UserIdTwo,
	Name:        "Example Item",
	Description: "This is an example item for testing purposes.",
	Status:      "active",
	Category: models.Category{
		ID:   CategoryIDOne,
		Name: "Mock",
	},
	Price: models.Price{
		Currency: models.Currency{
			ID:               "ARS",
			Symbol:           "$",
			DecimalDivider:   ",",
			ThousandsDivider: ".",
		},
		Amount: 19.99,
	},
	Images: []models.Image{
		"https://example.com/image1.jpg",
		"https://example.com/image2.jpg",
	},
	Attributes: models.Attributes{
		"Color": "Blue",
		"Size":  "M",
		"Brand": "Example Brand",
	},
	Eligible: []models.Eligible{
		{
			ID:         EligibleIdOne,
			Title:      "Eligible 1",
			Type:       "text",
			IsRequired: true,
			Options: []models.Option{
				"Option 1",
				"Option 2",
				"Option 3",
			},
		},
	},
}

var ItemsMock = models.Items{Items: []models.Item{
	ItemMockOne,
	ItemMockTwo,
}}

var ItemIds = models.ItemsIds{
	Items: []string{ItemIdOne, ItemIdTwo},
}

func ItemToJson() string {
	bytes, err := json.Marshal(ItemMockOne)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

func ItemsIdsToJson() string {
	bytes, err := json.Marshal(ItemIds)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}
