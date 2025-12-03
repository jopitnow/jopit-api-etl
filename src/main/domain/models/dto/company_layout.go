package dto

import "github.com/jopitnow/jopit-api-etl/src/main/domain/models"

type CompanyLayoutRequest struct {
	Name        string               `bson:"name" json:"name"`
	Request     models.RequestConfig `bson:"request" json:"request"`
	ItemMap     map[string]string    `bson:"item_map" json:"field_map"`
	CategoryMap map[string]string    `bson:"category_map" json:"category_map"`
}

func (c *CompanyLayoutRequest) ToModel() models.CompanyLayout {
	return models.CompanyLayout{
		Name:        c.Name,
		Request:     c.Request,
		ItemMap:     c.ItemMap,
		CategoryMap: c.CategoryMap,
	}

}
