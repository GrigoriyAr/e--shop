package rest

import (
	"jagodkiL0/store"
)

type API struct {
	db store.Storage
}

func NewAPI(s *store.Storage) *API {
	return &API{db: *s}
}
