package services

import (
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/utils"
	"go.uber.org/zap"
)

type ReviewServiceInterface interface {
	CreateReview(reviewPayload *dtos.CreateReviewDTO) *utils.AppError
	GetAllReviewsByHotelId(reviewPayload *dtos.GetAllReviewsByHotelIdDTO) ([]*models.ReviewModel, *utils.AppError)
	GetReviewById(reviewPayload *dtos.GetReviewByIdDTO) (*models.ReviewModel, *utils.AppError)
	UpdateReviewById(reviewId *dtos.UpdateReviewByIdParams, reviewPayload *dtos.UpdateReviewByIdDTO) *utils.AppError
	DeleteReviewById(reviewPayload *dtos.DeleteReviewByIdDTO) *utils.AppError
}

type ReviewService struct {
	ReviewRepository repositories.ReviewRepositoryInterface
	logger           *zap.Logger
	serverConfig     *config.ServerConfig
}

func (reviewService *ReviewService) CreateReview(reviewPayload *dtos.CreateReviewDTO) *utils.AppError {
	reviewService.logger.Info("Create review service called...")

	// call the create review repository
	repositoryErr := reviewService.ReviewRepository.CreateReview(reviewPayload)
	return repositoryErr
}

func (reviewService *ReviewService) GetAllReviewsByHotelId(reviewPayload *dtos.GetAllReviewsByHotelIdDTO) ([]*models.ReviewModel, *utils.AppError) {
	reviewService.logger.Info("Get all reviews service called...")

	// call the fetch all reviews repository
	reviewModels, repositoryErr := reviewService.ReviewRepository.GetAllReviewsByHotelId(reviewPayload)
	return reviewModels, repositoryErr
}

func (reviewService *ReviewService) GetReviewById(reviewPayload *dtos.GetReviewByIdDTO) (*models.ReviewModel, *utils.AppError) {
	reviewService.logger.Info("Get by id review service called...")

	// call the fetch review by id repository
	reviewModel, repositoryErr := reviewService.ReviewRepository.GetReviewById(reviewPayload)
	return reviewModel, repositoryErr
}

func (reviewService *ReviewService) UpdateReviewById(reviewId *dtos.UpdateReviewByIdParams, reviewPayload *dtos.UpdateReviewByIdDTO) *utils.AppError {
	reviewService.logger.Info("Update by id review service called...")

	// call the update review by id repository
	repositoryErr := reviewService.ReviewRepository.UpdateReviewById(reviewId, reviewPayload)
	return repositoryErr
}

func (reviewService *ReviewService) DeleteReviewById(reviewPayload *dtos.DeleteReviewByIdDTO) *utils.AppError {
	reviewService.logger.Info("Delete review service called...")

	// call the delete review by id repository
	repositoryErr := reviewService.ReviewRepository.DeleteReviewById(reviewPayload)
	return repositoryErr
}

func NewReviewService(reviewRepository repositories.ReviewRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) ReviewServiceInterface {
	newReviewService := &ReviewService{
		ReviewRepository: reviewRepository,
		logger:           logger,
		serverConfig:     serverConfig,
	}

	return newReviewService
}
