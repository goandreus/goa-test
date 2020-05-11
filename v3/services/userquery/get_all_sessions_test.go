package userqueryapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"gitlab.com/wiserskills/v3/services/userquery/gen/environment"
)

func TestGETGetAllSessions(t *testing.T) {

	os.Setenv("TOKEN_ACTIVE", "true")

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// We first sign in to get a token and create a new session
	orgID := "wiserskills"

	obj1 := e.POST("/signin").
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", orgID).
		WithHeader("X-AREA-ID", "DEV").
		WithBasicAuth("mpolo", "mysecret").
		Expect().
		Status(http.StatusOK).JSON().Object()

	obj1.Keys().ContainsOnly("token")
	tkn := obj1.Value("token").Raw().(string)

	// We wait for the event bus (the Token is revoked asynchronously)
	time.Sleep(time.Second * 3)

	e.GET("/sessions").
		WithHeader("Authorization", "Bearer "+tkn).
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", "wiserskills").
		WithHeader("X-AREA-ID", "DEV").
		Expect().
		Status(http.StatusOK).JSON().Array().NotEmpty()
}
