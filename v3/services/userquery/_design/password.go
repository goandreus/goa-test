package design

import (
	"strings"

	. "goa.design/goa/v3/dsl"
)

// PasswordInfo represents a password
var PasswordInfo = Type("Password", func() {

	Description("Password is the type used to represent a password.")

	Meta("type:generate:force", strings.ToLower(Domain)+"query", strings.ToLower(Domain)+"admin")
	Meta("type:category", "business")
	Meta("type:reference", "E0407")
	Meta("type:stored", "true")
	Meta("dgraph:stored", "false")

	Field(1, "id", String, "The UUID of the password.", func() {
		Meta("struct:tag:json", "id,omitempty")
		Meta("field:size", "36")

		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d1")
	})

	Field(2, "userId", String, "The UUID of the associated user.", func() {
		Meta("struct:tag:json", "userId,omitempty")
		Meta("field:size", "36")

		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d1")
	})

	Field(3, "encryptedPassword", String, "The encrypted password value.", func() {
		Meta("struct:tag:json", "encryptedPassword,omitempty")
		Meta("field:size", "255")

		MaxLength(255)
	})

	Field(4, "createdAt", String, "The date/time the record was created.", func() {
		Meta("struct:tag:json", "createdAt,omitempty")
		Meta("field:size", "20")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Required("id", "userId", "encryptedPassword")
})
