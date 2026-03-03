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
// @Summary      Create product
// @Description  Create a new product (admin only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      models.ProductCreateRequest  true  "Product data"
// @Success      201      {object}  models.Product
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Router       /admin/products [post]
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
// @Summary      List products
// @Description  List products with optional search and pagination
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        q          query     string  false  "search"
// @Param        min_price  query     int64   false  "min price (cents)"
// @Param        max_price  query     int64   false  "max price (cents)"
// @Param        page       query     int     false  "page"
// @Param        limit      query     int     false  "limit"
// @Success      200        {array}   models.Product
// @Router       /products [get]
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
// @Summary      Get product
// @Description  Get a product by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int64  true  "Product ID"
// @Success      200  {object}  models.Product
// @Failure      404  {object}  map[string]string
// @Router       /products/{id} [get]
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
// @Summary      Update product
// @Description  Update a product by ID (admin only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int64                        true  "Product ID"
// @Param        payload  body      models.ProductUpdateRequest  true  "Product fields"
// @Success      200      {object}  models.Product
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /admin/products/{id} [put]
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
// @Summary      Delete product
// @Description  Delete a product by ID (admin only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int64  true  "Product ID"
// @Success      204  "No Content"
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /admin/products/{id} [delete]
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
