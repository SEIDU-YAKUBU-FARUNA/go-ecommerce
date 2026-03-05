package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Checkpassword(hashedpassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
}
