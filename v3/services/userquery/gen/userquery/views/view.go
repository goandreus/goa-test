// Code generated by goa v3.0.9, DO NOT EDIT.
//
// userquery views
//
// Command:
// $ goa gen gitlab.com/wiserskills/v3/services/userquery/design

package views

import (
	"unicode/utf8"

	goa "goa.design/goa/v3/pkg"
)

// SessionCollection is the viewed result type that is projected based on a
// view.
type SessionCollection struct {
	// Type to project
	Projected SessionCollectionView
	// View to render
	View string
}

// RegisteredUserCollection is the viewed result type that is projected based
// on a view.
type RegisteredUserCollection struct {
	// Type to project
	Projected RegisteredUserCollectionView
	// View to render
	View string
}

// SessionCollectionView is a type that runs validations on a projected type.
type SessionCollectionView []*SessionView

// SessionView is a type that runs validations on a projected type.
type SessionView struct {
	// The UUID of the session.
	ID *string `json:"id,omitempty"`
	// The id of entity in the document database.
	Key *string `json:"_key,omitempty"`
	// The id of the associated organization.
	OrganizationID *string `json:"organizationId,omitempty"`
	// The id of the associated area.
	AreaID *string `json:"areaId,omitempty"`
	// The UUID of the user.
	UserID *string `json:"userId,omitempty"`
	// The UUID of the token.
	TokenID *string `json:"tokenId,omitempty"`
	// The date/time the record was created.
	CreatedAt *string `json:"createdAt,omitempty"`
	// The date/time the record was updated.
	UpdatedAt *string `updatedAt:"uri,omitempty"`
	// The date/time the session will expire.
	ExpiresAt *string `json:"expiresAt,omitempty"`
}

// RegisteredUserCollectionView is a type that runs validations on a projected
// type.
type RegisteredUserCollectionView []*RegisteredUserView

// RegisteredUserView is a type that runs validations on a projected type.
type RegisteredUserView struct {
	// The UUID of the underlying or referenced entity.
	ID *string `json:"id,omitempty"`
	// The first name of the user.
	FirstName *string `json:"firstName,omitempty"`
	// The last name of the user.
	LastName *string `json:"lastName,omitempty" storm:"index"`
	// Defines if the entity is active.
	Active *bool `json:"active"`
	// The birth name of the user.
	BirthName *string `json:"birthName,omitempty"`
	// The address of the user.
	Address *string `json:"address,omitempty"`
	// The id of the city.
	CityID *string `json:"cityId,omitempty"`
	// The ISO 3166-1 code of the country.
	CountryID *string `json:"countryId,omitempty"`
	// The latitude of the user's address.
	Latitude *float64 `json:"latitude,omitempty"`
	// The longitude of the user's address.
	Longitude *float64 `json:"longitude,omitempty"`
	// The birth date of the user.
	BirthDate *string `json:"birthDate,omitempty"`
	// The gender of the user.
	Gender *string `json:"gender,omitempty"`
	// The ISO code of the user's prefered language.
	LanguageID *string `json:"languageId,omitempty"`
	// The email of the user.
	Email *string `json:"email,omitempty"`
	// The login of the user.
	Login *string `json:"login,omitempty"`
	// The user mobile phone number.
	Mobile *string `json:"mobile,omitempty"`
	// Defines if this user is a B2C user.
	B2C *bool `json:"b2c,omitempty"`
	// The date/time the record was created.
	CreatedAt *string `json:"createdAt,omitempty"`
	// The date/time the record was last updated.
	UpdatedAt *string `json:"updatedAt,omitempty"`
	// The id of organization where the user was created.
	OrganizationID *string `json:"organizationId,omitempty"`
	// The list of the user roles.
	Roles []string `json:"roles,omitempty"`
}

var (
	// SessionCollectionMap is a map of attribute names in result type
	// SessionCollection indexed by view name.
	SessionCollectionMap = map[string][]string{
		"default": []string{
			"id",
			"key",
			"organizationId",
			"areaId",
			"userId",
			"tokenId",
			"createdAt",
			"updatedAt",
			"expiresAt",
		},
		"tiny": []string{
			"id",
			"userId",
			"createdAt",
			"updatedAt",
			"expiresAt",
		},
	}
	// RegisteredUserCollectionMap is a map of attribute names in result type
	// RegisteredUserCollection indexed by view name.
	RegisteredUserCollectionMap = map[string][]string{
		"default": []string{
			"id",
			"firstName",
			"lastName",
			"active",
			"birthName",
			"address",
			"cityId",
			"countryId",
			"latitude",
			"longitude",
			"birthDate",
			"gender",
			"languageId",
			"email",
			"login",
			"mobile",
			"B2C",
			"roles",
			"organizationId",
			"createdAt",
			"updatedAt",
		},
		"tiny": []string{
			"id",
			"firstName",
			"lastName",
			"active",
			"email",
			"login",
			"roles",
			"organizationId",
		},
	}
	// SessionMap is a map of attribute names in result type Session indexed by
	// view name.
	SessionMap = map[string][]string{
		"default": []string{
			"id",
			"key",
			"organizationId",
			"areaId",
			"userId",
			"tokenId",
			"createdAt",
			"updatedAt",
			"expiresAt",
		},
		"tiny": []string{
			"id",
			"userId",
			"createdAt",
			"updatedAt",
			"expiresAt",
		},
	}
	// RegisteredUserMap is a map of attribute names in result type RegisteredUser
	// indexed by view name.
	RegisteredUserMap = map[string][]string{
		"default": []string{
			"id",
			"firstName",
			"lastName",
			"active",
			"birthName",
			"address",
			"cityId",
			"countryId",
			"latitude",
			"longitude",
			"birthDate",
			"gender",
			"languageId",
			"email",
			"login",
			"mobile",
			"B2C",
			"roles",
			"organizationId",
			"createdAt",
			"updatedAt",
		},
		"tiny": []string{
			"id",
			"firstName",
			"lastName",
			"active",
			"email",
			"login",
			"roles",
			"organizationId",
		},
	}
)

