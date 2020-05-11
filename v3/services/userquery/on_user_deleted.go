package userqueryapi

import (
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// OnUser is executed when the User event is received
func (s *userquerysrvc) OnUserDeleted(e *events.Event, payload interface{}) error {

	s.logger.Info("UserDeleted event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of UserDeleted events received
	s.metrics.IncrCounter("events_userdeleted_total_received")

	// We cast the event payload
	user := payload.(*userquery.User)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(e.OrgID)

	if err != nil {
		return err
	}

	// and delete from the cache
	cacheKey := db.GetUserCacheKey(user)
	db.DeleteFromCache(cacheKey)

	s.logger.Info("UserDeleted event was treated.")

	return nil
}
