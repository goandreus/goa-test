package design

import (
	"strings"
)

// Service Constants ******************************************************************************

// ServiceName represents the name of the service
var ServiceName = "UserQuery"

// LowerServiceName represents the name of the service in lower case
var LowerServiceName = strings.ToLower(ServiceName)

// ServiceDescription represents the description of the service
var ServiceDescription = "API service to authenticate users. The service requires an ArangoDB server. Data synchronization with other services is done via a shared NATS Streaming event bus."

// VersionCode is the version of the API
var VersionCode = "1.0.0"

// Domain represents the domain of the microservice
var Domain = "User"

// DomainCode represents the code of the domain
var DomainCode = "04"

// ServiceType represents the type of the service (query or admin)
var ServiceType = "query"

// ServiceTypeCode represents the code of the service type
var ServiceTypeCode = "Q"

// ServiceNumber represents the number of the service within the domain
var ServiceNumber = "02"

// ServiceContact represents the contact for this service
var ServiceContact = "Christophe PEILLET (CTO)"

// ServiceContactEmail represents the email of the contact for this service
var ServiceContactEmail = "christophe.peillet@wiserskills.com"
