package shops

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/jopitnow/jopit-api-etl/src/main/api/config"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/clients"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/domain/mocks"

	"github.com/stretchr/testify/assert"
)

// ================================== Get Shop By UserID Tests ==================================

func TestClient_GetShopByUserID_Success(t *testing.T) {
	var model = mocks.Shop

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.InternalBaseShopsClient+clients.ShopsBaseEndpoint,
		func(req *http.Request) (*http.Response, error) {
			if !httpmock.HeaderExists("Authorization").Check(req) {
				return nil, errors.New("missing Authorization header")
			}

			response, err := httpmock.NewJsonResponse(200, model)

			return response, err
		},
	)

	api := clients.NewShopClient()

	body, err := api.GetShopByUserID(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, model, body)
}

func TestClient_GetShopByUserID_Response_Error(t *testing.T) {
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.InternalBaseShopsClient+clients.ShopsBaseEndpoint,
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("response error")
		},
	)

	api := clients.NewShopClient()

	_, err := api.GetShopByUserID(context.Background())

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, "unexpected error getting shop, url: /shops", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_GetShopByUserID_Not_Found_Error(t *testing.T) {
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.InternalBaseShopsClient+clients.ShopsBaseEndpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(404, nil)
		},
	)

	api := clients.NewShopClient()

	_, err := api.GetShopByUserID(context.Background())

	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, "shop not found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestClient_GetShopByUserID_Internal_Server_Error(t *testing.T) {
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.InternalBaseShopsClient+clients.ShopsBaseEndpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(505, nil)
		},
	)

	api := clients.NewShopClient()

	_, err := api.GetShopByUserID(context.Background())

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_GetShopByUserID_UnMarshall_Error(t *testing.T) {
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.InternalBaseShopsClient+clients.ShopsBaseEndpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, "fake")
		},
	)

	api := clients.NewShopClient()

	_, err := api.GetShopByUserID(context.Background())

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}
