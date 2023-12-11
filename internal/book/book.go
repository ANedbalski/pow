package book

import "context"

type Repository interface {
	GetRandomPhrase() (string, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetQuote(_ context.Context) (string, error) {
	return s.repo.GetRandomPhrase()
}
