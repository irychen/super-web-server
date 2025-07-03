package utils

import (
	"math/rand/v2"
	"time"
)

var (
	rs = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 1))
)

func GenerateRandomCode(length int) string {
	chars := "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rs.IntN(len(chars))]
	}
	return string(code)
}

func GenerateRandomCodeOnlyUppercase(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rs.IntN(len(chars))]
	}
	return string(code)
}

func GenerateRandomCodeOnlyLowercase(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rs.IntN(len(chars))]
	}
	return string(code)
}

func GenerateRandomCodeOnlyNumber(length int) string {
	chars := "1234567890"
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rs.IntN(len(chars))]
	}
	return string(code)
}

func GenerateRandomCodeOnlyNumberAndUppercase(length int) string {
	chars := "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rs.IntN(len(chars))]
	}
	return string(code)
}
