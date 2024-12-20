package utils

import (
	"golang.org/x/crypto/sha3"
)

func Sha3Sum(input []byte) []byte {
	sha := sha3.New512()
	_, _ = sha.Write(input)
	return sha.Sum(nil)
}

func Sha3SumString(input string) []byte {
	sha := sha3.New512()
	_, _ = sha.Write([]byte(input))
	return sha.Sum(nil)
}
