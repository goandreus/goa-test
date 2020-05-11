package userqueryapi

import (
	"fmt"
	"github.com/gavv/httpexpect"
	"gitlab.com/wiserskills/v3/services/userquery/gen/environment"
	userquerypb "gitlab.com/wiserskills/v3/services/userquery/gen/grpc/userquery/pb"
	"gitlab.com/wiserskills/v3/services/userquery/security"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestPOSTSignOut(t *testing.T) {

	os.Setenv("TOKEN_ACTIVE", "true")

	// run server using httptest
	server := httptest.NewServer(GetServerHandle())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// We first sign in to get a token
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

	// We wait for the event bus (the Token is saved asynchronously)
	time.Sleep(time.Second * 3)

	// We then sign out
	e.POST("/signout").
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", orgID).
		WithHeader("X-AREA-ID", "DEV").
		WithHeader("Authorization", "Bearer "+tkn).
		Expect().
		Status(http.StatusOK)

	// We wait for the event bus (the Token is revoked asynchronously)
	time.Sleep(time.Second * 3)

	// and check the token
	e.POST("/check").
		WithHeader("X-API-KEY", environment.GetApiKey()).
		WithHeader("X-ORG-ID", orgID).
		WithHeader("X-AREA-ID", "DEV").
		WithHeader("Authorization", "Bearer "+tkn).
		Expect().
		Status(http.StatusNotFound)
}

func TestPOSTSignOutgRPC(t *testing.T) {

	ctx,c,err := GRPCConnect("localhost:8090", "mpolo" , "mysecret", environment.GetApiKey() )
	if err != nil {
		t.Error(err.Error())
	}

	//resp, err := c.SignIn(ctx, &userquerypb.SignInRequest{
	//	OrgId:  "wiserskills",
	//	AreaId: "DEV",
	//})
	//if err != nil {
	//	t.Error(err.Error())
	//}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjNmMGIzYzUyLWM5NTMtNDFhZi1iYTY1LTc5ZGI4M2RlNTE5MyIsImxvZ2luIjoibXBvbG8iLCJmaXN0bmFtZSI6Ik1hcmNvIiwibGFzdG5hbWUiOiJQb2xvIiwiZW1haWwiOiJtYXJjby5wb2xvQGNoaW5hLmNvbSIsImxhbmd1YWdlIjoiIiwicm9sZXMiOm51bGwsImV4cCI6MTU4MDU5MjUwNiwianRpIjoiYjZiMDhhM2MtZGFiYy00YmUwLWFiMTUtOWEzZDZmMzk4YmZiIiwiaWF0IjoxNTgwNTA2MTA2LCJpc3MiOiJ3aXNlcnNraWxscyIsInN1YiI6Im1wb2xvIiwiU2NvcGUiOiJhcGk6aW50ZXJuYWwifQ.J3TOGUkR5D8HfzsgwJ-tp6WbOs1G7T7iQrHQ6L5CH6o"
	//fmt.Println(resp.Token)

	//fmt.Println(resp.Token)


	//md := metadata.Pairs( "token", token , "key" , environment.GetApiKey(), "authorization",  "Bearer "+token , "api_key",  "Bearer "+token  )
	//md := metadata.Pairs( "token", token , "key" , environment.GetApiKey(), "authorization",  "Bearer "+token , "api_key",  environment.GetApiKey()  )


	user := security.User{
		ID:        "",
		Login:     "",
		Firstname: "",
		Lastname:  "",
		Email:     "",
		Language:  "",
		Roles:     nil,
	}
	fmt.Println(user)

	//md := metadata.Pairs( "token", token , "key" , environment.GetApiKey(), "authorization", environment.GetApiKey() , "api_key",  "Bearer "+token ,"user", &user )
	md := metadata.Pairs( "token", token , "key" , environment.GetApiKey(), "authorization", environment.GetApiKey() , "api_key",  "Bearer "+token )

	ctx = metadata.NewOutgoingContext(ctx, md)

	respSignOut , err := c.SignOut(ctx,&userquerypb.SignOutRequest{
		OrgId:  "wiserskills",
		AreaId: "DEV",
	})
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(respSignOut)
}
