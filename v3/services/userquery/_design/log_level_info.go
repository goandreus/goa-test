package design

import (
	"strings"

	. "goa.design/goa/v3/dsl"
)

// LogLevelInfo represents information related to a log level
var LogLevelInfo = Type("LogLevelInfo", func() {

	Description("LogLevelInfo is the type used to specify a level of log.")

	Meta("type:generate:force", strings.ToLower(Domain)+"query", strings.ToLower(Domain)+"admin")
	Meta("type:category", "system")
	Meta("type:reference", "E0005")

	Field(1, "service", String, "The name of the service", func() {
		Example("userquery")
	})

	Field(2, "instanceId", String, "The id of the service instance", func() {
		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d9")
	})

	Field(3, "level", String, "The required log level. Possible values: DEBUG, INFO, WARN, ERROR", func() {
		Enum("DEBUG", "INFO", "WARN", "ERROR")
		Example("DEBUG")
	})

	Required("service", "instanceId", "level")
})
