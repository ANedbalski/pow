package app

import (
	"pow/internal/book"
	"time"
)

type Factory struct {
	pow         POW
	bookService *book.Service
	repo        POWRepository
	powTTL      time.Duration
}

func NewFactory(pow POW, book *book.Service, repo POWRepository, powTTL time.Duration) *Factory {
	return &Factory{
		pow:         pow,
		bookService: book,
		repo:        repo,
		powTTL:      powTTL,
	}
}

func (f *Factory) NewDispatcher(client string) *Handler {
	return New(f.pow, f.bookService, client, f.repo, f.powTTL)
}
