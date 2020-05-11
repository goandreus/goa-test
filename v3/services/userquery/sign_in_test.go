package userqueryapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/gavv/httpexpect"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gitlab.com/wiserskills/v3/services/userquery/gen/environment"
	userquerypb "gitlab.com/wiserskills/v3/services/userquery/gen/grpc/userquery/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestPostSignInFail(t *testing.T) {

	os.Setenv("TOKEN_ACTIVE", "true")

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	orgID := "wiserskills"

	e.POST("/signin").
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", orgID).
		WithHeader("X-AREA-ID", "DEV").
		WithBasicAuth("mpolo", "titi").
		Expect().
		Status(http.StatusForbidden)

	// We wait for the event bus
	time.Sleep(time.Second * 3)
}

func TestPostSignInSuccess(t *testing.T) {

	os.Setenv("TOKEN_ACTIVE", "true")

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	orgID := "wiserskills"

	e.POST("/signin").
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", orgID).
		WithHeader("X-AREA-ID", "DEV").
		WithBasicAuth("mpolo", "mysecret").
		Expect().
		Status(http.StatusOK)

	// We wait for the event bus
	time.Sleep(time.Second * 3)
}

func TestPostSignInPasswordExpired(t *testing.T) {

	os.Setenv("TOKEN_ACTIVE", "true")

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	orgID := "wiserskills"

	e.POST("/signin").
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", orgID).
		WithHeader("X-AREA-ID", "DEV").
		WithBasicAuth("jvaljean", "mysecret").
		Expect().
		Status(http.StatusForbidden)

	// We wait for the event bus
	time.Sleep(time.Second * 3)
}

func TestMuxMatching(t *testing.T) {

	r := mux.NewRouter()

	r.HandleFunc("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println("Handler was called")
	}).Methods("GET")

	rq := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/products/1234?view=default"},
	}
	var match mux.RouteMatch
	found := r.Match(rq, &match)

	assert.True(t, found)
}

func TestPostSignInFailgRPC(t *testing.T) {
	//server := httptest.NewServer(GetServerHandle())
	//defer server.Close()
	//
	//e := httpexpect.New(t, server.URL)
	//
	//e.String("")

	ctx,c,err := GRPCConnect("localhost:8090", "mpolo" , "xx", environment.GetApiKey() )
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := c.SignIn(ctx, &userquerypb.SignInRequest{
		OrgId:  "wiserskills",
		AreaId: "DEV",
	})
	if err != nil {
		fmt.Println(err.Error())
		//t.Error(err.Error())
	}
	//if resp != nil {
		fmt.Println(resp)
		//t.Error(err.Error())
	//}
	//if err == nil {
	//	fmt.Println("Ok")
	//}
	//if err == nil {
	//	log.Fatalf("error %v",err)
		//t.Error(err.Error())
	//}
	//fmt.Println(resp.GetToken())

}

func TestPostSignInSuccessgRPC(t *testing.T) {

	//server := httptest.NewServer(GetServerHandle())
	//defer server.Close()
	//
	//e := httpexpect.New(t, server.URL)
	//
	//e.String("")

	ctx,c,err := GRPCConnect("localhost:8090", "mpolo" , "mysecret", environment.GetApiKey() )
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := c.SignIn(ctx,&userquerypb.SignInRequest{
		OrgId:                "wiserskills",
		AreaId:               "DEV",
	})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(resp.GetToken())
	if resp.GetToken() == "" {
		t.Error("Token no fue recibido !!")
	}

}

func GRPCConnect( url, username, password, authorization string ) ( context.Context, userquerypb.UserqueryClient, error ) {

	ctx := context.Background()

	cc, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return ctx,nil,err
	}

	md := metadata.Pairs("username", username ,"password" ,password, "authorization", authorization )

	ctx = metadata.NewOutgoingContext(ctx, md)

	md, ok :=  metadata.FromOutgoingContext(ctx)
	if !ok {
		return ctx,nil,errors.New("cant convert metada")
	}
	usernameMd := md.Get("username")
	if usernameMd[0] != username {
		return ctx,nil,errors.New("error get values in metadata ")
	}

	return ctx,userquerypb.NewUserqueryClient(cc),nil
}


