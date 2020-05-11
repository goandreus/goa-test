package userqueryapi

import (
	"gitlab.com/wiserskills/v3/services/userquery/gen/caching"
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// OnSAMLConfig is executed when the SAMLConfig event is received
func (s *userquerysrvc) OnSAMLConfigDeleted(e *events.Event, payload interface{}) error {

	s.logger.Info("SAMLConfigDeleted event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of SAMLConfigDeleted events received
	s.metrics.IncrCounter("events_samlconfigdeleted_total_received")

	config := payload.(*userquery.SAMLConfig)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(e.OrgID)

	if err != nil {
		return err
	}

	// We delete the SAML config from cache
	kb := caching.NewKeyBuilder()
	kb.Add("samlconfig")
	kb.Add(config.Host)
	cacheKey := kb.Get()

	db.DeleteFromCache(cacheKey)

	s.logger.Info("SAMLConfigDeleted event was treated.")

	return nil
}
