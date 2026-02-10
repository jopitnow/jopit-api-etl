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
)

type MercadoLibreCredentialsHandler struct {
	service services.MercadoLibreCredentialsService
}

func NewMercadoLibreCredentialsHandler(
	service services.MercadoLibreCredentialsService,
) MercadoLibreCredentialsHandler {
	return MercadoLibreCredentialsHandler{
		service: service,
	}
}

// GetOAuthURL godoc
// @Summary Get MercadoLibre OAuth URL
// @Description Get MercadoLibre OAuth URL for authentication.
// @Tags MercadoLibre OAuth Credentials
// @Param Authorization header string true "Bearer token"
// @Produce json
// @Success 200 {object} models.MercadoLibreURL
// @Failure 401 "Unauthorized Firebase Token"
// @Failure 500 "Internal Server Error"
// @Router /mercadolibre/oauth [get]
func (h *MercadoLibreCredentialsHandler) GetOAuthURL(c *gin.Context) {
	userID, err := goauth.GetUserId(c)
	if err != nil {
		c.Error(err)
		c.JSON(err.Status(), err)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)

	url, err := h.service.GetOAuthURL(ctx)
	if err != nil {
		c.Error(err)
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, url)
}

// CreateOAuthCredentials godoc
// @Summary Create MercadoLibre OAuth Credentials
// @Description Exchange authorization code for MercadoLibre OAuth credentials and store them.
// @Tags MercadoLibre OAuth Credentials
// @Param Authorization header string true "Bearer token"
// @Param code body dto.MercadoLibreAuthRedirectDTO true "Authorization code from MercadoLibre"
// @Accept json
// @Produce json
// @Success 204 "No Content - Credentials created successfully"
// @Failure 400 "Bad Request - Invalid input"
// @Failure 401 "Unauthorized Firebase Token"
// @Failure 500 "Internal Server Error"
// @Router /mercadolibre/oauth [post]
func (h *MercadoLibreCredentialsHandler) CreateOAuthCredentials(c *gin.Context) {
	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)
	ctx = context.WithValue(ctx, goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	input := dto.MercadoLibreAuthRedirectDTO{}

	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		apiErr = apierrors.NewApiError(err.Error(), "bad_request", http.StatusBadRequest, apierrors.CauseList{})
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	apiErr = h.service.CreateOAuthCredentials(ctx, input)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCredentials godoc
// @Summary Get MercadoLibre Credentials
// @Description Get stored MercadoLibre credentials for the authenticated user.
// @Tags MercadoLibre OAuth Credentials
// @Param Authorization header string true "Bearer token"
// @Produce json
// @Success 200 {object} models.MercadoLibreCredential
// @Failure 401 "Unauthorized Firebase Token"
// @Failure 404 "Not Found - No credentials found"
// @Failure 500 "Internal Server Error"
// @Router /mercadolibre/credentials [get]
func (h *MercadoLibreCredentialsHandler) GetCredentials(c *gin.Context) {
	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)

	credentials, apiErr := h.service.GetCredentialsByUserID(ctx, userID)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Don't expose sensitive tokens in response
	credentials.AccessToken = ""
	credentials.RefreshToken = ""

	c.JSON(http.StatusOK, credentials)
}

// DeleteCredentials godoc
// @Summary Delete MercadoLibre Credentials
// @Description Delete stored MercadoLibre credentials for the authenticated user.
// @Tags MercadoLibre OAuth Credentials
// @Param Authorization header string true "Bearer token"
// @Produce json
// @Success 204 "No Content - Credentials deleted successfully"
// @Failure 401 "Unauthorized Firebase Token"
// @Failure 404 "Not Found - No credentials found"
// @Failure 500 "Internal Server Error"
// @Router /mercadolibre/credentials [delete]
func (h *MercadoLibreCredentialsHandler) DeleteCredentials(c *gin.Context) {
	userID, apiErr := goauth.GetUserId(c)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx := context.WithValue(c.Request.Context(), goauth.FirebaseUserID, userID)

	apiErr = h.service.DeleteCredentials(ctx, userID)
	if apiErr != nil {
		c.Error(apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}
