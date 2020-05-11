package datalayer

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	arangodb "github.com/arangodb/go-driver"
	arangodbhttp "github.com/arangodb/go-driver/http"
	jsoniter "github.com/json-iterator/go"
	caching "gitlab.com/wiserskills/v3/services/userquery/gen/caching"
	userquerysvc "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// Datastore represents the embedded database
type Datastore struct {
	client    arangodb.Client
	databases map[string]*Db
}

// Db represents a database
type Db struct {
	Name        string
	db          arangodb.Database
	collections map[string]arangodb.Collection
	cache       *caching.Cache
}

func finalizer(s *Datastore) {
	s.Close()
}

func getEnv(key string, defaultValue string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

// Open opens the datastore
func (s *Datastore) Open() error {

	var err error

	paths := getEnv("DB_CONNECTION", "http://localhost:8529")
	endpoints := strings.Split(paths, ";")

	login := getEnv("DB_LOGIN", "")
	pwd := getEnv("DB_PASSWORD", "")

	conn, err := arangodbhttp.NewConnection(arangodbhttp.ConnectionConfig{
		Endpoints: endpoints,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})

	if err != nil {
		return err
	}

	client, err := arangodb.NewClient(arangodb.ClientConfig{
		Connection:     conn,
		Authentication: arangodb.BasicAuthentication(login, pwd),
	})

	if s.databases == nil {
		s.databases = make(map[string]*Db)
	}

	if err != nil {
		return err
	}

	s.client = client

	return nil
}

// GetDatabase returns the database with the specified name
func (s *Datastore) GetDatabase(name string) (*Db, error) {

	//We first try to retrieve it from the map
	result := s.databases[name]
	var err error

	// If found, we return it
	if result != nil {
		return result, nil
	}

	// We create or retrieve the database
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	found, err := s.client.DatabaseExists(ctx, name)

	var db arangodb.Database

	if err != nil {
		return nil, err
	}

	if !found {

		cctx := arangodb.WithWaitForSync(ctx)
		db, err = s.client.CreateDatabase(cctx, name, nil)

		if err != nil {
			return nil, err
		}

	} else {

		db, err = s.client.Database(ctx, name)

		if err != nil {
			return nil, err
		}
	}

	c, err := caching.NewCache()

	if err != nil {
		return nil, err
	}

	result = &Db{
		Name:        name,
		db:          db,
		collections: make(map[string]arangodb.Collection),
		cache:       c,
	}

	// We populate the collections
	result.ensureCollections()

	s.databases[name] = result

	return result, nil
}

func (s *Db) ensureCollections() error {

	var err error
	err = s.ensureCollection("Session", s.SessionIsEdgeCollection())
	if err != nil {
		return err
	}
	err = s.ensureCollection("User", s.UserIsEdgeCollection())
	if err != nil {
		return err
	}
	err = s.ensureCollection("SAMLConfig", s.SAMLConfigIsEdgeCollection())
	if err != nil {
		return err
	}
	err = s.ensureCollection("Token", s.TokenIsEdgeCollection())
	if err != nil {
		return err
	}
	err = s.ensureCollection("Password", s.PasswordIsEdgeCollection())
	if err != nil {
		return err
	}

	return nil
}

// Close closes the datastore
func (s *Datastore) Close() error {

	return nil
}

// Seed populate the store with predefined data
func (s *Datastore) Seed() error {

	root := getEnv("DB_SEED_FOLDER", "seed")

	if _, err := os.Stat(root); !os.IsNotExist(err) {

		folders, err := ioutil.ReadDir(root)

		if err != nil {
			return err
		}

		// Each folder represents a database
		for _, folder := range folders {
			if folder.IsDir() {

				name := folder.Name()

				db, err := s.GetDatabase(name)

				if err != nil {
					return err
				}

				filepath.Walk(name, func(path string, info os.FileInfo, err error) error {

					if info.IsDir() == false && filepath.Ext(path) == ".json" {

						switch strings.ToLower(info.Name()) {
						case "sessions.json":
							return db.SeedSessions(path)
						case "users.json":
							return db.SeedUsers(path)
						case "samlconfigs.json":
							return db.SeedSAMLConfigs(path)
						case "tokens.json":
							return db.SeedTokens(path)
						case "passwords.json":
							return db.SeedPasswords(path)
						default:
							return nil
						}
					}

					return nil
				})
			}
		}
	}

	return nil
}

func (s *Db) ensureCollection(collectionName string, isRelation bool) error {

	col, err := s.GetDocumentCollection(collectionName, isRelation)

	if err != nil {
		return err
	}

	s.collections[collectionName] = col

	return nil
}

// GetDocumentCollection creates or/and returns the specified document collection
func (s *Db) GetDocumentCollection(collection string, isRelation bool) (arangodb.Collection, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cctx := arangodb.WithWaitForSync(ctx)

	found, err := s.db.CollectionExists(cctx, collection)

	if err != nil {

		return nil, err
	}

	if !found {

		var options *arangodb.CreateCollectionOptions

		if isRelation {
			options = &arangodb.CreateCollectionOptions{Type: arangodb.CollectionTypeEdge}
		} else {
			options = &arangodb.CreateCollectionOptions{Type: arangodb.CollectionTypeDocument}
		}

		col, err := s.db.CreateCollection(cctx, collection, options)

		if err != nil {
			return nil, err
		}

		return col, nil
	}

	col, err := s.db.Collection(cctx, collection)

	if err != nil {
		return nil, err
	}

	return col, nil
}

// DocumentExistsWithCollectionName defines if the document with the specified id exists in the passed collection
func (s *Db) DocumentExistsWithCollectionName(collectionName, id string) (bool, error) {

	return s.DocumentExists(s.collections[collectionName], id)
}

// DocumentExists defines if the document with the specified id exists in the passed collection
func (s *Db) DocumentExists(collection arangodb.Collection, documentID string) (bool, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cctx := arangodb.WithWaitForSync(ctx)

	result, err := collection.DocumentExists(cctx, documentID)

	if err != nil {
		return false, err
	}

	return result, nil
}

// CreateDocument creates the specified document within the passed collection
func (s *Db) CreateDocument(collection arangodb.Collection, document interface{}) (string, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cctx := arangodb.WithWaitForSync(ctx)

	meta, err := collection.CreateDocument(cctx, document)

	if err != nil {
		return "", err
	}

	return meta.Key, nil
}

// UpdateDocument updates the specified document within the passed collection
func (s *Db) UpdateDocument(collection arangodb.Collection, documentID string, document interface{}) error {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cctx := arangodb.WithWaitForSync(ctx)

	_, err := collection.UpdateDocument(cctx, documentID, document)

	if err != nil {
		return err
	}

	return nil
}

// DeleteDocument deletes the document with the specified id from the passed collection
func (s *Db) DeleteDocument(collection arangodb.Collection, documentID string) error {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cctx := arangodb.WithWaitForSync(ctx)

	_, err := collection.RemoveDocument(cctx, documentID)

	if err != nil {
		return err
	}

	return nil
}

// GetDocument returns the document with the specified id from the passed collection
func (s *Db) GetDocument(collection arangodb.Collection, documentID string, document interface{}) error {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cctx := arangodb.WithWaitForSync(ctx)

	_, err := collection.ReadDocument(cctx, documentID, document)

	if err != nil {
		return err
	}

	return nil
}

// CountDocuments returns the number of documents returned by the specified query
func (s *Db) CountDocuments(query string) (int64, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cctx := arangodb.WithQueryCount(ctx)

	cursor, err := s.db.Query(cctx, query, nil)

	if err != nil {
		return 0, err
	}

	defer cursor.Close()

	return cursor.Count(), nil
}

// CountDocumentsWithParams returns the number of documents returned by the specified query and parameters
func (s *Db) CountDocumentsWithParams(query string, vars map[string]interface{}) (int64, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cctx := arangodb.WithQueryCount(ctx)

	cursor, err := s.db.Query(cctx, query, vars)

	if err != nil {
		return 0, err
	}

	defer cursor.Close()

	return cursor.Count(), nil
}

// DeleteFromCache deletes the specified key from cache
func (s *Db) DeleteFromCache(key string) {
	s.cache.Delete(key)
}

// AddToCache adds the specified key/value to cache
func (s *Db) AddToCache(key string, value interface{}) bool {
	return s.cache.Set(key, value)
}

// GetSessionCacheKey returns the cache key for the passed Session
func (s *Db) GetSessionCacheKey(item *userquerysvc.Session) string {

	kb := caching.NewKeyBuilder()
	kb.Add("Session")
	kb.Add(item.UserID)

	return kb.Get()
}

// SessionIsEdgeCollection defines if the collection is an edge collection
func (s *Db) SessionIsEdgeCollection() bool {
	return false
}

// CreateSession creates a session in the store
func (s *Db) CreateSession(item *userquerysvc.Session) error {

	item.Key = &item.ID

	_, err := s.CreateDocument(s.collections["Session"], item)

	if err != nil {
		return err
	}

	return nil
}

// UpdateSession updates the passed session in the store
func (s *Db) UpdateSession(item *userquerysvc.Session) error {

	err := s.UpdateDocument(s.collections["Session"], item.ID, item)

	if err != nil {
		return err
	}

	return nil
}

// SessionExists defines if a Session with the specified id already exists in the store
func (s *Db) SessionExists(id string) (bool, error) {

	return s.DocumentExists(s.collections["Session"], id)
}

// GetSession returns the session with the specified id
func (s *Db) GetSession(id string) (*userquerysvc.Session, error) {

	var result userquerysvc.Session

	err := s.GetDocument(s.collections["Session"], id, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Db) GetAllSessions() ([]*userquerysvc.Session, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.Session

	query := fmt.Sprintf("FOR d IN %s RETURN d", "Session")

	cursor, err := s.db.Query(ctx, query, nil)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var doc userquerysvc.Session
		_, err := cursor.ReadDocument(ctx, &doc)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &doc)
	}

	return result, nil
}

