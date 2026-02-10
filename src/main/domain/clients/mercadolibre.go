package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/go-jopit-toolkit/rest"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const (
	meliAPIBaseURL = "https://api.mercadolibre.com"
)

var (
	MercadoLibreClientInstance = NewMercadoLibreClient()
	tracerMeliClient           = otel.Tracer("mercadolibre-client")
)

type MercadoLibreClient interface {
	GetItem(ctx context.Context, meliItemID string, accessToken string) (dto.MeliItemResponse, apierrors.ApiError)
	GetItems(ctx context.Context, meliItemIDs []string, accessToken string) ([]dto.MeliItemResponse, apierrors.ApiError)
	GetUserItems(ctx context.Context, meliUserID int64, accessToken string) (dto.MeliUserItemsSearchResponse, apierrors.ApiError)
	SearchItems(ctx context.Context, filters dto.MercadoLibreSearchFilters, accessToken string) (dto.MeliSearchResponse, apierrors.ApiError)
	GetSizeChart(ctx context.Context, chartID string, accessToken string) (dto.MeliSizeChartResponse, apierrors.ApiError)
}

type mercadoLibreClient struct {
	Builder *rest.RequestBuilder
}

func NewMercadoLibreClient() MercadoLibreClient {
	httpClient := http.Client{}
	httpClient.Transport = otelhttp.NewTransport(http.DefaultTransport)

	builder := &rest.RequestBuilder{
		BaseURL:        meliAPIBaseURL,
		Timeout:        10 * time.Second,
		ContentType:    rest.JSON,
		EnableCache:    false,
		DisableTimeout: false,
		CustomPool:     &rest.CustomPool{MaxIdleConnsPerHost: 100},
		FollowRedirect: true,
		MetricsConfig:  rest.MetricsReportConfig{TargetId: "meli-api"},
		Client:         &httpClient,
	}

	return &mercadoLibreClient{Builder: builder}
}

