package services

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/clients"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/utils"

	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

type EtlService interface {
	LoadApi(ctx context.Context) (string, apierrors.ApiError)
	LoadCsv(ctx context.Context, file *multipart.FileHeader) (string, apierrors.ApiError)
	LoadMercadoLibre(ctx context.Context) (*ETLResult, apierrors.ApiError)
	DeleteBatch(ctx context.Context, batchID string) apierrors.ApiError
}

// ETLResult contains the results of an ETL operation
type ETLResult struct {
	BatchID      string       `json:"batch_id"`
	TotalItems   int          `json:"total_items"`
	CreatedCount int          `json:"created_count"`
	UpdatedCount int          `json:"updated_count"`
	FailureCount int          `json:"failure_count"`
	FailedItems  []FailedItem `json:"failed_items,omitempty"`
}

// FailedItem represents an item that failed during ETL
type FailedItem struct {
	ExternalID   string `json:"external_id"`
	Title        string `json:"title,omitempty"`
	FailureStage string `json:"failure_stage"` // "transform" or "load"
	ErrorMessage string `json:"error_message"`
}

type etlService struct {
	companyConfigService CompanyLayoutService
	itemsClient          clients.ItemsClient
	httpClient           clients.FetchApiClient
	shopsClient          clients.ShopClient
	mercadoLibreService  MercadoLibreService
}

func NewEtlService(
	httpClient clients.FetchApiClient,
	itemsClient clients.ItemsClient,
	shopsClient clients.ShopClient,
	mercadoLibreService MercadoLibreService,
) EtlService {
	return &etlService{
		httpClient:          httpClient,
		itemsClient:         itemsClient,
		shopsClient:         shopsClient,
		mercadoLibreService: mercadoLibreService,
	}
}

func (s *etlService) LoadApi(ctx context.Context) (string, apierrors.ApiError) {

	shop, err := s.shopsClient.GetShopByUserID(ctx)
	if err != nil {
		return "", err
	}

	companyLayout, err := s.companyConfigService.GetByShopID(ctx, shop.ID)
	if err != nil {
		return "", err
	}

	//validate Layout ??

	response, err := s.httpClient.FetchAPI(ctx, companyLayout)
	if err != nil {
		return "", err
	}

	batchID, items := utils.Transform(response, companyLayout, fmt.Sprint(ctx.Value(goauth.FirebaseAuthHeader)))

	err = s.itemsClient.BulkCreateItems(ctx, items)
	if err != nil {
		return "", err
	}

	return batchID, nil
}

func (s *etlService) LoadCsv(ctx context.Context, file *multipart.FileHeader) (string, apierrors.ApiError) {

	shop, err := s.shopsClient.GetShopByUserID(ctx)
	if err != nil {
		return "", err
	}

	companyLayout, err := s.companyConfigService.GetByShopID(ctx, shop.ID)
	if err != nil {
		return "", err
	}

	//validate Layout ??

	data, err := utils.ExtractFromCSV(file)
	if err != nil {
		return "", err
	}

	batchID, items := utils.Transform(data, companyLayout, fmt.Sprint(ctx.Value(goauth.FirebaseAuthHeader)))

	err = s.itemsClient.BulkCreateItems(ctx, items)
	if err != nil {
		return "", err
	}

	return batchID, nil
}

func (s *etlService) DeleteBatch(ctx context.Context, batchID string) apierrors.ApiError {
	return s.itemsClient.BulkDeleteItems(ctx, batchID)
}

