package service

import (
	"context"

	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/auth"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/repository"
	pkg "github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

// UserService handles user business logic
type UserService struct {
	userRepo repository.UserRepository
	auth     *auth.AuthService
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, auth *auth.AuthService) *UserService {
	return &UserService{
		userRepo: userRepo,
		auth:     auth,
	}
}

// Signup creates a new user
func (s *UserService) Signup(ctx context.Context, req *model.UserSignupRequest) (*model.AuthResponse, error) {
	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && err != pkg.ErrRecordNotFound {
		return nil, pkg.ErrDatabase
	}
	if existing != nil {
		return nil, pkg.ErrUserExists
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, pkg.ErrValidation
	}

	// Create user
	user := &model.User{
		Email:        req.Email,
		Name:         req.Name,
		Role:         model.RoleUser,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.auth.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	return &model.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login authenticates a user
func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == pkg.ErrRecordNotFound {
			return nil, pkg.ErrInvalidPassword
		}
		return nil, pkg.ErrDatabase
	}

	// Check password
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		return nil, pkg.ErrInvalidPassword
	}

	// Generate token
	token, err := s.auth.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	return &model.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// List retrieves paginated users
func (s *UserService) List(ctx context.Context, page, limit int) ([]model.User, int64, error) {
	return s.userRepo.List(ctx, page, limit)
}

// Update updates a user
func (s *UserService) Update(ctx context.Context, user *model.User) error {
	// Check if email is already used by another user
	if user.Email != "" {
		existing, err := s.userRepo.GetByEmail(ctx, user.Email)
		if err == nil && existing.ID != user.ID {
			return pkg.ErrUserExists
		}
	}
	return s.userRepo.Update(ctx, user)
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

// ListAdmins retrieves all admin users
func (s *UserService) ListAdmins(ctx context.Context) ([]model.User, error) {
	return s.userRepo.ListAdmins(ctx)
}
