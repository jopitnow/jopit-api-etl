package dto

import (
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
)

type BulkCreateItemsRequest struct {
	Items []models.Item `json:"items" binding:"required"`
}

type BulkUpsertResponse struct {
	TotalItems   int64 `json:"total_items"`
	CreatedCount int64 `json:"created_count"`
	UpdatedCount int64 `json:"updated_count"`
}
