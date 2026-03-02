package repository

import (
	"database/sql"
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
)

// CartRepo is the repository for cart
type CartRepo struct {
	db *sql.DB
}

// NewCartRepo creates a new cart repository
func NewCartRepo(db *sql.DB) *CartRepo {
	return &CartRepo{db: db}
}

// Upsert inserts a new cart item or updates the quantity if it already exists
func (r *CartRepo) Upsert(userId int64, productId int64, qty int) error {
	_, err := r.db.Exec(
		`INSERT INTO cart_items (user_id, product_id, quantity)
         VALUES (?, ?, ?)
         ON CONFLICT(user_id, product_id)
         DO UPDATE SET quantity = quantity + excluded.quantity`,
		userId, productId, qty,
	)
	if err != nil {
		return fmt.Errorf("upserting cart item: %w", err)
	}
	return nil
}

// UpdateQuantity updates the quantity of a cart item
func (r *CartRepo) UpdateQuantity(userId, productId int64, qty int) error {
	res, err := r.db.Exec(
		`UPDATE cart_items
         SET quantity = ?
         WHERE user_id = ? AND product_id = ?`,
		qty, userId, productId,
	)
	if err != nil {
		return fmt.Errorf("updating cart quantity: %w", err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return models.ErrNotFound
	}

	return nil
}

// Remove removes a cart item
func (r *CartRepo) Remove(userId, productId int64) error {
	res, err := r.db.Exec(
		`DELETE FROM cart_items
         WHERE user_id = ? AND product_id = ?`,
		userId, productId,
	)
	if err != nil {
		return fmt.Errorf("removing cart item: %w", err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return models.ErrNotFound
	}

	return nil
}

// GetCart returns all cart items with their associated product data.
func (r *CartRepo) GetCart(userId int64) ([]models.CartItem, error) {
	rows, err := r.db.Query(
		`SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at,
                p.id, p.name, p.description, p.price, p.stock, p.image_url
         FROM cart_items ci
         JOIN products p ON p.id = ci.product_id
         WHERE ci.user_id = ?
         ORDER BY ci.created_at DESC`, userId,
	)
	if err != nil {
		return nil, fmt.Errorf("querying cart: %w", err)
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var ci models.CartItem
		var p models.Product
		if err := rows.Scan(
			&ci.ID, &ci.UserId, &ci.ProductId, &ci.Quantity, &ci.CreatedAt,
			&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("scanning cart item: %w", err)
		}
		ci.Product = &p
		items = append(items, ci)
	}

	return items, rows.Err()
}

// ClearCart removes all items for a user. Used after checkout.
func (r *CartRepo) ClearCart(tx *sql.Tx, userId int64) error {
	_, err := tx.Exec("DELETE FROM cart_items WHERE user_id = ?", userId)
	if err != nil {
		return fmt.Errorf("clearing cart: %w", err)
	}

	return nil
}

// GetCartTx returns all cart items with their associated product data from a transaction.
func (r *CartRepo) GetCartTx(tx *sql.Tx, userId int64) ([]models.CartItem, error) {
	rows, err := tx.Query(
		`SELECT ci.product_id, ci.quantity, p.price, p.stock
         FROM cart_items ci
         JOIN products p ON p.id = ci.product_id
         WHERE ci.user_id = ?`, userId,
	)
	if err != nil {
		return nil, fmt.Errorf("querying cart: %w", err)
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var ci models.CartItem
		var p models.Product
		if err := rows.Scan(&ci.ProductId, &ci.Quantity, &p.Price, &p.Stock); err != nil {
			return nil, fmt.Errorf("scanning cart item: %w", err)
		}
		ci.Product = &p
		items = append(items, ci)
	}

	return items, rows.Err()
}
