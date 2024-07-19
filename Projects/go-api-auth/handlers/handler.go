package handlers

import (
	"fullstackgo/services/auth"
	"fullstackgo/store"
)

type Handler struct {
	store *store.Storage
	auth  *auth.AuthService
}

func New(store *store.Storage, auth *auth.AuthService) *Handler {
	return &Handler{
		store: store,
		auth:  auth,
	}
}