// ValidateSessionCollection runs the validations defined on the viewed result
// type SessionCollection.
func ValidateSessionCollection(result SessionCollection) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateSessionCollectionView(result.Projected)
	case "tiny":
		err = ValidateSessionCollectionViewTiny(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default", "tiny"})
	}
	return
}

// ValidateRegisteredUserCollection runs the validations defined on the viewed
// result type RegisteredUserCollection.
func ValidateRegisteredUserCollection(result RegisteredUserCollection) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateRegisteredUserCollectionView(result.Projected)
	case "tiny":
		err = ValidateRegisteredUserCollectionViewTiny(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default", "tiny"})
	}
	return
}

// ValidateSessionCollectionView runs the validations defined on
// SessionCollectionView using the "default" view.
func ValidateSessionCollectionView(result SessionCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateSessionView(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateSessionCollectionViewTiny runs the validations defined on
// SessionCollectionView using the "tiny" view.
func ValidateSessionCollectionViewTiny(result SessionCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateSessionViewTiny(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateSessionView runs the validations defined on SessionView using the
// "default" view.
func ValidateSessionView(result *SessionView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.UserID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("userId", "result"))
	}
	if result.TokenID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("tokenId", "result"))
	}
	if result.OrganizationID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("organizationId", "result"))
	}
	if result.AreaID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("areaId", "result"))
	}
	if result.ID != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.id", *result.ID, goa.FormatUUID))
	}
	if result.Key != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.key", *result.Key, goa.FormatUUID))
	}
	if result.OrganizationID != nil {
		if utf8.RuneCountInString(*result.OrganizationID) > 255 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.organizationId", *result.OrganizationID, utf8.RuneCountInString(*result.OrganizationID), 255, false))
		}
	}
	if result.AreaID != nil {
		if !(*result.AreaID == "PROD" || *result.AreaID == "DEV" || *result.AreaID == "QA" || *result.AreaID == "UAT") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.areaId", *result.AreaID, []interface{}{"PROD", "DEV", "QA", "UAT"}))
		}
	}
	if result.AreaID != nil {
		if utf8.RuneCountInString(*result.AreaID) < 2 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.areaId", *result.AreaID, utf8.RuneCountInString(*result.AreaID), 2, true))
		}
	}
	if result.AreaID != nil {
		if utf8.RuneCountInString(*result.AreaID) > 4 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.areaId", *result.AreaID, utf8.RuneCountInString(*result.AreaID), 4, false))
		}
	}
	if result.UserID != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.userId", *result.UserID, goa.FormatUUID))
	}
	if result.TokenID != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.tokenId", *result.TokenID, goa.FormatUUID))
	}
	if result.CreatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.createdAt", *result.CreatedAt, goa.FormatDateTime))
	}
	if result.UpdatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.updatedAt", *result.UpdatedAt, goa.FormatDateTime))
	}
	if result.ExpiresAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.expiresAt", *result.ExpiresAt, goa.FormatDateTime))
	}
	return
}

// ValidateSessionViewTiny runs the validations defined on SessionView using
// the "tiny" view.
func ValidateSessionViewTiny(result *SessionView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.UserID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("userId", "result"))
	}
	if result.ID != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.id", *result.ID, goa.FormatUUID))
	}
	if result.UserID != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.userId", *result.UserID, goa.FormatUUID))
	}
	if result.CreatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.createdAt", *result.CreatedAt, goa.FormatDateTime))
	}
	if result.UpdatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.updatedAt", *result.UpdatedAt, goa.FormatDateTime))
	}
	if result.ExpiresAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.expiresAt", *result.ExpiresAt, goa.FormatDateTime))
	}
	return
}

