package utils

import (
	"crypto/rand"
	"math/big"
)

const lettersAndDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomAlias(length int) (string, error) {
	result := make([]byte, length)

	for i := range length {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersAndDigits))))
		if err != nil {
			return "", nil
		}
		result[i] = lettersAndDigits[idx.Int64()]
	}

	return string(result), nil
}
