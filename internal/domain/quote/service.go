package quote

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Service struct {
	quotes []string
}

func New(filePath string) (*Service, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	var quotes []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		quotes = append(quotes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return &Service{
		quotes: quotes,
	}, nil
}

// #nosec G404
func (s *Service) GetRandomQuote() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(s.quotes))

	return s.quotes[index]
}
