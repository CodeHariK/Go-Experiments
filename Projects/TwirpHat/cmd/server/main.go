package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"twirphat/docs"
	"twirphat/rpc/haberdasher"

	"github.com/twitchtv/twirp"
	"github.com/twitchtv/twirp/hooks/statsd"
)

type randomHaberdasher struct{}

// MakeHat is a function to make hats
//
//	@Summary		Post MakeHats
//	@Description	Post Hat Making
//	@Tags			Hat
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	haberdasher.Hat
//	@Failure		503	{object}	nil
//	@Router			/dash/twirphat.haberdasher.Haberdasher/MakeHat [post]
func (h *randomHaberdasher) MakeHat(ctx context.Context, size *haberdasher.Size) (*haberdasher.Hat, error) {
	if size.Inches <= 0 {
		return nil, twirp.InvalidArgumentError("Inches", "I can't make a hat that small!")
	}
	colors := []string{"white", "black", "brown", "red", "blue"}
	names := []string{"bowler", "baseball cap", "top hat", "derby"}
	return &haberdasher.Hat{
		Size:  size.Inches,
		Color: colors[rand.Intn(len(colors))],
		Name:  names[rand.Intn(len(names))],
	}, nil
}

// @title			Fiber Example API
// @version			1.0
// @description	This is a sample swagger for Fiber
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.email	fiber@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:8080
// @BasePath		/
func main() {
	hook := statsd.NewStatsdServerHooks(LoggingStatter{os.Stderr})
	dashServer := haberdasher.NewHaberdasherServer(
		&randomHaberdasher{},
		hook,
		twirp.WithServerPathPrefix("/dash"),
	)

	mux := http.NewServeMux()

	loggedMux := requestLogger(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: loggedMux,
	}

	process, err := os.FindProcess(os.Getpid())
	fmt.Printf("Server Pid %d\n", os.Getpid())
	go func() {
		time.Sleep(50 * time.Second)
		if err != nil {
			log.Fatalf("Failed to find process: %v", err)
		}
		if err := process.Signal(syscall.SIGINT); err != nil {
			log.Fatalf("Failed to send signal: %v", err)
		}
	}()

	mux.Handle(dashServer.PathPrefix(), dashServer)
	docs.SwaggerHandler(mux, "https://cdn.pixabay.com/photo/2017/03/16/21/18/logo-2150297_640.png")

	// Channel to listen for interrupt or termination signal from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal the completion of the shutdown process
	serverShutdown := make(chan struct{})

	go func() {
		fmt.Printf("Recieved Signal : %v", <-stop)

		// Initiate graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Error during shutdown: %v", err)
		}
		close(serverShutdown)
	}()

	log.Println("Server is starting...")

	// Start the server
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	// Wait for the shutdown process to complete
	<-serverShutdown
	log.Println("Server has shut down gracefully")
}

type LoggingStatter struct {
	io.Writer
}

func (ls LoggingStatter) Inc(metric string, val int64, rate float32) error {
	_, err := fmt.Fprintf(ls, "incr %s: %d @ %f\n", metric, val, rate)
	return err
}

func (ls LoggingStatter) TimingDuration(metric string, val time.Duration, rate float32) error {
	_, err := fmt.Fprintf(ls, "time %s: %s @ %f\n", metric, val, rate)
	return err
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request method and path
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
