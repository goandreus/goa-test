package design

import (
	"strings"

	. "goa.design/goa/v3/dsl"
)

// ServiceStatusInfo provides information about the status of a service
var ServiceStatusInfo = Type("ServiceStatusInfo", func() {

	Description("ServiceStatusInfo is the type used to hold the service status information.")

	Meta("type:generate:force", strings.ToLower(Domain)+"query", strings.ToLower(Domain)+"admin")
	Meta("type:category", "system")
	Meta("type:reference", "E0001")

	Field(1, "service", String, "The name of the service", func() {
		Example("userquery")
	})

	Field(2, "instanceId", String, "The id of the service instance", func() {
		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d9")
	})

	Field(3, "status", String, "The new status of the service. Possible values: Initializing, Connecting, Pending, Syncing, OK", func() {
		Enum("Initializing", "Connecting", "Pending", "Syncing", "OK")
		Example("OK")
	})

	Required("service", "instanceId", "status")
})
