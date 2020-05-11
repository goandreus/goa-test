package design

import (
	"strings"

	. "goa.design/goa/v3/dsl"
)

// Session represents a session object
var Session = ResultType("Session", func() {

	TypeName("Session")
	Description("Session is the type used to represents an application session.")

	Meta("type:generate:force", strings.ToLower(Domain)+"query", strings.ToLower(Domain)+"admin")
	Meta("type:stored", "true")
	Meta("dgraph:stored", "false")
	Meta("type:category", "business")
	Meta("type:reference", "E0405")

	Attributes(func() {

		Field(1, "id", String, "The UUID of the session.", func() {
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

		Field(3, "organizationId", String, "The id of the associated organization.", func() {
			Meta("struct:tag:json", "organizationId,omitempty")

			MaxLength(255)
			Example("WiserSKILLS")
		})

		Field(4, "areaId", String, "The id of the associated area.", func() {
			Meta("struct:tag:json", "areaId,omitempty")
			Meta("field:size", "20")

			Enum("PROD", "DEV", "QA", "UAT")
			MinLength(2)
			MaxLength(4)
		})

		Field(5, "userId", String, "The UUID of the user.", func() {
			Meta("struct:tag:json", "userId,omitempty")
			Meta("field:size", "36")
			Meta("cache:key", "1")

			Format(FormatUUID)
			Example("1f052270-dec6-4930-a042-eaec722407d1")
		})

		Field(6, "tokenId", String, "The UUID of the token.", func() {
			Meta("struct:tag:json", "tokenId,omitempty")
			Meta("field:size", "36")

			Format(FormatUUID)
			Example("1f052270-dec6-4930-a042-eaec722407d1")
		})

		Field(7, "createdAt", String, "The date/time the record was created.", func() {
			Meta("struct:tag:json", "createdAt,omitempty")
			Meta("field:size", "20")

			Format(FormatDateTime)
			Example("2019-01-01T12:52:24Z")
		})

		Field(8, "updatedAt", String, "The date/time the record was updated.", func() {
			Meta("struct:tag:updatedAt", "uri,omitempty")
			Meta("field:size", "20")

			Format(FormatDateTime)
			Example("2019-01-01T12:52:24Z")
		})

		Field(9, "expiresAt", String, "The date/time the session will expire.", func() {
			Meta("struct:tag:json", "expiresAt,omitempty")
			Meta("field:size", "20")

			Format(FormatDateTime)
			Example("2019-01-01T12:52:24Z")
		})

		Required("id", "userId", "tokenId", "organizationId", "areaId")
	})

	View("default", func() {
		Attribute("id")
		Attribute("key")
		Attribute("organizationId")
		Attribute("areaId")
		Attribute("userId")
		Attribute("tokenId")
		Attribute("createdAt")
		Attribute("updatedAt")
		Attribute("expiresAt")
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("userId")
		Attribute("createdAt")
		Attribute("updatedAt")
		Attribute("expiresAt")
	})
})

// AuthorizePayload ...
var AuthorizePayload = Type("AuthorizePayload", func() {

	Meta("type:category", "payload")
	Meta("type:reference", "P0403")

	TokenField(1, "token", String, func() {
		Description("JWT used for authentication")
	})

	Field(2, "path", String, func() {
		Description("The path to authorize")
	})

	Field(3, "orgId", String, func() {
		Description("The id of the organization")
		MaxLength(255)
	})

	Field(4, "areaId", String, func() {
		Description("The id of the area")
		Enum("PROD", "DEV", "QA", "UAT")
	})

	Required("path", "orgId", "areaId")
})

// AllSessionsPayload ...
var AllSessionsPayload = Type("AllSessionsPayload", func() {

	Meta("type:category", "payload")
	Meta("type:reference", "P0404")

	APIKeyField(1, "api_key", "key", String, func() {
		Description("API key")
		Example("abcdef12345")
	})

	TokenField(2, "token", String, func() {
		Description("JWT used for authentication")
		MaxLength(1024)
	})

	Field(3, "orgId", String, func() {
		Description("The id of the organization")
		MaxLength(255)
	})

	Field(4, "areaId", String, func() {
		Description("The id of the area")
		Enum("PROD", "DEV", "QA", "UAT")
	})

	Field(5, "view", String, func() {
		Description("View to render.")
		Enum("default", "tiny")
	})

	Required("key", "token", "orgId", "areaId")
})
