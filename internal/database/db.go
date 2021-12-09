package database

import (
	"context"
	"errors"
	"sync"

	"superhuman-social/internal/models"
)

var (
	NotFoundErr = errors.New("not found")
)

type DB interface {
	GetPerson(ctx context.Context, email string) (models.Person, error)
	SavePerson(ctx context.Context, p *models.Person) error
	ListAllPersons(ctx context.Context) ([]models.Person, error)
}

type dbClient struct {
	m  map[string]models.Person
	mu sync.RWMutex
}

func NewDB() (DB, error) {
	return &dbClient{m: make(map[string]models.Person)}, nil
}

func (dc *dbClient) GetPerson(ctx context.Context, email string) (models.Person, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	if p, ok := dc.m[email]; ok {
		return p, nil
	}
	return models.Person{}, NotFoundErr
}

func (dc *dbClient) SavePerson(ctx context.Context, p *models.Person) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.m[p.Email] = *p
	return nil
}

func (dc *dbClient) ListAllPersons(ctx context.Context) ([]models.Person, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	ps := make([]models.Person, len(dc.m))
	i := 0
	for _, v := range dc.m {
		ps[i] = v
		i++
	}
	return ps, nil
}
