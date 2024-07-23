package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"sandslash/store/query/user"
	"sandslash/types"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Username    string           `json:"username"`
	// Email       string           `json:"email"`
	// IsAdmin     pgtype.Bool      `json:"is_admin"`
	// DateOfBirth pgtype.Date      `json:"date_of_birth"`
	// PhoneNumber string           `json:"phone_number"`
	// LastLogin   pgtype.Timestamp `json:"last_login"`
	// Location    pgtype.Int4      `json:"location"`

	_, err := h.store.UserStore.CreateUser(context.Background(), user.CreateUserParams{
		Username: "Hello",
		Email:    "hello@hello.com",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	users, err := h.store.UserStore.ListAllUsers(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, _ := json.Marshal(users)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(content))
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(types.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Info:\nID: %s\nUsername: %s\nEmail: %s", user.ID, user.Username, user.Email)
}
