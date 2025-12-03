package dto

import (
	"net/http"
	"testing"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"github.com/stretchr/testify/assert"
)

func TestValidateDimensions_Valid(t *testing.T) {
	item := dto.ItemDTO{
		Dimensions: dto.DimensionsDTO{
			Height: 100,
			Length: 200,
			Weight: 100,
			Width:  300,
		},
	}

	err := item.ValidateDimensions()
	assert.Nil(t, err)
}

func TestValidateDimensions_InvalidHeight(t *testing.T) {
	item := dto.ItemDTO{
		Dimensions: dto.DimensionsDTO{Height: -1},
	}

	err := item.ValidateDimensions()
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Status())
}

func TestValidateDimensions_InvalidLength(t *testing.T) {
	item := dto.ItemDTO{
		Dimensions: dto.DimensionsDTO{Length: 1600},
	}

	err := item.ValidateDimensions()
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Status())
}

func TestValidateDimensions_InvalidWeight(t *testing.T) {
	item := dto.ItemDTO{
		Dimensions: dto.DimensionsDTO{Weight: 2000},
	}

	err := item.ValidateDimensions()
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Status())
}

func TestValidateDimensions_InvalidWidth(t *testing.T) {
	item := dto.ItemDTO{
		Dimensions: dto.DimensionsDTO{Width: -5},
	}

	err := item.ValidateDimensions()
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Status())
}
