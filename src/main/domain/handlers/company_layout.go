package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/services"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/utils"
)

type CompanyLayoutHandler struct {
	Service services.CompanyLayoutService
}

func NewCompanyLayoutHandler(service services.CompanyLayoutService) CompanyLayoutHandler {
	return CompanyLayoutHandler{
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
func (h CompanyLayoutHandler) Create(c *gin.Context) {

	var input dto.CompanyLayoutRequest

	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		apiErr := apierrors.NewApiError(err.Error(), "bad_request", http.StatusBadRequest, apierrors.CauseList{})
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)
	ctx = context.WithValue(c.Request.Context(), goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	apiErr = h.Service.Create(ctx, input)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, "")
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
func (h CompanyLayoutHandler) Get(c *gin.Context) {

	itemID := c.Param("id")
	err := utils.ValidateHexID([]string{itemID})
	if err != nil {
		c.Error(err)
		c.JSON(err.Status(), err)
		return
	}

	companyLayout, apiErr := h.Service.Get(c.Request.Context(), itemID)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, companyLayout)
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
func (h CompanyLayoutHandler) GetByUserID(c *gin.Context) {

	itemID := c.Param("id")
	err := utils.ValidateHexID([]string{itemID})
	if err != nil {
		c.Error(err)
		c.JSON(err.Status(), err)
		return
	}

	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)
	ctx = context.WithValue(c.Request.Context(), goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	companyLayout, apiErr := h.Service.GetByShopID(ctx, itemID)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, companyLayout)
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
func (h CompanyLayoutHandler) Update(c *gin.Context) {

	var input dto.CompanyLayoutRequest

	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		apiErr := apierrors.NewApiError(err.Error(), "bad_request", http.StatusBadRequest, apierrors.CauseList{})
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)
	ctx = context.WithValue(c.Request.Context(), goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	apiErr = h.Service.Update(ctx, input)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, "")
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
func (h CompanyLayoutHandler) Delete(c *gin.Context) {

	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)
	ctx = context.WithValue(c.Request.Context(), goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	err := h.Service.Delete(ctx)
	if err != nil {
		c.Error(err)
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusNoContent)
}
