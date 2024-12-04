package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type Resolver struct {
	Difficulty int
}

func NewPoWResolver(difficulty int) *Resolver {
	return &Resolver{Difficulty: difficulty}
}

func (pr *Resolver) Solve(nonce string) (string, error) {
	for attempt := 0; ; attempt++ {
		data := fmt.Sprintf("%s%d", nonce, attempt)

		hash := sha256.Sum256([]byte(data))

		if pr.isValidHash(hash[:]) {
			return strconv.Itoa(attempt), nil
		}
	}
}

func (pr *Resolver) isValidHash(hash []byte) bool {
	hashHex := hex.EncodeToString(hash)

	return strings.HasPrefix(hashHex, strings.Repeat("0", pr.Difficulty))
}
