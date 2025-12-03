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
}