// LoadMercadoLibre performs full ETL from MercadoLibre to Jopit Items
func (s *etlService) LoadMercadoLibre(ctx context.Context) (*ETLResult, apierrors.ApiError) {
	// Get shop and user info
	shop, err := s.shopsClient.GetShopByUserID(ctx)
	if err != nil {
		return nil, err
	}

	userID := fmt.Sprint(ctx.Value(goauth.FirebaseUserID))
	batchID := fmt.Sprintf("meli-%s", userID)

	// STEP 1: EXTRACT - Get all MercadoLibre items with pagination
	meliItems, err := s.mercadoLibreService.GetUserItemsDetailsWithPagination(ctx, 50) // 50 items per page
	if err != nil {
		return nil, err
	}

	if len(meliItems) == 0 {
		return nil, apierrors.NewApiError("no items found from MercadoLibre", "not_found", 404, apierrors.CauseList{})
	}

	// Write MercadoLibre items to JSON file
	if meliJSON, err := json.MarshalIndent(meliItems, "", "  "); err == nil {
		os.WriteFile("meli-items-extracted.json", meliJSON, 0644)
	}

	// STEP 2: TRANSFORM - Convert MercadoLibre items to Jopit format
	jopitItems := make([]models.Item, 0, len(meliItems))
	failedItems := make([]FailedItem, 0)

	for _, meliItem := range meliItems {
		// Transform with error handling
		jopitItem, transformErr := s.transformMeliItem(ctx, meliItem, shop.ID, userID, batchID)

		if transformErr != nil {
			// Log failure and continue
			failedItems = append(failedItems, FailedItem{
				ExternalID:   meliItem.ID,
				Title:        meliItem.Title,
				FailureStage: "transform",
				ErrorMessage: transformErr.Error(),
			})
			continue
		}

		jopitItems = append(jopitItems, jopitItem)
	}

	// Write Jopit items to JSON file
	if jopitJSON, err := json.MarshalIndent(jopitItems, "", "  "); err == nil {
		os.WriteFile("jopit-items-transformed.json", jopitJSON, 0644)
	}

	// Write failed items to JSON file if any
	if len(failedItems) > 0 {
		if failedJSON, err := json.MarshalIndent(failedItems, "", "  "); err == nil {
			os.WriteFile("jopit-items-failed.json", failedJSON, 0644)
		}
	}

	// STEP 3: LOAD - Bulk upsert items into Jopit Items API
	var createdCount int64
	var updatedCount int64

	if len(jopitItems) > 0 {
		upsertResponse, upsertErr := s.itemsClient.BulkUpsertItems(ctx, jopitItems)
		if upsertErr != nil {
			// If bulk upsert fails entirely, mark all as failed
			for _, item := range jopitItems {
				failedItems = append(failedItems, FailedItem{
					ExternalID:   item.Source.ExternalID,
					Title:        item.Name,
					FailureStage: "load",
					ErrorMessage: upsertErr.Message(),
				})
			}
			jopitItems = []models.Item{} // Clear successful items
		} else {
			createdCount = upsertResponse.CreatedCount
			updatedCount = upsertResponse.UpdatedCount
		}
	}

	result := &ETLResult{
		BatchID:      batchID,
		TotalItems:   len(meliItems),
		CreatedCount: int(createdCount),
		UpdatedCount: int(updatedCount),
		FailureCount: len(failedItems),
		FailedItems:  failedItems,
	}

	// Return error only if ALL items failed
	successCount := result.CreatedCount + result.UpdatedCount
	if successCount == 0 && result.FailureCount > 0 {
		return result, apierrors.NewApiError(
			fmt.Sprintf("all %d items failed to load", result.FailureCount),
			"etl_failed",
			500,
			apierrors.CauseList{},
		)
	}

	// If at least one succeeded, return 200 (handled in handler)
	return result, nil
}

// transformMeliItem safely transforms a single MercadoLibre item with error recovery
func (s *etlService) transformMeliItem(
	ctx context.Context,
	meliItem dto.MeliItemResponse,
	shopID, userID, batchID string,
) (models.Item, error) {
	var result models.Item
	var transformErr error

	// Recover from panics during transformation
	defer func() {
		if r := recover(); r != nil {
			transformErr = fmt.Errorf("panic during transformation: %v", r)
		}
	}()

	// Extract size chart ID if exists
	sizeChartID := utils.ExtractSizeChartID(meliItem.Attributes)
	var sizeChart *dto.MeliSizeChartResponse

	// Fetch size chart if ID exists
	if sizeChartID != "" {
		chart, chartErr := s.mercadoLibreService.GetSizeChart(ctx, sizeChartID)
		if chartErr == nil {
			sizeChart = &chart
		}
		// Continue even if size chart fetch fails
	}

	// Transform item
	result = utils.TransformMeliItemToJopitItem(meliItem, shopID, userID, batchID, sizeChart)

	return result, transformErr
}
