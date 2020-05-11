package design

import (
	. "goa.design/goa/v3/dsl"
)

// Security Constants ******************************************************************************

// BasicAuth defines a security scheme using basic authentication. The scheme
// protects the "signin" action used to create JWTs.
var BasicAuth = BasicAuthSecurity("basic", func() {
	Description("Basic authentication used to authenticate security principal during signin.")
})

// JWTToken defines the credentials to use for authenticating to service methods.
var JWTToken = Type("JWTToken", func() {
	Field(1, "token", String, "JWT token", func() {
		Example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
	})
	Required("token")
})

// JWTAuth defines a security scheme that uses JWT tokens.
var JWTAuth = JWTSecurity("jwt", func() {
	Description(`Secures endpoints by requiring a valid JWT token retrieved via the signin endpoint.`)
	Scope("api:internal")    // Internal users scope
	Scope("api:b2c")         // B2C users
	Scope("api:b2b:[ORGID]") // B2B users
})

// APIKeyAuth defines a security scheme that uses API keys.
var APIKeyAuth = APIKeySecurity("api_key", func() {
	Description("Secures endpoints by requiring an API key.")
})

// API Server Definition *******************************************************************************

// API describes the global properties of the API server.
var _ = API(LowerServiceName, func() {
	Title(ServiceName)
	Description(ServiceDescription)
	Version(VersionCode)
	Contact(func() {
		Name(ServiceContact)
		Email(ServiceContactEmail)
	})

	Meta("api:domain", Domain)
	Meta("api:reference", "S"+DomainCode+ServiceNumber)
	Meta("api:repository", "https://gitlab.com/wiserskills/v3/services/"+LowerServiceName)
	Meta("api:namespace", "Services")
	Meta("api:storage", "NoSQL (Ristretto Cache + ArangoDB)")

	// Server describes a single process listening for client requests. The DSL
	// defines the set of services that the server hosts as well as hosts details.
	Server(LowerServiceName, func() {

		Description(LowerServiceName + " hosts the " + ServiceName + " service.")

		// List the services hosted by this server.
		Services(LowerServiceName)

		// List the Hosts and their transport URLs.
		Host(LowerServiceName, func() {
			// Transport specific URLs, supported schemes are:
			// 'http', 'https', 'grpc' and 'grpcs' with the respective default
			// ports: 80, 443, 8080, 8443.
			URI("http://0.0.0.0:8080")
			URI("grpc://0.0.0.0:8090")
		})
	})
})
