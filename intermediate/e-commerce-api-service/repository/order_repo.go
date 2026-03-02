package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
)

// OrderRepo is the repository for order operations
type OrderRepo struct {
	db *sql.DB
}

// NewOrderRepo creates a new instance of OrderRepo
func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// CreateOrder creates a new order in the database
func (r *OrderRepo) CreateOrder(tx *sql.Tx, userId int64, total int64, stripeID string) (int64, error) {
	res, err := tx.Exec(
		`INSERT INTO orders (user_id, total, status, stripe_payment_id) VALUES (?, ?, 'pending', ?)`,
		userId, total, stripeID,
	)
	if err != nil {
		return 0, fmt.Errorf("creating order: %w", err)
	}

	return res.LastInsertId()
}

// CreateOrderItem creates a new order item in the database
func (r *OrderRepo) CreateOrderItem(tx *sql.Tx, orderID, productID int64, qty int, price int64) error {
	_, err := tx.Exec(
		`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`,
		orderID, productID, qty, price,
	)

	return err
}

// UpdateStatus updates the status of an order
func (r *OrderRepo) UpdateStatus(orderId int64, status string) error {
	_, err := r.db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, orderId)

	return err
}

// FindById retrieves an order by its ID
func (r *OrderRepo) FindById(orderId int64) (*models.Order, error) {
	o := &models.Order{}
	err := r.db.QueryRow(
		`SELECT id, user_id, total, status, stripe_payment_id, created_at
         FROM orders WHERE id = ?`, orderId,
	).Scan(&o.ID, &o.UserId, &o.Total, &o.Status, &o.StripePaymentId, &o.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Load order items
	rows, err := r.db.Query(
		`SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
                p.name, p.description, p.image_url
         FROM order_items oi
         JOIN products p ON p.id = oi.product_id
         WHERE oi.order_id = ?`, orderId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		var p models.Product
		if err := rows.Scan(&item.ID, &item.OrderId, &item.ProductId,
			&item.Quantity, &item.Price, &p.Name, &p.Description, &p.ImageURL); err != nil {
			return nil, err
		}
		p.ID = item.ProductId
		p.Price = item.Price
		item.Product = &p
		o.Items = append(o.Items, item)
	}

	return o, nil
}

// ListByUser retrieves all orders for a given user ID
func (r *OrderRepo) ListByUser(userId int64) ([]models.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, total, status, stripe_payment_id, created_at
         FROM orders WHERE user_id = ? ORDER BY created_at DESC`, userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserId, &o.Total, &o.Status,
			&o.StripePaymentId, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

// BeginTx exposes the DB's transaction capability to the service layer.
func (r *OrderRepo) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}
