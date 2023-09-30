module github.com/sherlingre/orders-api

go 1.21.1

// Dependencies for the project
require (
	// Package for fast hashing using the xxHash algorithm
	github.com/cespare/xxhash/v2 v2.2.0

	// Package for consistent hashing
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f

	// Routing package for building RESTful APIs
	github.com/go-chi/chi/v5 v5.0.10

	// Redis client library for Go
	github.com/redis/go-redis/v9 v9.2.0

	// Package for generating and working with UUIDs
	github.com/google/uuid v1.3.1
)