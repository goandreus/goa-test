package design

import (
	"strings"

	. "goa.design/goa/v3/dsl"
)

// User represents a user entity
var User = Type("User", func() {

	Description("User is the type used to represents a system user.")

	Meta("type:generate:force", strings.ToLower(Domain)+"query", strings.ToLower(Domain)+"admin")
	Meta("type:category", "business")
	Meta("type:reference", "E0401")

	Meta("type:stored", "true")

	Field(1, "id", String, "The UUID of the underlying or referenced entity.", func() {
		Meta("struct:tag:json", "id,omitempty")
		Meta("field:size", "36")

		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d1")
	})

	Field(2, "uid", String, "The UID of the underlying or referenced entity in DGraph.", func() {
		Meta("struct:tag:json", "uid,omitempty")
		Meta("field:size", "40")

		Example("Ox123")
	})

	Field(3, "key", String, "The UUID of the underlying or referenced entity in the document database.", func() {
		Meta("struct:tag:json", "_key,omitempty")
		Meta("field:size", "36")

		Format(FormatUUID)
		Example("1f052270-dec6-4930-a042-eaec722407d1")
	})

	Field(4, "createdAt", String, "The date/time the record was created.", func() {
		Meta("struct:tag:json", "createdAt,omitempty")
		Meta("field:size", "20")
		Meta("dgraph:type", "datetime")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Field(5, "updatedAt", String, "The date/time the record was last updated.", func() {
		Meta("struct:tag:json", "updatedAt,omitempty")
		Meta("field:size", "20")
		Meta("dgraph:type", "datetime")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Field(6, "active", Boolean, "Defines if the entity is active.", func() {
		Meta("struct:tag:json", "active")
		Meta("dgraph:type", "bool")

		Default(true)
	})

	Field(7, "firstName", String, "The first name of the user.", func() {
		Meta("struct:tag:json", "firstName,omitempty")
		Meta("field:size", "255")
		Meta("dgraph:stored", "false")

		MaxLength(255)
		Example("John")
	})

	Field(8, "lastName", String, "The last name of the user.", func() {
		Meta("struct:tag:json", "lastName,omitempty")
		Meta("struct:tag:storm", "index")
		Meta("field:size", "255")
		Meta("dgraph:stored", "false")

		MaxLength(255)
		Example("SMITH")
	})

	Field(9, "birthName", String, "The birth name of the user.", func() {
		Meta("struct:tag:json", "birthName,omitempty")
		Meta("field:size", "255")
		Meta("rpc:tag", "9")
		Meta("dgraph:stored", "false")

		MaxLength(255)
		Example("SMITH")
	})

	Field(10, "address", String, "The address of the user.", func() {
		Meta("struct:tag:json", "address,omitempty")
		Meta("field:size", "1024")
		Meta("dgraph:stored", "false")

		MaxLength(1024)
		Example("8, rue de l'Hôtel de Ville")
	})

	Field(11, "cityId", String, "The id of the city.", func() {
		Meta("struct:tag:json", "cityId,omitempty")
		Meta("field:size", "255")

		Example("FR92220")
	})

	Field(12, "countryId", String, "The ISO 3166-1 code of the country.", func() {
		Meta("struct:tag:json", "countryId,omitempty")
		Meta("field:size", "2")

		MinLength(2)
		MaxLength(2)
		Example("FR")
	})

	Field(13, "latitude", Float64, "The latitude of the user's address.", func() {
		Meta("struct:tag:json", "latitude,omitempty")

		Example(41.40338)
	})

	Field(14, "longitude", Float64, "The longitude of the user's address.", func() {
		Meta("struct:tag:json", "longitude,omitempty")

		Example(2.17403)
	})

	Field(15, "birthDate", String, "The birth date of the user.", func() {
		Meta("struct:tag:json", "birthDate,omitempty")
		Meta("field:size", "10")

		Format(FormatDate)
		Example("1982-01-01")
	})

	Field(16, "gender", String, "The gender of the user.", func() {
		Meta("struct:tag:json", "gender,omitempty")
		Enum("M", "F")

		MinLength(1)
		MaxLength(1)
		Example("M")
	})

	Field(17, "languageId", String, "The ISO code of the user's prefered language.", func() {
		Meta("struct:tag:json", "languageId,omitempty")
		Meta("field:size", "2")

		MinLength(2)
		MaxLength(2)
		Example("en")
	})

	Field(18, "email", String, "The email of the user.", func() {
		Meta("struct:tag:json", "email,omitempty")
		Meta("field:size", "512")
		Meta("dgraph:stored", "false")

		Format(FormatEmail)
		MaxLength(512)
		Example("john.smith@google.com")
	})

	Field(19, "login", String, "The login of the user.", func() {
		Meta("struct:tag:json", "login,omitempty")
		Meta("dgraph:stored", "false")
		Meta("cache:key", "1")

		MaxLength(255)
	})

	Field(20, "encryptedPassword", String, "The encrypted password of the user.", func() {
		Meta("struct:tag:json", "encryptedPassword,omitempty")
		Meta("field:size", "255")
		Meta("dgraph:stored", "false")

		MaxLength(255)
	})

	Field(21, "passwordExpiresAt", String, "The date/time the password will expire.", func() {
		Meta("struct:tag:json", "passwordExpiresAt,omitempty")
		Meta("field:size", "20")
		Meta("dgraph:type", "datetime")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Field(22, "emailValidatedAt", String, "The date/time the email was validated.", func() {
		Meta("struct:tag:json", "emailValidatedAt,omitempty")
		Meta("field:size", "20")
		Meta("dgraph:type", "datetime")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Field(23, "suspendedUpTo", String, "The date/time the user is suspended up to.", func() {

		Meta("struct:tag:json", "suspendedUpTo,omitempty")
		Meta("field:size", "20")
		Meta("dgraph:type", "datetime")
		Meta("dgraph:stored", "false")

		Format(FormatDateTime)
		Example("2019-01-01T12:52:24Z")
	})

	Field(24, "failedAttempts", Int, "The number of failed authentication attempts.", func() {
		Meta("struct:tag:json", "fialedAttempts,omitempty")
		Meta("rpc:tag", "26")
		Meta("dgraph:stored", "false")
	})

	Field(25, "mobile", String, "The user mobile phone number.", func() {
		Meta("struct:tag:json", "mobile,omitempty")
		Meta("field:size", "20")
		Meta("dgraph:stored", "false")

		MaxLength(20)
		Example("+33645675637")
	})

	Field(26, "B2C", Boolean, "Defines if this user is a B2C user.", func() {
		Meta("struct:tag:json", "b2c,omitempty")

		Default(false)
	})

	Field(27, "organizationId", String, "The id of organization where the user was created.", func() {
		Meta("struct:tag:json", "organizationId,omitempty")
		Meta("rpc:tag", "30")

		Example("WiserSKILLS")
	})

	Field(28, "roles", ArrayOf(String), "The list of the user roles.", func() {
		Meta("struct:tag:json", "roles,omitempty")
	})

	Required("id", "firstName", "lastName", "birthName", "email", "login", "encryptedPassword", "organizationId")
})

// ManyUserIDPayload ...
var ManyUserIDPayload = Type("ManyUserIDPayload", func() {

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
		MaxLength(20)
	})

	Field(5, "ids", ArrayOf(String), func() {
		Description("The ids of the users to retrieve")
	})

	Field(6, "view", String, func() {
		Description("The view used for the returned users")

		Default("default")
	})

	Field(7, "activeOnly", Boolean, func() {
		Description("defines if only active users must be returned")

		Default(true)
	})

	Required("key", "token", "orgId", "areaId", "ids")
})
