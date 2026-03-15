package repository

import (
	"context"
	"database/sql"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

const (
	StatusActive    = "active"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)

const (
	RoleUser  model.UserRole = "user"
	RoleAdmin model.UserRole = "admin"
)

type userRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, email, name, role, password_hash, created_at, updated_at
			  FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrRecordNotFound
		}
		return nil, pkg.ErrDatabase
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, name, role, password_hash, created_at, updated_at
			  FROM users WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)

	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrRecordNotFound
		}
		return nil, pkg.ErrDatabase
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password_hash, name, role)
			  VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Role,
	).Scan(&user.ID)

	if err != nil {
		return pkg.ErrDuplicateEmail
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET email = $1, name = $2, role = $3, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Name,
		user.Role,
		user.ID,
	)

	if err != nil {
		return pkg.ErrDatabase
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkg.ErrDatabase
	}

	if rowsAffected == 0 {
		return pkg.ErrRecordNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return pkg.ErrDatabase
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkg.ErrDatabase
	}

	if rowsAffected == 0 {
		return pkg.ErrRecordNotFound
	}

	return nil
}

func (r *userRepository) List(ctx context.Context, page, limit int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}

	// Get users
	query := `SELECT id, email, name, role, password_hash, created_at, updated_at
			  FROM users ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Name,
			&u.Role,
			&u.PasswordHash,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, pkg.ErrDatabase
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *userRepository) ListAdmins(ctx context.Context) ([]model.User, error) {
	query := `SELECT id, email, name, role, password_hash, created_at, updated_at
			  FROM users WHERE role = $1 ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, RoleAdmin)
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	defer rows.Close()

	var admins []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Name,
			&u.Role,
			&u.PasswordHash,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, pkg.ErrDatabase
		}
		admins = append(admins, u)
	}

	return admins, nil
}
