package pow

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"cangrejo_gigante/internal/logger"
)

type Service struct {
	difficulty int
	logger     logger.Logger
}

func New(difficulty int, logger logger.Logger) *Service {
	return &Service{
		difficulty: difficulty,
		logger:     logger,
	}
}

func (s *Service) GenerateChallenge() (*Challenge, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	nonce := fmt.Sprintf("%x", r.Int63())

	return &Challenge{
		Nonce:      nonce,
		Difficulty: s.difficulty,
	}, nil
}

func (s *Service) VerifySolution(nonce, solution string) bool {
	data := fmt.Sprintf("%s%s", nonce, solution)
	hash := sha256.Sum256([]byte(data))
	leadingZeros := countLeadingZeros(hash[:])

	s.logger.Infof("Verifying solution: Nonce='%s', Solution='%s', Hash='%x', LeadingZeros=%d, Difficulty=%d\n",
		nonce, solution, hash, leadingZeros, s.difficulty)

	return leadingZeros >= s.difficulty
}
