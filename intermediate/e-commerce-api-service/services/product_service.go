package services

import (
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/repository"
)

// ProductService handles product-related business logic
type ProductService struct {
	repo *repository.ProductRepo
}

// NewProductService creates a new instance of ProductService
func NewProductService(repo *repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

// Create creates a new product
func (s *ProductService) Create(req models.ProductCreateRequest) (*models.Product, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("%w: name is required", models.ErrBadRequest)
	}
	if req.Price <= 0 {
		return nil, fmt.Errorf("%w: price must be positive", models.ErrBadRequest)
	}
	return s.repo.Create(&req)
}

// GetById retrieves a product by its ID
func (s *ProductService) GetById(id int64) (*models.Product, error) {
	return s.repo.FindById(id)
}

// List retrieves a list of products based on the query parameters
func (s *ProductService) List(q models.ProductQuery) ([]models.Product, error) {
	return s.repo.List(q)
}

// Update updates a product by its ID
func (s *ProductService) Update(id int64, req models.ProductUpdateRequest) (*models.Product, error) {
	return s.repo.Update(id, &req)
}

// Delete deletes a product by its ID
func (s *ProductService) Delete(id int64) error {
	return s.repo.Delete(id)
}
