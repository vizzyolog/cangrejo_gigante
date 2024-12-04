package pow

import (
	"crypto/sha256"
	"fmt"

	"cangrejo_gigante/internal/logger"
	"cangrejo_gigante/internal/utils"
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
	r, err := utils.GenerateCryptoRandomInt()
	if err != nil {
		return &Challenge{}, fmt.Errorf("%w: failed to generate Challenge, crypto random ", err)
	}

	nonce := fmt.Sprintf("%x", r)

	return &Challenge{
		Nonce:      nonce,
		Difficulty: s.difficulty,
	}, nil
}

func (s *Service) VerifySolution(nonce, solution string) bool {
	data := fmt.Sprintf("%s%s", nonce, solution)
	hash := sha256.Sum256([]byte(data))
	leadingZeros := utils.CountLeadingZeros(hash[:])

	s.logger.Infof("Verifying solution: Nonce='%s', Solution='%s', Hash='%x', LeadingZeros=%d, Difficulty=%d\n",
		nonce, solution, hash, leadingZeros, s.difficulty)

	return leadingZeros >= s.difficulty
}
