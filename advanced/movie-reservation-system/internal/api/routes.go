package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/api/handlers"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/auth"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/service"
)

// SetupRoutes sets up all API routes
func SetupRoutes(router *gin.Engine, services *service.Services, auth *auth.AuthService) {
	// Public routes
	authGroup := router.Group("/api/auth")
	{
		userHandler := handlers.NewUserHandler(services.UserService)
		authGroup.POST("/signup", userHandler.Signup)
		authGroup.POST("/login", userHandler.Login)
		authGroup.POST("/refresh", userHandler.RefreshToken)
	}

	// Protected routes
	protected := router.Group("/")
	protected.Use(authMiddleware(auth))

	// User routes
	users := protected.Group("/api/users")
	{
		userHandler := handlers.NewUserHandler(services.UserService)
		users.GET("/:id", userHandler.GetByID)
		users.GET("", userHandler.List)

		// Admin-only routes
		adminUsers := users.Group("")
		adminUsers.Use(adminMiddleware())
		{
			adminUsers.POST("/:id/promote", userHandler.PromoteToAdmin)
		}
	}

	// Movie routes
	movies := protected.Group("/api/movies")
	{
		movieHandler := handlers.NewMovieHandler(services.MovieService)
		movies.GET("", movieHandler.GetAll)
		movies.GET("/genres", movieHandler.GetAllGenres)
		movies.GET("/genre/:genreId", movieHandler.FindByGenre)
		movies.GET("/:id", movieHandler.GetByID)

		// Admin-only routes
		adminMovies := movies.Group("")
		adminMovies.Use(adminMiddleware())
		{
			adminMovies.POST("", movieHandler.Create)
			adminMovies.PUT("/:id", movieHandler.Update)
			adminMovies.DELETE("/:id", movieHandler.Delete)
		}
	}

	// Showtime routes
	showtimes := protected.Group("/api/showtimes")
	{
		showtimeHandler := handlers.NewShowtimeHandler(services.ShowtimeService)
		showtimes.GET("", showtimeHandler.GetAll)
		showtimes.GET("/by-date", showtimeHandler.GetByDate)
		showtimes.GET("/movie/:movieId", showtimeHandler.GetByMovieID)
		showtimes.GET("/:id", showtimeHandler.GetByID)

		// Admin-only routes
		adminShowtimes := showtimes.Group("")
		adminShowtimes.Use(adminMiddleware())
		{
			adminShowtimes.POST("", showtimeHandler.Create)
			adminShowtimes.PUT("/:id", showtimeHandler.Update)
			adminShowtimes.DELETE("/:id", showtimeHandler.Delete)
		}
	}

	// Reservation routes
	reservations := protected.Group("/api/reservations")
	{
		reservationHandler := handlers.NewReservationHandler(services.ReservationService)
		reservations.GET("/me", reservationHandler.GetAllByUserID)
		reservations.GET("/:id", reservationHandler.GetByID)
		reservations.POST("", reservationHandler.Create)
		reservations.DELETE("/:id", reservationHandler.Cancel)

		// Admin-only routes
		adminReservations := reservations.Group("")
		adminReservations.Use(adminMiddleware())
		{
			adminReservations.GET("", reservationHandler.GetAll)
			adminReservations.GET("/stats", reservationHandler.GetStats)
		}
	}
}

// authMiddleware validates JWT tokens
func authMiddleware(auth *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := auth.GetTokenFromHeader(authHeader)

		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing authorization header"})
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// adminMiddleware checks if user is an admin
func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		if roleStr, ok := role.(string); ok && roleStr == "admin" {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "forbidden: admin access required"})
	}
}
