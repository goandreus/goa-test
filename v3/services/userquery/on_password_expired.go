package userqueryapi

import (
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// OnPassword is executed when the Password event is received
func (s *userquerysrvc) OnPasswordExpired(e *events.Event, payload interface{}) error {

	s.logger.Info("PasswordExpired event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of PasswordExpired events received
	s.metrics.IncrCounter("events_passwordexpired_total_received")

	p := payload.(*userquery.Password)

	// We retrieve the database for the passed organization id
	db, err := s.store.GetDatabase(e.OrgID)

	if err != nil {
		return err
	}

	// We call the store to save the received object
	err = db.CreatePassword(p)

	if err != nil {
		s.logger.Error(err)
	}

	s.logger.Info("PasswordExpired event was treated.")

	return nil
}
