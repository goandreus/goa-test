package userqueryapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestGETMetrics(t *testing.T) {

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/metrics").
		Expect().
		Status(http.StatusOK)
}
