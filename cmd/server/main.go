package main

import (
	api "TulaHackDays-Backend/internal/api"
	"TulaHackDays-Backend/internal/api/handlers"
	"TulaHackDays-Backend/internal/api/middleware"
	"TulaHackDays-Backend/internal/config"
	"TulaHackDays-Backend/internal/pkg/auth"
	"TulaHackDays-Backend/internal/repository/postgres"
	"TulaHackDays-Backend/internal/services"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	dbPool, err := postgres.NewPostgresPool(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbPool.Close()

	authService := auth.NewAuthService(&config.JWTConfig{
		Secret:     os.Getenv("JWT_SECRET"),
		Expiration: 24,
	})

	txManager := postgres.NewTxManager(dbPool)
	userRepo := postgres.NewUserRepository(dbPool)
	achievementRepo := postgres.NewAchievementRepository(dbPool)
	placeRepo := postgres.NewPlaceRepository(dbPool)
	reviewRepo := postgres.NewReviewRepository(dbPool)

	userSvc := services.NewUserService(userRepo, txManager, authService)
	achievementSvc := services.NewAchievementService(achievementRepo, txManager)
	placeSvc := services.NewPlaceService(placeRepo, txManager)
	reviewSvc := services.NewReviewService(reviewRepo, txManager)

	userHandler := handlers.NewUserHandler(userSvc, achievementSvc)
	achievementHandler := handlers.NewAchievementHandler(achievementSvc)
	placeHandler := handlers.NewPlaceHandler(placeSvc)
	reviewHandler := handlers.NewReviewHandler(reviewSvc)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	gin.SetMode(gin.ReleaseMode)
	router := api.NewRouter(authMiddleware, userHandler, achievementHandler, placeHandler, reviewHandler)

	// Start server
	serverAddr := os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
