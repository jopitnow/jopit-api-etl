package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/clients"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
)

type MercadoLibreService interface {
	GetItem(ctx context.Context, meliItemID string) (dto.MeliItemResponse, apierrors.ApiError)
	GetItems(ctx context.Context, meliItemIDs []string) ([]dto.MeliItemResponse, apierrors.ApiError)
	GetUserItems(ctx context.Context) (dto.MeliUserItemsSearchResponse, apierrors.ApiError)
	SearchItemsBySeller(ctx context.Context, siteID string) (dto.MeliUserItemsSearchResponse, apierrors.ApiError)
	GetSizeChart(ctx context.Context, chartID string) (dto.MeliSizeChartResponse, apierrors.ApiError)
	GetUserItemsDetails(ctx context.Context) ([]dto.MeliItemResponse, apierrors.ApiError)
	GetUserItemsDetailsWithPagination(ctx context.Context, pageSize int) ([]dto.MeliItemResponse, apierrors.ApiError)
}

type mercadoLibreService struct {
	meliClient         clients.MercadoLibreClient
	credentialsService MercadoLibreCredentialsService
}

func NewMercadoLibreService(
	meliClient clients.MercadoLibreClient,
	credentialsService MercadoLibreCredentialsService,
) MercadoLibreService {
	return &mercadoLibreService{
		meliClient:         meliClient,
		credentialsService: credentialsService,
	}
}

func (s *mercadoLibreService) GetItem(ctx context.Context, meliItemID string) (dto.MeliItemResponse, apierrors.ApiError) {
	userID := fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	// Get credentials with auto-refresh
	credentials, err := s.credentialsService.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		return dto.MeliItemResponse{}, err
	}

	// Call MercadoLibre API
	item, err := s.meliClient.GetItem(ctx, meliItemID, credentials.AccessToken)
	if err != nil {
		return dto.MeliItemResponse{}, err
	}

	return item, nil
}

func (s *mercadoLibreService) GetItems(ctx context.Context, meliItemIDs []string) ([]dto.MeliItemResponse, apierrors.ApiError) {
	userID := fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	// Get credentials with auto-refresh
	credentials, err := s.credentialsService.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Call MercadoLibre API
	items, err := s.meliClient.GetItems(ctx, meliItemIDs, credentials.AccessToken)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *mercadoLibreService) SearchItemsBySeller(ctx context.Context, siteID string) (dto.MeliUserItemsSearchResponse, apierrors.ApiError) {
	userID := fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	// Get credentials with auto-refresh
	credentials, err := s.credentialsService.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		return dto.MeliUserItemsSearchResponse{}, err
	}

	if credentials.UserIDMeli == 0 {
		return dto.MeliUserItemsSearchResponse{}, apierrors.NewApiError("seller_id not found in credentials", "bad_request", http.StatusBadRequest, apierrors.CauseList{})
	}

	// Call MercadoLibre API using the authenticated user items endpoint
	result, err := s.meliClient.GetUserItems(ctx, credentials.UserIDMeli, credentials.AccessToken)
	if err != nil {
		return dto.MeliUserItemsSearchResponse{}, err
	}

	return result, nil
}

func (s *mercadoLibreService) GetUserItems(ctx context.Context) (dto.MeliUserItemsSearchResponse, apierrors.ApiError) {
	userID := fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	// Get credentials with auto-refresh
	credentials, err := s.credentialsService.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		return dto.MeliUserItemsSearchResponse{}, err
	}

	if credentials.UserIDMeli == 0 {
		return dto.MeliUserItemsSearchResponse{}, apierrors.NewApiError("seller_id not found in credentials", "bad_request", http.StatusBadRequest, apierrors.CauseList{})
	}

	// Call MercadoLibre API using the authenticated user items endpoint
	result, err := s.meliClient.GetUserItems(ctx, credentials.UserIDMeli, credentials.AccessToken)
	if err != nil {
		return dto.MeliUserItemsSearchResponse{}, err
	}

	return result, nil
}

func (s *mercadoLibreService) GetUserItemsDetails(ctx context.Context) ([]dto.MeliItemResponse, apierrors.ApiError) {

	userID := fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	// Get credentials with auto-refresh
	credentials, err := s.credentialsService.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		return []dto.MeliItemResponse{}, err
	}

	if credentials.UserIDMeli == 0 {
		return []dto.MeliItemResponse{}, apierrors.NewApiError("seller_id not found in credentials", "bad_request", http.StatusBadRequest, apierrors.CauseList{})
	}

	// Step 1: Search for user items to get the list of item IDs
	searchResult, err := s.meliClient.GetUserItems(ctx, credentials.UserIDMeli, credentials.AccessToken)
	if err != nil {
		return []dto.MeliItemResponse{}, err
	}

	if len(searchResult.Results) == 0 {
		return []dto.MeliItemResponse{}, nil
	}

	// Step 2: Batch fetch full details for all items using MercadoLibre's multi-get endpoint
	itemsDetails, err := s.meliClient.GetItems(ctx, searchResult.Results, credentials.AccessToken)
	if err != nil {
		return []dto.MeliItemResponse{}, err
	}

	return itemsDetails, nil

}

func (s *mercadoLibreService) GetUserItemsDetailsWithPagination(ctx context.Context, pageSize int) ([]dto.MeliItemResponse, apierrors.ApiError) {
	userID := fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	// Get credentials with auto-refresh
	credentials, err := s.credentialsService.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		return []dto.MeliItemResponse{}, err
	}

	if credentials.UserIDMeli == 0 {
		return []dto.MeliItemResponse{}, apierrors.NewApiError("seller_id not found in credentials", "bad_request", http.StatusBadRequest, apierrors.CauseList{})
	}

	// Step 1: Fetch all item IDs with pagination
	var allItemIDs []string
	offset := 0

	for {
		searchResult, err := s.meliClient.GetUserItemsWithPagination(ctx, credentials.UserIDMeli, credentials.AccessToken, offset, pageSize)
		if err != nil {
			return []dto.MeliItemResponse{}, err
		}

		if len(searchResult.Results) == 0 {
			break
		}

		allItemIDs = append(allItemIDs, searchResult.Results...)
		offset += pageSize

		// Check if we've received fewer results than requested (last page)
		if len(searchResult.Results) < pageSize {
			break
		}
	}

	if len(allItemIDs) == 0 {
		return []dto.MeliItemResponse{}, nil
	}

	// Step 2: Batch fetch full details for all items using MercadoLibre's multi-get endpoint
	itemsDetails, err := s.meliClient.GetItems(ctx, allItemIDs, credentials.AccessToken)
	if err != nil {
		return []dto.MeliItemResponse{}, err
	}

	return itemsDetails, nil
}

func (s *mercadoLibreService) GetSizeChart(ctx context.Context, chartID string) (dto.MeliSizeChartResponse, apierrors.ApiError) {
	userID := fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	// Get credentials with auto-refresh
	credentials, err := s.credentialsService.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		return dto.MeliSizeChartResponse{}, err
	}

	// Call MercadoLibre API to get size chart
	sizeChart, err := s.meliClient.GetSizeChart(ctx, chartID, credentials.AccessToken)
	if err != nil {
		return dto.MeliSizeChartResponse{}, err
	}

	return sizeChart, nil
}
