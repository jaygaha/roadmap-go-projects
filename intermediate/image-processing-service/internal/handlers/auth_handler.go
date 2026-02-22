package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/errors"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/logger"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// SignUpHandler handles user sign up requests
func (h *Handler) SignUpHandler(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := utils.ValidateName(req.Name); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	exists, err := h.userService.EmailExists(req.Email)
	if err != nil {
		h.logger.Error("Failed to check email existence", logger.Error(err))
		c.Error(errors.NewDatabaseError("email check", err))
		return
	}

	if exists {
		c.Error(errors.NewConflictError("Email already exists"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Failed to hash password", logger.Error(err))
		c.Error(errors.NewInternalError("Failed to hash password", err))
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.userService.CreateUser(user); err != nil {
		h.logger.Error("Failed to create user", logger.Error(err))
		c.Error(errors.NewDatabaseError("user creation", err))
		return
	}

	tokenString, err := h.generateJWT(user)
	if err != nil {
		h.logger.Error("Failed to generate JWT", logger.Error(err))
		c.Error(errors.NewInternalError("Failed to generate token", err))
		return
	}

	user.Password = ""
	utils.SuccessResponse(c, http.StatusCreated, "User created successfully", gin.H{
		"user":  user,
		"token": tokenString,
	})
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		h.logger.Error("Failed to get user by email", logger.Error(err))
		c.Error(errors.NewDatabaseError("user retrieval", err))
		return
	}

	if user == nil {
		c.Error(errors.NewUnauthorizedError("Invalid email or password"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.Error(errors.NewUnauthorizedError("Invalid email or password"))
		return
	}

	tokenString, err := h.generateJWT(user)
	if err != nil {
		h.logger.Error("Failed to generate JWT", logger.Error(err))
		c.Error(errors.NewInternalError("Failed to generate token", err))
		return
	}

	user.Password = ""
	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"user":  user,
		"token": tokenString,
	})
}

func (h *Handler) generateJWT(user *models.User) (string, error) {
	expiresAt := time.Now().Add(time.Duration(h.cfg.JWTExpirationMinute) * time.Minute)

	claims := jwt.RegisteredClaims{
		Subject:   user.ID.Hex(),
		Issuer:    h.cfg.JWTIssuer,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(h.cfg.JWTSecret))
}
