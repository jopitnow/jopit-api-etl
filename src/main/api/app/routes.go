package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/jopit-api-etl/src/main/api/dependencies"
)

func RouterMapper(router *gin.Engine, h dependencies.HandlersStruct) {

	// Items
	router.POST("/etl/company-layout", goauth.AuthWithFirebase(), h.CompanyLayout.Create)
	router.GET("/etl/company-layout/:id", goauth.AuthWithFirebase(), h.CompanyLayout.Get)
	router.GET("/etl/company-layout", goauth.AuthWithFirebase(), h.CompanyLayout.GetByUserID)
	router.PUT("/etl/company-layout", goauth.AuthWithFirebase(), h.CompanyLayout.Update)
	router.DELETE("/etl/company-layout", goauth.AuthWithFirebase(), h.CompanyLayout.Delete)

	// MercadoLibre Credentials
	router.GET("/etl/mercadolibre/oauth", goauth.AuthWithFirebase(), h.MercadoLibreCredentials.GetOAuthURL)
	router.POST("/etl/mercadolibre/oauth", goauth.AuthWithFirebase(), h.MercadoLibreCredentials.CreateOAuthCredentials)
	//router.GET("/etl/mercadolibre/credentials", goauth.AuthWithFirebase(), h.MercadoLibreCredentials.GetCredentials)
	//router.DELETE("/etl/mercadolibre/credentials", goauth.AuthWithFirebase(), h.MercadoLibreCredentials.DeleteCredentials)

	// MercadoLibre ETL
	router.POST("/etl/mercadolibre/load", goauth.AuthWithFirebase(), h.Etl.LoadMercadoLibre)
	router.GET("/etl/mercadolibre/item/:item_id", goauth.AuthWithFirebase(), h.Etl.GetMercadoLibreItem)
	router.GET("/etl/mercadolibre/items", goauth.AuthWithFirebase(), h.Etl.GetMercadoLibreItems)
}
