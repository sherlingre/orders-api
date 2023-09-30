package main

import (
	"github.com/sherlingre/orders-api/application"
	"os/signal"
	"context"
	"fmt"
	"os"
)

func main() {
	// Load the configuration for the application
	app := application.New(application.LoadConfig())

	// Setup a context with cancellation on interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Start the application, passing the context for possible cancellation
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app", err)
	}
}