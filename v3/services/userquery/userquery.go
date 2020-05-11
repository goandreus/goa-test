package userqueryapi

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.com/wiserskills/v3/services/userquery/datalayer"
	"gitlab.com/wiserskills/v3/services/userquery/gen/environment"
	events "gitlab.com/wiserskills/v3/services/userquery/gen/events"
	log "gitlab.com/wiserskills/v3/services/userquery/gen/log"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// userquery service example implementation.
// The example methods log the requests and return zero values.
type userquerysrvc struct {
	logger  *log.Logger
	metrics *userquery.Metrics
	status  userquery.ServiceStatus
	bus     events.EventbusClient
	store   *datalayer.Datastore
}

// NewUserquery returns the userquery service implementation.
func NewUserquery(logger *log.Logger) userquery.Service {

	svc := &userquerysrvc{logger: logger}
	svc.SetServiceStatus(userquery.CONNECTING)

	svc.metrics = userquery.NewMetrics()
	svc.CreateCustomMetrics()
	// initializing plugins
	initBus(svc)
	initDatastore(svc)

	// plugins initialized
	svc.ServiceInit()

	svc.SetServiceStatus(userquery.READY)
	return svc
}

// Health status endpoint.
func (s *userquerysrvc) Health(ctx context.Context) (res *userquery.HealthResult, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "userquery.Health")
	defer span.Finish()
	res = &userquery.HealthResult{Status: string(s.status)}
	s.logger.Info("userquery.health")
	return
}

// initDatastore initializes & opens the datastore
func initDatastore(svc *userquerysrvc) {

	// Initializing & opening datastore
	svc.store = &datalayer.Datastore{}
	err := svc.store.Open()

	if err != nil {
		svc.logger.Fatalf("error while connecting to store: ", err)
	}

	err = svc.store.Seed()

	if err != nil {
		svc.logger.Fatalf("error while seeding store: ", err)
	}
}

// PublishServiceStatus publishes the service status to the event bus
func (s *userquerysrvc) PublishServiceStatus(status userquery.ServiceStatus) error {

	s.logger.Info("Publishing service status to event bus...")

	if s.bus != nil {

		payload := &userquery.ServiceStatusInfo{
			Service:    "userquery",
			InstanceID: events.GetInstanceID(),
			Status:     string(status),
		}

		e, err := events.NewEvent(context.Background(), "StatusChanged", 0, "ServiceStatus", payload)

		if err == nil {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, "[CLUSTER].services", e, payload, true)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of StatusChanged events published
				s.metrics.IncrCounter("events_statuschanged_total_published")

				s.logger.Infof("StatusChanged event published to '%s'.", parsedTopic)
				return nil
			}
			return err
		}
		return err
	}

	return nil
}

// PublishErrorInfo sends an error to the event bus
func (s *userquerysrvc) PublishErrorInfo(err error) error {

	if s.bus != nil {

		s.logger.Info("Publishing error to service bus...")

		payload := &userquery.ErrorInfo{
			Service:    "userquery",
			InstanceID: events.GetInstanceID(),
			Message:    err.Error(),
		}

		e, err := events.NewEvent(context.Background(), "ErrorInfo", 0, "ErrorInfo", payload)

		if err == nil {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, "[CLUSTER].services", e, payload, true)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of StatusChanged events published
				s.metrics.IncrCounter("events_errors_total_published")

				s.logger.Infof("ErrorInfo event published to '%s'", parsedTopic)
				return nil
			}
			return err
		}
		return err
	}

	return nil
}

// SetServiceStatus updates the service status
func (s *userquerysrvc) SetServiceStatus(status userquery.ServiceStatus) error {

	s.status = status

	return s.PublishServiceStatus(status)
}

// PublishTokenCreated sends a TokenCreated event to the event bus
func (s *userquerysrvc) PublishTokenCreated(ctx context.Context, payload *userquery.Token, async bool, additionalTopics ...string) error {

	s.logger.Info("Publishing TokenCreated...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_tokencreated_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_tokencreated_total_published", "Total number of TokenCreated events published.")
	}

	e, err := events.NewEvent(ctx, "TokenCreated", 0, "Token", payload)

	if err == nil {

		topics := []string{"[CLUSTER].user.query"}

		if len(additionalTopics) > 0 {
			topics = append(topics, additionalTopics...)
		}

		for _, topic := range topics {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, topic, e, payload, async)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of TokenCreated events published
				if published != nil {
					published.Inc()
				}

				s.logger.Infof("TokenCreated published to '%s'.", parsedTopic)

			} else {
				return err
			}
		}
		return nil
	}
	return err
}

