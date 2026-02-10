package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/clients"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/repositories"
)

type MercadoLibreCredentialsService interface {
	GetCredentialsByShopID(ctx context.Context, shopID string) (models.MercadoLibreCredential, apierrors.ApiError)
	GetCredentialsByUserID(ctx context.Context, userID string) (models.MercadoLibreCredential, apierrors.ApiError)
	GetOAuthURL(ctx context.Context) (models.MercadoLibreURL, apierrors.ApiError)
	CreateOAuthCredentials(ctx context.Context, input dto.MercadoLibreAuthRedirectDTO) apierrors.ApiError
	DeleteCredentials(ctx context.Context, userID string) apierrors.ApiError
}

type mercadoLibreCredentialsService struct {
	repository repositories.MercadoLibreCredentialsRepository
	shopClient clients.ShopClient
	authClient clients.MercadoLibreAuthClient
}

func NewMercadoLibreCredentialsService(
	repository repositories.MercadoLibreCredentialsRepository,
	shopClient clients.ShopClient,
	authClient clients.MercadoLibreAuthClient,
) MercadoLibreCredentialsService {
	return &mercadoLibreCredentialsService{
		repository: repository,
		shopClient: shopClient,
		authClient: authClient,
	}
}

func (s *mercadoLibreCredentialsService) GetCredentialsByShopID(ctx context.Context, shopID string) (models.MercadoLibreCredential, apierrors.ApiError) {
	credentials, err := s.repository.GetCredentialsByShopID(ctx, shopID)
	if err != nil {
		return models.MercadoLibreCredential{}, err
	}

	return s.checkAndRefreshToken(ctx, credentials)
}

func (s *mercadoLibreCredentialsService) GetCredentialsByUserID(ctx context.Context, userID string) (models.MercadoLibreCredential, apierrors.ApiError) {
	credentials, err := s.repository.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		return models.MercadoLibreCredential{}, err
	}

	return s.checkAndRefreshToken(ctx, credentials)
}

func (s *mercadoLibreCredentialsService) checkAndRefreshToken(ctx context.Context, credentials models.MercadoLibreCredential) (models.MercadoLibreCredential, apierrors.ApiError) {
	now := time.Now().UTC()
	expirationTime := time.Duration(credentials.ExpiresIn) * time.Second
	// Refresh 1 hour before expiration
	expirationDate := credentials.UpdatedAt.Add(expirationTime - 1*time.Hour)

	if now.After(expirationDate) {
		refresh, err := s.refreshOAuthCredentials(ctx, credentials)
		if err != nil {
			return models.MercadoLibreCredential{}, err
		}

		return refresh, nil
	}

	return credentials, nil
}

func (s *mercadoLibreCredentialsService) GetOAuthURL(ctx context.Context) (models.MercadoLibreURL, apierrors.ApiError) {
	return s.authClient.GetOAuthURL(ctx)
}

func (s *mercadoLibreCredentialsService) CreateOAuthCredentials(ctx context.Context, input dto.MercadoLibreAuthRedirectDTO) apierrors.ApiError {
	userID := fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	shop, err := s.shopClient.GetShopByUserID(ctx)
	if err != nil {
		return err
	}

	// Exchange code for tokens
	response, err := s.authClient.GetOAuthCredentials(ctx, input.Code)
	if err != nil {
		return err
	}

	model := models.MercadoLibreCredential{
		ShopID:       shop.ID,
		UserID:       userID,
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		TokenType:    response.TokenType,
		ExpiresIn:    response.ExpiresIn,
		Scope:        response.Scope,
		UserIDMeli:   response.UserID,
		UpdatedAt:    time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
	}

	// Check if credentials already exist
	credentials, err := s.repository.GetCredentialsByUserID(ctx, userID)
	if err != nil && err.Status() != http.StatusNotFound {
		return err
	} else if err != nil {
		// Create new credentials
		err = s.repository.CreateCredentials(ctx, model)
		if err != nil {
			return err
		}
	} else {
		// Update existing credentials
		model.ID = credentials.ID
		model.CreatedAt = credentials.CreatedAt
		err = s.repository.UpdateCredentials(ctx, model)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *mercadoLibreCredentialsService) DeleteCredentials(ctx context.Context, userID string) apierrors.ApiError {
	return s.repository.DeleteCredentials(ctx, userID)
}

func (s *mercadoLibreCredentialsService) refreshOAuthCredentials(ctx context.Context, credentials models.MercadoLibreCredential) (models.MercadoLibreCredential, apierrors.ApiError) {
	response, err := s.authClient.RefreshOAuthCredentials(ctx, credentials.RefreshToken)
	if err != nil {
		return models.MercadoLibreCredential{}, err
	}

	model := models.MercadoLibreCredential{
		ID:           credentials.ID,
		ShopID:       credentials.ShopID,
		UserID:       credentials.UserID,
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		TokenType:    response.TokenType,
		ExpiresIn:    response.ExpiresIn,
		Scope:        response.Scope,
		UserIDMeli:   response.UserID,
		UpdatedAt:    time.Now().UTC(),
		CreatedAt:    credentials.CreatedAt,
	}

	err = s.repository.UpdateCredentials(ctx, model)
	if err != nil {
		return models.MercadoLibreCredential{}, err
	}

	return model, nil
}
