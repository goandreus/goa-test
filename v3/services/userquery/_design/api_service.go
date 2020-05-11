package design

import (
	. "goa.design/goa/v3/dsl"

	_ "goa.design/plugins/v3/environment"
	. "goa.design/plugins/v3/environment/dsl"
	_ "goa.design/plugins/v3/health"
	_ "goa.design/plugins/v3/opentracing"
	_ "goa.design/plugins/v3/prometheus"
	_ "goa.design/plugins/v3/zaplogger"

	_ "goa.design/plugins/v3/secure"

	_ "goa.design/plugins/v3/events"
	. "goa.design/plugins/v3/events/dsl"
	"goa.design/plugins/v3/events/expr"

	_ "goa.design/plugins/v3/nats"

	_ "goa.design/plugins/v3/arangodb"

	_ "goa.design/plugins/v3/apitests"
	_ "goa.design/plugins/v3/buildfile"
	_ "goa.design/plugins/v3/doc"
	_ "goa.design/plugins/v3/dockerfile"
	_ "goa.design/plugins/v3/swagger"
)

// Service describes a service
var _ = Service(LowerServiceName, func() {

	Description(ServiceDescription)

	Meta("endpoint:code:root", "EP"+DomainCode+ServiceTypeCode)

	// Environment defines the environment variables used by the service
	Environment("CLUSTER", String, "The name of the cluster the service is running on.", "fr0", true)
	Environment("NAMESPACE", String, "The name of the namespace the service is running in.", "services", true)
	Environment("BUSINESS_DOMAIN", String, "The name of the business domain that the service is covering.", Domain, true)
	Environment("SERVICE_TYPE", String, "Defines the type of the service.", ServiceType, true)
	Environment("SERVICE_NAME", String, "The name of the service.", LowerServiceName, true)
	Environment("NODE", String, "The name of the node the service is running on.", "", false)

	Environment("EVENTBUS_CLUSTER", String, "The name of the event bus cluster.", "wiserskills", false)
	Environment("EVENTBUS_URL", String, "The URL of the event bus.", "nats://localhost:4222", false)
	Environment("EVENTBUS_DEFAULT_TOPIC", String, "The event bus default topic.", "", false)
	Environment("EVENTBUS_WRITE_TIMEOUT", Int, "The timeout to write to the event bus in microseconds.", "5000000", false)
	Environment("EVENTBUS_RETRY_COUNT", Int, "The number of retry to connect to the event bus in case of connection failure.", "5", false)
	Environment("EVENTBUS_RETRY_DELAY", Int, "The delay in seconds before retrying to connect to the event bus in case of failure.", "5", false)
	Environment("EVENTBUS_CLIENT_CERT_PATH", String, "The path to the client certificate file in PEM format used to connect to NATS.", "", false)
	Environment("EVENTBUS_CLIENT_KEY_PATH", String, "The path to the client private key file in PEM format used to connect to NATS.", "", false)

	Environment("INTERNAL_ORG_ID", String, "The id of the organization that uses the internal scope.", "wiserskills", false)

	Environment("API_KEY", String, "The api key.", "wsk3z", false)
	Environment("TOKEN_KEY", String, "The key used to sign the token.", "sdzT3wxcA", false)
	Environment("TOKEN_ACTIVE", Boolean, "Defines if security by token is active.", "false", false)

	Environment("TOKEN_LIFE_TIME", Int, "The life time of a token in minutes.", "1440", false)
	Environment("SESSION_LIFE_TIME", Int, "The life time of a session in minutes.", "20", false)
	Environment("LOGIN_MAX_ATTEMPT", Int, "The maximum number of failed login attempts.", "5", false)
	Environment("LOGIN_SUSPENSION_TIME", Int, "The number of minutes a login can be suspended when the maximum failed number of attempts is reached.", "120", false)

	Environment("DB_CONNECTION", String, "The URL of the database server.", "http://localhost:8529", false)
	Environment("DB_LOGIN", String, "The login used to connect to the database server.", "", false)
	Environment("DB_PASSWORD", String, "The password used to connect to the database server.", "", false)
	Environment("DB_SEED_FOLDER", String, "The path of the folder that contain seed files.", "seed", false)
	Environment("DB_ENCRYPT", Boolean, "Defines if personal data fields must be encrypted.", "false", false)

	// List of errors raised by the service
	Error("not_found", ErrorResult, "Entity was not found.")
	Error("bad_argument", ErrorResult, "Unvalid parameter(s).")
	Error("not_authorized", ErrorResult, "Access denied.")
	Error("unexpected_error", ErrorResult, "An unexpected error was raised.")
	Error("invalid_scope", ErrorResult, "Token scope is invalid.")

	Error("token_invalid", ErrorResult, "Token invalid.")
	Error("token_expired", ErrorResult, "Token has expired.")
	Error("password_expired", ErrorResult, "Password has expired.")
	Error("invalid_credentials", ErrorResult, "Username or password is invalid.")
	Error("login_blocked", ErrorResult, "Login is temporarly blocked due to too many failed attempts.")

	// SubscribeEvent defines the event(s) the service subscribes to
	// Topic options define the durability and replay stategy of the event subscriptions
	// NB: OnlyOnce is set to true. If multiple instances of the service are running, events will only be treated by one of the instances
	TopicOptions("[CLUSTER].services.userquery", false, true, expr.LASTRECEIVED, nil)
	TopicOptions("[CLUSTER].user.admin", true, true, expr.LASTRECEIVED, nil)
	TopicOptions("[CLUSTER].user.query", true, true, expr.LASTRECEIVED, nil)

	// Service specific events
	SubscribeEvent("LogLevelUpdated", "EV0002", LogLevelInfo, "LogLevelUpdated is received when the level of log must be changed.", "[CLUSTER].services.userquery")

	// Events related to entities stored in cache only (not owned by the service)

	SubscribeEvent("UserCreated", "EV0401", User, "UserCreated is received when a new user has been created.", "[CLUSTER].user.admin")
	SubscribeEvent("UserUpdated", "EV0402", User, "UserUpdated is received when an existing user has been updated.", "[CLUSTER].user.admin")
	SubscribeEvent("UserDeleted", "EV0403", User, "UserDeleted is received when an existing user has been deleted.", "[CLUSTER].user.admin")
	SubscribeEvent("UserEnabled", "EV0404", User, "UserEnabled is received when an existing user has been enabled.", "[CLUSTER].user.admin")
	SubscribeEvent("UserDisabled", "EV0405", User, "UserDisabled is received when an existing user has been disabled.", "[CLUSTER].user.admin")

	SubscribeEvent("SAMLConfigCreated", "EV0422", SAMLConfig, "SAMLConfigCreated is received when a new SAML config has been created.", "[CLUSTER].user.admin")
	SubscribeEvent("SAMLConfigUpdated", "EV0423", SAMLConfig, "SAMLConfigUpdated is received when a SAML config has been updated.", "[CLUSTER].user.admin")
	SubscribeEvent("SAMLConfigDeleted", "EV0424", SAMLConfig, "SAMLConfigDeleted is received when a SAML config has been deleted.", "[CLUSTER].user.admin")

	// Events related to entities owned by the service

	SubscribeEvent("TokenCreated", "EV0406", TokenInfo, "TokenCreated is received when a new token has been created.", "[CLUSTER].user.query")
	SubscribeEvent("TokenExpired", "EV0407", TokenInfo, "TokenExpired is received when a token has expired.", "[CLUSTER].user.query")
	SubscribeEvent("TokenRevoked", "EV0408", TokenInfo, "TokenRevoked is received when a token has been revoked.", "[CLUSTER].user.query")

	SubscribeEvent("SessionCreated", "EV0409", Session, "SessionCreated is received when a new session has been created.", "[CLUSTER].user.query")
	SubscribeEvent("SessionUpdated", "EV0410", Session, "SessionUpdated is received when a session has been updated.", "[CLUSTER].user.query")
	SubscribeEvent("SessionExpired", "EV0411", Session, "SessionExpired is received when a session has expired.", "[CLUSTER].user.query")

	SubscribeEvent("AuthFailureCreated", "EV0417", AuthFailure, "AuthFailureCreated is received when an authentication failure was created.", "[CLUSTER].user.query")

	SubscribeEvent("PasswordExpired", "EV0418", PasswordInfo, "PasswordExpired is received when a password has expired.", "[CLUSTER].user.admin")
	SubscribeEvent("PasswordChanged", "EV0419", PasswordInfo, "PasswordChanged is received when a password has been changed.", "[CLUSTER].user.admin")

	// Endpoints definitions ******************************************************************************

	// EP04Q001 - SignIn
	Method("SignIn", func() {

		Description("Signs a user in.")

		Meta("endpoint:code", "EP04Q001")

		// The signin endpoint is secured via basic auth
		Security(BasicAuth, APIKeyAuth)

		Payload(func() {

			Description("Credentials used to authenticate and retrieve a JWT token")

			APIKeyField(1, "api_key", "key", String, func() {
				Description("API key")

				Example("abcdef12345")
			})

			UsernameField(2, "username", String, "Username used to perform signin", func() {
				Example("user")
			})

			PasswordField(3, "password", String, "Password used to perform signin", func() {
				Example("password")
			})

			Field(4, "orgId", String, func() {
				Description("The id of the organization.")

				Example("WiserSKILLS")
			})

			Field(5, "areaId", String, func() {
				Description("The id of the area.")

				Enum("PROD", "UAT", "QA", "DEV")
				Example("PROD")
			})

			Required("username", "password", "orgId", "areaId")
		})

		Result(JWTToken)

		// PublishEvent defines the event(s) that is/are raised by the method
		PublishEvent("TokenCreated", "EV0406", TokenInfo, "TokenCreated is raised when a new token is created.", "[CLUSTER].user.query", expr.CREATED)
		PublishEvent("TokenExpired", "EV0407", TokenInfo, "TokenExpired is raised when a token is expired.", "[CLUSTER].user.query", expr.DELETED)

		PublishEvent("SessionCreated", "EV0409", Session, "SessionCreated is raised when a new session is created.", "[CLUSTER].user.query", expr.CREATED)
		PublishEvent("SessionUpdated", "EV0410", Session, "SessionUpdated is raised when a session is updated.", "[CLUSTER].user.query", expr.UPDATED)
		PublishEvent("SessionExpired", "EV0411", Session, "SessionExpired is raised when a session is expired.", "[CLUSTER].user.query", expr.DELETED)

		PublishEvent("PasswordExpired", "EV0418", PasswordInfo, "PasswordExpired is raised when a user's password has expired.", "[CLUSTER].user.admin", expr.DELETED)

		PublishEvent("AuthFailureCreated", "EV0417", AuthFailure, "AuthFailureCreated is raised when a authentication attempt has failed.", "[CLUSTER].user.query", expr.CREATED)

		// Note: this event will be propagated to the data server
		PublishEvent("UserUpdated", "EV0402", User, "UserUpdated is raised when a user has been updated.", "[CLUSTER].user.admin;[CLUSTER].[$Event.OrgID].user", expr.UPDATED)

		// HTTP describes the HTTP transport mapping.
		HTTP(func() {

			// Requests to the service consist of HTTP POST requests.
			POST("/signin")

			Header("key:X-API-KEY")    // key passed in "X-API-KEY" header
			Header("orgId:X-ORG-ID")   // The id of the organization if any
			Header("areaId:X-AREA-ID") // The id of the area if any

			// Responses use a "200 OK" HTTP status.
			// The result is encoded in the response body.
			Response(StatusOK)
			Response("bad_argument", StatusBadRequest)
			Response("not_authorized", StatusForbidden)
			Response("invalid_credentials", StatusForbidden)
			Response("password_expired", StatusForbidden)
			Response("login_blocked", StatusForbidden)
			Response("unexpected_error", StatusInternalServerError)
		})

		GRPC(func() {
			Response(CodeOK)
			Response("bad_argument", CodeInvalidArgument)
			Response("not_authorized", CodePermissionDenied)
			Response("invalid_credentials", CodePermissionDenied)
			Response("password_expired", CodePermissionDenied)
			Response("login_blocked", CodePermissionDenied)
			Response("unexpected_error", CodeInternal)
		})
	})

	// EP04Q002 - SignOut
	Method("SignOut", func() {

		Description("Signs a user out.")

		Meta("endpoint:code", "EP04Q002")

		Security(JWTAuth, APIKeyAuth)

		Payload(TokenPayload)

		// PublishEvent defines the event(s) that is/are raised by the method
		PublishEvent("TokenRevoked", "EV0408", TokenInfo, "TokenRevoked is raised when a token is revoked.", "[CLUSTER].user.query", expr.DELETED)
		PublishEvent("SessionExpired", "EV0411", Session, "SessionExpired is raised when a session is expired.", "[CLUSTER].user.query", expr.DELETED)

		HTTP(func() {

			POST("/signout")

			Header("token:Authorization") // JWT token passed in "Authorization" header
			Header("key:X-API-KEY")       // key passed in "X-API-KEY" header
			Header("orgId:X-ORG-ID")      // The id of the organization if any
			Header("areaId:X-AREA-ID")    // The id of the area if any

			Response(StatusOK)
			Response("not_found", StatusNotFound)
			Response("token_expired", StatusForbidden)
			Response("token_invalid", StatusForbidden)
			Response("bad_argument", StatusBadRequest)
			Response("unexpected_error", StatusInternalServerError)
			Response("not_authorized", StatusForbidden)
		})

		GRPC(func() {
			Response(CodeOK)
			Response("not_found", CodeNotFound)
			Response("token_expired", CodePermissionDenied)
			Response("token_invalid", CodePermissionDenied)
			Response("bad_argument", CodeInvalidArgument)
			Response("unexpected_error", CodeInternal)
			Response("not_authorized", CodePermissionDenied)
		})
	})

	// EP04Q003 - GetAllSessions
	Method("GetAllSessions", func() {

		Description("Returns the currently active sessions.")
		Meta("endpoint:code", "EP04Q003")

		Security(JWTAuth, APIKeyAuth)

		Payload(AllSessionsPayload)

		Result(CollectionOf(Session))

		HTTP(func() {

			GET("/sessions")

			Header("token:Authorization") // JWT token passed in "Authorization" header
			Header("key:X-API-KEY")       // key passed in "X-API-KEY" header
			Header("orgId:X-ORG-ID")      // The id of the organization
			Header("areaId:X-AREA-ID")    // The id of the area

			Param("view")

			Response(StatusOK)
			Response("token_expired", StatusForbidden)
			Response("token_invalid", StatusForbidden)
			Response("bad_argument", StatusBadRequest)
			Response("unexpected_error", StatusInternalServerError)
		})

		GRPC(func() {
			Response(CodeOK)
			Response("not_found", CodeNotFound)
			Response("token_expired", CodePermissionDenied)
			Response("token_invalid", CodePermissionDenied)
			Response("bad_argument", CodeInvalidArgument)
			Response("unexpected_error", CodeInternal)
		})
	})

	// EP04Q004 - GetIDPURL
	Method("GetIDPURL", func() {

		Description("Returns the URL of the IDP to redirect the user to.")

		Meta("endpoint:code", "EP04Q004")

		// The signin endpoint is secured via api key
		Security(APIKeyAuth)

		Payload(HostPayload)

		Result(RedirectResult)

		HTTP(func() {

			GET("/saml")

			Header("key:X-API-KEY")    // key passed in "X-API-KEY" header
			Header("orgId:X-ORG-ID")   // The id of the organization
			Header("areaId:X-AREA-ID") // The id of the area
			Header("host:Host")        // The host

			Response(func() {
				Code(StatusPermanentRedirect)
				Header("Location")
			})
			Response("bad_argument", StatusBadRequest)
			Response("unexpected_error", StatusInternalServerError)
		})

		GRPC(func() {
			Response(CodeOK)
			Response("bad_argument", CodeInvalidArgument)
			Response("unexpected_error", CodeInternal)
		})
	})

	// EP04Q005 - SamlSignIn
	Method("SamlSignIn", func() {

		Description("Call back endpoint called by the IDP once the user is authenticated.")

		Meta("endpoint:code", "EP04Q005")

		Payload(String)

		Result(RedirectResult)

		// PublishEvent defines the event(s) that is/are raised by the method
		PublishEvent("TokenCreated", "EV0406", TokenInfo, "TokenCreated is raised when a new token is created.", "[CLUSTER].user.query", expr.CREATED)
		PublishEvent("TokenExpired", "EV0407", TokenInfo, "TokenExpired is raised when a token is expired.", "[CLUSTER].user.query", expr.DELETED)

		PublishEvent("SessionCreated", "EV0409", Session, "SessionCreated is raised when a new session is created.", "[CLUSTER].user.query", expr.CREATED)
		PublishEvent("SessionUpdated", "EV0410", Session, "SessionUpdated is raised when a session is updated.", "[CLUSTER].user.query", expr.UPDATED)
		PublishEvent("SessionExpired", "EV0411", Session, "SessionExpired is raised when a session is expired.", "[CLUSTER].user.query", expr.DELETED)

		PublishEvent("PasswordExpired", "EV0418", PasswordInfo, "PasswordExpired is raised when a user's password has expired.", "[CLUSTER].user.admin", expr.DELETED)

		PublishEvent("AuthFailureCreated", "EV0417", AuthFailure, "AuthFailureCreated is raised when a authentication attempt has failed.", "[CLUSTER].user.query", expr.CREATED)

		// Note: this event will be propagated to the data server
		PublishEvent("UserUpdated", "EV0402", User, "UserUpdated is raised when a user has been updated.", "[CLUSTER].user.admin;[CLUSTER].[$Event.OrgID].user", expr.UPDATED)

		HTTP(func() {

			POST("/saml/signin")

			Response(func() {
				Code(StatusPermanentRedirect)
				Header("Location")
			})
			Response("not_authorized", StatusForbidden)
			Response("bad_argument", StatusBadRequest)
			Response("unexpected_error", StatusInternalServerError)
		})
	})

	// EP04Q006 - CheckToken
	Method("CheckToken", func() {

		Description("Checks if the passed token is valid.")

		Meta("endpoint:code", "EP04Q006")

		Security(JWTAuth, APIKeyAuth)

		Payload(TokenPayload)

		HTTP(func() {

			POST("/check")

			Header("token:Authorization") // JWT token passed in "Authorization" header
			Header("key:X-API-KEY")       // key passed in "X-API-KEY" header
			Header("orgId:X-ORG-ID")      // The id of the organization
			Header("areaId:X-AREA-ID")    // The id of the area

			Response(StatusOK)
			Response("not_found", StatusNotFound)
			Response("token_expired", StatusForbidden)
			Response("token_invalid", StatusForbidden)
			Response("bad_argument", StatusBadRequest)
			Response("unexpected_error", StatusInternalServerError)
			Response("not_authorized", StatusForbidden)
		})

		GRPC(func() {
			Response(CodeOK)
			Response("not_found", CodeNotFound)
			Response("token_expired", CodePermissionDenied)
			Response("token_invalid", CodePermissionDenied)
			Response("bad_argument", CodeInvalidArgument)
			Response("unexpected_error", CodeInternal)
			Response("not_authorized", CodePermissionDenied)
		})
	})

	// EP04Q007 - GetUsersByID
	Method("GetUsersByID", func() {

		Description("Returns the users with the specified ids.")

		Meta("endpoint:code", "EP04Q007")

		Security(JWTAuth, APIKeyAuth)

		Payload(ManyUserIDPayload)

		Result(CollectionOf(RegisteredUser))

		HTTP(func() {

			GET("/users/id")

			Header("token:Authorization") // JWT token passed in "Authorization" header
			Header("key:X-API-KEY")       // key passed in "X-API-KEY" header
			Header("orgId:X-ORG-ID")      // The id of the organization
			Header("areaId:X-AREA-ID")    // The id of the area

			Body("ids")

			Param("view")

			Response(StatusOK)
			Response("not_found", StatusNotFound)
			Response("token_expired", StatusForbidden)
			Response("token_invalid", StatusForbidden)
			Response("bad_argument", StatusBadRequest)
			Response("unexpected_error", StatusInternalServerError)
			Response("not_authorized", StatusForbidden)
		})

		GRPC(func() {
			Response(CodeOK)
			Response("not_found", CodeNotFound)
			Response("token_expired", CodePermissionDenied)
			Response("token_invalid", CodePermissionDenied)
			Response("bad_argument", CodeInvalidArgument)
			Response("unexpected_error", CodeInternal)
			Response("not_authorized", CodePermissionDenied)
		})
	})

})
