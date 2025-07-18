package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(inputStr string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(inputStr), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
