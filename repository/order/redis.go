package order

import (
	"github.com/sherlingre/orders-api/model"
	"github.com/redis/go-redis/v9"
	"encoding/json"
	"context"
	"errors"
	"fmt"
)

// RedisRepo is a repository implementation for orders using Redis as the storage backend.
type RedisRepo struct {
	Client *redis.Client // Client is the Redis client used by the repository.
}

// orderIDKey generates a key for storing an order in Redis based on its ID.
func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

// Insert adds a new order to the Redis repository.
func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	// Serialize the order data to JSON
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	// Generate the Redis key for the order
	key := orderIDKey(order.OrderID)

	// Start a transaction pipeline
	txn := r.Client.TxPipeline()

	// Set the order data in Redis with a non-existence condition (SetNX)
	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	// Add the order key to the "orders" set in Redis
	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to orders set: %w", err)
	}

	// Execute the transaction
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

// ErrNotExist is an error indicating that the order does not exist.
var ErrNotExist = errors.New("order does not exist")

// FindByID retrieves an order from the Redis repository based on its ID.
func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	// Generate the Redis key for the order
	key := orderIDKey(id)

	// Retrieve the order data from Redis
	value, err := r.Client.Get(ctx, key).Result()

	// Check if the order does not exist in Redis
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		// Return an error if there was an issue retrieving the order
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	// Deserialize the order data from JSON
	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		// Return an error if there was an issue decoding the order JSON
		return model.Order{}, fmt.Errorf("failed to decode order json: %w", err)
	}

	// Return the retrieved order
	return order, nil
}

// DeleteById deletes an order from the Redis repository based on its ID.
func (r *RedisRepo) DeleteById(ctx context.Context, id uint64) error {
	// Generate the Redis key for the order
	key := orderIDKey(id)

	// Start a transaction pipeline
	txn := r.Client.TxPipeline()

	// Delete the order from Redis
	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		// If the order does not exist, discard the transaction and return ErrNotExist
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		// Return an error if there was an issue deleting the order
		txn.Discard()
		return fmt.Errorf("get order: %w", err)
	}

	// Remove the order key from the "orders" set in Redis
	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		// Return an error if there was an issue removing the order key from the set
		txn.Discard()
		return fmt.Errorf("failed to remove from orders set: %w", err)
	}

	// Execute the transaction
	if _, err := txn.Exec(ctx); err != nil {
		// Return an error if there was an issue executing the transaction
		return fmt.Errorf("failed to exec: %w", err)
	}

	// Return nil if the deletion is successful
	return nil
}

// Update modifies an existing order in the Redis repository.
func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	// Serialize the updated order data to JSON
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	// Generate the Redis key for the order
	key := orderIDKey(order.OrderID)

	// Use SetXX to update the order in Redis only if it already exists
	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		// If the order does not exist, return ErrNotExist
		return ErrNotExist
	} else if err != nil {
		// Return an error if there was an issue setting the updated order
		return fmt.Errorf("set order: %w", err)
	}

	// Return nil if the update is successful
	return nil
}

type FindAllPage struct {
	Size uint
	Offset uint
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

// FindAll retrieves a paginated list of orders from the Redis repository.
func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	// Use SScan to get a paginated list of order keys from the "orders" set in Redis
	res := r.Client.SScan(ctx, "orders", uint64(page.Offset), "*", int64(page.Size))

	// Retrieve the order keys and cursor from the SScan result
	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get order ids: %w", err)
	}

	// If no order keys were found, return an empty result
	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
			Cursor: cursor,
		}, nil
	}

	// Use MGet to retrieve the order data for the found keys
	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}

	// Process the retrieved order data
	orders := make([]model.Order, len(xs))

	for i, x := range xs {
		x := x.(string)
		var order model.Order

		// Deserialize the order data from JSON
		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json: %w", err)
		}

		orders[i] = order
	}

	// Return the paginated list of orders along with the cursor
	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}