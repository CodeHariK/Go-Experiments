package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"Checkout-System/config"

	"Checkout-System/internal/infra/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config: %v", err)
	}

	ctxDBTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	mongoUrl := "mongodb+srv://" + cfg.Mongo.UserName + ":" + cfg.Mongo.Password +
		"@" + cfg.Mongo.Host + "/?retryWrites=true&w=majority&appName=" + cfg.Mongo.AppName
	db, err := database.Connect(ctxDBTimeout, mongoUrl)
	if err != nil {
		slog.Error("Failed to connect to database %s", err)
	}
	defer db.Disconnect()
	defer cancel()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	fmt.Printf("Server running on port : %x", cfg.Server.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%x", cfg.Server.Host, cfg.Server.Port), r)
}
