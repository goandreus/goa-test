package datalayer

import (
	caching "gitlab.com/wiserskills/v3/services/userquery/gen/caching"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// GetTokenByUserID returns the token with the specified id
func (s *Db) GetTokenByUserID(userID string) (*userquery.Token, bool, error) {

	var result *userquery.Token
	fromCache := false

	// We first try to get it from the cache
	kb := caching.NewKeyBuilder()
	kb.Add("token")
	kb.Add(userID)
	cacheKey := kb.Get()

	r := s.cache.Get(cacheKey)

	if r != nil {
		t := r.(userquery.Token)
		result = &t
		fromCache = true
	} else {

		// We retrieve it from the database
		vars := make(map[string]interface{})
		vars["userID"] = userID

		items, err := s.ExecuteTokenQuery("FOR t IN Token FILTER t.userID==@userID RETURN t", vars)

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
