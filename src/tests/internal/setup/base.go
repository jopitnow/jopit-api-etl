package setup

import (
	"github.com/gin-gonic/gin"
	"github.com/jopitnow/jopit-api-etl/src/main/api/app"
	"github.com/jopitnow/jopit-api-etl/src/main/api/dependencies"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// BuildRouter Helper function to create a router during testing.
func BuildRouter(depend dependencies.HandlersStruct) *gin.Engine {
	router := app.ConfigureRouter()
	router.Use(otelgin.Middleware("mock-otelgin"))
	mockRouteMapper(router, depend)
	return router
}

func mockRouteMapper(router *gin.Engine, h dependencies.HandlersStruct) {

	// Items
	router.GET("/items", mockAuthFirebase("01-USER-TEST"), h.Items.GetItemsByUserID)
	router.GET("/items/:id", h.Items.GetItemByID)
	router.GET("/items/shop/:id", h.Items.GetItemsByShopID)
	router.GET("/items/shop/:id/category/:category_id", h.Items.GetItemsByShopCategoryID)
	router.POST("/items/list", h.Items.GetItemsByIDs)
	router.POST("/items", mockAuthFirebase("01-USER-TEST"), h.Items.CreateItem)
	router.DELETE("/items/:id", h.Items.DeleteItem)
	router.PUT("/items/:id", mockAuthFirebase("01-USER-TEST"), h.Items.UpdateItem)

	//Categories
	router.PUT("/items/category", h.Categories.Update)
	router.DELETE("/items/category/:id_category", h.Categories.Delete)
	router.POST("/items/category", h.Categories.Create)
	router.GET("/items/category/:id_category", h.Categories.Get)
	router.GET("/items/categories", h.Categories.GetAllCategories)
}

func mockAuthFirebase(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
	}
}
