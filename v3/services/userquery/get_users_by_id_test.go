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

func TestGETGetUsersByID(t *testing.T) {

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

	ids := [2]string{"3f0b3c52-c953-41af-ba65-79db83de5193", "4f0b3c52-c953-41af-ba65-79db83de5193"}

	// We wait for the event bus (the token and sessions are created asynchronously)
	time.Sleep(time.Second * 3)

	e.GET("/users/id").
		WithHeader("Authorization", "Bearer "+tkn).
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", "wiserskills").
		WithHeader("X-AREA-ID", "DEV").
		WithJSON(ids).
		Expect().
		Status(http.StatusOK).JSON().Array().Length().Equal(2)
}
