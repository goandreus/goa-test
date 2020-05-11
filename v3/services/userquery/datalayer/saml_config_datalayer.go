package datalayer

import (
	caching "gitlab.com/wiserskills/v3/services/userquery/gen/caching"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// GetSAMLConfigByHost returns the SAML config associated with the specified host
func (s *Db) GetSAMLConfigByHost(host string) (*userquery.SAMLConfig, bool, error) {

	fromCache := false
	var result *userquery.SAMLConfig

	// We first try to get it from the cache
	kb := caching.NewKeyBuilder()
	kb.Add("samlconfig")
	kb.Add(host)
	cacheKey := kb.Get()

	r := s.cache.Get(cacheKey)

	if r != nil {
		t := r.(userquery.SAMLConfig)
		result = &t
		fromCache = true
	} else {

		// We retrieve the user from the database
		vars := make(map[string]interface{})
		vars["host"] = host

		items, err := s.ExecuteSAMLConfigQuery("FOR c IN SAMLConfig FILTER c.host==@host AND c.active==true RETURN c", vars)

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
