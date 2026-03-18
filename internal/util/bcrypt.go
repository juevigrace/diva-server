package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	pass := string(encpw)

	return pass, nil
}

func ValidatePassword(reqPass, encPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encPass), []byte(reqPass)) == nil
}
