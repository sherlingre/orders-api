package application

import (
	"github.com/redis/go-redis/v9"
	"net/http"
	"context"
	"time"
	"fmt"
)

// App represents the main application structure.
type App struct {
	router http.Handler     // HTTP router for handling requests
	rdb    *redis.Client    // Redis client for interacting with Redis
	config Config           // Configuration for the application
}

// New creates a new instance of the App.
func New(config Config) *App {
	// Initialize the App struct with the provided configuration
	app := &App{
		rdb:    redis.NewClient(&redis.Options{Addr: config.RedisAddress}),
		config: config,
	}

	// Configure and load routes for the application
	app.loadRoutes()

	// Return the initialized App
	return app
}

// Start starts the application, including the HTTP server and connecting to Redis.
func (a *App) Start(ctx context.Context) error {
	// Create an HTTP server with the configured router
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	// Check if the application can connect to Redis
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Defer closing the Redis client to ensure it's closed when the application exits
	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close Redis", err)
		}
	}()

	fmt.Println("Starting server")

	// Create a channel to receive potential errors from the server startup
	ch := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	// Use a select statement to handle either an error from server startup or a cancellation signal
	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		// If the context is canceled, attempt to gracefully shut down the server
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}