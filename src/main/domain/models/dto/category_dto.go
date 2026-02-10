package dto

import "github.com/jopitnow/jopit-api-etl/src/main/domain/models"

type CategoryDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name" binding:"required"`
}

type CategoriesDTO struct {
	CategoryDTO []models.ItemCategory `json:"categories"`
}
