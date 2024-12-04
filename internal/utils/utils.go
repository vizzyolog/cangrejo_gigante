package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateCryptoRandomInt() (int64, error) {
	maxInt := big.NewInt(1<<63 - 1) // Max value for int63

	randomBigInt, err := rand.Int(rand.Reader, maxInt)
	if err != nil {
		return 0, fmt.Errorf("failed to generate random int: %w", err)
	}

	return randomBigInt.Int64(), nil
}

func CountLeadingZeros(hash []byte) int {
	zeros := 0

	for _, b := range hash {
		for i := 7; i >= 0; i-- {
			if (b>>i)&1 == 0 {
				zeros++
			} else {
				return zeros
			}
		}
	}

	return zeros
}
