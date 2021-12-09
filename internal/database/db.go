package database

import (
	"errors"
	"fmt"

	"superhuman-social/internal/models"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

const (
	personEntityKind = "person"
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
	datastoreClient *datastore.Client
}

func NewDB(projectID string) (DB, error) {
	c, err := datastore.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to make a datastore client: %v", err)
	}
	return &dbClient{datastoreClient: c}, nil
}

func (dc *dbClient) GetPerson(ctx context.Context, email string) (models.Person, error) {
	q := datastore.NewQuery(personEntityKind).Filter("Email =", email)
	var ps []models.Person
	keys, err := dc.datastoreClient.GetAll(ctx, q, &ps)
	if err != nil {
		return models.Person{}, err
	} else if len(keys) == 0 {
		return models.Person{}, NotFoundErr
	}
	// TODO what if there were more than one?
	return ps[0], nil
}

func (dc *dbClient) SavePerson(ctx context.Context, p *models.Person) error {
	if p.Key != nil {
		_, err := dc.datastoreClient.Put(ctx, p.Key, p)
		if err != nil {
			return err
		}
	} else {
		key, err := dc.datastoreClient.Put(ctx, datastore.IncompleteKey(personEntityKind, nil), p)
		if err != nil {
			return err
		}
		p.Key = key
	}
	return nil
}

func (dc *dbClient) ListAllPersons(ctx context.Context) ([]models.Person, error) {
	q := datastore.NewQuery(personEntityKind).Order("-TimesLookedUp")
	var ps []models.Person
	if _, err := dc.datastoreClient.GetAll(ctx, q, &ps); err != nil {
		return nil, err
	}
	return ps, nil
}
