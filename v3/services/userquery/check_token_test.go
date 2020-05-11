package userqueryapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"gitlab.com/wiserskills/v3/services/userquery/gen/environment"
)

func TestPOSTCheckToken(t *testing.T) {

	os.Setenv("TOKEN_ACTIVE", "true")

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// We first sign in to get a valid token and create a new session
	orgID := "wiserskills"

	obj := e.POST("/signin").
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", orgID).
		WithHeader("X-AREA-ID", "DEV").
		WithBasicAuth("mpolo", "mysecret").
		Expect().
		Status(http.StatusOK).JSON().Object()

	obj.Keys().ContainsOnly("token")
	tkn := obj.Value("token").Raw().(string)

	e.POST("/check").
		WithHeader("Authorization", "Bearer "+tkn).
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", orgID).
		WithHeader("X-AREA-ID", "DEV").
		Expect().
		Status(http.StatusOK)
}
