package userqueryapi

import (
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// OnToken is executed when the Token event is received
func (s *userquerysrvc) OnTokenCreated(e *events.Event, payload interface{}) error {

	s.logger.Info("TokenCreated event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of TokenCreated events received
	s.metrics.IncrCounter("events_tokencreated_total_received")

	token := payload.(*userquery.Token)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(e.OrgID)

	if err != nil {
		return err
	}

	// We save it to the database
	db.CreateToken(token)

	// and add it to cache
	cacheKey := db.GetTokenCacheKey(token)
	db.AddToCache(cacheKey, *token)

	s.logger.Info("TokenCreated event was treated.")

	return nil
}
