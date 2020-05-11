package userqueryapi

import (
	"context"
	"fmt"
	"net/url"

	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.com/wiserskills/v3/services/userquery/gen/constants"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
	security "gitlab.com/wiserskills/v3/services/userquery/security"
)

// Call back endpoint called by the IDP once the user is authenticated.
func (s *userquerysrvc) SamlSignIn(ctx context.Context, p string) (res *userquery.RedirectResult, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SamlSignIn")
	defer span.Finish()

	// Instrumenting the endpoint with prometheus metrics
	defer s.metrics.EndpointStarted("SamlSignIn").Dec()

	m, _ := url.ParseQuery(p)

	relayState := m["RelayState"][0]
	samlResponse := m["SAMLResponse"][0]

	if relayState == "" {
		s.logger.Error("RelayState parameter cannot be empty")
		s.metrics.EndpointFailure("GetIDPURL", "500")
		return res, userquery.InternalServerError("RelayState parameter cannot be empty")
	}

	if samlResponse == "" {
		s.logger.Error("SAMLResponse parameter cannot be empty")
		s.metrics.EndpointFailure("GetIDPURL", "500")
		return res, userquery.InternalServerError("SAMLResponse parameter cannot be empty")
	}

	// We retrieve the orgId from the context
	orgID, _ := ctx.Value(constants.ContextKeyOrgID).(string)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(orgID)

	if err != nil {
		s.logger.Error(err)
		s.metrics.EndpointFailure("EnableUserByEmail", "500")
		return nil, userquery.MakeUnexpectedError(err)
	}

	// We retrieve the SAML config from the store
	samlConfig, fromCache, err := db.GetSAMLConfigByHost(relayState)

	if fromCache {
		s.logger.Info("SAMLConfig was retrieved from cache")
	}

	if err != nil || samlConfig == nil {
		s.logger.Error("SAML configuration not found for host: ", relayState)
		s.metrics.EndpointFailure("GetIDPURL", "500")
		return res, userquery.InternalServerError("SAML configuration not found for specified host")
	}

	sp, err := security.GetServiceProvider(samlConfig)

	if err != nil {
		s.logger.Error(err.Error())
		s.metrics.EndpointFailure("SamlSignIn", "500")
		return res, userquery.InternalServerError("error ")
	}

	username, err := security.GetValueFromSAMLResponse(sp, samlResponse, samlConfig.IDKey)

	if err != nil {
		s.logger.Error(err.Error())
		s.metrics.EndpointFailure("SamlSignIn", "500")
		return res, userquery.InternalServerError("error ")
	}

	// We retrieve the user associated with the username
	user, err := s.getUserByLogin(ctx, db, username, samlConfig.OrganizationID, samlConfig.AreaID)

	if err != nil {
		s.logger.Error(err.Error())
		s.metrics.EndpointFailure("SamlSignIn", "500")
		return res, userquery.InternalServerError("error ")
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
				s.metrics.EndpointFailure("SamlSignIn", "500")
				return nil, userquery.MakeUnexpectedError(err)
			}

		} else {

			// We create a new session
			session = s.createNewSession(user.ID, tkn.ID, samlConfig.OrganizationID, samlConfig.AreaID)

			// We publish the event to the bus
			// The session will be saved in the handler
			err := s.PublishSessionCreated(ctx, session, false)

			if err != nil {
				s.logger.Error(err)
				s.metrics.EndpointFailure("SamlSignIn", "500")
				return nil, userquery.MakeUnexpectedError(err)
			}

			s.logger.Info("A new session has been created")
		}

	} else {

		s.logger.Info("No active token found")

		// We create a new token
		tkn, err = s.createNewToken(db, user, samlConfig.OrganizationID, samlConfig.AreaID)

		if err != nil {
			s.logger.Error(err)
			s.metrics.EndpointFailure("SamlSignIn", "500")
			return nil, userquery.MakeUnexpectedError(err)
		}

		// We publish the event to the bus
		// The token will be saved in the handler
		err := s.PublishTokenCreated(ctx, tkn, false)

		if err != nil {
			s.logger.Error(err)
			s.metrics.EndpointFailure("SamlSignIn", "500")
			return nil, userquery.MakeUnexpectedError(err)
		}

		// We create a new session
		session := s.createNewSession(user.ID, tkn.ID, samlConfig.OrganizationID, samlConfig.AreaID)

		// We publish the event to the bus
		// The session will be saved in the handler
		err = s.PublishSessionCreated(ctx, session, false)

		if err != nil {
			s.logger.Error(err)
			s.metrics.EndpointFailure("SamlSignIn", "500")
			return nil, userquery.MakeUnexpectedError(err)
		}

		s.logger.Info("A new token and a new session have been created")
	}

	// We return the redirect url
	url := fmt.Sprintf(samlConfig.RedirectURL, tkn.Token)
	res = &userquery.RedirectResult{Location: url}

	s.metrics.EndpointSuccess("SamlSignIn")

	return
}
