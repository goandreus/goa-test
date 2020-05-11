package security

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"

	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
	"gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
)

// GetServiceProvider returns the SAML service provider for the specified parameters
func GetServiceProvider(samlConfig *userquery.SAMLConfig) (*saml2.SAMLServiceProvider, error) {

	metadataIDP := []byte(samlConfig.IdpMetadata)

	metadata := &types.EntityDescriptor{}
	err := xml.Unmarshal(metadataIDP, metadata)

	if err != nil {
		return nil, err
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				err = fmt.Errorf("Metadata certificate(%d) must not be empty", idx)
				break
			}

			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				break
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				break
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}

		if err != nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	// We sign the AuthnRequest with a random key because
	randomKeyStore := dsig.RandomKeyStoreForTest()

	sp := &saml2.SAMLServiceProvider{}
	sp.IdentityProviderSSOURL = metadata.IDPSSODescriptor.SingleSignOnServices[0].Location
	sp.IdentityProviderIssuer = metadata.EntityID
	sp.SignAuthnRequests = true
	sp.IDPCertificateStore = &certStore
	sp.SPKeyStore = randomKeyStore
	sp.AssertionConsumerServiceURL = samlConfig.CallbackURL

	return sp, nil
}

// GetURLToAuthenticate this method create custom URL for request authenticate
func GetURLToAuthenticate(samlConfig *userquery.SAMLConfig) (url string, err error) {

	sp, err := GetServiceProvider(samlConfig)

	if err != nil {
		return
	}

	url, err = sp.BuildAuthURL(samlConfig.Host)
	if err != nil {
		return
	}

	return
}

// GetValueFromSAMLResponse returns the value from the passed IDP data
func GetValueFromSAMLResponse(sp *saml2.SAMLServiceProvider, samlResponse string, key string) (result string, err error) {

	assertionInfo, err := sp.RetrieveAssertionInfo(samlResponse)
	if err != nil {
		return
	}

	if assertionInfo.WarningInfo.InvalidTime {

		err = errors.New("invalid time in SAML response")
		return
	}

	if assertionInfo.WarningInfo.NotInAudience {
		err = errors.New("received SAML response not in expected audience")
		return
	}

	result = assertionInfo.Values.Get(key)

	if result == "" {
		err = fmt.Errorf("expected '%s' value is not present in the SAML response", key)
		return
	}

	return
}
