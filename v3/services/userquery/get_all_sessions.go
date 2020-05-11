package userqueryapi

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// Returns the currently active sessions.
func (s *userquerysrvc) GetAllSessions(ctx context.Context, p *userquery.AllSessionsPayload) (res userquery.SessionCollection, view string, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetAllSessions")
	defer span.Finish()

	// Instrumenting the endpoint with prometheus metrics
	defer s.metrics.EndpointStarted("GetAllSessions").Dec()

	// We assign the passed view value or use "default"
	if p.View == nil {
		view = "default"
	} else {
		view = *p.View
	}

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(p.OrgID)

	if err != nil {
		s.logger.Error(err)
		s.metrics.EndpointFailure("GetAllSessions", "500")
		err = userquery.MakeUnexpectedError(err)
		return
	}

	// We call the database
	res, err = db.GetSessionsByOrg(p.OrgID, p.AreaID)

	if err != nil {
		s.logger.Error(err)
		s.metrics.EndpointFailure("GetAllSessions", "404")
		err = userquery.NotFound("no session found")
		return
	}

	// We marked the call as successful
	s.metrics.EndpointSuccess("GetAllSessions")

	return
}
