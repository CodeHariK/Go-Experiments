package handler

import (
	slick "temple"

	"templeapp/model"
	"templeapp/views/profile"
)

func HandleUserProfile(c *slick.Context) error {
	user := model.User{
		FirstName: "Go",
		LastName:  "Experiments",
		Email:     "go@exp",
	}
	return c.Render(profile.Index(user))
}