// DeleteSession deletes the passed session from the datastore
func (s *Db) DeleteSession(item *userquerysvc.Session) error {

	err := s.DeleteDocument(s.collections["Session"], item.ID)

	if err != nil {
		return err
	}

	return nil
}

// ExecuteSessionQuery returns the Session documents for the specified query
func (s *Db) ExecuteSessionQuery(query string, vars map[string]interface{}) ([]*userquerysvc.Session, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.Session

	cursor, err := s.db.Query(ctx, query, vars)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var item userquerysvc.Session
		_, err := cursor.ReadDocument(ctx, &item)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	return result, nil
}

// SeedSessions populates the datastore with sessions from a json file
func (s *Db) SeedSessions(path string) error {

	jsonFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	var items []userquerysvc.Session
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(byteValue, &items)

	if err != nil {
		return err
	}

	for _, item := range items {
		err := s.CreateSession(&item)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetUserCacheKey returns the cache key for the passed User
func (s *Db) GetUserCacheKey(item *userquerysvc.User) string {

	kb := caching.NewKeyBuilder()
	kb.Add("User")
	kb.Add(item.Login)

	return kb.Get()
}

// UserIsEdgeCollection defines if the collection is an edge collection
func (s *Db) UserIsEdgeCollection() bool {
	return false
}

// CreateUser creates a user in the store
func (s *Db) CreateUser(item *userquerysvc.User) error {

	item.Key = &item.ID

	_, err := s.CreateDocument(s.collections["User"], item)

	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates the passed user in the store
func (s *Db) UpdateUser(item *userquerysvc.User) error {

	err := s.UpdateDocument(s.collections["User"], item.ID, item)

	if err != nil {
		return err
	}

	return nil
}

// UserExists defines if a User with the specified id already exists in the store
func (s *Db) UserExists(id string) (bool, error) {

	return s.DocumentExists(s.collections["User"], id)
}

// GetUser returns the user with the specified id
func (s *Db) GetUser(id string) (*userquerysvc.User, error) {

	var result userquerysvc.User

	err := s.GetDocument(s.collections["User"], id, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAllUsers returns all the users from the datastore
func (s *Db) GetAllUsers(activeOnly bool) ([]*userquerysvc.User, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.User

	var query string

	if activeOnly == true {
		query = fmt.Sprintf("FOR d IN %s FILTER d.active == true RETURN d", "User")
	} else {
		query = fmt.Sprintf("FOR d IN %s RETURN d", "User")
	}

	cursor, err := s.db.Query(ctx, query, nil)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var doc userquerysvc.User
		_, err := cursor.ReadDocument(ctx, &doc)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &doc)
	}

	return result, nil
}

// DeleteUser deletes the passed user from the datastore
func (s *Db) DeleteUser(item *userquerysvc.User) error {

	err := s.DeleteDocument(s.collections["User"], item.ID)

	if err != nil {
		return err
	}

	return nil
}

// ExecuteUserQuery returns the User documents for the specified query
func (s *Db) ExecuteUserQuery(query string, vars map[string]interface{}) ([]*userquerysvc.User, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.User

	cursor, err := s.db.Query(ctx, query, vars)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var item userquerysvc.User
		_, err := cursor.ReadDocument(ctx, &item)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	return result, nil
}

// SeedUsers populates the datastore with users from a json file
func (s *Db) SeedUsers(path string) error {

	jsonFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	var items []userquerysvc.User
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(byteValue, &items)

	if err != nil {
		return err
	}

	for _, item := range items {
		err := s.CreateUser(&item)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetSAMLConfigCacheKey returns the cache key for the passed SAMLConfig
func (s *Db) GetSAMLConfigCacheKey(item *userquerysvc.SAMLConfig) string {

	kb := caching.NewKeyBuilder()
	kb.Add("SAMLConfig")
	kb.Add(item.Host)

	return kb.Get()
}

// SAMLConfigIsEdgeCollection defines if the collection is an edge collection
func (s *Db) SAMLConfigIsEdgeCollection() bool {
	return false
}

// CreateSAMLConfig creates a samlconfig in the store
func (s *Db) CreateSAMLConfig(item *userquerysvc.SAMLConfig) error {

	item.Key = &item.ID

	_, err := s.CreateDocument(s.collections["SAMLConfig"], item)

	if err != nil {
		return err
	}

	return nil
}

// UpdateSAMLConfig updates the passed samlconfig in the store
func (s *Db) UpdateSAMLConfig(item *userquerysvc.SAMLConfig) error {

	err := s.UpdateDocument(s.collections["SAMLConfig"], item.ID, item)

	if err != nil {
		return err
	}

	return nil
}

// SAMLConfigExists defines if a SAMLConfig with the specified id already exists in the store
func (s *Db) SAMLConfigExists(id string) (bool, error) {

	return s.DocumentExists(s.collections["SAMLConfig"], id)
}

// GetSAMLConfig returns the samlconfig with the specified id
func (s *Db) GetSAMLConfig(id string) (*userquerysvc.SAMLConfig, error) {

	var result userquerysvc.SAMLConfig

	err := s.GetDocument(s.collections["SAMLConfig"], id, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAllSAMLConfigs returns all the samlconfigs from the datastore
func (s *Db) GetAllSAMLConfigs(activeOnly bool) ([]*userquerysvc.SAMLConfig, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.SAMLConfig

	var query string

	if activeOnly == true {
		query = fmt.Sprintf("FOR d IN %s FILTER d.active == true RETURN d", "SAMLConfig")
	} else {
		query = fmt.Sprintf("FOR d IN %s RETURN d", "SAMLConfig")
	}

	cursor, err := s.db.Query(ctx, query, nil)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var doc userquerysvc.SAMLConfig
		_, err := cursor.ReadDocument(ctx, &doc)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &doc)
	}

	return result, nil
}

// DeleteSAMLConfig deletes the passed samlconfig from the datastore
func (s *Db) DeleteSAMLConfig(item *userquerysvc.SAMLConfig) error {

	err := s.DeleteDocument(s.collections["SAMLConfig"], item.ID)

	if err != nil {
		return err
	}

	return nil
}

// ExecuteSAMLConfigQuery returns the SAMLConfig documents for the specified query
func (s *Db) ExecuteSAMLConfigQuery(query string, vars map[string]interface{}) ([]*userquerysvc.SAMLConfig, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.SAMLConfig

	cursor, err := s.db.Query(ctx, query, vars)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var item userquerysvc.SAMLConfig
		_, err := cursor.ReadDocument(ctx, &item)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	return result, nil
}

// SeedSAMLConfigs populates the datastore with samlconfigs from a json file
func (s *Db) SeedSAMLConfigs(path string) error {

	jsonFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	var items []userquerysvc.SAMLConfig
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(byteValue, &items)

	if err != nil {
		return err
	}

	for _, item := range items {
		err := s.CreateSAMLConfig(&item)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetTokenCacheKey returns the cache key for the passed Token
func (s *Db) GetTokenCacheKey(item *userquerysvc.Token) string {

	kb := caching.NewKeyBuilder()
	kb.Add("Token")
	kb.Add(item.UserID)

	return kb.Get()
}

// TokenIsEdgeCollection defines if the collection is an edge collection
func (s *Db) TokenIsEdgeCollection() bool {
	return false
}

// CreateToken creates a token in the store
func (s *Db) CreateToken(item *userquerysvc.Token) error {

	item.Key = &item.ID

	_, err := s.CreateDocument(s.collections["Token"], item)

	if err != nil {
		return err
	}

	return nil
}

// UpdateToken updates the passed token in the store
func (s *Db) UpdateToken(item *userquerysvc.Token) error {

	err := s.UpdateDocument(s.collections["Token"], item.ID, item)

	if err != nil {
		return err
	}

	return nil
}

// TokenExists defines if a Token with the specified id already exists in the store
func (s *Db) TokenExists(id string) (bool, error) {

	return s.DocumentExists(s.collections["Token"], id)
}

// GetToken returns the token with the specified id
func (s *Db) GetToken(id string) (*userquerysvc.Token, error) {

	var result userquerysvc.Token

	err := s.GetDocument(s.collections["Token"], id, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Db) GetAllTokens() ([]*userquerysvc.Token, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.Token

	query := fmt.Sprintf("FOR d IN %s RETURN d", "Token")

	cursor, err := s.db.Query(ctx, query, nil)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var doc userquerysvc.Token
		_, err := cursor.ReadDocument(ctx, &doc)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &doc)
	}

	return result, nil
}

// DeleteToken deletes the passed token from the datastore
func (s *Db) DeleteToken(item *userquerysvc.Token) error {

	err := s.DeleteDocument(s.collections["Token"], item.ID)

	if err != nil {
		return err
	}

	return nil
}

// ExecuteTokenQuery returns the Token documents for the specified query
func (s *Db) ExecuteTokenQuery(query string, vars map[string]interface{}) ([]*userquerysvc.Token, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.Token

	cursor, err := s.db.Query(ctx, query, vars)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var item userquerysvc.Token
		_, err := cursor.ReadDocument(ctx, &item)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	return result, nil
}

// SeedTokens populates the datastore with tokens from a json file
func (s *Db) SeedTokens(path string) error {

	jsonFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	var items []userquerysvc.Token
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(byteValue, &items)

	if err != nil {
		return err
	}

	for _, item := range items {
		err := s.CreateToken(&item)

		if err != nil {
			return err
		}
	}

	return nil
}

// PasswordIsEdgeCollection defines if the collection is an edge collection
func (s *Db) PasswordIsEdgeCollection() bool {
	return false
}

// CreatePassword creates a password in the store
func (s *Db) CreatePassword(item *userquerysvc.Password) error {

	_, err := s.CreateDocument(s.collections["Password"], item)

	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword updates the passed password in the store
func (s *Db) UpdatePassword(item *userquerysvc.Password) error {

	err := s.UpdateDocument(s.collections["Password"], item.ID, item)

	if err != nil {
		return err
	}

	return nil
}

// PasswordExists defines if a Password with the specified id already exists in the store
func (s *Db) PasswordExists(id string) (bool, error) {

	return s.DocumentExists(s.collections["Password"], id)
}

// GetPassword returns the password with the specified id
func (s *Db) GetPassword(id string) (*userquerysvc.Password, error) {

	var result userquerysvc.Password

	err := s.GetDocument(s.collections["Password"], id, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Db) GetAllPasswords() ([]*userquerysvc.Password, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.Password

	query := fmt.Sprintf("FOR d IN %s RETURN d", "Password")

	cursor, err := s.db.Query(ctx, query, nil)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var doc userquerysvc.Password
		_, err := cursor.ReadDocument(ctx, &doc)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &doc)
	}

	return result, nil
}

// DeletePassword deletes the passed password from the datastore
func (s *Db) DeletePassword(item *userquerysvc.Password) error {

	err := s.DeleteDocument(s.collections["Password"], item.ID)

	if err != nil {
		return err
	}

	return nil
}

// ExecutePasswordQuery returns the Password documents for the specified query
func (s *Db) ExecutePasswordQuery(query string, vars map[string]interface{}) ([]*userquerysvc.Password, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*userquerysvc.Password

	cursor, err := s.db.Query(ctx, query, vars)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var item userquerysvc.Password
		_, err := cursor.ReadDocument(ctx, &item)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	return result, nil
}

// SeedPasswords populates the datastore with passwords from a json file
func (s *Db) SeedPasswords(path string) error {

	jsonFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	var items []userquerysvc.Password
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(byteValue, &items)

	if err != nil {
		return err
	}

	for _, item := range items {
		err := s.CreatePassword(&item)

		if err != nil {
			return err
		}
	}

	return nil
}
