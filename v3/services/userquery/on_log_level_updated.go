package userqueryapi

import (
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
	log "gitlab.com/wiserskills/v3/services/userquery/gen/log"
)

// OnLogLevelInfo is executed when the LogLevelInfo event is received
func (s *userquerysrvc) OnLogLevelUpdated(e *events.Event, payload interface{}) error {

	s.logger.Info("LogLevelUpdated event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of LogLevelUpdated events received
	s.metrics.IncrCounter("events_loglevelupdated_total_received")

	newLevel := payload.(*userquery.LogLevelInfo).Level

	switch newLevel {
	case "DEBUG":
		s.logger.SetLevel(log.DEBUG)
	case "INFO":
		s.logger.SetLevel(log.INFO)
	case "WARN":
		s.logger.SetLevel(log.WARNING)
	case "ERROR":
		s.logger.SetLevel(log.ERROR)
	}

	s.logger.Info("LogLevelUpdated event was treated.")

	return nil
}
