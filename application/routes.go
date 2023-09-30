package application

import (
	"github.com/sherlingre/orders-api/repository/order"
	"github.com/sherlingre/orders-api/handler"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// loadRoutes configures and loads routes for the application.
func (a *App) loadRoutes() {
	// Create a new Chi router
	router := chi.NewRouter()

	// Use the Chi middleware Logger for logging HTTP requests
	router.Use(middleware.Logger)

	// Define a simple endpoint for testing
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Route to the /orders path and delegate to loadOrderRouter
	router.Route("/orders", a.loadOrderRouter)

	// Assign the configured router to the App
	a.router = router
}

// loadOrderRouter configures and loads routes related to orders.
func (a *App) loadOrderRouter(router chi.Router) {
	// Create an order handler with a Redis repository
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{
			Client: a.rdb,
		},
	}

	// Define routes for order-related operations
	router.Post("/", orderHandler.Create)         // Create a new order
	router.Get("/", orderHandler.List)            // List all orders
	router.Get("/{id}", orderHandler.GetById)     // Get order by ID
	router.Put("/{id}", orderHandler.UpdateById)  // Update order by ID
	router.Delete("/{id}", orderHandler.DeleteById) // Delete order by ID
}