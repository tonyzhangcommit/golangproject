package utils

import (
	"math/rand"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const charsetInt = "0123456789"

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func GenerateRandomIntString(length int) string {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子
	result := make([]byte, length)
	for i := range result {
		result[i] = charsetInt[rand.Intn(len(charsetInt))]
	}
	return string(result)
}
