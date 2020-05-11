package userqueryapi

import (
	"context"
	"fmt"
	"net/http"
	"os"

	stand "github.com/nats-io/nats-streaming-server/server"
	"gitlab.com/wiserskills/v3/services/userquery/datalayer"
	userquerysvcsvr "gitlab.com/wiserskills/v3/services/userquery/gen/http/userquery/server"
	log "gitlab.com/wiserskills/v3/services/userquery/gen/log"
	"gitlab.com/wiserskills/v3/services/userquery/gen/tracing"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
	security "gitlab.com/wiserskills/v3/services/userquery/security"
	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
)

var server *userquerysvcsvr.Server
var handler http.Handler
var nats *stand.StanServer

func GetServerHandle() http.Handler {

	if nats == nil {

		os.Setenv("CLUSTER", "fr0")
		os.Setenv("EVENTBUS_CLUSTER", "wiserskills")
		os.Setenv("EVENTBUS_URL", "nats://localhost:4222")

		// create a NATS streaming server instance
		var err error
		//nats, err = stand.RunServer("wiserskills")

		if err != nil {
			fmt.Println(fmt.Errorf("Error with NATS Streaming server: %s", err.Error()))
		}
	}

	if server == nil {
		// create an instance of the server
		logger := log.New("userquery", false)
		svc := NewUserquery(logger)
		db := svc.(Dataservice).Store()
		InitializeStore(db)
		endpoints := userquery.NewEndpoints(svc)
		dec := goahttp.RequestDecoder
		enc := goahttp.ResponseEncoder
		m := goahttp.NewMuxer()
		handler = m
		handler = tracing.OpenTracing()(handler)
		handler = httpmdlwr.RequestID()(handler)
		handler = httpmdlwr.PopulateRequestContext()(handler)

		eh := func(ctx context.Context, w http.ResponseWriter, err error) {
			w.Write([]byte("ERROR: " + err.Error()))
		}
		server = userquerysvcsvr.New(endpoints, m, dec, enc, eh, nil)
		userquerysvcsvr.Mount(m, server)
	}

	// returns the handler
	return handler
}

