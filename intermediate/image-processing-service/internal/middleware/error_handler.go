package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/config"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/errors"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/utils"
	"go.uber.org/zap"
)

// ErrorHandler middleware handles application errors
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Handle custom application errors
			if appErr, ok := err.(*errors.AppError); ok {
				logger.Error("Application error",
					zap.String("error", appErr.Error()),
					zap.Int("status_code", appErr.Code),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				utils.ErrorResponse(c, appErr.Code, appErr.Message)
				return
			}

			// Handle generic errors
			logger.Error("Unhandled error",
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)

			utils.InternalErrorResponse(c, "An unexpected error occurred")
		}
	}
}

// RecoveryHandler middleware handles panics
func RecoveryHandler(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered",
			zap.Any("panic", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		utils.InternalErrorResponse(c, "Internal server error")
		c.Abort()
	})
}

// RequestLogger middleware logs incoming requests
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	})
}

// CORS middleware handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

type rateLimitEntry struct {
	windowStart time.Time
	count       int
}

func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	entries := make(map[string]rateLimitEntry)

	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			userID = c.ClientIP()
		}
		if userID == "" {
			c.Next()
			return
		}

		now := time.Now()
		mu.Lock()
		entry := entries[userID]
		if entry.windowStart.IsZero() || now.Sub(entry.windowStart) > window {
			entry.windowStart = now
			entry.count = 0
		}
		entry.count++
		entries[userID] = entry
		count := entry.count
		mu.Unlock()

		if count > limit {
			c.Error(errors.NewAppError(http.StatusTooManyRequests, "Rate limit exceeded", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Error(errors.NewUnauthorizedError("Invalid or missing authorization header"))
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(authHeader[len("Bearer "):])

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewUnauthorizedError("Unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			c.Error(errors.NewUnauthorizedError("Invalid token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			c.Error(errors.NewUnauthorizedError("Invalid token"))
			c.Abort()
			return
		}

		if claims.Subject == "" {
			c.Error(errors.NewUnauthorizedError("Invalid token subject"))
			c.Abort()
			return
		}

		if claims.Issuer != "" && claims.Issuer != cfg.JWTIssuer {
			c.Error(errors.NewUnauthorizedError("Invalid token issuer"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.Subject)
		c.Next()
	}
}
