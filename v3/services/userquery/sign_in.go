package userqueryapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.com/wiserskills/v3/services/userquery/datalayer"
	"gitlab.com/wiserskills/v3/services/userquery/gen/constants"
	"gitlab.com/wiserskills/v3/services/userquery/gen/environment"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
	"gitlab.com/wiserskills/v3/services/userquery/security"
)

// Signs a user in.
func (s *userquerysrvc) SignIn(ctx context.Context, p *userquery.SignInPayload) (res *userquery.JWTToken, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SignIn")
	defer span.Finish()

	// Instrumenting the endpoint with prometheus metrics
	defer s.metrics.EndpointStarted("SignIn").Dec()

	// We retrieve the "real" user from the context
	// This is specific for the BasicAuth result.
	// In other endpoints, a secutil user (user from the token) is returned.
	user := ctx.Value(constants.ContextKeyUser).(*userquery.User)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(p.OrgID)

	if err != nil {
		s.logger.Error(err)
		s.metrics.EndpointFailure("SignIn", "500")
		return nil, userquery.MakeUnexpectedError(err)
	}

	// We search for an existing and valid token
	tkn, _ := s.getToken(db, user.ID)

	if tkn != nil {

		s.logger.Info("An existing token was found")

		// We look for an exisitng session for the existing token
		session, _ := s.getSession(db, user.ID)

		if session != nil {

			s.logger.Info("An existing session was found")

			// We extend the existing session
			session = s.extendSession(session)

			// We publish the event to the bus
			// The session will be saved in the handler
			err = s.PublishSessionUpdated(ctx, session, false)

			if err != nil {
				s.logger.Error(err)
				s.metrics.EndpointFailure("SignIn", "500")
				return nil, userquery.MakeUnexpectedError(err)
			}

			res = &userquery.JWTToken{Token: tkn.Token}

		} else {

			// We create a new session
			session = s.createNewSession(user.ID, tkn.ID, p.OrgID, p.AreaID)

			// We publish the event to the bus
			// The session will be saved in the handler
			err := s.PublishSessionCreated(ctx, session, false)

			if err != nil {
				s.logger.Error(err)
				s.metrics.EndpointFailure("SignIn", "500")
				return nil, userquery.MakeUnexpectedError(err)
			}

			s.logger.Info("A new session has been created")

			// We return the existing token
			res = &userquery.JWTToken{Token: tkn.Token}
		}

	} else {

		s.logger.Info("No active token found")

		// We create a new token
		tkn, err = s.createNewToken(db, user, p.OrgID, p.AreaID)

		if err != nil {
			s.logger.Error(err)
			s.metrics.EndpointFailure("SignIn", "500")
			return nil, userquery.MakeUnexpectedError(err)
		}

		// We publish the event to the bus
		// The token will be saved in the handler
		err := s.PublishTokenCreated(ctx, tkn, false)

		if err != nil {
			s.logger.Error(err)
			s.metrics.EndpointFailure("SignIn", "500")
			return nil, userquery.MakeUnexpectedError(err)
		}

		// We create a new session
		session := s.createNewSession(user.ID, tkn.ID, p.OrgID, p.AreaID)

		// We publish the event to the bus
		// The session will be saved in the handler
		err = s.PublishSessionCreated(ctx, session, false)

		if err != nil {
			s.logger.Error(err)
			s.metrics.EndpointFailure("SignIn", "500")
			return nil, userquery.MakeUnexpectedError(err)
		}

		s.logger.Info("A new token and a new session have been created")

		// We return the token
		res = &userquery.JWTToken{Token: tkn.Token}
	}

	s.metrics.EndpointSuccess("SignIn")

	return
}

func (s *userquerysrvc) createNewSession(userID string, tokenID string, orgID string, areaID string) *userquery.Session {

	now := time.Now().Local()
	nowUTC := now.UTC()

	duration := environment.GetSessionLifeTime()
	sessionExpiresAt := now.Add(time.Minute * time.Duration(duration)).UTC()

	nowStr := nowUTC.Format(time.RFC3339)
	sessionExpiresAtStr := sessionExpiresAt.Format(time.RFC3339)

	sessionID, _ := uuid.NewV4()
	result := &userquery.Session{
		ID:             sessionID.String(),
		UserID:         userID,
		TokenID:        tokenID,
		CreatedAt:      &nowStr,
		ExpiresAt:      &sessionExpiresAtStr,
		OrganizationID: orgID,
		AreaID:         areaID,
	}

	return result
}

func (s *userquerysrvc) createNewToken(db *datalayer.Db, user *userquery.User, orgID string, areaID string) (*userquery.Token, error) {

	signingkey := []byte(environment.GetTokenKey())

	tkUser := &security.User{
		ID:        user.ID,
		Login:     user.Login,
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		Roles:     user.Roles,
	}

	if user.LanguageID != nil {
		tkUser.Language = *user.LanguageID
	}

	var scope string

	if orgID == "" {
		scope = "api:b2c"
	} else {
		if orgID == environment.GetInternalOrgId() {
			scope = "api:internal"
		} else {
			scope = fmt.Sprintf("api:b2b:%s:%s", orgID, areaID)
		}
	}

	duration, _ := time.ParseDuration(fmt.Sprintf("%vm", environment.GetTokenLifeTime()))
	token, tokenID, err := security.CreateToken(signingkey, "wiserskills", scope, tkUser, time.Duration(duration))

	if err != nil {
		return nil, err
	}

	now := time.Now().Local()
	nowUTC := now.UTC()
	tokenExpiresAt := now.Add(time.Duration(duration)).UTC()

	nowStr := nowUTC.Format(time.RFC3339)
	tokenExpiresAtStr := tokenExpiresAt.Format(time.RFC3339)

	result := &userquery.Token{
		ID:             tokenID,
		Type:           1, // Session Token
		UserID:         user.ID,
		OrganizationID: orgID,
		AreaID:         areaID,
		Token:          token,
		CreatedAt:      &nowStr,
		ExpiresAt:      &tokenExpiresAtStr,
	}

	return result, nil
}

