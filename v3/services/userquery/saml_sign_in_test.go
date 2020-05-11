package userqueryapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"gitlab.com/wiserskills/v3/services/userquery/gen/environment"
)

func TestPOSTSamlSignIn(t *testing.T) {

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.POST("/saml/signin").
		WithHeader("Authorization", "Bearer asdasdasdsadsa").
		WithHeader("X-API-KEY", environment.GetApiKey()).
		Expect().
		Status(http.StatusCreated)

}
