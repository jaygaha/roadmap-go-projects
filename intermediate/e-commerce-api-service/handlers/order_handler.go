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
// @Summary      Checkout
// @Description  Create an order from the user's cart and a Stripe PaymentIntent
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      201  {object}  models.CheckoutResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /checkout [post]
func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.Checkout(getUserId(r))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// GET /orders
// @Summary      List orders
// @Description  List all orders for the authenticated user
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Order
// @Failure      401  {object}  map[string]string
// @Router       /orders [get]
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.ListOrders(getUserId(r))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

// GET /orders/{id}
// @Summary      Get order
// @Description  Retrieve a single order by ID (owned by the user)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int64  true  "Order ID"
// @Success      200  {object}  models.Order
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /orders/{id} [get]
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
