package handler

import (
	"context"
	"fmt"
	"net/http"

	"sandslash/store/query/user"
	"sandslash/types"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	h.store.UserStore.CreateUser(context.Background(), user.CreateUserParams{
		Username: "superuser",
	})
	users, _ := h.store.UserStore.ListUsers(context.Background())

	fmt.Fprintln(w, users)
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(types.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Info:\nID: %s\nUsername: %s\nEmail: %s", user.ID, user.Username, user.Email)
}
