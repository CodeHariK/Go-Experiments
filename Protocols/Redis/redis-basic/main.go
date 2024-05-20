package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/redis/go-redis/v9"
)

func main() {
	app := fiber.New()

	// Redis configuration
	redis := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.3:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		// Set and retrieve a value from Redis
		err := redis.Set(c.Context(), "key", "Hello, Fiber & Redis", 0).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		value, err := redis.Get(c.Context(), "key").Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendString(fmt.Sprintf("Value from Redis: %s", value))
	})

	app.Get("/map", func(c *fiber.Ctx) error {
		session := map[string]string{"name": "Red", "surname": "Goddess", "company": "Redis", "age": "27"}
		for k, v := range session {
			err := redis.HSet(c.Context(), "user-session:123", k, v).Err()
			if err != nil {
				panic(err)
			}
		}

		userSession := redis.HGetAll(c.Context(), "user-session:123").Val()
		fmt.Println(userSession)
		return c.SendString(fmt.Sprintf("Value from Redis: %s", userSession))
	})

	app.Get("/quote", func(c *fiber.Ctx) error {
		cache, err := redis.JSONGet(c.Context(), "quote").Result()
		if err != nil {
			resp, _ := http.Get("https://dummyjson.com/quotes")

			resBody, _ := io.ReadAll(resp.Body)

			err := redis.JSONSet(c.Context(), "quote", "$", resBody).Err()

			redis.Expire(c.Context(), "quote", time.Second*30)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}

			return c.Send(resBody)
		}

		return c.SendString(cache)
	})

	// Start the Fiber app
	port := 3000
	log.Printf("Server started on :%d...", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
