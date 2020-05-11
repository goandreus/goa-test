package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// Encrypt encrypts the specified text with the specified key
func Encrypt(text string) (*string, error) {

	key := []byte("thisisa32byteslongencryptionkey!")
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	result := base64.URLEncoding.EncodeToString(ciphertext)
	return &result, err
}

// Decrypt decrypts the specified text with the specified key
func Decrypt(cryptoText string) (*string, error) {

	key := []byte("thisisa32byteslongencryptionkey!")
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	result := string(ciphertext)

	return &result, nil
}

// HashText creates a hash value of the passed text
func HashText(text string) string {

	key := []byte("thisisa32byteslongencryptionkey!")
	t := []byte(text)

	hash := hmac.New(sha256.New, key)
	hash.Write(t)

	hex.EncodeToString(hash.Sum(nil))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
