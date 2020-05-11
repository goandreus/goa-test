package design

import (
	"strings"

	. "goa.design/goa/v3/dsl"
)

// SAMLConfig represents an SAML configuration
var SAMLConfig = Type("SAMLConfig", func() {

	Description("SAMLConfig is the type used to represents a SAML configuration.")

	Meta("type:generate:force", strings.ToLower(Domain)+"query", strings.ToLower(Domain)+"admin")
	Meta("type:stored", "true")
	Meta("type:category", "business")
	Meta("type:reference", "E0402")

	Field(1, "id", String, "The UUID of the underlying or referenced entity.", func() {
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

	Field(3, "createdAt", String, "The date/time the record was created.", func() {
		Meta("struct:tag:json", "createdAt,omitempty")
		Meta("field:size", "20")
		Meta("rpc:tag", "4")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Field(4, "updatedAt", String, "The date/time the record was last updated.", func() {
		Meta("struct:tag:json", "updatedAt,omitempty")
		Meta("field:size", "20")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Field(5, "active", Boolean, "Defines if the entity is active.", func() {
		Meta("struct:tag:json", "active,omitempty")
		Meta("dgraph:type", "bool")

		Default(true)
	})

	Field(6, "organizationId", String, "The id of the associated organization.", func() {
		Meta("struct:tag:json", "organizationId,omitempty")

		MaxLength(255)
		Example("WiserSKILLS")
	})

	Field(7, "areaId", String, "The id of the associated area.", func() {
		Meta("struct:tag:json", "areaId,omitempty")
		Enum("PROD", "DEV", "QA", "UAT")

		Example("PROD")
	})

	Field(8, "host", String, "The associated host/domain.", func() {
		Meta("struct:tag:json", "host,omitempty")
		Meta("cache:key", "1")

		Example("wiserskills.io")
	})

	Field(9, "idpMetadata", String, "The associated IDP metadata in XML.", func() {
		Meta("struct:tag:json", "idpMetadata")
	})

	Field(10, "idKey", String, "The key used to manage the user identification.", func() {
		Meta("struct:tag:json", "idKey")
	})

	Field(11, "callbackURL", String, "The URL to be called back by the IDP.", func() {
		Meta("struct:tag:json", "callbackURL")
	})

	Field(12, "redirectURL", String, "The URL the service redirects to once the user is authenticated.", func() {
		Meta("struct:tag:json", "redirectURL")
	})

	Required("id", "organizationId", "areaId", "host", "active", "idpMetadata", "idKey", "callbackURL", "redirectURL")
})

// HostPayload ...
var HostPayload = Type("HostPayload", func() {

	APIKeyField(1, "api_key", "key", String, func() {
		Description("API key")
		Example("abcdef12345")
	})

	Field(2, "host", String, func() {
		Description("The target host.")
	})

	Field(3, "orgId", String, func() {
		Description("The id of the organization")
		MaxLength(255)
	})

	Field(4, "areaId", String, func() {
		Description("The id of the area")
		Enum("PROD", "DEV", "QA", "UAT")
		MaxLength(20)
	})

	Required("key", "host", "orgId", "areaId")
})

// RedirectResult ...
var RedirectResult = Type("RedirectResult", func() {

	Field(1, "Location", String, func() {
		Description("The url to redirect the caller to.")
	})

	Required("Location")
})
