package handler

import (
	"fmt"
	"net/http"

	"sandslash/types"
)

// type Item struct {
// 	ID          int    `db:"id"`
// 	Name        string `db:"name"`
// 	Description string `db:"description"`
// }

// items := []Item{
// 	{Name: "Item 1", Description: "Description for item 1"},
// 	{Name: "Item 2", Description: "Description for item 2"},
// 	{Name: "Item 3", Description: "Description for item 3"},
// 	{Name: "Item 4", Description: "Description for item 4"},
// }

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello How are you")
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(types.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Info:\nID: %s\nUsername: %s\nEmail: %s", user.ID, user.Username, user.Email)
}
