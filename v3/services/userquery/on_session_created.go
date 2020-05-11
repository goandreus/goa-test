package userqueryapi

import (
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// OnSession is executed when the Session event is received
func (s *userquerysrvc) OnSessionCreated(e *events.Event, payload interface{}) error {

	s.logger.Info("SessionCreated event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of SessionCreated events received
	s.metrics.IncrCounter("events_sessioncreated_total_received")

	session := payload.(*userquery.Session)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(e.OrgID)

	if err != nil {
		return err
	}

	// We save it to the database
	db.CreateSession(session)

	// and add it to cache
	cacheKey := db.GetSessionCacheKey(session)
	db.AddToCache(cacheKey, *session)

	s.logger.Info("SessionCreated event was treated.")

	return nil
}
