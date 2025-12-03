package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jopitnow/go-jopit-toolkit/rest"
	"github.com/jopitnow/jopit-api-etl/src/main/api/config"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

var FetchApiClientInstance = newFetchApiClient()
var tracerClientEtlFetchAPI = otel.Tracer("etl-http-client") // Tracer for this package

type FetchApiClient interface {
	FetchAPI(ctx context.Context, layout models.CompanyLayout) ([]map[string]string, apierrors.ApiError)
}

type fetchApiClient struct {
	Builder *rest.RequestBuilder
}

func newFetchApiClient() FetchApiClient {

	httpclient := http.Client{}
	httpclient.Transport = otelhttp.NewTransport(http.DefaultTransport)

	builder := &rest.RequestBuilder{
		BaseURL:        config.InternalBaseShopsClient,
		Timeout:        5 * time.Second,
		ContentType:    rest.JSON,
		EnableCache:    false,
		DisableTimeout: false,
		CustomPool:     &rest.CustomPool{MaxIdleConnsPerHost: 100},
		FollowRedirect: true,
		MetricsConfig:  rest.MetricsReportConfig{TargetId: "shops-api"},
		Client:         &httpclient,
	}

	return &fetchApiClient{Builder: builder}
}

func (client *fetchApiClient) FetchAPI(ctx context.Context, layout models.CompanyLayout) ([]map[string]string, apierrors.ApiError) {

	ctx, span := tracerClientEtlFetchAPI.Start(ctx, "ExtractFromAPI")
	defer span.End()

	// Prepare headers with tracing context
	headers := http.Header{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))

	// Build request
	req, err := http.NewRequest(layout.Request.Method, layout.Request.Endpoint, nil)
	if err != nil {
		return nil, apierrors.NewApiError(
			fmt.Sprintf("error building request for %s API. %s", layout.Name, err.Error()),
			"internal_server_error",
			http.StatusInternalServerError,
			apierrors.CauseList{})
	}

	// Add headers
	for k, v := range layout.Request.Headers {
		req.Header.Set(k, v)
	}

	// Add query params
	q := req.URL.Query()
	for k, v := range layout.Request.QueryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Add body if present
	if layout.Request.Method == "POST" || layout.Request.Method == "PUT" {
		bodyBytes, _ := json.Marshal(layout.Request.Body)
		req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, apierrors.NewApiError(
			fmt.Sprintf("error fetching the %s API. %s", layout.Name, err.Error()),
			"internal_server_error",
			http.StatusInternalServerError,
			apierrors.CauseList{})
	}
	defer resp.Body.Close()

	var records []map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, apierrors.NewApiError(
			fmt.Sprintf("error decoding from the %s API. %s", layout.Name, err.Error()),
			"internal_server_error",
			http.StatusInternalServerError,
			apierrors.CauseList{})
	}

	return records, nil
}
