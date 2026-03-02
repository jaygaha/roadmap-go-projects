package handlers

import (
	"net/http"
	"strconv"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/services"
)

type ProductHandler struct {
	svc *services.ProductService
}

func NewProductHandler(svc *services.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// POST /admin/products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.ProductCreateRequest
	if err := readJSON(r, &req); err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}
	product, err := h.svc.Create(req)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, product)
}

// GET /products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	minPrice, _ := strconv.ParseInt(q.Get("min_price"), 10, 64)
	maxPrice, _ := strconv.ParseInt(q.Get("max_price"), 10, 64)

	products, err := h.svc.List(models.ProductQuery{
		Name:     q.Get("q"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Page:     page,
		Limit:    limit,
	})
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, products)
}

// GET /products/{id}
func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := urlParamInt64(r, "id")
	if err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}
	product, err := h.svc.GetById(id)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, product)
}

// PUT /admin/products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := urlParamInt64(r, "id")
	if err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}
	var req models.ProductUpdateRequest
	if err := readJSON(r, &req); err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}
	product, err := h.svc.Update(id, req)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, product)
}

// DELETE /admin/products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := urlParamInt64(r, "id")
	if err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}
	if err := h.svc.Delete(id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
