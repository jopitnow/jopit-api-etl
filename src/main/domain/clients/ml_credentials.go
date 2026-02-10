package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/go-jopit-toolkit/rest"
	"github.com/jopitnow/jopit-api-etl/src/main/api/config"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

const (
	// MercadoLibre OAuth URLs - Argentina
	meliOAuthURL     = "https://auth.mercadolibre.com.ar/authorization"
	meliTokenBaseURL = "https://api.mercadolibre.com"
	meliOAuthToken   = "/oauth/token"
)

var (
	MercadoLibreAuthClientInstance = NewMercadoLibreAuthClient()
	tracerMeliAuth                 = otel.Tracer("mercadolibre-auth-client")
)

type MercadoLibreAuthClient interface {
	GetOAuthURL(ctx context.Context) (models.MercadoLibreURL, apierrors.ApiError)
	GetOAuthCredentials(ctx context.Context, code string) (dto.MercadoLibreAuthResponse, apierrors.ApiError)
	RefreshOAuthCredentials(ctx context.Context, refreshToken string) (dto.MercadoLibreAuthResponse, apierrors.ApiError)
}

type mercadoLibreAuthClient struct {
	Builder *rest.RequestBuilder
}

func NewMercadoLibreAuthClient() MercadoLibreAuthClient {
	httpClient := http.Client{}
	httpClient.Transport = otelhttp.NewTransport(http.DefaultTransport)

	builder := &rest.RequestBuilder{
		BaseURL:        meliTokenBaseURL,
		Timeout:        10 * time.Second,
		ContentType:    rest.JSON,
		EnableCache:    false,
		DisableTimeout: false,
		CustomPool:     &rest.CustomPool{MaxIdleConnsPerHost: 100},
		FollowRedirect: true,
		MetricsConfig:  rest.MetricsReportConfig{TargetId: "meli-auth-api"},
		Client:         &httpClient,
	}

	return &mercadoLibreAuthClient{Builder: builder}
}

func (c *mercadoLibreAuthClient) GetOAuthURL(ctx context.Context) (models.MercadoLibreURL, apierrors.ApiError) {
	ctx, span := tracerMeliAuth.Start(ctx, "GetOAuthURL")
	defer span.End()

	clientID := getMeliClientID()
	redirectURI := getMeliCallbackURL()

	u, err := url.Parse(meliOAuthURL)
	if err != nil {
		return models.MercadoLibreURL{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf("error parsing OAuth URL: %s", err.Error()), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	// PKCE disabled for now; we can re-enable later if needed.
	params.Add("redirect_uri", redirectURI)

	u.RawQuery = params.Encode()

	return models.MercadoLibreURL{URL: u.String()}, nil
}

func getMeliCallbackURL() string {
	env := os.Getenv("DEPLOY_ENVIRONMENT")
	var base = "https://jopit.com.ar"

	if strings.Contains(env, "stgn") {
		base = "https://staging.jopit.com.ar"
	}

	return base + "/user/shop"
}

func getMeliClientID() string {
	if config.ConfMap.MercadolibreClientId != "" {
		return config.ConfMap.MercadolibreClientId
	}

	clientID := os.Getenv("MERCADOLIBRE_CLIENT_ID")
	if clientID == "" {
		clientID = config.ConfMap.APIRestUsername // placeholder
	}

	return clientID
}

func getMeliClientSecret() string {
	if config.ConfMap.MercadolibreClientSecret != "" {
		return config.ConfMap.MercadolibreClientSecret
	}

	clientSecret := os.Getenv("MERCADOLIBRE_CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = config.ConfMap.APIRestPassword // placeholder
	}

	return clientSecret
}

func getMeliAPISignature() string {
	// TODO: Add to config
	signature := os.Getenv("MERCADOLIBRE_API_SIGNATURE")
	if signature == "" {
		signature = "jopit-meli-signature" // placeholder
	}
	return signature
}

func (c *mercadoLibreAuthClient) GetOAuthCredentials(ctx context.Context, code string) (dto.MercadoLibreAuthResponse, apierrors.ApiError) {
	var auth dto.MercadoLibreAuthResponse

	ctx, span := tracerMeliAuth.Start(ctx, "GetOAuthCredentials")
	defer span.End()

	redirectURI := getMeliCallbackURL()
	codeVerifier := getMeliAPISignature()

	request := dto.MercadoLibreAuthRequestDTO{
		ClientID:     getMeliClientID(),
		ClientSecret: getMeliClientSecret(),
		GrantType:    "authorization_code",
		Code:         &code,
		CodeVerifier: &codeVerifier,
		RedirectURI:  &redirectURI,
	}

	response := c.Builder.Post(meliOAuthToken, request, rest.Context(ctx))

	if response.Err != nil || response.Response == nil || response.StatusCode != http.StatusOK {
		return dto.MercadoLibreAuthResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("MercadoLibre credentials handshake failed.", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{response}))
	}

	if err := json.Unmarshal(response.Bytes(), &auth); err != nil {
		return dto.MercadoLibreAuthResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("unexpected error unmarshalling auth json response. value: "+string(response.Bytes()), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	return auth, nil
}

func (c *mercadoLibreAuthClient) RefreshOAuthCredentials(ctx context.Context, refreshToken string) (dto.MercadoLibreAuthResponse, apierrors.ApiError) {
	var auth dto.MercadoLibreAuthResponse

	ctx, span := tracerMeliAuth.Start(ctx, "RefreshOAuthCredentials")
	defer span.End()

	request := dto.MercadoLibreAuthRequestDTO{
		ClientID:     getMeliClientID(),
		ClientSecret: getMeliClientSecret(),
		GrantType:    "refresh_token",
		RefreshToken: &refreshToken,
	}

	response := c.Builder.Post(meliOAuthToken, request, rest.Context(ctx))

	if response.Err != nil || response.Response == nil || response.StatusCode != http.StatusOK {
		return dto.MercadoLibreAuthResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("MercadoLibre credentials refresh failed.", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{response}))
	}

	if err := json.Unmarshal(response.Bytes(), &auth); err != nil {
		return dto.MercadoLibreAuthResponse{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("unexpected error unmarshalling refresh token response. value: "+string(response.Bytes()), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	return auth, nil
}
