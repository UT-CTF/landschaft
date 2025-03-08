package util

import (
	"math/rand"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

var titleColor = color.New(color.FgYellow).Add(color.Underline)

func PrintSectionTitle(title string) {
	titleColor.Println(title)
}

const allowedCharacters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.?!+=:^()"

func GenerateRandomPassword(length uint) string {
	randomPassword := make([]byte, length)
	for i := range randomPassword {
		randomPassword[i] = allowedCharacters[rand.Intn(len(allowedCharacters))]
	}
	return string(randomPassword)
}

func GenerateStrictRandomPassword(length uint) string {
	password := GenerateRandomPassword(length)
	for !isValidPassword(password) {
		password = GenerateRandomPassword(length)
	}
	return password
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