// ValidateRegisteredUserCollectionView runs the validations defined on
// RegisteredUserCollectionView using the "default" view.
func ValidateRegisteredUserCollectionView(result RegisteredUserCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateRegisteredUserView(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateRegisteredUserCollectionViewTiny runs the validations defined on
// RegisteredUserCollectionView using the "tiny" view.
func ValidateRegisteredUserCollectionViewTiny(result RegisteredUserCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateRegisteredUserViewTiny(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateRegisteredUserView runs the validations defined on
// RegisteredUserView using the "default" view.
func ValidateRegisteredUserView(result *RegisteredUserView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.FirstName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("firstName", "result"))
	}
	if result.LastName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("lastName", "result"))
	}
	if result.Login == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("login", "result"))
	}
	if result.Email == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("email", "result"))
	}
	if result.BirthName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("birthName", "result"))
	}
	if result.OrganizationID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("organizationId", "result"))
	}
	if result.ID != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.id", *result.ID, goa.FormatUUID))
	}
	if result.FirstName != nil {
		if utf8.RuneCountInString(*result.FirstName) > 255 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.firstName", *result.FirstName, utf8.RuneCountInString(*result.FirstName), 255, false))
		}
	}
	if result.LastName != nil {
		if utf8.RuneCountInString(*result.LastName) > 255 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.lastName", *result.LastName, utf8.RuneCountInString(*result.LastName), 255, false))
		}
	}
	if result.BirthName != nil {
		if utf8.RuneCountInString(*result.BirthName) > 255 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.birthName", *result.BirthName, utf8.RuneCountInString(*result.BirthName), 255, false))
		}
	}
	if result.Address != nil {
		if utf8.RuneCountInString(*result.Address) > 1024 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.address", *result.Address, utf8.RuneCountInString(*result.Address), 1024, false))
		}
	}
	if result.CountryID != nil {
		if utf8.RuneCountInString(*result.CountryID) < 2 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.countryId", *result.CountryID, utf8.RuneCountInString(*result.CountryID), 2, true))
		}
	}
	if result.CountryID != nil {
		if utf8.RuneCountInString(*result.CountryID) > 2 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.countryId", *result.CountryID, utf8.RuneCountInString(*result.CountryID), 2, false))
		}
	}
	if result.BirthDate != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.birthDate", *result.BirthDate, goa.FormatDate))
	}
	if result.Gender != nil {
		if !(*result.Gender == "M" || *result.Gender == "F") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.gender", *result.Gender, []interface{}{"M", "F"}))
		}
	}
	if result.Gender != nil {
		if utf8.RuneCountInString(*result.Gender) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.gender", *result.Gender, utf8.RuneCountInString(*result.Gender), 1, true))
		}
	}
	if result.Gender != nil {
		if utf8.RuneCountInString(*result.Gender) > 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.gender", *result.Gender, utf8.RuneCountInString(*result.Gender), 1, false))
		}
	}
	if result.LanguageID != nil {
		if utf8.RuneCountInString(*result.LanguageID) < 2 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.languageId", *result.LanguageID, utf8.RuneCountInString(*result.LanguageID), 2, true))
		}
	}
	if result.LanguageID != nil {
		if utf8.RuneCountInString(*result.LanguageID) > 2 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.languageId", *result.LanguageID, utf8.RuneCountInString(*result.LanguageID), 2, false))
		}
	}
	if result.Email != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.email", *result.Email, goa.FormatEmail))
	}
	if result.Email != nil {
		if utf8.RuneCountInString(*result.Email) > 512 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.email", *result.Email, utf8.RuneCountInString(*result.Email), 512, false))
		}
	}
	if result.Login != nil {
		if utf8.RuneCountInString(*result.Login) > 255 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.login", *result.Login, utf8.RuneCountInString(*result.Login), 255, false))
		}
	}
	if result.Mobile != nil {
		if utf8.RuneCountInString(*result.Mobile) > 20 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.mobile", *result.Mobile, utf8.RuneCountInString(*result.Mobile), 20, false))
		}
	}
	if result.CreatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.createdAt", *result.CreatedAt, goa.FormatDateTime))
	}
	if result.UpdatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.updatedAt", *result.UpdatedAt, goa.FormatDateTime))
	}
	return
}

// ValidateRegisteredUserViewTiny runs the validations defined on
// RegisteredUserView using the "tiny" view.
func ValidateRegisteredUserViewTiny(result *RegisteredUserView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.FirstName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("firstName", "result"))
	}
	if result.LastName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("lastName", "result"))
	}
	if result.Login == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("login", "result"))
	}
	if result.Email == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("email", "result"))
	}
	if result.OrganizationID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("organizationId", "result"))
	}
	if result.ID != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.id", *result.ID, goa.FormatUUID))
	}
	if result.FirstName != nil {
		if utf8.RuneCountInString(*result.FirstName) > 255 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.firstName", *result.FirstName, utf8.RuneCountInString(*result.FirstName), 255, false))
		}
	}
	if result.LastName != nil {
		if utf8.RuneCountInString(*result.LastName) > 255 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.lastName", *result.LastName, utf8.RuneCountInString(*result.LastName), 255, false))
		}
	}
	if result.Email != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.email", *result.Email, goa.FormatEmail))
	}
	if result.Email != nil {
		if utf8.RuneCountInString(*result.Email) > 512 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.email", *result.Email, utf8.RuneCountInString(*result.Email), 512, false))
		}
	}
	if result.Login != nil {
		if utf8.RuneCountInString(*result.Login) > 255 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("result.login", *result.Login, utf8.RuneCountInString(*result.Login), 255, false))
		}
	}
	return
}
