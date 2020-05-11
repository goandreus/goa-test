// Code generated by goa v3.0.9, DO NOT EDIT.
//
// userquery gRPC client CLI support package
//
// Command:
// $ goa gen gitlab.com/wiserskills/v3/services/userquery/design

package client

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	userquerypb "gitlab.com/wiserskills/v3/services/userquery/gen/grpc/userquery/pb"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// BuildSignInPayload builds the payload for the userquery SignIn endpoint from
// CLI flags.
func BuildSignInPayload(userquerySignInMessage string, userquerySignInUsername string, userquerySignInPassword string, userquerySignInKey string) (*userquery.SignInPayload, error) {
	var err error
	var message userquerypb.SignInRequest
	{
		if userquerySignInMessage != "" {
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err = json.Unmarshal([]byte(userquerySignInMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, example of valid JSON:\n%s", "'{\n      \"areaId\": \"PROD\",\n      \"orgId\": \"WiserSKILLS\"\n   }'")
			}
		}
	}
	var username string
	{
		username = userquerySignInUsername
	}
	var password string
	{
		password = userquerySignInPassword
	}
	var key *string
	{
		if userquerySignInKey != "" {
			key = &userquerySignInKey
		}
	}
	v := &userquery.SignInPayload{
		OrgID:  message.OrgId,
		AreaID: message.AreaId,
	}
	v.Username = username
	v.Password = password
	v.Key = key
	return v, nil
}

// BuildSignOutPayload builds the payload for the userquery SignOut endpoint
// from CLI flags.
func BuildSignOutPayload(userquerySignOutMessage string, userquerySignOutToken string, userquerySignOutKey string) (*userquery.TokenPayload, error) {
	var err error
	var message userquerypb.SignOutRequest
	{
		if userquerySignOutMessage != "" {
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err = json.Unmarshal([]byte(userquerySignOutMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, example of valid JSON:\n%s", "'{\n      \"areaId\": \"ss\",\n      \"orgId\": \"fka\"\n   }'")
			}
		}
	}
	var token string
	{
		token = userquerySignOutToken
	}
	var key string
	{
		key = userquerySignOutKey
	}
	v := &userquery.TokenPayload{
		OrgID:  message.OrgId,
		AreaID: message.AreaId,
	}
	v.Token = token
	v.Key = key
	return v, nil
}

// BuildGetAllSessionsPayload builds the payload for the userquery
// GetAllSessions endpoint from CLI flags.
func BuildGetAllSessionsPayload(userqueryGetAllSessionsMessage string, userqueryGetAllSessionsToken string, userqueryGetAllSessionsKey string) (*userquery.AllSessionsPayload, error) {
	var err error
	var message userquerypb.GetAllSessionsRequest
	{
		if userqueryGetAllSessionsMessage != "" {
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err = json.Unmarshal([]byte(userqueryGetAllSessionsMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, example of valid JSON:\n%s", "'{\n      \"areaId\": \"PROD\",\n      \"orgId\": \"65e\",\n      \"view\": \"tiny\"\n   }'")
			}
		}
	}
	var token string
	{
		token = userqueryGetAllSessionsToken
	}
	var key string
	{
		key = userqueryGetAllSessionsKey
	}
	v := &userquery.AllSessionsPayload{
		OrgID:  message.OrgId,
		AreaID: message.AreaId,
	}
	if message.View != "" {
		v.View = &message.View
	}
	v.Token = token
	v.Key = key
	return v, nil
}

// BuildGetIDPURLPayload builds the payload for the userquery GetIDPURL
// endpoint from CLI flags.
func BuildGetIDPURLPayload(userqueryGetIDPURLMessage string, userqueryGetIDPURLKey string) (*userquery.HostPayload, error) {
	var err error
	var message userquerypb.GetIDPURLRequest
	{
		if userqueryGetIDPURLMessage != "" {
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err = json.Unmarshal([]byte(userqueryGetIDPURLMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, example of valid JSON:\n%s", "'{\n      \"areaId\": \"3zr\",\n      \"host\": \"Voluptas eaque perferendis accusantium.\",\n      \"orgId\": \"jzo\"\n   }'")
			}
		}
	}
	var key string
	{
		key = userqueryGetIDPURLKey
	}
	v := &userquery.HostPayload{
		Host:   message.Host,
		OrgID:  message.OrgId,
		AreaID: message.AreaId,
	}
	v.Key = key
	return v, nil
}

// BuildCheckTokenPayload builds the payload for the userquery CheckToken
// endpoint from CLI flags.
func BuildCheckTokenPayload(userqueryCheckTokenMessage string, userqueryCheckTokenToken string, userqueryCheckTokenKey string) (*userquery.TokenPayload, error) {
	var err error
	var message userquerypb.CheckTokenRequest
	{
		if userqueryCheckTokenMessage != "" {
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err = json.Unmarshal([]byte(userqueryCheckTokenMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, example of valid JSON:\n%s", "'{\n      \"areaId\": \"ad\",\n      \"orgId\": \"gaz\"\n   }'")
			}
		}
	}
	var token string
	{
		token = userqueryCheckTokenToken
	}
	var key string
	{
		key = userqueryCheckTokenKey
	}
	v := &userquery.TokenPayload{
		OrgID:  message.OrgId,
		AreaID: message.AreaId,
	}
	v.Token = token
	v.Key = key
	return v, nil
}

// BuildGetUsersByIDPayload builds the payload for the userquery GetUsersByID
// endpoint from CLI flags.
func BuildGetUsersByIDPayload(userqueryGetUsersByIDMessage string, userqueryGetUsersByIDToken string, userqueryGetUsersByIDKey string) (*userquery.ManyUserIDPayload, error) {
	var err error
	var message userquerypb.GetUsersByIDRequest
	{
		if userqueryGetUsersByIDMessage != "" {
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err = json.Unmarshal([]byte(userqueryGetUsersByIDMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, example of valid JSON:\n%s", "'{\n      \"activeOnly\": true,\n      \"areaId\": \"8w8\",\n      \"ids\": [\n         \"Ab modi magni saepe quia.\",\n         \"Ratione voluptatem accusamus atque nostrum.\",\n         \"Quae velit aliquam fugiat et ipsa.\"\n      ],\n      \"orgId\": \"bu3\",\n      \"view\": \"Adipisci rerum adipisci.\"\n   }'")
			}
		}
	}
	var token string
	{
		token = userqueryGetUsersByIDToken
	}
	var key string
	{
		key = userqueryGetUsersByIDKey
	}
	v := &userquery.ManyUserIDPayload{
		OrgID:      message.OrgId,
		AreaID:     message.AreaId,
		View:       message.View,
		ActiveOnly: message.ActiveOnly,
	}
	if message.Ids != nil {
		v.Ids = make([]string, len(message.Ids))
		for i, val := range message.Ids {
			v.Ids[i] = val
		}
	}
	if message.View == "" {
		v.View = "default"
	}
	v.Token = token
	v.Key = key
	return v, nil
}
