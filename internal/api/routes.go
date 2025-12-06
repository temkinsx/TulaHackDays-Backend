package api

import (
	"net/http"

	"TulaHackDays-Backend/internal/api/handlers"
	"TulaHackDays-Backend/internal/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter создаёт и настраивает Gin роутер со всеми HTTP-группами.
func NewRouter(authMiddleware *middleware.AuthMiddleware, userHandler *handlers.UserHandler, placeHandler *handlers.PlaceHandler, reviewHandler *handlers.ReviewHandler) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	{
		registerAuthRoutes(v1, userHandler)
		registerUserRoutes(v1, authMiddleware, userHandler)
		registerPlaceRoutes(v1, authMiddleware, placeHandler, reviewHandler)
		registerReviewRoutes(v1, authMiddleware, reviewHandler)
	}

	return router
}

func registerAuthRoutes(v1 *gin.RouterGroup, userHandler *handlers.UserHandler) {
	auth := v1.Group("/auth")
	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)
}

func registerUserRoutes(v1 *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, userHandler *handlers.UserHandler) {
	users := v1.Group("/users")
	users.GET("/leaderboard", userHandler.GetLeaderboard)

	usersAuthenticated := users.Group("")
	usersAuthenticated.Use(authMiddleware.AuthRequired())
	{
		usersAuthenticated.GET("/profile", userHandler.GetProfile)
		usersAuthenticated.PUT("/profile", userHandler.UpdateProfile)
		usersAuthenticated.POST("/change-password", userHandler.ChangePassword)
		usersAuthenticated.GET("/achievements", userHandler.GetUserAchievements)
	}
}

func registerAchievementRoutes(v1 *gin.RouterGroup, achievementHandler *handlers.AchievementHandler) {
	achievement := v1.Group("/achievement")
	achievement.GET("/", achievementHandler.ListActive)
}

func registerPlaceRoutes(v1 *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, placeHandler *handlers.PlaceHandler, reviewHandler *handlers.ReviewHandler) {
	places := v1.Group("/places")
	places.GET("/search", placeHandler.SearchPlaces)
	places.GET("/nearby", placeHandler.GetNearbyPlaces)
	places.GET("/:id", placeHandler.GetPlace)
	places.GET("/:place_id/reviews", reviewHandler.GetReviewsByPlace)

	protected := places.Group("")
	protected.Use(authMiddleware.AuthRequired())
	{
		protected.POST("", placeHandler.CreatePlace)
		protected.PUT("/:id", placeHandler.UpdatePlace)
		protected.DELETE("/:id", placeHandler.DeletePlace)
	}
}

func registerReviewRoutes(v1 *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, reviewHandler *handlers.ReviewHandler) {
	reviews := v1.Group("/reviews")
	reviews.GET("/:id", reviewHandler.GetReview)
	reviews.GET("/:id/comments", reviewHandler.GetCommentsByReview)

	protected := reviews.Group("")
	protected.Use(authMiddleware.AuthRequired())
	{
		protected.POST("", reviewHandler.CreateReview)
		protected.PUT("/:id", reviewHandler.UpdateReview)
		protected.DELETE("/:id", reviewHandler.DeleteReview)
		protected.POST("/:id/comments", reviewHandler.AddComment)
	}
}
