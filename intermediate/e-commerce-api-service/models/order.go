package models

import "time"

// Order represents an order placed by a user
type Order struct {
	ID              int64       `json:"id"`
	UserId          int64       `json:"user_id"`
	Total           int64       `json:"total"` // total price in cents
	StripePaymentId string      `json:"stripe_payment_id,omitempty"`
	Items           []OrderItem `json:"items,omitempty"`
	Status          string      `json:"status"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        int64    `json:"id"`
	OrderId   int64    `json:"order_id"`
	ProductId int64    `json:"product_id"`
	Quantity  int      `json:"quantity"`
	Price     int64    `json:"price"`
	Product   *Product `json:"product,omitempty"`
}

// CheckoutResponse represents the response payload for checkout
type CheckoutResponse struct {
	OrderID         int64  `json:"order_id"`
	ClientSecret    string `json:"client_secret"`
	StripePaymentId string `json:"stripe_payment_id"`
}
