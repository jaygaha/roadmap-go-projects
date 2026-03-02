package models

import "time"

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	ProductId int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Product   *Product  `json:"product,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddToCartRequest represents the request payload for adding a product to the cart
type AddToCartRequest struct {
	ProductId int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// UpdateCartItemRequest represents the request payload for updating a cart item
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
}

// CartResponse represents the response payload for cart retrieval
type CartResponse struct {
	Items []CartItem `json:"items"`
	Total int64      `json:"total"` // total price in cents
}
