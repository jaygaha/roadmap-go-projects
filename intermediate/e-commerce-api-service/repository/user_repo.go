package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
)

// UserRepo represents the repository for user operations
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo creates a new instance of UserRepo
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(email, hashedPassword, role string) (*models.User, error) {
	res, err := r.db.Exec(
		`INSERT INTO users (email, password, role) VALUES (?, ?, ?)`,
		email, hashedPassword, role,
	)
	if err != nil {
		// SQLite unique constraint violation
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, models.ErrConflict
		}
		return nil, fmt.Errorf("inserting user: %w", err)
	}

	id, _ := res.LastInsertId()

	return r.FindById(id)
}

// FindById finds a user by their ID
func (r *UserRepo) FindById(id int64) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, role, created_at, updated_at FROM users WHERE id = ?`,
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by id: %w", err)
	}

	return u, nil
}

// FindByEmail finds a user by their email
func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, role, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by email: %w", err)
	}
	return u, nil
}
