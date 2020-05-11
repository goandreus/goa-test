package design

import (
	"strings"

	. "goa.design/goa/v3/dsl"
)

// ErrorInfo represents information related to an error
var ErrorInfo = Type("ErrorInfo", func() {

	Description("ErrorInfo is the type used to represent an error.")

	Meta("type:generate:force", strings.ToLower(Domain)+"query", strings.ToLower(Domain)+"admin")
	Meta("type:category", "system")
	Meta("type:reference", "E0002")

	Field(1, "service", String, "The name of the service.", func() {
		Example("userquery")
	})

	Field(2, "instanceId", String, "The id of the service instance.", func() {
		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d9")
	})

	Field(3, "message", String, "The error message.", func() {
		Meta("field:size", "1024")
		MaxLength(1024)
		Example("Unsufficient disk space.")
	})

	Field(4, "details", String, "The details about the error.", func() {
		Meta("field:size", "8000")

		MaxLength(8000)
		Example("Operation cannot be completed because only 1 Mb is available on drive.")
	})

	Required("service", "instanceId", "message")
})
