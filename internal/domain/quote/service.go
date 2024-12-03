package quote

import (
	"bufio"
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
		return nil, err
	}
	defer file.Close()

	var quotes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		quotes = append(quotes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Service{
		quotes: quotes,
	}, nil
}

func (s *Service) GetRandomQuote() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(s.quotes))
	return s.quotes[index]
}
