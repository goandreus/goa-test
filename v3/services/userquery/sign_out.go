package userqueryapi

import (
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// Signs a user out.
func (s *userquerysrvc) SignOut(ctx context.Context, p *userquery.TokenPayload) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SignOut")
	defer span.Finish()

	// Instrumenting the endpoint with prometheus metrics
	defer s.metrics.EndpointStarted("SignOut").Dec()

	// We retrieve the user from the context
	user := s.GetUserFromContext(ctx)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(p.OrgID)

	if err != nil {
		s.logger.Error(err)
		s.metrics.EndpointFailure("SignIn", "500")
		return userquery.MakeUnexpectedError(err)
	}

	if user != nil {

		// We get the token
		tkn, _ := s.getToken(db, user.ID)

		// We retrieve the token and revoke it
		if tkn != nil {
			s.PublishTokenRevoked(ctx, tkn, false)
		}

		// We retrieve the session and expires it
		session, fromCache, _ := db.GetSessionByUserID(user.ID)

		if session != nil {

			if fromCache {
				s.logger.Info("Session was retrieved from cache")
			}

			nowUTC := time.Now().UTC()
			nowStr := nowUTC.Format(time.RFC3339)

			// we publish the event to force the expiration of the session
			session.UpdatedAt = &nowStr
			session.ExpiresAt = &nowStr

			s.PublishSessionExpired(ctx, session, false)
		}
	}

	s.metrics.EndpointSuccess("SignOut")

	return
}
