package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/hello", basicHandler)

	server := &http.Server{
		Addr: ":3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen to server", err)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.MethodGet {
	// 	if r.URL.Path == "/foo" {
	// 		// Handle get foo
	// 	}
	// }

	// if r.Method == http.MethodPost {
	// 	// Handle POST
	// }
	w.Write([]byte("Hello, world!"))
}