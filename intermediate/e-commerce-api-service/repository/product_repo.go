package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
)

// ProductRepo is the repository for products
type ProductRepo struct {
	db *sql.DB
}

// NewProductRepo creates a new instance of ProductRepo
func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

// Create creates a new product
func (r *ProductRepo) Create(p *models.ProductCreateRequest) (*models.Product, error) {
	res, err := r.db.Exec(
		`INSERT INTO products (name, description, price, stock, image_url) VALUES (?, ?, ?, ?, ?)`,
		p.Name, p.Description, p.Price, p.Stock, p.ImageURL,
	)
	if err != nil {
		return nil, fmt.Errorf("creating product: %w", err)
	}

	id, _ := res.LastInsertId()

	return r.FindById(id)
}

// FindById finds a product by its ID
func (r *ProductRepo) FindById(id int64) (*models.Product, error) {
	p := &models.Product{}
	err := r.db.QueryRow(
		`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = ?`,
		id,
	).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}

	return p, nil
}

// List supports filtering, search, and pagination.
// All query parameters are safely parameterized — no SQL injection risk.
func (r *ProductRepo) List(q models.ProductQuery) ([]models.Product, error) {
	var (
		where []string
		args  []any
	)

	if q.Name != "" {
		where = append(where, "(name LIKE ? OR description LIKE ?)")
		s := "%" + q.Name + "%"
		args = append(args, s, s)
	}
	if q.MinPrice > 0 {
		where = append(where, "price >= ?")
		args = append(args, q.MinPrice)
	}
	if q.MaxPrice > 0 {
		where = append(where, "price <= ?")
		args = append(args, q.MaxPrice)
	}

	query := "SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products"
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY created_at DESC"

	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	offset := (q.Page - 1) * q.Limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", q.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price,
			&p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning product: %w", err)
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// Update updates a product
func (r *ProductRepo) Update(id int64, req *models.ProductUpdateRequest) (*models.Product, error) {
	var sets []string
	var args []any

	if req.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *req.Description)
	}
	if req.Price != nil {
		sets = append(sets, "price = ?")
		args = append(args, *req.Price)
	}
	if req.Stock != nil {
		sets = append(sets, "stock = ?")
		args = append(args, *req.Stock)
	}
	if req.ImageURL != nil {
		sets = append(sets, "image_url = ?")
		args = append(args, *req.ImageURL)
	}

	if len(sets) == 0 {
		return r.FindById(id)
	}

	sets = append(sets, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE products SET %s WHERE id = ?", strings.Join(sets, ", "))
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("updating product: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, models.ErrNotFound
	}

	return r.FindById(id)
}

// Delete deletes a product by its ID
func (r *ProductRepo) Delete(id int64) error {
	res, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return models.ErrNotFound
	}

	return nil
}

// DecrementStock atomically reduces stock. Used during checkout.
func (r *ProductRepo) DecrementStock(tx *sql.Tx, productID int64, qty int) error {
	res, err := tx.Exec(
		"UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?",
		qty, productID, qty,
	)
	if err != nil {
		return fmt.Errorf("decrementing stock: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return models.ErrInsufficientStock
	}

	return nil
}