// PublishTokenExpired sends a TokenExpired event to the event bus
func (s *userquerysrvc) PublishTokenExpired(ctx context.Context, payload *userquery.Token, async bool, additionalTopics ...string) error {

	s.logger.Info("Publishing TokenExpired...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_tokenexpired_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_tokenexpired_total_published", "Total number of TokenExpired events published.")
	}

	e, err := events.NewEvent(ctx, "TokenExpired", 2, "Token", payload)

	if err == nil {

		topics := []string{"[CLUSTER].user.query"}

		if len(additionalTopics) > 0 {
			topics = append(topics, additionalTopics...)
		}

		for _, topic := range topics {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, topic, e, payload, async)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of TokenExpired events published
				if published != nil {
					published.Inc()
				}

				s.logger.Infof("TokenExpired published to '%s'.", parsedTopic)

			} else {
				return err
			}
		}
		return nil
	}
	return err
}

// PublishSessionCreated sends a SessionCreated event to the event bus
func (s *userquerysrvc) PublishSessionCreated(ctx context.Context, payload *userquery.Session, async bool, additionalTopics ...string) error {

	s.logger.Info("Publishing SessionCreated...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_sessioncreated_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_sessioncreated_total_published", "Total number of SessionCreated events published.")
	}

	e, err := events.NewEvent(ctx, "SessionCreated", 0, "Session", payload)

	if err == nil {

		topics := []string{"[CLUSTER].user.query"}

		if len(additionalTopics) > 0 {
			topics = append(topics, additionalTopics...)
		}

		for _, topic := range topics {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, topic, e, payload, async)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of SessionCreated events published
				if published != nil {
					published.Inc()
				}

				s.logger.Infof("SessionCreated published to '%s'.", parsedTopic)

			} else {
				return err
			}
		}
		return nil
	}
	return err
}

// PublishSessionUpdated sends a SessionUpdated event to the event bus
func (s *userquerysrvc) PublishSessionUpdated(ctx context.Context, payload *userquery.Session, async bool, additionalTopics ...string) error {

	s.logger.Info("Publishing SessionUpdated...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_sessionupdated_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_sessionupdated_total_published", "Total number of SessionUpdated events published.")
	}

	e, err := events.NewEvent(ctx, "SessionUpdated", 1, "Session", payload)

	if err == nil {

		topics := []string{"[CLUSTER].user.query"}

		if len(additionalTopics) > 0 {
			topics = append(topics, additionalTopics...)
		}

		for _, topic := range topics {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, topic, e, payload, async)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of SessionUpdated events published
				if published != nil {
					published.Inc()
				}

				s.logger.Infof("SessionUpdated published to '%s'.", parsedTopic)

			} else {
				return err
			}
		}
		return nil
	}
	return err
}

// PublishSessionExpired sends a SessionExpired event to the event bus
func (s *userquerysrvc) PublishSessionExpired(ctx context.Context, payload *userquery.Session, async bool, additionalTopics ...string) error {

	s.logger.Info("Publishing SessionExpired...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_sessionexpired_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_sessionexpired_total_published", "Total number of SessionExpired events published.")
	}

	e, err := events.NewEvent(ctx, "SessionExpired", 2, "Session", payload)

	if err == nil {

		topics := []string{"[CLUSTER].user.query"}

		if len(additionalTopics) > 0 {
			topics = append(topics, additionalTopics...)
		}

		for _, topic := range topics {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, topic, e, payload, async)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of SessionExpired events published
				if published != nil {
					published.Inc()
				}

				s.logger.Infof("SessionExpired published to '%s'.", parsedTopic)

			} else {
				return err
			}
		}
		return nil
	}
	return err
}

// PublishPasswordExpired sends a PasswordExpired event to the event bus
func (s *userquerysrvc) PublishPasswordExpired(ctx context.Context, payload *userquery.Password, async bool, additionalTopics ...string) error {

	s.logger.Info("Publishing PasswordExpired...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_passwordexpired_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_passwordexpired_total_published", "Total number of PasswordExpired events published.")
	}

	e, err := events.NewEvent(ctx, "PasswordExpired", 2, "Password", payload)

	if err == nil {

		topics := []string{"[CLUSTER].user.admin"}

		if len(additionalTopics) > 0 {
			topics = append(topics, additionalTopics...)
		}

		for _, topic := range topics {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, topic, e, payload, async)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of PasswordExpired events published
				if published != nil {
					published.Inc()
				}

				s.logger.Infof("PasswordExpired published to '%s'.", parsedTopic)

			} else {
				return err
			}
		}
		return nil
	}
	return err
}

