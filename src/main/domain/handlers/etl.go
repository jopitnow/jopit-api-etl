package handlers

import (
	"context"
	"net/http"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/services"

	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"

	"github.com/gin-gonic/gin"
)

type EtlHandler struct {
	Service             services.EtlService
	MercadoLibreService services.MercadoLibreService
}

func NewEtlsHandler(service services.EtlService, mercadoLibreService services.MercadoLibreService) EtlHandler {
	return EtlHandler{
		Service:             service,
		MercadoLibreService: mercadoLibreService,
	}
}

// CreateItem godoc
// @Summary Create Item
// @Description Create item in db
// @Tags Items
// @Accept  json
// @Produce  json
// @Param item body dto.ItemDTO true "Add item"
// @Success 201
// @Router /items [post]
func (h EtlHandler) LoadApi(c *gin.Context) {

	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)

	response, apiErr := h.Service.LoadApi(ctx)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// CreateItem godoc
// @Summary Create Item
// @Description Create item in db
// @Tags Items
// @Accept  json
// @Produce  json
// @Param item body dto.ItemDTO true "Add item"
// @Success 201
// @Router /items [post]
func (h EtlHandler) LoadCsv(c *gin.Context) {

	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file required"})
		return
	}

	response, apiErr := h.Service.LoadCsv(ctx, fileHeader)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// DeleteItem godoc
// @Summary Delete item
// @Description Delete item by ID
// @Tags Items
// @Accept  json
// @Produce  json
// @Param id path string true "Item ID"
// @Success 204
// @Router /items/{id} [delete]
func (h EtlHandler) Delete(c *gin.Context) {

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	batchID := c.Param("id")
	if batchID == "" {
		err := apierrors.NewApiError("empty required BatchID", "bad_request", http.StatusBadRequest, apierrors.CauseList{})
		c.Error(err)
		c.JSON(err.Status(), err)

	}

	err := h.Service.DeleteBatch(ctx, batchID)
	if err != nil {
		c.Error(err)
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusNoContent)
}

// LoadMercadoLibre godoc
// @Summary Load items from MercadoLibre
// @Description Fetch item IDs from MercadoLibre search, then batch fetch full item details
// @Tags ETL
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /etl/mercadolibre/load [post]
func (h EtlHandler) LoadMercadoLibre(c *gin.Context) {
	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)
	ctx = context.WithValue(ctx, goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	// Step 1: Search for user items to get the list of item IDs
	searchResult, apiErr := h.MercadoLibreService.GetUserItems(ctx)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	if len(searchResult.Results) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":    "No items found for this seller",
			"item_count": 0,
			"items":      []interface{}{},
		})
		return
	}

	// Step 2: Batch fetch full details for all items using MercadoLibre's multi-get endpoint
	itemsDetails, apiErr := h.MercadoLibreService.GetItems(ctx, searchResult.Results)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "MercadoLibre items retrieved and processed successfully",
		"item_count": len(itemsDetails),
		"items":      itemsDetails,
	})
}

// GetMercadoLibreItem godoc
// @Summary Get single item from MercadoLibre
// @Description Test fetching a single item detail from MercadoLibre API
// @Tags ETL
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param item_id path string true "MercadoLibre Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /etl/mercadolibre/item/{item_id} [get]
func (h EtlHandler) GetMercadoLibreItem(c *gin.Context) {
	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)
	ctx = context.WithValue(ctx, goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	itemID := c.Param("item_id")
	if itemID == "" {
		err := apierrors.NewApiError("item_id is required", "bad_request", http.StatusBadRequest, apierrors.CauseList{})
		c.Error(err)
		c.JSON(err.Status(), err)
		return
	}

	// Get item detail
	item, apiErr := h.MercadoLibreService.GetItem(ctx, itemID)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "MercadoLibre item retrieved successfully",
		"data":    item,
	})
}

// GetMercadoLibreItems godoc
// @Summary Get user items from MercadoLibre
// @Description Fetch all item IDs for the authenticated seller
// @Tags ETL
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /etl/mercadolibre/items [get]
func (h EtlHandler) GetMercadoLibreItems(c *gin.Context) {
	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)
	ctx = context.WithValue(ctx, goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	// Get user items (search endpoint) - returns just the IDs
	result, apiErr := h.MercadoLibreService.GetUserItems(ctx)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "MercadoLibre user items retrieved successfully",
		"data":    result,
	})
}
