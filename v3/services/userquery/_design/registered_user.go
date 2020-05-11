package design

import . "goa.design/goa/v3/dsl"

// RegisteredUser represents a limited view of the user entity
var RegisteredUser = ResultType("RegisteredUser", func() {

	Description("Represents a registered user.")

	Reference(User)

	TypeName("RegisteredUser")

	Meta("type:category", "business")
	Meta("type:reference", "E0409")

	Attributes(func() {
		Attribute("id")
		Attribute("firstName")
		Attribute("lastName")
		Attribute("active")
		Attribute("birthName")
		Attribute("address")
		Attribute("cityId")
		Attribute("countryId")
		Attribute("latitude")
		Attribute("longitude")
		Attribute("birthDate")
		Attribute("gender")
		Attribute("languageId")
		Attribute("email")
		Attribute("login")
		Attribute("mobile")
		Attribute("B2C")
		Attribute("createdAt")
		Attribute("updatedAt")
		Attribute("organizationId")
		Attribute("roles")
	})

	View("default", func() {
		Attribute("id")
		Attribute("firstName")
		Attribute("lastName")
		Attribute("active")
		Attribute("birthName")
		Attribute("address")
		Attribute("cityId")
		Attribute("countryId")
		Attribute("latitude")
		Attribute("longitude")
		Attribute("birthDate")
		Attribute("gender")
		Attribute("languageId")
		Attribute("email")
		Attribute("login")
		Attribute("mobile")
		Attribute("B2C")
		Attribute("roles")
		Attribute("organizationId")
		Attribute("createdAt")
		Attribute("updatedAt")
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("firstName")
		Attribute("lastName")
		Attribute("active")
		Attribute("email")
		Attribute("login")
		Attribute("roles")
		Attribute("organizationId")
	})

	Required("id", "firstName", "lastName", "login", "email")
})
