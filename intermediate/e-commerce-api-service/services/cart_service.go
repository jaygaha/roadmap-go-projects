package services

import (
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/repository"
)

// CartService handles cart operations
type CartService struct {
	cartRepo    *repository.CartRepo
	productRepo *repository.ProductRepo
}

func NewCartService(cr *repository.CartRepo, pr *repository.ProductRepo) *CartService {
	return &CartService{cartRepo: cr, productRepo: pr}
}

func (s *CartService) AddItem(userID int64, req models.AddToCartRequest) error {
	if req.Quantity <= 0 {
		return fmt.Errorf("%w: quantity must be positive", models.ErrBadRequest)
	}

	// Verify product exists and has stock
	product, err := s.productRepo.FindById(req.ProductId)
	if err != nil {
		return err
	}
	if product.Stock < req.Quantity {
		return models.ErrInsufficientStock
	}

	return s.cartRepo.Upsert(userID, req.ProductId, req.Quantity)
}

func (s *CartService) UpdateItem(userID, productId int64, req models.UpdateCartItemRequest) error {
	if req.Quantity <= 0 {
		return s.cartRepo.Remove(userID, productId)
	}
	return s.cartRepo.UpdateQuantity(userID, productId, req.Quantity)
}

func (s *CartService) RemoveItem(userId, productId int64) error {
	return s.cartRepo.Remove(userId, productId)
}

// GetCart retrieves the user's cart with product details
func (s *CartService) GetCart(userId int64) (*models.CartResponse, error) {
	items, err := s.cartRepo.GetCart(userId)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, item := range items {
		if item.Product != nil {
			total += item.Product.Price * int64(item.Quantity)
		}
	}

	return &models.CartResponse{Items: items, Total: total}, nil
}
