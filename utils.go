package coze

import (
	"crypto/rand"
	"encoding/json"
)

func ptrValue[T any](s *T) T {
	if s != nil {
		return *s
	}
	var empty T
	return empty
}

func ptr[T any](s T) *T {
	return &s
}

func ptrNotZero[T comparable](s T) *T {
	var zero T
	if s == zero {
		return nil
	}
	return &s
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return bytesToHex(bytes), nil
}

func bytesToHex(bytes []byte) string {
	hex := make([]byte, len(bytes)*2)
	for i, b := range bytes {
		hex[i*2] = hexChar(b >> 4)
		hex[i*2+1] = hexChar(b & 0xF)
	}
	return string(hex)
}

func hexChar(b byte) byte {
	if b < 10 {
		return '0' + b
	}
	return 'a' + (b - 10)
}

func mustToJson(obj any) string {
	jsonArray, err := json.Marshal(obj)
	if err != nil {
		return "{}"
	}
	return string(jsonArray)
}
