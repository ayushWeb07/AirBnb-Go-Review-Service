package routers

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/controllers"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/middlewares"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type RouterInterface interface {
	Register(r *chi.Mux)
}

func RegisterRouters(logger *zap.Logger, db *sql.DB, serverConfig *config.ServerConfig) *chi.Mux {
	// create the router instance
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middlewares.RateLimiter(serverConfig))

	// register health router
	//SetupHealthRouter(router)

	// register review router
	reviewRepository := repositories.NewReviewRepository(logger, db, serverConfig)
	reviewService := services.NewReviewService(reviewRepository, logger, serverConfig)
	reviewController := controllers.NewReviewController(reviewService, logger, serverConfig)
	reviewRouter := NewReviewRouter(reviewController, logger, serverConfig)

	reviewRouter.Register(router)

	return router
}
