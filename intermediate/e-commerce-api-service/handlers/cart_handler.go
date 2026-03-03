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
// @Summary      Get cart
// @Description  Retrieve the authenticated user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.CartResponse
// @Failure      401  {object}  map[string]string
// @Router       /cart [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	cart, err := h.svc.GetCart(getUserId(r))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cart)
}

// POST /cart/items
// @Summary      Add item to cart
// @Description  Add a product to the authenticated user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      models.AddToCartRequest  true  "Cart item"
// @Success      201      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Router       /cart/items [post]
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
// @Summary      Update cart item
// @Description  Update quantity or remove an item from the cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        productId  path      int64                        true  "Product ID"
// @Param        payload    body      models.UpdateCartItemRequest true  "Quantity (0 deletes)"
// @Success      200        {object}  map[string]string
// @Failure      400        {object}  map[string]string
// @Failure      401        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Router       /cart/items/{productId} [put]
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
// @Summary      Remove cart item
// @Description  Remove a product from the authenticated user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        productId  path  int64  true  "Product ID"
// @Success      204        "No Content"
// @Failure      401        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Router       /cart/items/{productId} [delete]
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
