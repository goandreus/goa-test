package userqueryapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"gitlab.com/wiserskills/v3/services/userquery/gen/environment"
)

func TestGETGetIDPURL(t *testing.T) {

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	orgID := "wiserskills"

	location := e.GET("/saml").
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", orgID).
		WithHeader("X-AREA-ID", "DEV").
		WithHeader("Host", "wiserskills.io").
		Expect().
		Status(http.StatusPermanentRedirect).Header("Location").Raw()

	fmt.Println("Location: " + location)
}
