package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

func CheckPasswordHash(hashedPassword, plainText string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainText))
	return err
}
