package userqueryapi

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.com/wiserskills/v3/services/userquery/gen/constants"
	"gitlab.com/wiserskills/v3/services/userquery/gen/environment"
	secutil "gitlab.com/wiserskills/v3/services/userquery/gen/security"
	userquerysvc "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
	"goa.design/goa/v3/security"
)

// GetUserFromContext returns the user from the call context
func (s *userquerysrvc) GetUserFromContext(ctx context.Context) *secutil.User {
	return ctx.Value(constants.ContextKeyUser).(*secutil.User)
}

// GetTokenIDFromContext returns the id of the token from the call context
func (s *userquerysrvc) GetTokenIDFromContext(ctx context.Context) string {
	return ctx.Value(constants.ContextKeyTokenID).(string)
}

// BasicAuth implements the authorization logic for service "userquery" for the
// "basic" security scheme.
func (s *userquerysrvc) BasicAuth(ctx context.Context, user, pass string, scheme *security.BasicScheme) (context.Context, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "userquery.BasicAuth")
	defer span.Finish()

	// We check the credentials
	if len(user) == 0 || len(strings.TrimSpace(user)) == 0 {
		s.metrics.EndpointFailure("BasicAuth", "400")
		return ctx, userquerysvc.BadArgument("username cannot be empty")
	}

	if len(pass) == 0 || len(strings.TrimSpace(pass)) == 0 {
		s.metrics.EndpointFailure("BasicAuth", "400")
		return ctx, userquerysvc.BadArgument("password cannot be empty")
	}

	// We ensure that the OrgID was provided
	if ctx.Value(constants.ContextKeyOrgID) == nil {
		return ctx, userquerysvc.BadArgument("OrgID was not provided")
	}

	// We retrieve the orgID and the areaID from the request headers
	orgID, _ := ctx.Value(constants.ContextKeyOrgID).(string)
	areaID, _ := ctx.Value(constants.ContextKeyAreaID).(string)

	s.logger.Infof("Looking for user: '%s' in org: '%s'", user, orgID)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(orgID)

	if err != nil {
		s.logger.Error(err)
		s.metrics.EndpointFailure("BasicAuth", "500")
		return ctx, userquerysvc.MakeUnexpectedError(err)
	}

	// We look for the user in the database
	u, err := s.getUser(ctx, db, user, pass, orgID, areaID)

	if err != nil {
		s.metrics.EndpointFailure("BasicAuth", "401")
		return ctx, userquerysvc.MakeNotAuthorized(err)
	}

	// Login was successful, we reset the failed attempts
	if u.FailedAttempts != nil {
		u.FailedAttempts = nil
		s.PublishUserUpdated(ctx, u, true)
	}

	// We test if the user password has expired or not
	if u.PasswordExpiresAt != nil {

		nowUTC := time.Now().UTC()
		nowStr := nowUTC.Format(time.RFC3339)
		expirationDate, err := time.Parse(time.RFC3339, *u.PasswordExpiresAt)

		if err == nil && nowUTC.After(expirationDate) {

			s.logger.Info("Password has expired.")

			id, _ := uuid.NewV4()

			pinfo := &userquerysvc.Password{
				ID:                id.String(),
				UserID:            u.ID,
				EncryptedPassword: u.EncryptedPassword,
				CreatedAt:         &nowStr,
			}

			s.PublishPasswordExpired(ctx, pinfo, true)

			return ctx, userquerysvc.MakePasswordExpired(errors.New("password has expired"))
		}
	}

	ctx = context.WithValue(ctx, constants.ContextKeyUser, u)

	return ctx, nil
}

// APIKeyAuth implements the authorization logic for service "userquery" for
// the "api_key" security scheme.
func (s *userquerysrvc) APIKeyAuth(ctx context.Context, key string, scheme *security.APIKeyScheme) (context.Context, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "userquery.APIKeyAuth")
	defer span.Finish()

	if key != environment.GetApiKey() {
		return ctx, userquerysvc.BadArgument("wrong API key")
	}

	return ctx, nil
}

// JWTAuth implements the authorization logic for service "userquery" for the
// "jwt" security scheme.
func (s *userquerysrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "userquery.JWTAuth")
	defer span.Finish()

	if environment.GetTokenActive() {

		// We uncrypt the token to retrieve the user
		signingKey := []byte(environment.GetTokenKey())
		claims, err := secutil.ParseToken(token, signingKey)

		if err != nil || claims == nil || claims.User == nil {
			return ctx, userquerysvc.BadArgument("invalid token")
		}

		// We validate if the token is not expired
		if claims.Valid() != nil {
			return ctx, userquerysvc.BadArgument("expired token")
		}

		// We ensure that the OrgID was provided
		if ctx.Value(constants.ContextKeyOrgID) == nil {
			return ctx, userquerysvc.BadArgument("OrgID was not provided")
		}

		// We parse the required scopes and validate the provided one
		orgID, _ := ctx.Value(constants.ContextKeyOrgID).(string)

		for i, scope := range scheme.RequiredScopes {
			scheme.RequiredScopes[i] = strings.Replace(scope, "[ORGID]", orgID, -1)
		}

		scopes := []string{claims.Scope}
		err = scheme.Validate(scopes)

		if err != nil {
			return ctx, userquerysvc.MakeInvalidScope(errors.New("invalid token scope"))
		}

		// We add the user/userId to the context
		ctx = context.WithValue(ctx, constants.ContextKeyTokenID, claims.StandardClaims.Id)
		ctx = context.WithValue(ctx, constants.ContextKeyUser, claims.User)
		ctx = context.WithValue(ctx, constants.ContextKeyUserID, claims.User.ID)
	}

	return ctx, nil
}
