// Code generated by goa v3.0.9, DO NOT EDIT.
//
// userquery HTTP client CLI support package
//
// Command:
// $ goa gen gitlab.com/wiserskills/v3/services/userquery/design

package client

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// BuildSignInPayload builds the payload for the userquery SignIn endpoint from
// CLI flags.
func BuildSignInPayload(userquerySignInKey string, userquerySignInOrgID string, userquerySignInAreaID string, userquerySignInUsername string, userquerySignInPassword string) (*userquery.SignInPayload, error) {
	var key *string
	{
		if userquerySignInKey != "" {
			key = &userquerySignInKey
		}
	}
	var orgID string
	{
		orgID = userquerySignInOrgID
	}
	var areaID string
	{
		areaID = userquerySignInAreaID
	}
	var username string
	{
		username = userquerySignInUsername
	}
	var password string
	{
		password = userquerySignInPassword
	}
	payload := &userquery.SignInPayload{
		Key:      key,
		OrgID:    orgID,
		AreaID:   areaID,
		Username: username,
		Password: password,
	}
	return payload, nil
}

// BuildSignOutPayload builds the payload for the userquery SignOut endpoint
// from CLI flags.
func BuildSignOutPayload(userquerySignOutToken string, userquerySignOutKey string, userquerySignOutOrgID string, userquerySignOutAreaID string) (*userquery.TokenPayload, error) {
	var token string
	{
		token = userquerySignOutToken
	}
	var key string
	{
		key = userquerySignOutKey
	}
	var orgID string
	{
		orgID = userquerySignOutOrgID
	}
	var areaID string
	{
		areaID = userquerySignOutAreaID
	}
	payload := &userquery.TokenPayload{
		Token:  token,
		Key:    key,
		OrgID:  orgID,
		AreaID: areaID,
	}
	return payload, nil
}

// BuildGetAllSessionsPayload builds the payload for the userquery
// GetAllSessions endpoint from CLI flags.
func BuildGetAllSessionsPayload(userqueryGetAllSessionsView string, userqueryGetAllSessionsToken string, userqueryGetAllSessionsKey string, userqueryGetAllSessionsOrgID string, userqueryGetAllSessionsAreaID string) (*userquery.AllSessionsPayload, error) {
	var view *string
	{
		if userqueryGetAllSessionsView != "" {
			view = &userqueryGetAllSessionsView
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
	var orgID string
	{
		orgID = userqueryGetAllSessionsOrgID
	}
	var areaID string
	{
		areaID = userqueryGetAllSessionsAreaID
	}
	payload := &userquery.AllSessionsPayload{
		View:   view,
		Token:  token,
		Key:    key,
		OrgID:  orgID,
		AreaID: areaID,
	}
	return payload, nil
}

// BuildGetIDPURLPayload builds the payload for the userquery GetIDPURL
// endpoint from CLI flags.
func BuildGetIDPURLPayload(userqueryGetIDPURLKey string, userqueryGetIDPURLOrgID string, userqueryGetIDPURLAreaID string, userqueryGetIDPURLHost string) (*userquery.HostPayload, error) {
	var key string
	{
		key = userqueryGetIDPURLKey
	}
	var orgID string
	{
		orgID = userqueryGetIDPURLOrgID
	}
	var areaID string
	{
		areaID = userqueryGetIDPURLAreaID
	}
	var host string
	{
		host = userqueryGetIDPURLHost
	}
	payload := &userquery.HostPayload{
		Key:    key,
		OrgID:  orgID,
		AreaID: areaID,
		Host:   host,
	}
	return payload, nil
}

// BuildCheckTokenPayload builds the payload for the userquery CheckToken
// endpoint from CLI flags.
func BuildCheckTokenPayload(userqueryCheckTokenToken string, userqueryCheckTokenKey string, userqueryCheckTokenOrgID string, userqueryCheckTokenAreaID string) (*userquery.TokenPayload, error) {
	var token string
	{
		token = userqueryCheckTokenToken
	}
	var key string
	{
		key = userqueryCheckTokenKey
	}
	var orgID string
	{
		orgID = userqueryCheckTokenOrgID
	}
	var areaID string
	{
		areaID = userqueryCheckTokenAreaID
	}
	payload := &userquery.TokenPayload{
		Token:  token,
		Key:    key,
		OrgID:  orgID,
		AreaID: areaID,
	}
	return payload, nil
}

// BuildGetUsersByIDPayload builds the payload for the userquery GetUsersByID
// endpoint from CLI flags.
func BuildGetUsersByIDPayload(userqueryGetUsersByIDBody string, userqueryGetUsersByIDView string, userqueryGetUsersByIDToken string, userqueryGetUsersByIDKey string, userqueryGetUsersByIDOrgID string, userqueryGetUsersByIDAreaID string) (*userquery.ManyUserIDPayload, error) {
	var err error
	var body []string
	{
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal([]byte(userqueryGetUsersByIDBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'[\n      \"Ullam error rerum et.\",\n      \"Aut dignissimos enim dolor magnam.\",\n      \"Aspernatur hic ad nihil illum alias.\",\n      \"Animi laboriosam ad et tempora magnam corporis.\"\n   ]'")
		}
	}
	var view string
	{
		if userqueryGetUsersByIDView != "" {
			view = userqueryGetUsersByIDView
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
	var orgID string
	{
		orgID = userqueryGetUsersByIDOrgID
	}
	var areaID string
	{
		areaID = userqueryGetUsersByIDAreaID
	}
	v := make([]string, len(body))
	for i, val := range body {
		v[i] = val
	}
	res := &userquery.ManyUserIDPayload{
		Ids: v,
	}
	res.View = view
	res.Token = token
	res.Key = key
	res.OrgID = orgID
	res.AreaID = areaID
	return res, nil
}
