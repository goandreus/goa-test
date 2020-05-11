package design

import (
	"strings"

	. "goa.design/goa/v3/dsl"
)

// AuthFailure represents information about a failed authentication
var AuthFailure = Type("AuthFailure", func() {

	Description("AuthFailure is the type used to describe a failed authentication.")

	Meta("type:generate:force", strings.ToLower(Domain)+"query", strings.ToLower(Domain)+"admin")
	Meta("type:category", "business")
	Meta("type:reference", "E0408")

	Field(1, "login", String, "The login used for the authentication.", func() {
		Meta("field:size", "255")

		MaxLength(255)
		Example("mpolo")
	})

	Field(2, "organizationId", String, "The id of the associated organization.", func() {
		Meta("struct:tag:json", "organizationId,omitempty")
		Meta("field:size", "255")

		MaxLength(255)
		Example("wiserskills")
	})

	Field(3, "areaId", String, "The id of the associated area.", func() {
		Meta("struct:tag:json", "areaId,omitempty")
		Meta("field:size", "20")

		Enum("PROD", "DEV", "QA", "UAT")
	})

	Field(4, "createdAt", String, "The date/time the record was created.", func() {
		Meta("field:size", "20")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Required("login", "organizationId", "areaId", "createdAt")
})
