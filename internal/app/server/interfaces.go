package server

import "cangrejo_gigante/internal/domain/pow"

type PowService interface {
	GenerateChallenge() (*pow.Challenge, error)
	VerifySolution(nonce, solution string) bool
}

type QuoteService interface {
	GetRandomQuote() string
}
