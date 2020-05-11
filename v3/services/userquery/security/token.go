package security

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

// WSClaims represents a custom claim used by WS
type WSClaims struct {
	*User
	*jwt.StandardClaims
	Scope string
}

// CreateToken creates a new token with the specified parameters
func CreateToken(signingkey []byte, issuer string, scope string, user *User, validity time.Duration) (string, string, error) {

	// Generate the ID of the token
	id, _ := uuid.NewV4()
	tokenID := id.String()

	// Create the Claims
	stdClaims := &jwt.StandardClaims{
		Id:        tokenID,
		ExpiresAt: time.Now().Add(validity).Unix(),
		Issuer:    issuer,
		Subject:   user.Login,
		IssuedAt:  time.Now().UTC().Unix(),
	}

	claims := &WSClaims{
		User:           user,
		StandardClaims: stdClaims,
		Scope:          scope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	value, err := token.SignedString(signingkey)

	return value, tokenID, err
}

// ParseToken parses the specified string token
func ParseToken(tokenstring string, signingkey []byte) (*WSClaims, error) {

	if tokenstring == "" {
		return nil, errors.New("passed token string is empty")
	}

	if signingkey == nil || len(signingkey) == 0 {
		return nil, errors.New("passed signing key is nil or zero length")
	}

	token, err := jwt.ParseWithClaims(tokenstring, &WSClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingkey, nil
	})

	return token.Claims.(*WSClaims), err
}
