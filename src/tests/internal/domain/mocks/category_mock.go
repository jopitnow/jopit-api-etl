package mocks

import (
	"encoding/json"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CategoryNameOne   = "Name One"
	CategoryNameTwo   = "Name Two"
	CategoryNameThree = "Name Three"
)

var Categories = []models.Category{CategoryOne, CategoryTwo, CategoryThree}

var CategoryOne = models.Category{
	ID:   primitive.NewObjectID().Hex(),
	Name: CategoryNameOne,
}

var CategoryTwo = models.Category{
	ID:   primitive.NewObjectID().Hex(),
	Name: CategoryNameTwo,
}

var CategoryThree = models.Category{
	ID:   primitive.NewObjectID().Hex(),
	Name: CategoryNameThree,
}

func CategoryToJson(item models.Category) string {
	bytes, _ := json.Marshal(item)

	return string(bytes)
}
