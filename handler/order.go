package handler

import (
	"github.com/sherlingre/orders-api/repository/order"
	"github.com/sherlingre/orders-api/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"errors"
	"time"
	"fmt"
)

// Order represents the HTTP handler for order-related operations.
type Order struct {
	Repo *order.RedisRepo // Repository for order data (e.g., Redis implementation)
}

// Create handles the creation of a new order.
func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}

	// Decode the request body into the 'body' struct
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	// Create an order with random OrderID (for simplicity)
	order := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	// Insert the order into the repository
	err := o.Repo.Insert(r.Context(), order)
	if err != nil {
		fmt.Println("failed to insert order:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with the created order
	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshal order:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

// List handles the retrieval of a list of orders.
func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const base = 10
	const bitSize = 64

	// Parse the cursor from the URL query
	cursor, err := strconv.ParseUint(cursorStr, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50

	// Retrieve a list of orders from the repository
	res, err := o.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: uint(cursor),
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare the response with orders and the next cursor
	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor

	// Marshal the response into JSON
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with the JSON data
	w.Write(data)
}

// GetById handles the retrieval of an order by its ID.
func (h *Order) GetById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	// Parse the order ID from the URL parameter
	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve the order by ID from the repository
	o, err := h.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with the order in JSON format
	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateById handles the update of an order by its ID.
func (o *Order) UpdateById(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	// Decode the request body into the 'body' struct
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	// Parse the order ID from the URL parameter
	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve the existing order by ID from the repository
	theOrder, err := o.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"
	now := time.Now().UTC()

	// Update the order status based on the request body
	switch body.Status {
	case shippedStatus:
		if theOrder.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAt = &now
	case completedStatus:
		if theOrder.CompletedAt != nil || theOrder.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the order in the repository
	err = o.Repo.Update(r.Context(), theOrder)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with the updated order in JSON format
	if err := json.NewEncoder(w).Encode(theOrder); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// DeleteById handles the deletion of an order by its ID.
func (o *Order) DeleteById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	// Parse the order ID from the URL parameter
	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete the order by ID from the repository
	err = o.Repo.DeleteById(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}