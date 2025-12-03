package clients

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jopitnow/jopit-api-etl/src/main/api/config"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/go-jopit-toolkit/rest"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var (
	ItemsClientInstance = newItemsClient()
	tracerClientItems   = otel.Tracer("items-client") // Tracer for this package
)

const (
	GetIntegrity = "/items/list" //to-do
)

type itemsClient struct {
	Client *rest.RequestBuilder
}

type ItemsClient interface {
	BulkCreateItems(ctx context.Context, items []models.Item) apierrors.ApiError
	BulkDeleteItems(ctx context.Context, batchID string) apierrors.ApiError
}

func newItemsClient() *itemsClient {

	httpClient := http.Client{}
	httpClient.Transport = otelhttp.NewTransport(http.DefaultTransport)

	customPool := &rest.CustomPool{
		MaxIdleConnsPerHost: 100,
	}
	restClientItems := &rest.RequestBuilder{
		BaseURL:        config.InternalBaseItemsClient,
		Timeout:        5 * time.Second,
		ContentType:    rest.JSON,
		EnableCache:    false,
		DisableTimeout: false,
		CustomPool:     customPool,
		FollowRedirect: true,
		//RetryStrategy:  retry.NewSimpleRetryStrategy(3, 30*time.Millisecond),
		MetricsConfig: rest.MetricsReportConfig{TargetId: "items-api"},
		Client:        &httpClient,
	}
	return &itemsClient{Client: restClientItems}
}

func (c *itemsClient) BulkCreateItems(ctx context.Context, items []models.Item) apierrors.ApiError {

	ctx, span := tracerClientItems.Start(ctx, "BulkCreateItems")
	defer span.End()

	headers := http.Header{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))

	var response *rest.Response

	reqBody := dto.BulkCreateItemsRequest{Items: items}

	endpoint := GetIntegrity
	response = c.Client.Post(endpoint, reqBody, rest.Context(ctx), rest.Headers(headers))

	if response.Err != nil || response.Response == nil || (response.StatusCode != http.StatusNotFound && response.StatusCode != http.StatusOK) {
		return apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprint("Unexpected error hitting items api, url: "+endpoint, "\nresponse: ", response), "error hitting Items Api", http.StatusInternalServerError, apierrors.CauseList{response}))
	}

	if response.StatusCode != http.StatusCreated {
		return apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprint("Unexpected error hitting items api, url: "+endpoint, "\nresponse: ", response), "error hitting Items Api", http.StatusInternalServerError, apierrors.CauseList{response}))
	}

	return nil
}

func (c *itemsClient) BulkDeleteItems(ctx context.Context, batchID string) apierrors.ApiError {

	return nil

}
