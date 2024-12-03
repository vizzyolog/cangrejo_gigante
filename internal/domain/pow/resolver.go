package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type PoWResolver struct {
	Difficulty int
}

func NewPoWResolver(difficulty int) *PoWResolver {
	return &PoWResolver{Difficulty: difficulty}
}

func (pr *PoWResolver) Solve(nonce string) (string, error) {
	for i := 0; ; i++ {
		data := fmt.Sprintf("%s%d", nonce, i)
		hash := sha256.Sum256([]byte(data))
		if pr.isValidHash(hash[:]) {
			return fmt.Sprintf("%d", i), nil
		}
	}
}

func (pr *PoWResolver) isValidHash(hash []byte) bool {
	hashHex := hex.EncodeToString(hash)
	return strings.HasPrefix(hashHex, strings.Repeat("0", pr.Difficulty))
}
