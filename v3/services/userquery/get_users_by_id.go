package userqueryapi

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// Returns the users with the specified ids.
func (s *userquerysrvc) GetUsersByID(ctx context.Context, p *userquery.ManyUserIDPayload) (res userquery.RegisteredUserCollection, view string, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetUsersByID")
	defer span.Finish()

	// Instrumenting the endpoint with prometheus metrics
	defer s.metrics.EndpointStarted("GetUsersByID").Dec()

	view = p.View

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(p.OrgID)

	if err != nil {
		s.logger.Error(err)
		s.metrics.EndpointFailure("GetUsersByID", "500")
		return nil, view, userquery.MakeUnexpectedError(err)
	}

	// We call the database
	res, err = db.GetRegisteredUsersByID(p.Ids, p.ActiveOnly)

	s.logger.Info("userquery.GetUsersByID")
	return
}
