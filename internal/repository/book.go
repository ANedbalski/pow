package repository

import "math/rand"

type Book struct {
	ID     int64
	Phrase string
}

type InMemoBookRepo struct {
	books []Book
}

func NewInMemoBookRepo() *InMemoBookRepo {
	return &InMemoBookRepo{
		//quotes from word of wisdom
		books: []Book{
			{ID: 1, Phrase: "When the going gets rough - turn to wonder"},
			{ID: 2, Phrase: "If you have knowledge, let others light their candles in it."},
			{ID: 3, Phrase: "The best way to predict the future is to create it."},
			{ID: 4, Phrase: "The only way to do great work is to love what you do."},
		},
	}
}

func (r *InMemoBookRepo) GetRandomPhrase() (string, error) {
	i := rand.Intn(len(r.books))
	return r.books[i].Phrase, nil
}
