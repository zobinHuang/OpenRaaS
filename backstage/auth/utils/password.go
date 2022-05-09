package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

/*
	func: HashPassword
	description: hash password with 32-bytes salt
*/
func HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	
	// passwd + salt = hash passwd
	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

/*
	func: ComparePassword
	description: compare hash password with same 32-bytes salt
*/
func ComparePassword(storedPassword, suppliedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")

	salt, err := hex.DecodeString(pwsalt[1])
	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	// passwd + salt = hash passwd
	shash, _ := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)

	return hex.EncodeToString(shash) == pwsalt[0], nil
}
