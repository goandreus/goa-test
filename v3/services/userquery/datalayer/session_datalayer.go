package datalayer

import (
	caching "gitlab.com/wiserskills/v3/services/userquery/gen/caching"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// GetSessionsByOrg returns all the sessions from the datastore for the specified orgID and areaID
func (s *Db) GetSessionsByOrg(orgID string, areaID string) ([]*userquery.Session, error) {

	// We retrieve it from the database
	vars := make(map[string]interface{})
	vars["orgID"] = orgID
	vars["areaID"] = areaID

	result, err := s.ExecuteSessionQuery("FOR s IN Session FILTER s.orgID==@orgID AND s.areaID==@areaID  RETURN s", vars)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetSessionByUserID returns the session associated with the specified userId
func (s *Db) GetSessionByUserID(userID string) (*userquery.Session, bool, error) {

	fromCache := false
	var result *userquery.Session

	// We first try to get it from the cache
	kb := caching.NewKeyBuilder()
	kb.Add("session")
	kb.Add(userID)
	cacheKey := kb.Get()

	r := s.cache.Get(cacheKey)

	if r != nil {
		t := r.(userquery.Session)
		result = &t
		fromCache = true
	} else {

		// We retrieve it from the database
		vars := make(map[string]interface{})
		vars["userID"] = userID

		items, err := s.ExecuteSessionQuery("FOR s IN Session FILTER s.userID==@userID RETURN s", vars)

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
