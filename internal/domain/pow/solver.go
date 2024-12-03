package pow

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
)

func SolveChallenge(ctx context.Context, challenge *Challenge) (*Solution, error) {
	var solution string

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
			solution = fmt.Sprintf("%x", rand.Int63())
			data := fmt.Sprintf("%s%s", challenge.Nonce, solution)
			hash := sha256.Sum256([]byte(data))

			if countLeadingZeros(hash[:]) >= challenge.Difficulty {
				return &Solution{
					Nonce:    challenge.Nonce,
					Response: solution,
				}, nil
			}
		}
	}
}

func countLeadingZeros(hash []byte) int {
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
