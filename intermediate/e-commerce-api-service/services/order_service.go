package services

import (
	"fmt"
	"log"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/repository"
)

// OrderService handles order creation and payment processing.
type OrderService struct {
	orderRepo   *repository.OrderRepo
	cartRepo    *repository.CartRepo
	productRepo *repository.ProductRepo
	paymentSvc  *PaymentService
}

// NewOrderService creates a new OrderService with the given repositories and payment service.
func NewOrderService(
	or *repository.OrderRepo,
	cr *repository.CartRepo,
	pr *repository.ProductRepo,
	ps *PaymentService,
) *OrderService {
	return &OrderService{
		orderRepo: or, cartRepo: cr,
		productRepo: pr, paymentSvc: ps,
	}
}

// Checkout is the core business operation. It runs inside a single
// transaction to guarantee atomicity: either everything succeeds
// (stock decremented, order created, cart cleared) or nothing changes.
func (s *OrderService) Checkout(userId int64) (*models.CheckoutResponse, error) {
	tx, err := s.orderRepo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("[ORDER] rollback failed: %v", rbErr)
			}
		}
	}()

	// 1. Read cart within the transaction
	items, err := s.cartRepo.GetCartTx(tx, userId)
	if err != nil {
		return nil, fmt.Errorf("reading cart: %w", err)
	}
	if len(items) == 0 {
		err = models.ErrEmptyCart
		return nil, err
	}

	// 2. Calculate total & validate stock
	var total int64
	for _, item := range items {
		total += item.Product.Price * int64(item.Quantity)
	}

	// 3. Create Stripe PaymentIntent
	pi, piErr := s.paymentSvc.CreatePaymentIntent(total, "usd")
	if piErr != nil {
		err = fmt.Errorf("creating payment intent: %w", piErr)
		return nil, err
	}

	// 4. Create order record
	orderID, err := s.orderRepo.CreateOrder(tx, userId, total, pi.ID)
	if err != nil {
		return nil, err
	}

	// 5. Decrement stock & create order items
	for _, item := range items {
		if stockErr := s.productRepo.DecrementStock(tx, item.ProductId, item.Quantity); stockErr != nil {
			err = stockErr
			return nil, err
		}
		if itemErr := s.orderRepo.CreateOrderItem(tx, orderID, item.ProductId,
			item.Quantity, item.Product.Price); itemErr != nil {
			err = itemErr
			return nil, err
		}
	}

	// 6. Clear the cart
	if err = s.cartRepo.ClearCart(tx, userId); err != nil {
		return nil, fmt.Errorf("clearing cart: %w", err)
	}

	// 7. Commit
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	log.Printf("[ORDER] Created order %d for user %d, total: %d cents", orderID, userId, total)

	return &models.CheckoutResponse{
		OrderID:         orderID,
		ClientSecret:    pi.ClientSecret,
		StripePaymentId: pi.ID,
	}, nil
}

// GetOrder returns a single order by ID, ensuring the user has permission to view it.
func (s *OrderService) GetOrder(userId, orderID int64) (*models.Order, error) {
	order, err := s.orderRepo.FindById(orderID)
	if err != nil {
		return nil, err
	}
	// Ensure users can only view their own orders
	if order.UserId != userId {
		return nil, models.ErrForbidden
	}
	return order, nil
}

// ListOrders returns all orders for a user.
func (s *OrderService) ListOrders(userId int64) ([]models.Order, error) {
	return s.orderRepo.ListByUser(userId)
}

// ConfirmPayment is called via Stripe webhook or manual confirmation.
func (s *OrderService) ConfirmPayment(orderId int64) error {
	return s.orderRepo.UpdateStatus(orderId, "paid")
}
