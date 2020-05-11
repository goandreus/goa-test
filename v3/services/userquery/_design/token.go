package design

import (
	"strings"

	. "goa.design/goa/v3/dsl"
)

// TokenInfo represents a session token
var TokenInfo = Type("Token", func() {

	Description("Token is the type used to represents a JWT token.")

	Meta("type:generate:force", strings.ToLower(Domain)+"query", strings.ToLower(Domain)+"admin")
	Meta("type:category", "business")
	Meta("type:reference", "E0404")

	Meta("type:stored", "true")

	Field(1, "id", String, "The UUID of the token.", func() {
		Meta("struct:tag:json", "id,omitempty")
		Meta("field:size", "36")

		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d1")
	})

	Field(2, "key", String, "The id of entity in the document database.", func() {
		Meta("struct:tag:json", "_key,omitempty")
		Meta("field:size", "36")

		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d1")
	})

	Field(3, "type", Int, "The type of the token.", func() {
		Meta("struct:tag:json", "type,omitempty")

		Default(1)

		Example(1)
	})

	Field(4, "userId", String, "The UUID of the associated user.", func() {
		Meta("struct:tag:json", "userId,omitempty")
		Meta("field:size", "36")
		Meta("cache:key", "1")

		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d2")
	})

	Field(5, "organizationId", String, "The id of the associated organization.", func() {
		Meta("struct:tag:json", "organizationId,omitempty")
	})

	Field(6, "areaId", String, "The id of the associated area.", func() {
		Meta("struct:tag:json", "areaId,omitempty")
		Meta("field:size", "20")

		Enum("PROD", "UAT", "QA", "DEV")
		Example("PROD")
	})

	Field(7, "token", String, "The value of the token.", func() {
		Meta("struct:tag:storm", "index")
		Meta("struct:tag:json", "token")
		Meta("field:size", "1024")

		MaxLength(1024)
	})

	Field(8, "createdAt", String, "The date/time the record was created.", func() {
		Meta("struct:tag:json", "createdAt,omitempty")
		Meta("field:size", "20")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Field(9, "expiresAt", String, "The date/time the token will expire.", func() {
		Meta("struct:tag:json", "expiresAt,omitempty")
		Meta("field:size", "20")

		Format(FormatDateTime)
		Example("2019-01-03T12:52:24Z")
	})

	Required("id", "type", "userId", "organizationId", "areaId", "token")
})

// TokenPayload ...
var TokenPayload = Type("TokenPayload", func() {

	Meta("type:category", "payload")
	Meta("type:reference", "P0401")

	APIKeyField(1, "api_key", "key", String, func() {
		Description("API key")
		Example("abcdef12345")
	})

	TokenField(2, "token", String, func() {
		Description("JWT used for authentication")
	})

	Field(3, "orgId", String, func() {
		Description("The id of the organization")
		MaxLength(255)
	})

	Field(4, "areaId", String, func() {
		Description("The id of the area")
		Enum("PROD", "DEV", "QA", "UAT")
		MinLength(2)
		MaxLength(4)
	})

	Required("key", "token", "orgId", "areaId")
})
