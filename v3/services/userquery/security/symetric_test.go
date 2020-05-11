package security

import (
	"testing"
)

// TestArrayContains is testing the contains function
func TestEncryptAndDecrypt(t *testing.T) {

	text := "helloworld"

	r, err := Encrypt(text)

	if err != nil {
		t.Error("Error while encrypting.")
	}

	if r == nil {
		t.Error("Encrypted text is nil.")
	}

	o, err := Decrypt(*r)

	if err != nil {
		t.Error("Error while decrypting.")
	}

	value := *o

	if value != text {
		t.Error("Decrypted text is different from original text.")
	}

}
