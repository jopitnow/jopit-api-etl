package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/rest"
	"github.com/jopitnow/jopit-api-etl/src/main/api/config"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

const (
	ShopsBaseEndpoint = "/shops"
)

var ShopsClientInstance = newShopClient()

var tracerClientShops = otel.Tracer("shops-client") // Tracer for this package

type ShopClient interface {
	GetShopByUserID(ctx context.Context) (models.Shop, apierrors.ApiError)
}

type shopClient struct {
	Builder *rest.RequestBuilder
}

func newShopClient() ShopClient {

	httpClient := http.Client{}
	httpClient.Transport = otelhttp.NewTransport(http.DefaultTransport)

	builder := &rest.RequestBuilder{
		BaseURL:        config.InternalBaseShopsClient,
		Timeout:        5 * time.Second,
		ContentType:    rest.JSON,
		EnableCache:    false,
		DisableTimeout: false,
		CustomPool:     &rest.CustomPool{MaxIdleConnsPerHost: 100},
		FollowRedirect: true,
		MetricsConfig:  rest.MetricsReportConfig{TargetId: "shops-api"},
		Client:         &httpClient,
	}

	return &shopClient{Builder: builder}
}

func (client shopClient) GetShopByUserID(ctx context.Context) (models.Shop, apierrors.ApiError) {

	var shop models.Shop

	ctx, span := tracerClientShops.Start(ctx, "GetShopByUserID")
	defer span.End()

	headers := http.Header{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))

	headers.Add("Authorization", fmt.Sprint(ctx.Value(goauth.FirebaseAuthHeader)))

	response := client.Builder.Get(ShopsBaseEndpoint, rest.Headers(headers))

	if response.Response == nil {
		return models.Shop{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprint("unexpected error getting shop, url: "+ShopsBaseEndpoint), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if response.StatusCode == http.StatusNotFound {
		return models.Shop{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("shop not found", "not_found", http.StatusNotFound, apierrors.CauseList{}))
	}

	if response.StatusCode != http.StatusOK {
		return models.Shop{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("unexpected response code from shops api", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if err := json.Unmarshal(response.Bytes(), &shop); err != nil {
		return models.Shop{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("unexpected error decoding response body from shops api", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	return shop, nil
}
