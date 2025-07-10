package helpers

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password using bcrypt with cost 14
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// VerifyPassword compares a bcrypt hashed password with its possible plain-text version
func VerifyPassword(hashedPwd, plainPwd string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
    return err == nil
}
