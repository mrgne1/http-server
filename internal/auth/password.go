package auth

import "golang.org/x/crypto/bcrypt"

func HashPasword(password string) (string, error) {
	hash, err :=  bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
