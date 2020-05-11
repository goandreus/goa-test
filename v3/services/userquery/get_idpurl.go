package userqueryapi

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
	security "gitlab.com/wiserskills/v3/services/userquery/security"
)

// Returns the URL of the IDP to redirect the user to.
func (s *userquerysrvc) GetIDPURL(ctx context.Context, p *userquery.HostPayload) (res *userquery.RedirectResult, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetIDPURL")
	defer span.Finish()

	// Instrumenting the endpoint with prometheus metrics
	defer s.metrics.EndpointStarted("GetIDPURL").Dec()

	// We retrieve the SAML config from the store
	s.logger.Info("Searching SAML Config for host: ", p.Host)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(p.OrgID)

	if err != nil {
		s.logger.Error(err)
		s.metrics.EndpointFailure("GetIDPURL", "500")
		return nil, userquery.MakeUnexpectedError(err)
	}

	samlConfig, fromCache, err := db.GetSAMLConfigByHost(p.Host)

	if fromCache {
		s.logger.Info("SAMLConfig was retrieved from cache")
	}

	if err != nil || samlConfig == nil {
		s.metrics.EndpointFailure("GetIDPURL", "500")
		s.logger.Warn("SAML configuration not found for host: ", p.Host)
		return res, userquery.InternalServerError("SAML configuration not found for specified host")
	}

	url, err2 := security.GetURLToAuthenticate(samlConfig)

	if err2 != nil {
		s.logger.Error(err2.Error())
		s.metrics.EndpointFailure("GetIDPURL", "500")
		return res, userquery.InternalServerError("error ")
	}

	res = &userquery.RedirectResult{Location: url}

	s.metrics.EndpointSuccess("GetIDPURL")

	return
}
