package handlers

import (
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/services"
)

type OrderHandler struct {
	svc *services.OrderService
}

func NewOrderHandler(svc *services.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// POST /checkout
func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.Checkout(getUserId(r))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// GET /orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.ListOrders(getUserId(r))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

// GET /orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := urlParamInt64(r, "id")
	if err != nil {
		handleError(w, err)
		return
	}
	order, err := h.svc.GetOrder(getUserId(r), orderID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, order)
}
