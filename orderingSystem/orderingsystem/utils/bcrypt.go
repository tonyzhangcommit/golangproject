package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// 加密
func BcryptMake(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// 验证
func BcryptMakeCheck(pwd []byte, hashePwd string) bool {
	byteHash := []byte(hashePwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(hashePwd))
	if err != nil {
		return false
	}
	return true
}