// PublishAuthFailureCreated sends a AuthFailureCreated event to the event bus
func (s *userquerysrvc) PublishAuthFailureCreated(ctx context.Context, payload *userquery.AuthFailure, async bool, additionalTopics ...string) error {

	s.logger.Info("Publishing AuthFailureCreated...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_authfailurecreated_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_authfailurecreated_total_published", "Total number of AuthFailureCreated events published.")
	}

	e, err := events.NewEvent(ctx, "AuthFailureCreated", 0, "AuthFailure", payload)

	if err == nil {

		topics := []string{"[CLUSTER].user.query"}

		if len(additionalTopics) > 0 {
			topics = append(topics, additionalTopics...)
		}

		for _, topic := range topics {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, topic, e, payload, async)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of AuthFailureCreated events published
				if published != nil {
					published.Inc()
				}

				s.logger.Infof("AuthFailureCreated published to '%s'.", parsedTopic)

			} else {
				return err
			}
		}
		return nil
	}
	return err
}

// PublishUserUpdated sends a UserUpdated event to the event bus
func (s *userquerysrvc) PublishUserUpdated(ctx context.Context, payload *userquery.User, async bool, additionalTopics ...string) error {

	s.logger.Info("Publishing UserUpdated...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_userupdated_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_userupdated_total_published", "Total number of UserUpdated events published.")
	}

	e, err := events.NewEvent(ctx, "UserUpdated", 1, "User", payload)

	if err == nil {

		topics := []string{"[CLUSTER].user.admin;[CLUSTER].[$Event.OrgID].user"}

		if len(additionalTopics) > 0 {
			topics = append(topics, additionalTopics...)
		}

		for _, topic := range topics {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, topic, e, payload, async)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of UserUpdated events published
				if published != nil {
					published.Inc()
				}

				s.logger.Infof("UserUpdated published to '%s'.", parsedTopic)

			} else {
				return err
			}
		}
		return nil
	}
	return err
}

// PublishTokenRevoked sends a TokenRevoked event to the event bus
func (s *userquerysrvc) PublishTokenRevoked(ctx context.Context, payload *userquery.Token, async bool, additionalTopics ...string) error {

	s.logger.Info("Publishing TokenRevoked...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_tokenrevoked_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_tokenrevoked_total_published", "Total number of TokenRevoked events published.")
	}

	e, err := events.NewEvent(ctx, "TokenRevoked", 2, "Token", payload)

	if err == nil {

		topics := []string{"[CLUSTER].user.query"}

		if len(additionalTopics) > 0 {
			topics = append(topics, additionalTopics...)
		}

		for _, topic := range topics {
			var parsedTopic string
			err, parsedTopic = events.PublishEvent(s.bus, topic, e, payload, async)
			if err == nil {

				// We increment the total number of published events
				s.metrics.IncrCounter("events_total_published")

				// We increment the total number of TokenRevoked events published
				if published != nil {
					published.Inc()
				}

				s.logger.Infof("TokenRevoked published to '%s'.", parsedTopic)

			} else {
				return err
			}
		}
		return nil
	}
	return err
}

