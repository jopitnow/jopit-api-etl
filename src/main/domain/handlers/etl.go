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
	Service services.EtlService
}

func NewEtlsHandler(service services.EtlService) EtlHandler {
	return EtlHandler{
		Service: service,
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
