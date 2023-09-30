package model

import (
	"github.com/google/uuid"
	"time"
)

// Order represents an order in the system.
type Order struct {
	OrderID     uint64      `json:"order_id"`     // Unique identifier for the order
	CustomerID  uuid.UUID   `json:"customer_id"`  // Unique identifier for the customer
	LineItems   []LineItem  `json:"line_items"`   // List of line items in the order
	CreatedAt   *time.Time  `json:"created_at"`   // Timestamp when the order was created
	ShippedAt   *time.Time  `json:"shipped_at"`   // Timestamp when the order was shipped
	CompletedAt *time.Time  `json:"completed_at"` // Timestamp when the order was completed
}

// LineItem represents an item within an order.
type LineItem struct {
	ItemID   uuid.UUID `json:"item_id"`   // Unique identifier for the item
	Quantity uint      `json:"quantity"`  // Quantity of the item in the order
	Price    uint      `json:"price"`     // Price of the item
}