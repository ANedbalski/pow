package repository

import "time"

type InMemoPOWRepo struct {
	data map[string]time.Time
}

func NewInMemoPOWRepo() *InMemoPOWRepo {
	return &InMemoPOWRepo{
		data: make(map[string]time.Time),
	}
}

func (p *InMemoPOWRepo) Add(key string, ttl time.Duration) error {
	p.data[key] = time.Now().Add(ttl)
	return nil
}

func (p *InMemoPOWRepo) Exists(key string) (bool, error) {
	if _, ok := p.data[key]; !ok || p.data[key].Before(time.Now()) {
		return false, nil
	}
	return true, nil
}
