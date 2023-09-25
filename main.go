package main

import (
	"fmt"
	"context"

	"github.com/sherlingre/orders-api/application"
	// "net/http"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app", err)
	}
}