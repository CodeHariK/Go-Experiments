package main

import (
	"fmt"
	slick "temple"

	"templeapp/handler"

	"github.com/google/uuid"
)

func main() {
	app := slick.New()

	app.Plug(WithAuth, WithRequestID)

	app.Get("/profile", handler.HandleUserProfile)
	app.Get("/dashboard", handler.HandleDashboard)

	app.Start(":3000")
}

func WithAuth(h slick.Handler) slick.Handler {
	return func(c *slick.Context) error {
		c.Set("email", "Go@Exp")
		fmt.Println("Hello from Auth")
		return h(c)
	}
}

func WithRequestID(h slick.Handler) slick.Handler {
	return func(c *slick.Context) error {
		c.Set("requestID", uuid.New())
		return h(c)
	}
}
