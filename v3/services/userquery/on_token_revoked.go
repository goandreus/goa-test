package userqueryapi

import (
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// OnToken is executed when the Token event is received
func (s *userquerysrvc) OnTokenRevoked(e *events.Event, payload interface{}) error {

	s.logger.Info("TokenRevoked event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of TokenRevoked events received
	s.metrics.IncrCounter("events_tokenrevoked_total_received")

	token := payload.(*userquery.Token)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(e.OrgID)

	if err != nil {
		return err
	}

	// We delete the token
	s.deleteToken(db, token)

	s.logger.Info("TokenRevoked event was treated.")

	return nil
}
