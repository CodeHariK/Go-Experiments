package main

import (
	"log"
	"net/http"
	"time"

	// "fmt"
	// "io"
	// "log"
	"goexperiments/database"
	"goexperiments/handlers"
	"goexperiments/router"

	"github.com/gofiber/fiber/v2"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

const idleTimeout = 5 * time.Second

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.ListFacts)

	app.Post("/fact", handlers.CreateFact)
}

func main() {
	app := fiber.New()
	// app := fiber.New(fiber.Config{
	// 	Prefork:           true,
	// 	CaseSensitive:     true,
	// 	StrictRouting:     true,
	// 	ServerHeader:      "Fiber",
	// 	EnablePrintRoutes: true,
	// 	AppName:           "App Name",
	// })
	// app.Use(cors.New())

	database.ConnectDb()

	router.SetupRoutes(app)

	// setupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}

// func main() {
// 	app := fiber.New(fiber.Config{
// 		AppName: `goexperiments`,
// 		// Prefork:           true,
// 		DisableStartupMessage: true,
// 		ColorScheme:           fiber.DefaultColors,
// 		// EnablePrintRoutes: true,
// 		IdleTimeout: idleTimeout,
// 	})

// 	if !fiber.IsChild() {
// 		fmt.Println("I'm the parent process")
// 	} else {
// 		fmt.Println("I'm a child process")
// 	}

// 	app.Use(func(c *fiber.Ctx) error {
// 		c.Bind(fiber.Map{
// 			"Title": "Hello, World!",
// 		})
// 		return c.Next()
// 	})

// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.SendString("Hello, Word!")
// 	}).Name(`Home`)

// 	app.Static("/", "./public", fiber.Static{
// 		Compress:      true,
// 		ByteRange:     true,
// 		Browse:        true,
// 		CacheDuration: 10 * time.Second,
// 		MaxAge:        3600,
// 	})

// 	app.Get("/user/:id", func(c *fiber.Ctx) error {
// 		return c.SendString(c.Params("id"))
// 	}).Name("user.show")

// 	app.Get("/test", func(c *fiber.Ctx) error {
// 		location, _ := c.GetRouteURL("user.show", fiber.Map{"id": 1})
// 		return c.SendString(location)
// 	})

// 	app.Use([]string{"/api", "/home"}, func(c *fiber.Ctx) error {
// 		c.Set("X-Custom-Header", `random.String(32)`)
// 		return c.Next()
// 	}, func(c *fiber.Ctx) error {
// 		return c.Next()
// 	})

// 	app.Get("/json", func(c *fiber.Ctx) error {
// 		resp, err := client.Get("https://dummyjson.com/products/1")
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
// 				"success": false,
// 				"error":   err.Error(),
// 			})
// 		}

// 		defer resp.Body.Close()
// 		if resp.StatusCode != http.StatusOK {
// 			return c.Status(resp.StatusCode).JSON(&fiber.Map{
// 				"success": false,
// 				"error":   err.Error(),
// 			})
// 		}

// 		if _, err := io.Copy(c.Response().BodyWriter(), resp.Body); err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
// 				"success": false,
// 				"error":   err.Error(),
// 			})
// 		}
// 		return c.SendStatus(fiber.StatusOK)
// 	})

// 	app.Get("/struct", func(c *fiber.Ctx) error {
// 		// Create data struct:
// 		data := SomeStruct{
// 			Name: "Grame",
// 			Age:  20,
// 		}

// 		return c.JSON(data)
// 		// => Content-Type: application/json
// 		// => "{"Name": "Grame", "Age": 20}"

// 		return c.JSON(fiber.Map{
// 			"name": "Grame",
// 			"age":  20,
// 		})
// 		// => Content-Type: application/json
// 		// => "{"name": "Grame", "age": 20}"

// 		return c.JSON(fiber.Map{
// 			"type":     "https://example.com/probs/out-of-credit",
// 			"title":    "You do not have enough credit.",
// 			"status":   403,
// 			"detail":   "Your current balance is 30, but that costs 50.",
// 			"instance": "/account/12345/msgs/abc",
// 		}, "application/problem+json")
// 		// => Content-Type: application/problem+json
// 		// => "{
// 		// =>     "type": "https://example.com/probs/out-of-credit",
// 		// =>     "title": "You do not have enough credit.",
// 		// =>     "status": 403,
// 		// =>     "detail": "Your current balance is 30, but that costs 50.",
// 		// =>     "instance": "/account/12345/msgs/abc",
// 		// => }"
// 	})

// 	// curl "http://localhost:3000/person/?name=hi&pass=ls&products=shoe,hat"
// 	app.Get("/person", func(c *fiber.Ctx) error {
// 		p := new(Person)

// 		if err := c.QueryParser(p); err != nil {
// 			return err
// 		}

// 		log.Println(p.Name)
// 		log.Println(p.Pass)
// 		log.Println(p.Products)

// 		return c.SendStatus(fiber.StatusOK)
// 	})

// 	app.Get("/coffee", func(c *fiber.Ctx) error {
// 		return c.Redirect("/teapot")
// 	})

// 	app.Get("/teapot", func(c *fiber.Ctx) error {
// 		return c.Status(fiber.StatusTeapot).SendString("🍵 short and stout 🍵")
// 	})

// 	api := app.Group("/api", handler) // /api

// 	v1 := api.Group("/v1", handler) // /api/v1
// 	v1.Get("/list", handler)        // /api/v1/list
// 	v1.Get("/user", handler)        // /api/v1/user

// 	v2 := api.Group("/v2", handler) // /api/v2
// 	v2.Get("/list", handler)        // /api/v2/list
// 	v2.Get("/user", handler)

// 	app.Route("/test", func(api fiber.Router) {
// 		api.Get("/foo", handler).Name("foo") // /test/foo (name: test.foo)
// 		api.Get("/bar", handler).Name("bar") // /test/bar (name: test.bar)
// 	}, "test.")

// 	// 404 Handler
// 	app.Use(func(c *fiber.Ctx) error {
// 		return c.SendStatus(404) // => 404 "Not Found"
// 	})

// 	// data, _ := json.MarshalIndent(app.Stack(), "", "  ")
// 	// fmt.Println(string(data))

// 	fmt.Println("Hello, World @ http://localhost:8080")

// 	log.Fatal(app.Listen(":8080"))
// }

func handler(c *fiber.Ctx) error {
	return c.SendString(c.Route().Path)
}

type SomeStruct struct {
	Name string
	Age  uint8
}

// Field names should start with an uppercase letter
type Person struct {
	Name     string   `query:"name"`
	Pass     string   `query:"pass"`
	Products []string `query:"products"`
}
