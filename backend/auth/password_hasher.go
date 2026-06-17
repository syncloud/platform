package auth

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
)

type PasswordHasher struct{}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (h *PasswordHasher) Hash(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	salt := make([]byte, 4)
	_, err := rand.Read(salt)
	if err != nil {
		log.Printf("unable to generate password salt: %s", err)
		salt = []byte("salt")
	}
	hasher.Write(salt)
	hash := hasher.Sum(nil)
	hashWithSalt := append(hash, salt...)
	encodedHash := base64.StdEncoding.EncodeToString(hashWithSalt)
	return fmt.Sprintf("{SSHA}%s", encodedHash)
}
