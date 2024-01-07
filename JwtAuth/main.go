package main

import (
	"log"

	// "fmt"
	// "io"
	// "log"
	"goexperiments/database"
	"goexperiments/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.ConnectDb()

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
