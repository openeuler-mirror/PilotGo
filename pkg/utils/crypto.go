package utils

import "golang.org/x/crypto/bcrypt"

func CryptoPassword(data string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
}

func ComparePassword(hash, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}
