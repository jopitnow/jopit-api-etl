package mocks

import (
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
)

var ItemDTO = dto.ItemDTO{
	Name:        "Example Item",
	Description: "This is an example item for testing purposes.",
	Status:      "active",
	Category: dto.CategoryDTO{
		Name: "Mock",
	},
	Dimensions: dto.DimensionsDTO{
		Weight: 500,
		Length: 30,
		Height: 30,
		Width:  30,
	},
	Price: dto.PriceDTO{
		Currency: dto.CurrencyDTO{
			ID: "USD",
		},
		Amount: 19.99,
	},
	Images: []dto.ImageDTO{
		"https://example.com/image1.jpg",
		"https://example.com/image2.jpg",
	},
	Attributes: dto.AttributesDTO{
		"Color": "Blue",
		"Size":  "M",
		"Brand": "Example Brand",
	},
	Eligible: []dto.EligibleDTO{
		{
			ID:         "1",
			Title:      "Eligible 1",
			Type:       "text",
			IsRequired: true,
			Options: []dto.OptionDTO{
				"Option 1",
				"Option 2",
				"Option 3",
			},
		},
		{
			ID:         "2",
			Title:      "Eligible 2",
			Type:       "color",
			IsRequired: false,
			Options: []dto.OptionDTO{
				"Option A",
				"Option B",
				"Option C",
			},
		},
	},
}