// InitializeStore initializes the datastore for the tests
func InitializeStore(store *datalayer.Datastore) {

	orgID := "wiserskills"
	db, err := store.GetDatabase(orgID)

	if err != nil {
		fmt.Println("Error while retrieving database: " + err.Error())
		return
	}

	userID := "3f0b3c52-c953-41af-ba65-79db83de5193"
	userID2 := "4f0b3c52-c953-41af-ba65-79db83de5193"
	userID3 := "5f0b3c52-c953-41af-ba65-79db83de5193"

	pwd, _ := security.HashPassword([]byte("mysecret"))

	exists, _ := db.DocumentExistsWithCollectionName("User", userID)
	if !exists {
		err = db.CreateUser(&userquery.User{
			ID:                userID,
			FirstName:         "Marco",
			LastName:          "Polo",
			BirthName:         "Polo",
			Active:            true,
			B2C:               true,
			Email:             "marco.polo@china.com",
			Login:             "mpolo",
			EncryptedPassword: string(pwd),
			OrganizationID:    orgID,
			CreatedAt:         &[]string{"2019-02-27T19:50:23.980Z"}[0],
			UpdatedAt:         &[]string{"2019-02-27T19:50:23.980Z"}[0],
		})

		if err != nil {
			fmt.Println("Error while saving user: " + err.Error())
		}
	}

	exists, _ = db.DocumentExistsWithCollectionName("User", userID2)
	if !exists {
		err = db.CreateUser(&userquery.User{
			ID:                userID2,
			FirstName:         "Jean",
			LastName:          "Valjean",
			BirthName:         "Valjean",
			Active:            true,
			Email:             "jean.valjean@paris.com",
			Login:             "jvaljean",
			OrganizationID:    orgID,
			EncryptedPassword: string(pwd),
			PasswordExpiresAt: &[]string{"2019-02-27T19:50:23.980Z"}[0],
			CreatedAt:         &[]string{"2019-02-27T19:50:23.980Z"}[0],
			UpdatedAt:         &[]string{"2019-02-27T19:50:23.980Z"}[0],
		})

		if err != nil {
			fmt.Println("Error while saving user: " + err.Error())
		}
	}

	exists, _ = db.DocumentExistsWithCollectionName("User", userID3)
	if !exists {
		err = db.CreateUser(&userquery.User{
			ID:                userID3,
			FirstName:         "Jules",
			LastName:          "Verne",
			BirthName:         "Verne",
			Active:            true,
			Email:             "jules.verne@paris.com",
			Login:             "jverne",
			EncryptedPassword: string(pwd),
			OrganizationID:    orgID,
			PasswordExpiresAt: &[]string{"2019-02-27T19:50:23.980Z"}[0],
			CreatedAt:         &[]string{"2019-02-27T19:50:23.980Z"}[0],
			UpdatedAt:         &[]string{"2019-02-27T19:50:23.980Z"}[0],
		})

		if err != nil {
			fmt.Println("Error while saving user: " + err.Error())
		}
	}

	exists, _ = db.DocumentExistsWithCollectionName("SAMLConfig", "3f0b3c52-c953-41af-ba65-79db83de5199")
	if !exists {
		err = db.CreateSAMLConfig(&userquery.SAMLConfig{
			ID:             "3f0b3c52-c953-41af-ba65-79db83de5199",
			OrganizationID: orgID,
			AreaID:         "DEV",
			Host:           "wiserskills.io",
			Active:         true,
			IdpMetadata:    `<?xml version="1.0"?><md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://capriza.github.io/samling/samling.html"><md:IDPSSODescriptor WantAuthnRequestsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol"><md:KeyDescriptor use="signing"><ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:X509Data><ds:X509Certificate>MIICpzCCAhACCQDuFX0Db5iljDANBgkqhkiG9w0BAQsFADCBlzELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExEjAQBgNVBAcMCVBhbG8gQWx0bzEQMA4GA1UECgwHU2FtbGluZzEPMA0GA1UECwwGU2FsaW5nMRQwEgYDVQQDDAtjYXByaXphLmNvbTEmMCQGCSqGSIb3DQEJARYXZW5naW5lZXJpbmdAY2Fwcml6YS5jb20wHhcNMTgwNTE1MTgxMTEwWhcNMjgwNTEyMTgxMTEwWjCBlzELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExEjAQBgNVBAcMCVBhbG8gQWx0bzEQMA4GA1UECgwHU2FtbGluZzEPMA0GA1UECwwGU2FsaW5nMRQwEgYDVQQDDAtjYXByaXphLmNvbTEmMCQGCSqGSIb3DQEJARYXZW5naW5lZXJpbmdAY2Fwcml6YS5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAJEBNDJKH5nXr0hZKcSNIY1l4HeYLPBEKJLXyAnoFTdgGrvi40YyIx9lHh0LbDVWCgxJp21BmKll0CkgmeKidvGlr3FUwtETro44L+SgmjiJNbftvFxhNkgA26O2GDQuBoQwgSiagVadWXwJKkodH8tx4ojBPYK1pBO8fHf3wOnxAgMBAAEwDQYJKoZIhvcNAQELBQADgYEACIylhvh6T758hcZjAQJiV7rMRg+Omb68iJI4L9f0cyBcJENR+1LQNgUGyFDMm9Wm9o81CuIKBnfpEE2Jfcs76YVWRJy5xJ11GFKJJ5T0NEB7txbUQPoJOeNoE736lF5vYw6YKp8fJqPW0L2PLWe9qTn8hxpdnjo3k6r5gXyl8tk=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></md:KeyDescriptor><md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress</md:NameIDFormat><md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://capriza.github.io/samling/samling.html"/><md:SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://capriza.github.io/samling/samling.html"/></md:IDPSSODescriptor></md:EntityDescriptor>`,
			IDKey:          "login",
			CallbackURL:    "http://localhost:8080/saml/signin",
			RedirectURL:    "http://localhost:8080/token?token=%s",
		})

		if err != nil {
			fmt.Println("Error while saving SAML config: " + err.Error())
		}
	}
}
