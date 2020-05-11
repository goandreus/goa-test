package userqueryapi

import (
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// OnUser is executed when the User event is received
func (s *userquerysrvc) OnUserUpdated(e *events.Event, payload interface{}) error {

	s.logger.Info("UserUpdated event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of UserUpdated events received
	s.metrics.IncrCounter("events_userupdated_total_received")

	user := payload.(*userquery.User)

	if e.OrgID != "" {

		// We retrieve the database for the passed organization id
		db, err := s.store.GetDatabase(e.OrgID)

		if err != nil {
			return err
		}

		s.logger.Info("adding to cache...")

		// and add the user to the cache
		cacheKey := db.GetUserCacheKey(user)
		db.AddToCache(cacheKey, *user)
	}

	s.logger.Info("UserUpdated event was treated.")

	return nil
}
