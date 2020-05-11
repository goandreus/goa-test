package userqueryapi

import (
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// OnSession is executed when the Session event is received
func (s *userquerysrvc) OnSessionExpired(e *events.Event, payload interface{}) error {

	s.logger.Info("SessionExpired event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of SessionExpired events received
	s.metrics.IncrCounter("events_sessionexpired_total_received")

	session := payload.(*userquery.Session)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(e.OrgID)

	if err != nil {
		return err
	}

	// We delete the session from the database
	db.DeleteSession(session)

	// and from the cache
	cacheKey := db.GetSessionCacheKey(session)
	db.DeleteFromCache(cacheKey)

	s.logger.Info("SessionExpired event was treated.")

	return nil
}
