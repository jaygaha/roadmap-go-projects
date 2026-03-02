package handlers

import (
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/services"
)

type CartHandler struct {
	svc *services.CartService
}

func NewCartHandler(svc *services.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

// GET /cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	cart, err := h.svc.GetCart(getUserId(r))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cart)
}

// POST /cart/items
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var req models.AddToCartRequest
	if err := readJSON(r, &req); err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}
	if err := h.svc.AddItem(getUserId(r), req); err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "item added to cart"})
}

// PUT /cart/items/{productId}
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	productId, err := urlParamInt64(r, "productId")
	if err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}
	var req models.UpdateCartItemRequest
	if err := readJSON(r, &req); err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}
	if err := h.svc.UpdateItem(getUserId(r), productId, req); err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "cart updated"})
}

// DELETE /cart/items/{productId}
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	productId, err := urlParamInt64(r, "productId")
	if err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}
	if err := h.svc.RemoveItem(getUserId(r), productId); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
