package model

import "crypto/rand"

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenerateToken(length int) string {
	ll := len(chars)
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%ll]
	}
	return string(b)
}
