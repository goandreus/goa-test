package userqueryapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestGETHealth(t *testing.T) {

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	obj := e.GET("/health").
		Expect().
		Status(http.StatusOK).JSON().Object()

	obj.Keys().ContainsOnly("status")
	obj.Value("status").Equal("OK")

}
