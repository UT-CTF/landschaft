package util

import (
	"crypto/rand"
	"math/big"
	"strings"
	"unicode"
)

// Strict mode means that the password must contain at least one lowercase letter,
// one uppercase letter, one number, and one symbol.
func GenerateRandomPassword(length uint, allowedCharacters string, strictMode bool) string {
	var randomPassword []byte
	charsetLen := big.NewInt(int64(len(allowedCharacters)))

	for randomPassword == nil || (strictMode && !isValidPassword(string(randomPassword))) {
		randomPassword = make([]byte, length)
		for i := range randomPassword {
			n, err := rand.Int(rand.Reader, charsetLen)
			if err != nil {
				panic("failed to generate secure random number: " + err.Error())
			}
			randomPassword[i] = allowedCharacters[n.Int64()]
		}
	}
	return string(randomPassword)
}

func GenerateRandomDefaultPassword(length uint) string {
	allowedCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_!"
	return GenerateRandomPassword(length, allowedCharacters, true)
}

func isValidPassword(password string) bool {
	hasLower, hasUpper, hasNumber, hasSymbol := false, false, false, false

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune("-_.?!+=:^()", char):
			hasSymbol = true
		}
	}

	return hasLower && hasUpper && hasNumber && hasSymbol
}
