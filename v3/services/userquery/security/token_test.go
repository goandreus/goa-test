package security

import (
	"testing"
	"time"
)

// TestArrayContains is testing the contains function
func TestCreateAndParseToken(t *testing.T) {

	mySigningKey := []byte("mysecret")

	roles := []string{"admin"}

	user := &User{
		ID:        "u1234556",
		Login:     "mpolo",
		Firstname: "Marco",
		Lastname:  "Polo",
		Email:     "marco@polo.com",
		Language:  "fr",
		Roles:     roles,
	}

	// We create a token

	duration := time.Minute * 5

	token, tokenID, err := CreateToken(mySigningKey, "ws-test", "api:internal", user, duration)

	if err != nil {
		t.Error("No error was expected.")
	}

	t.Log("Generated Token: ", token)

	jwttoken, err := ParseToken(token, mySigningKey)

	if err != nil {
		t.Error("No error was expected.")
	}

	if jwttoken.Valid() != nil {
		t.Error("Valid token expected.")
	}

	if jwttoken.Subject != "mpolo" {
		t.Error("Wrong subject value.")
	}

	if jwttoken.Id != tokenID {
		t.Error("Wrong claim id.")
	}
}
