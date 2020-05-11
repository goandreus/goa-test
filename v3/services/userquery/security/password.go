package security

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns an hashed/salted password
func HashPassword(pwd []byte) ([]byte, error) {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		return nil, err
	}

	return hash, nil
}

// ComparePasswords compares a hashed password with a plain text one
func ComparePasswords(hashedPwd []byte, plainPwd []byte) (bool, error) {

	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	if err != nil {
		return false, err
	}

	return true, nil
}