func (s *userquerysrvc) getUser(ctx context.Context, db *datalayer.Db, username string, pwd string, orgID string, areaID string) (*userquery.User, error) {

	result, err := s.getUserByLogin(ctx, db, username, orgID, areaID)

	if err != nil {
		return nil, err
	}

	if result == nil {
		err = s.createAuthError(ctx, username, orgID, areaID)
		return nil, err
	}

	valid, err := security.ComparePasswords([]byte(result.EncryptedPassword), []byte(pwd))

	s.logger.Infof("Password comparison result: %t", valid)

	if err != nil || valid == false {

		err = s.createAuthError(ctx, username, orgID, areaID)
		return nil, err
	}

	return result, nil
}

func (s *userquerysrvc) getUserByLogin(ctx context.Context, db *datalayer.Db, username string, orgID string, areaID string) (*userquery.User, error) {

	result, fromCache, _ := db.GetUserByLogin(username)

	if result == nil {
		err := s.createAuthError(ctx, username, orgID, areaID)
		return nil, err
	}

	if fromCache {
		s.logger.Info("User was retrieved from cache")
	}

	// We test if the user is suspended or not
	if result.SuspendedUpTo != nil {

		nowUTC := time.Now().UTC()
		expirationDate, err := time.Parse(time.RFC3339, *result.SuspendedUpTo)

		if err == nil && nowUTC.Before(expirationDate) {

			err := userquery.MakeLoginBlocked(errors.New("login is temporarly blocked due to too many failed attempts"))

			return nil, err
		}

		// We reset the SuspendedUpTo and FailedAttempts fields
		result.SuspendedUpTo = nil
		result.FailedAttempts = nil

		// And publish the user updated event
		s.PublishUserUpdated(ctx, result, true)
	}

	s.logger.Infof("User retrieved: %s", result.ID)

	return result, nil
}

func (s *userquerysrvc) createAuthError(ctx context.Context, username string, orgID string, areaID string) error {

	timeUTC := time.Now().UTC()
	timeStr := timeUTC.Format(time.RFC3339)

	f := userquery.AuthFailure{
		Login:          username,
		OrganizationID: orgID,
		AreaID:         areaID,
		CreatedAt:      timeStr,
	}

	s.PublishAuthFailureCreated(ctx, &f, true)

	return userquery.MakeInvalidCredentials(errors.New("invalid login or password"))
}

func (s *userquerysrvc) getToken(db *datalayer.Db, userID string) (*userquery.Token, error) {

	s.logger.Infof("Looking token for user: %s", userID)

	// We search for an existing and valid token
	tkn, fromCache, _ := db.GetTokenByUserID(userID)

	if tkn != nil {

		if fromCache {
			s.logger.Info("Token was retrieved from cache")
		}

		s.logger.Info("Token found. Testing expiration date...")

		// if the current token is expired, we need to remove it
		now := time.Now().UTC()
		expirationDate, err := time.Parse(time.RFC3339, *tkn.ExpiresAt)

		if err != nil {
			return nil, err
		}

		if now.After(expirationDate) {

			s.logger.Info("Token expired. Deleting token...")

			//the token is expired we need to remove it from the store
			err := db.DeleteToken(tkn)

			if err == nil {

				// We also delete it from the cache
				cacheKey := db.GetTokenCacheKey(tkn)
				db.DeleteFromCache(cacheKey)

				s.logger.Info("Token has been deleted.")
			}

			return nil, err
		}

		return tkn, nil
	}

	return nil, nil
}

func (s *userquerysrvc) getSession(db *datalayer.Db, userID string) (*userquery.Session, error) {

	// We search for an existing and valid session
	session, fromCache, err := db.GetSessionByUserID(userID)

	if err != nil {
		return nil, err
	}

	if session != nil {

		if fromCache {
			s.logger.Info("Session was retrieved from cache")
		}

		// if the current token is expired, we need to remove it
		now := time.Now().UTC()
		expirationDate, err := time.Parse(time.RFC3339, *session.ExpiresAt)

		if err != nil {
			return nil, err
		}

		if now.After(expirationDate) {
			return nil, nil
		}

		return session, nil
	}

	return nil, nil
}

func (s *userquerysrvc) extendSession(session *userquery.Session) *userquery.Session {

	// We extend the existing session
	now := time.Now().Local()
	nowUTC := now.UTC()

	duration := environment.GetSessionLifeTime()
	sessionExpiresAt := now.Add(time.Minute * time.Duration(duration)).UTC()

	nowStr := nowUTC.Format(time.RFC3339)
	sessionExpiresAtStr := sessionExpiresAt.Format(time.RFC3339)

	session.ExpiresAt = &sessionExpiresAtStr
	session.UpdatedAt = &nowStr

	return session
}

func (s *userquerysrvc) deleteToken(db *datalayer.Db, token *userquery.Token) {

	// We delete the token from the database
	db.DeleteToken(token)

	// and from the cache
	cacheKey := db.GetTokenCacheKey(token)
	db.DeleteFromCache(cacheKey)
}
