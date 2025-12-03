package dto

import (
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
)

type BulkCreateItemsRequest struct {
	Items []models.Item `json:"items" binding:"required"`
}