// getTopicSubscriptions returns the list of the subscribed topics
func (s *userquerysrvc) getTopicSubscriptions() []*events.TopicSubscription {

	// We register the topic options
	events.RegisterTopicOptions("[CLUSTER].services.userquery", false, "", true, 1, nil)
	events.RegisterTopicOptions("[CLUSTER].user.admin", true, "userquery", true, 1, nil)
	events.RegisterTopicOptions("[CLUSTER].user.query", true, "userquery", true, 1, nil)

	// We register to the endpoints requested event
	events.RegisterEventHandler("[CLUSTER].services.endpoints", "EndpointsRequested", s.OnEndpointsRequested, []*events.EndpointInfo{})

	// We register with the other events
	events.RegisterEventHandler("[CLUSTER].services.userquery", "LogLevelUpdated", s.OnLogLevelUpdated, userquery.LogLevelInfo{})
	s.metrics.NewCounter("events_loglevelupdated_total_received", "Total number of LogLevelUpdated events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "UserCreated", s.OnUserCreated, userquery.User{})
	s.metrics.NewCounter("events_usercreated_total_received", "Total number of UserCreated events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "UserUpdated", s.OnUserUpdated, userquery.User{})
	s.metrics.NewCounter("events_userupdated_total_received", "Total number of UserUpdated events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "UserDeleted", s.OnUserDeleted, userquery.User{})
	s.metrics.NewCounter("events_userdeleted_total_received", "Total number of UserDeleted events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "UserEnabled", s.OnUserEnabled, userquery.User{})
	s.metrics.NewCounter("events_userenabled_total_received", "Total number of UserEnabled events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "UserDisabled", s.OnUserDisabled, userquery.User{})
	s.metrics.NewCounter("events_userdisabled_total_received", "Total number of UserDisabled events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "SAMLConfigCreated", s.OnSAMLConfigCreated, userquery.SAMLConfig{})
	s.metrics.NewCounter("events_samlconfigcreated_total_received", "Total number of SAMLConfigCreated events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "SAMLConfigUpdated", s.OnSAMLConfigUpdated, userquery.SAMLConfig{})
	s.metrics.NewCounter("events_samlconfigupdated_total_received", "Total number of SAMLConfigUpdated events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "SAMLConfigDeleted", s.OnSAMLConfigDeleted, userquery.SAMLConfig{})
	s.metrics.NewCounter("events_samlconfigdeleted_total_received", "Total number of SAMLConfigDeleted events received.")
	events.RegisterEventHandler("[CLUSTER].user.query", "TokenCreated", s.OnTokenCreated, userquery.Token{})
	s.metrics.NewCounter("events_tokencreated_total_received", "Total number of TokenCreated events received.")
	events.RegisterEventHandler("[CLUSTER].user.query", "TokenExpired", s.OnTokenExpired, userquery.Token{})
	s.metrics.NewCounter("events_tokenexpired_total_received", "Total number of TokenExpired events received.")
	events.RegisterEventHandler("[CLUSTER].user.query", "TokenRevoked", s.OnTokenRevoked, userquery.Token{})
	s.metrics.NewCounter("events_tokenrevoked_total_received", "Total number of TokenRevoked events received.")
	events.RegisterEventHandler("[CLUSTER].user.query", "SessionCreated", s.OnSessionCreated, userquery.Session{})
	s.metrics.NewCounter("events_sessioncreated_total_received", "Total number of SessionCreated events received.")
	events.RegisterEventHandler("[CLUSTER].user.query", "SessionUpdated", s.OnSessionUpdated, userquery.Session{})
	s.metrics.NewCounter("events_sessionupdated_total_received", "Total number of SessionUpdated events received.")
	events.RegisterEventHandler("[CLUSTER].user.query", "SessionExpired", s.OnSessionExpired, userquery.Session{})
	s.metrics.NewCounter("events_sessionexpired_total_received", "Total number of SessionExpired events received.")
	events.RegisterEventHandler("[CLUSTER].user.query", "AuthFailureCreated", s.OnAuthFailureCreated, userquery.AuthFailure{})
	s.metrics.NewCounter("events_authfailurecreated_total_received", "Total number of AuthFailureCreated events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "PasswordExpired", s.OnPasswordExpired, userquery.Password{})
	s.metrics.NewCounter("events_passwordexpired_total_received", "Total number of PasswordExpired events received.")
	events.RegisterEventHandler("[CLUSTER].user.admin", "PasswordChanged", s.OnPasswordChanged, userquery.Password{})
	s.metrics.NewCounter("events_passwordchanged_total_received", "Total number of PasswordChanged events received.")

	return events.GetTopicSubscriptions()
}

// OnEndpointsRequested is executed when the EndpointsRequested event is received
func (s *userquerysrvc) OnEndpointsRequested(e *events.Event, payload interface{}) error {

	s.logger.Info("EndpointsRequested event was received...")

	// We increment the total number of events received
	s.metrics.IncrCounter("events_total_received")

	// We increment the number of EndpointsRequested events received
	s.metrics.IncrCounter("events_endpointsrequested_total_received")

	// We populate the endpoints array
	endpoints := []*events.EndpointInfo{}
	namespace := environment.GetEnv("NAMESPACE", "")
	endpoints = append(endpoints, &events.EndpointInfo{Name: "SignIn", Code: "EP04Q001", Description: "Signs a user in.", Verb: "POST", Path: "/signin", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "SignOut", Code: "EP04Q002", Description: "Signs a user out.", Verb: "POST", Path: "/signout", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "GetAllSessions", Code: "EP04Q003", Description: "Returns the currently active sessions.", Verb: "GET", Path: "/sessions", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "GetIDPURL", Code: "EP04Q004", Description: "Returns the URL of the IDP to redirect the user to.", Verb: "GET", Path: "/saml", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "SamlSignIn", Code: "EP04Q005", Description: "Call back endpoint called by the IDP once the user is authenticated.", Verb: "POST", Path: "/saml/signin", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "CheckToken", Code: "EP04Q006", Description: "Checks if the passed token is valid.", Verb: "POST", Path: "/check", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "GetUsersByID", Code: "EP04Q007", Description: "Returns the users with the specified ids.", Verb: "GET", Path: "/users/id", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "health", Code: "EP04Q801", Description: "Health status endpoint.", Verb: "GET", Path: "/health", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "metrics", Code: "EP04Q802", Description: "Prometheus metrics endpoint.", Verb: "GET", Path: "/metrics", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "swagger", Code: "EP04Q901", Description: "Swagger service description endpoint.", Verb: "GET", Path: "/swagger", Service: "userquery", Namespace: namespace})
	endpoints = append(endpoints, &events.EndpointInfo{Name: "swagger-ui", Code: "EP04Q902", Description: "Swagger UI service endpoint.", Verb: "GET", Path: "/swaggerui", Service: "userquery", Namespace: namespace})

	// We publish the response event
	ctx := events.CreateContextFromEvent(e)
	s.PublishEndpointsInfo(ctx, endpoints, true)

	s.logger.Info("EndpointsRequested event was treated.")

	return nil
}

// PublishEndpointsInfo sends an EndpointsInfo event to the event bus
func (s *userquerysrvc) PublishEndpointsInfo(ctx context.Context, payload []*events.EndpointInfo, async bool) error {

	s.logger.Info("Publishing EndpointsInfo...")

	// Creating endpoint total calls counter if doesn't exist already
	published := s.metrics.GetCounter("events_endpointsinfo_total_published")
	if published == nil {
		published = s.metrics.NewCounter("events_endpointsinfo_total_published", "Total number of EndpointsInfo events published.")
	}

	e, err := events.NewEvent(ctx, "EndpointsInfo", 0, "[]EndpointInfo", payload)

	if err == nil {
		var parsedTopic string
		err, parsedTopic = events.PublishEvent(s.bus, "[CLUSTER].services.endpoints", e, payload, async)
		if err == nil {

			// We increment the total number of published events
			s.metrics.IncrCounter("events_total_published")

			// We increment the total number of EndpointsInfo events published
			if published != nil {
				published.Inc()
			}

			s.logger.Infof("EndpointsInfo published to '%s'.", parsedTopic)
			return nil
		}
		return err
	}
	return err
}

// Dataservice defines a service that has a datastore
type Dataservice interface {
	Store() *datalayer.Datastore
}

// Store returns the inner datastore
func (s *userquerysrvc) Store() *datalayer.Datastore {
	return s.store
}

// Reloads the topic subscriptions
func (s *userquerysrvc) ReloadTopicSubscriptions() {

	if s.bus != nil {
		s.bus.Reload(s.getTopicSubscriptions())
	}
}

// initBus initializes the event bus
func initBus(svc *userquerysrvc) {

	// Creating event bus client and register with NATS
	busConfig, err := events.GetBusConfig()

	if err != nil {
		svc.logger.Fatalf("error while instanciating event bus configuration: ", err)
	}

	busConfig.Subscriptions = svc.getTopicSubscriptions()
	bus, err := events.NewNATSClient(svc.logger, busConfig)

	if err != nil {
		svc.logger.Fatalf("error while instanciating event bus client: ", err)
	}

	svc.bus = bus

	err = svc.bus.Connect()

	if err != nil {
		svc.logger.Fatalf("error while connecting to event bus: ", err)
	}

	// Creating events related metrics
	svc.metrics.NewCounter("events_total_received", "Total number of received events since service startup.")
	svc.metrics.NewCounter("events_total_published", "Total number of published events since service startup.")
	svc.metrics.NewCounter("events_statuschanged_total_published", "Total number of statuschanged events published since service startup.")
	svc.metrics.NewCounter("events_errors_total_published", "Total number of errors published since service startup.")
}
