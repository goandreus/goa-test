package datalayer

import (
	caching "gitlab.com/wiserskills/v3/services/userquery/gen/caching"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// GetUserByLogin returns the user associated with the specified login and password
func (s *Db) GetUserByLogin(login string) (*userquery.User, bool, error) {

	fromCache := false
	var result *userquery.User

	// We first try to get it from the cache
	kb := caching.NewKeyBuilder()
	kb.Add("user")
	kb.Add(login)
	cacheKey := kb.Get()

	r := s.cache.Get(cacheKey)

	if r != nil {
		t := r.(userquery.User)
		result = &t
		fromCache = true
	} else {

		// We retrieve it from the database
		vars := make(map[string]interface{})
		vars["login"] = login

		items, err := s.ExecuteUserQuery("FOR u IN User FILTER u.login==@login AND u.active == true RETURN u", vars)

		if err != nil {
			return nil, fromCache, err
		}

		if len(items) > 0 {
			result = items[0]

			// We add the result to cache
			s.cache.Set(cacheKey, *result)
		}
	}

	return result, fromCache, nil
}

// GetRegisteredUsersByID returns the registered users associated with the specified ids
func (s *Db) GetRegisteredUsersByID(ids []string, activeOnly bool) (userquery.RegisteredUserCollection, error) {

	// We retrieve it from the database
	vars := make(map[string]interface{})
	vars["ids"] = ids
	query := "FOR u IN User FILTER u.id IN @ids AND u.active == true RETURN u"

	if !activeOnly {
		query = "FOR u IN User FILTER u.id IN @ids RETURN u"
	}

	users, err := s.ExecuteUserQuery(query, vars)

	if err != nil {
		return nil, err
	}

	result := make(userquery.RegisteredUserCollection, len(users))
	for i, user := range users {
		result[i] = UserToRegisteredUser(user)
	}

	return result, nil
}

// UserToRegisteredUser converts a user to a registered user
func UserToRegisteredUser(user *userquery.User) *userquery.RegisteredUser {

	return &userquery.RegisteredUser{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		BirthName:  user.BirthName,
		Active:     user.Active,
		Address:    user.Address,
		CityID:     user.CityID,
		CountryID:  user.CountryID,
		Latitude:   user.Latitude,
		Longitude:  user.Longitude,
		BirthDate:  user.BirthDate,
		Gender:     user.Gender,
		LanguageID: user.LanguageID,
		Login:      user.Login,
		Mobile:     user.Mobile,
		B2C:        user.B2C,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		Roles:      user.Roles,
	}
}