func (c *mercadoLibreClient) GetItem(ctx context.Context, meliItemID string, accessToken string) (dto.MeliItemResponse, apierrors.ApiError) {
	ctx, span := tracerMeliClient.Start(ctx, "GetItem")
	defer span.End()

	if strings.TrimSpace(meliItemID) == "" {
		return dto.MeliItemResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("meli item id is required", "bad_request", http.StatusBadRequest, apierrors.CauseList{}))
	}

	headers := http.Header{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))
	if accessToken != "" {
		headers.Add("Authorization", "Bearer "+accessToken)
	}

	endpoint := fmt.Sprintf("/items/%s", meliItemID)
	response := c.Builder.Get(endpoint, rest.Context(ctx), rest.Headers(headers))

	if response.Response == nil {
		return dto.MeliItemResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("unexpected error calling MercadoLibre items endpoint", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if response.StatusCode != http.StatusOK {
		return dto.MeliItemResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf("unexpected response from MercadoLibre items endpoint, status: %d", response.StatusCode), "bad_gateway", http.StatusBadGateway, apierrors.CauseList{response}))
	}

	var item dto.MeliItemResponse
	if err := json.Unmarshal(response.Bytes(), &item); err != nil {
		return dto.MeliItemResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("error decoding MercadoLibre item response", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	return item, nil
}

func (c *mercadoLibreClient) GetItems(ctx context.Context, meliItemIDs []string, accessToken string) ([]dto.MeliItemResponse, apierrors.ApiError) {
	ctx, span := tracerMeliClient.Start(ctx, "GetItems")
	defer span.End()

	if len(meliItemIDs) == 0 {
		return nil, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("meli item ids are required", "bad_request", http.StatusBadRequest, apierrors.CauseList{}))
	}

	headers := http.Header{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))
	if accessToken != "" {
		headers.Add("Authorization", "Bearer "+accessToken)
	}

	query := url.Values{}
	query.Set("ids", strings.Join(meliItemIDs, ","))
	endpoint := fmt.Sprintf("/items?%s", query.Encode())

	response := c.Builder.Get(endpoint, rest.Context(ctx), rest.Headers(headers))
	if response.Response == nil {
		return nil, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("unexpected error calling MercadoLibre batch items endpoint", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if response.StatusCode != http.StatusOK {
		return nil, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf("unexpected response from MercadoLibre batch items endpoint, status: %d", response.StatusCode), "bad_gateway", http.StatusBadGateway, apierrors.CauseList{response}))
	}

	// MercadoLibre returns array of objects with code and body fields for batch requests
	type batchResponse struct {
		Code int                  `json:"code"`
		Body dto.MeliItemResponse `json:"body"`
	}

	var batchItems []batchResponse
	if err := json.Unmarshal(response.Bytes(), &batchItems); err != nil {
		return nil, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("error decoding MercadoLibre batch items response", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	items := make([]dto.MeliItemResponse, 0, len(batchItems))
	for _, batch := range batchItems {
		if batch.Code == http.StatusOK {
			items = append(items, batch.Body)
		}
	}

	return items, nil
}

func (c *mercadoLibreClient) GetUserItems(ctx context.Context, meliUserID int64, accessToken string) (dto.MeliUserItemsSearchResponse, apierrors.ApiError) {
	ctx, span := tracerMeliClient.Start(ctx, "GetUserItems")
	defer span.End()

	if meliUserID == 0 {
		return dto.MeliUserItemsSearchResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("meli_user_id is required for MercadoLibre user items search", "bad_request", http.StatusBadRequest, apierrors.CauseList{}))
	}

	headers := http.Header{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))
	if accessToken != "" {
		headers.Add("Authorization", "Bearer "+accessToken)
	}

	endpoint := fmt.Sprintf("/users/%d/items/search", meliUserID)
	response := c.Builder.Get(endpoint, rest.Context(ctx), rest.Headers(headers))

	if response.Response == nil {
		return dto.MeliUserItemsSearchResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("unexpected error calling MercadoLibre user items endpoint", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if response.StatusCode != http.StatusOK {
		return dto.MeliUserItemsSearchResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf("unexpected response from MercadoLibre user items endpoint, status: %d, body: %s", response.StatusCode, string(response.Bytes())), "bad_gateway", http.StatusBadGateway, apierrors.CauseList{response}))
	}

	var searchResult dto.MeliUserItemsSearchResponse
	if err := json.Unmarshal(response.Bytes(), &searchResult); err != nil {
		return dto.MeliUserItemsSearchResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("error decoding MercadoLibre user items search response", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	return searchResult, nil
}

func (c *mercadoLibreClient) SearchItems(ctx context.Context, filters dto.MercadoLibreSearchFilters, accessToken string) (dto.MeliSearchResponse, apierrors.ApiError) {
	ctx, span := tracerMeliClient.Start(ctx, "SearchItems")
	defer span.End()

	if strings.TrimSpace(filters.SiteID) == "" {
		return dto.MeliSearchResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("site_id is required for MercadoLibre search", "bad_request", http.StatusBadRequest, apierrors.CauseList{}))
	}

	headers := http.Header{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))
	if accessToken != "" {
		headers.Add("Authorization", "Bearer "+accessToken)
	}

	query := url.Values{}
	if filters.SellerID != "" {
		query.Set("seller_id", filters.SellerID)
	}
	if filters.Query != "" {
		query.Set("q", filters.Query)
	}
	for k, v := range filters.Params {
		if k == "" {
			continue
		}
		query.Set(k, v)
	}

	endpoint := fmt.Sprintf("/sites/%s/search?%s", filters.SiteID, query.Encode())
	response := c.Builder.Get(endpoint, rest.Context(ctx), rest.Headers(headers))

	if response.Response == nil {
		return dto.MeliSearchResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("unexpected error calling MercadoLibre search endpoint", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if response.StatusCode != http.StatusOK {
		return dto.MeliSearchResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf("unexpected response from MercadoLibre search endpoint, status: %d", response.StatusCode), "bad_gateway", http.StatusBadGateway, apierrors.CauseList{response}))
	}

	var searchResult dto.MeliSearchResponse
	if err := json.Unmarshal(response.Bytes(), &searchResult); err != nil {
		return dto.MeliSearchResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("error decoding MercadoLibre search response", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	return searchResult, nil
}

func (c *mercadoLibreClient) GetSizeChart(ctx context.Context, chartID string, accessToken string) (dto.MeliSizeChartResponse, apierrors.ApiError) {
	ctx, span := tracerMeliClient.Start(ctx, "GetSizeChart")
	defer span.End()

	if strings.TrimSpace(chartID) == "" {
		return dto.MeliSizeChartResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("chart_id is required for MercadoLibre size chart", "bad_request", http.StatusBadRequest, apierrors.CauseList{}))
	}

	headers := http.Header{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))
	if accessToken != "" {
		headers.Add("Authorization", "Bearer "+accessToken)
	}

	endpoint := fmt.Sprintf("/catalog/charts/%s", chartID)
	response := c.Builder.Get(endpoint, rest.Context(ctx), rest.Headers(headers))

	if response.Response == nil {
		return dto.MeliSizeChartResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("unexpected error calling MercadoLibre size chart endpoint", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if response.StatusCode != http.StatusOK {
		return dto.MeliSizeChartResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf("unexpected response from MercadoLibre size chart endpoint, status: %d", response.StatusCode), "bad_gateway", http.StatusBadGateway, apierrors.CauseList{response}))
	}

	var sizeChart dto.MeliSizeChartResponse
	if err := json.Unmarshal(response.Bytes(), &sizeChart); err != nil {
		return dto.MeliSizeChartResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("error decoding MercadoLibre size chart response", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	return sizeChart, nil
}
