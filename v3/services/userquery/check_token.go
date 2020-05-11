package userqueryapi

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// Check if the passed token is valid.
func (s *userquerysrvc) CheckToken(ctx context.Context, p *userquery.TokenPayload) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CheckToken")
	defer span.Finish()

	// Instrumenting the endpoint with prometheus metrics
	defer s.metrics.EndpointStarted("CheckToken").Dec()

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(p.OrgID)

	if err != nil {
		s.logger.Error(err)
		s.metrics.EndpointFailure("CheckToken", "500")
		return userquery.MakeUnexpectedError(err)
	}

	// We retrieve the token id from the context
	tokenID := s.GetTokenIDFromContext(ctx)

	s.logger.Infof("Looking for token with id: %s", tokenID)

	// We look for the token in the database
	tkn, err := db.GetToken(tokenID)

	if err != nil || tkn == nil {

		if err != nil {
			s.logger.Error(err)
		}
		s.metrics.EndpointFailure("CheckToken", "404")
		return userquery.MakeNotFound(err)
	}

	return
}
