package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

const (
	minLength = 43
	maxLength = 128
)

func main() {
	codeVerifier := createCodeVerifier(43)
	codeChallange := createCodeChallange(codeVerifier)
	fmt.Printf("code_verifier  : %s\n", codeVerifier)
	fmt.Printf("code_challange : %s\n", codeChallange)
}

func createCodeVerifier(digits int) string {
	// check the range of digits
	// ref: https://datatracker.ietf.org/doc/html/rfc7636#section-4.1
	switch {
	case digits < minLength:
		log.Printf("getCodeVerifier: insufficient digits given. change the digits to %d", minLength)
		digits = minLength
	case digits > maxLength:
		log.Printf("getCodeVerifier: excessive digits given. change the digits to %d", maxLength)
		digits = maxLength
	}

	// ref: https://datatracker.ietf.org/doc/html/rfc7636#section-4.1
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYabcdefghijklmnopqrstuvwxyz0123456789-._~"

	b := make([]byte, digits)
	rand.Read(b)

	var result string
	for _, v := range b {
		result += string(chars[int(v)%len(chars)])
	}

	return result
}

func createCodeChallange(codeVerifier string) string {
	hashedByte := getHashedByteSHA256(codeVerifier)
	return getEncodedStringBASE64woPadding(hashedByte)
}

func getHashedByteSHA256(s string) []byte {
	r := sha256.Sum256([]byte(s))
	return r[:]
}

func getEncodedStringBASE64woPadding(b []byte) string {
	s := base64.StdEncoding.EncodeToString(b)

	// ref: https://datatracker.ietf.org/doc/html/rfc7636#appendix-A
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)
	s = strings.Replace(s, "=", "", -1)

	return s
}
