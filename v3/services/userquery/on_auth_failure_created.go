package userqueryapi

import (
	"time"

	environment "gitlab.com/wiserskills/v3/services/userquery/gen/environment"
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// OnAuthFailure is executed when the AuthFailure event is received
func (s *userquerysrvc) OnAuthFailureCreated(e *events.Event, payload interface{}) error {

	s.logger.Info("AuthFailureCreated event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of AuthFailureCreated events received
	s.metrics.IncrCounter("events_authfailurecreated_total_received")

	p := payload.(*userquery.AuthFailure)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(e.OrgID)

	if err != nil {
		return err
	}

	user, fromCache, _ := db.GetUserByLogin(p.Login)

	if user != nil {

		if fromCache {
			s.logger.Info("User was retrieved from cache")
		}

		maxAttempts := environment.GetLoginMaxAttempt()
		delay := environment.GetLoginSuspensionTime()

		var counter int

		if user.FailedAttempts != nil {

			counter = *user.FailedAttempts

			if counter+1 > maxAttempts {

				timeUTC := time.Now().UTC().Add(time.Minute * time.Duration(delay))
				timeStr := timeUTC.Format(time.RFC3339)

				user.SuspendedUpTo = &timeStr

			} else {
				counter++
				*user.FailedAttempts = counter
			}
		} else {
			user.FailedAttempts = new(int)
			*user.FailedAttempts = 1
		}

		// We populate the context with data from the initial event
		ctx := events.CreateContextFromEvent(e)

		s.logger.Infof("User '%s' failed attempts = %v", user.ID, *user.FailedAttempts)

		s.PublishUserUpdated(ctx, user, true)
	}

	s.logger.Info("AuthFailureCreated event was treated.")

	return nil
}
