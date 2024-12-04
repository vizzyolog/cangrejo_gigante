package pow

import (
	"context"
	"crypto/sha256"
	"fmt"

	"cangrejo_gigante/internal/utils"
)

func SolveChallenge(ctx context.Context, challenge *Challenge) (*Solution, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
			randomInt, err := utils.GenerateCryptoRandomInt()
			if err != nil {
				return nil, fmt.Errorf("failed to generate random number: %w", err)
			}

			solution := fmt.Sprintf("%x", randomInt)
			data := fmt.Sprintf("%s%s", challenge.Nonce, solution)
			hash := sha256.Sum256([]byte(data))

			if utils.CountLeadingZeros(hash[:]) >= challenge.Difficulty {
				return &Solution{
					Nonce:    challenge.Nonce,
					Response: solution,
				}, nil
			}
		}
	}
}
