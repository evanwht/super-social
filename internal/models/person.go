package models

import (
	"cloud.google.com/go/datastore"
)

type Person struct {
	Key           *datastore.Key `json:"-" datastore:"__key__"`
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	TimesLookedUp int64          `json:"times_looked_up"`
}
