package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/utils"
	"go.uber.org/zap"
)

type ReviewServiceInterface interface {
	CreateReview(reviewPayload *dtos.CreateReviewDTO) *utils.AppError
	GetAllReviewsByHotelId(reviewParams *dtos.GetAllReviewsByHotelIdDTO) ([]*models.ReviewModel, *utils.AppError)
	GetReviewById(reviewParams *dtos.GetReviewByIdDTO) (*models.ReviewModel, *utils.AppError)
	UpdateReviewById(reviewParams *dtos.UpdateReviewByIdParams, reviewPayload *dtos.UpdateReviewByIdDTO) *utils.AppError
	DeleteReviewById(reviewParams *dtos.DeleteReviewByIdDTO) *utils.AppError
}

type ReviewService struct {
	ReviewRepository repositories.ReviewRepositoryInterface
	logger           *zap.Logger
	serverConfig     *config.ServerConfig
}

func (reviewService *ReviewService) CreateReview(reviewPayload *dtos.CreateReviewDTO) *utils.AppError {
	reviewService.logger.Info("Create review service called...")

	// check if the user exists
	apiGatewayUrl := reviewService.serverConfig.ApiGatewayBaseUrl + "/users/" + strconv.Itoa(reviewPayload.UserID)

	resp, err := http.Get(apiGatewayUrl)
	if err != nil {
		reviewService.logger.Error("Failed to make request to api gateway: " + err.Error())
		return utils.InternalServerError("Failed to make request to api gateway: " + err.Error())
	}

	resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {

		reviewService.logger.Error("Such user not found",
			zap.Int("user_id", reviewPayload.UserID))

		return utils.NotFound("User with such id not found")

	} else if resp.StatusCode != http.StatusOK {

		reviewService.logger.Error("Something went wrong while checking if the user exists",
			zap.Int("user_id", reviewPayload.UserID))

		return utils.InternalServerError("Something went wrong while checking if the user exists")

	}

	// check if the hotel exists
	hotelServiceUrl := reviewService.serverConfig.HotelServiceBaseUrl + "/hotels/" + strconv.Itoa(reviewPayload.HotelID)

	resp, err = http.Get(hotelServiceUrl)
	if err != nil {
		reviewService.logger.Error("Failed to make request to hotel service: " + err.Error())
		return utils.InternalServerError("Failed to make request to hotel service: " + err.Error())
	}

	resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {

		reviewService.logger.Error("Such hotel not found",
			zap.Int("hotel_id", reviewPayload.HotelID))

		return utils.NotFound("Hotel with such id not found")

	} else if resp.StatusCode != http.StatusOK {

		reviewService.logger.Error("Something went wrong while checking if the hotel exists",
			zap.Int("hotel_id", reviewPayload.HotelID))

		return utils.InternalServerError("Something went wrong while checking if the hotel exists")

	}

	// check if the booking exists
	bookingServiceUrl := reviewService.serverConfig.BookingServiceBaseUrl + "/bookings/" + strconv.Itoa(reviewPayload.BookingID)

	resp, err = http.Get(bookingServiceUrl)
	if err != nil {
		reviewService.logger.Error("Failed to make request to booking service: " + err.Error())
		return utils.InternalServerError("Failed to make request to booking service: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {

		reviewService.logger.Error("Such booking not found",
			zap.Int("booking_id", reviewPayload.BookingID))

		return utils.NotFound("Booking with such id not found")

	} else if resp.StatusCode != http.StatusOK {

		reviewService.logger.Error("Something went wrong while checking if the booking exists",
			zap.Int("booking_id", reviewPayload.BookingID))

		return utils.InternalServerError("Something went wrong while checking if the booking exists")

	}

	// check if the booking has confirmed status
	var fetchBookingResp *dtos.FetchBookingDTO

	if err := json.NewDecoder(resp.Body).Decode(&fetchBookingResp); err != nil {
		reviewService.logger.Error("Something went wrong while checking the booking status",
			zap.Int("booking_id", reviewPayload.BookingID))

		return utils.InternalServerError("Something went wrong while checking the booking status")
	}

	if fetchBookingResp.Data.Status != "confirmed" {
		reviewService.logger.Error("Your booking must be confirmed, in order to give a review",
			zap.Int("booking_id", reviewPayload.BookingID),
			zap.String("booking_status", fetchBookingResp.Data.Status))

		return utils.Forbidden("Your booking must be confirmed, in order to give a review")
	}

	// call the create review repository
	repositoryErr := reviewService.ReviewRepository.CreateReview(reviewPayload)
	return repositoryErr
}

func (reviewService *ReviewService) GetAllReviewsByHotelId(reviewParams *dtos.GetAllReviewsByHotelIdDTO) ([]*models.ReviewModel, *utils.AppError) {
	reviewService.logger.Info("Get all reviews service called...")

	// check if the hotel exists
	url := reviewService.serverConfig.HotelServiceBaseUrl + "/hotels/" + strconv.Itoa(reviewParams.HotelID)

	resp, err := http.Get(url)
	if err != nil {
		reviewService.logger.Error("Failed to make request to hotel service: " + err.Error())
		return nil, utils.InternalServerError("Failed to make request to hotel service: " + err.Error())
	}

	if resp.StatusCode == http.StatusNotFound {

		reviewService.logger.Error("Such hotel not found",
			zap.Int("hotel_id", reviewParams.HotelID))

		return nil, utils.NotFound("Hotel with such id not found")

	} else if resp.StatusCode != http.StatusOK {

		reviewService.logger.Error("Something went wrong while checking if the hotel exists",
			zap.Int("hotel_id", reviewParams.HotelID))

		return nil, utils.InternalServerError("Something went wrong while checking if the hotel exists")

	}

	// call the fetch all reviews repository
	reviewModels, repositoryErr := reviewService.ReviewRepository.GetAllReviewsByHotelId(reviewParams)
	return reviewModels, repositoryErr
}

func (reviewService *ReviewService) GetReviewById(reviewParams *dtos.GetReviewByIdDTO) (*models.ReviewModel, *utils.AppError) {
	reviewService.logger.Info("Get by id review service called...")

	// call the fetch review by id repository
	reviewModel, repositoryErr := reviewService.ReviewRepository.GetReviewById(reviewParams)
	return reviewModel, repositoryErr
}

func (reviewService *ReviewService) UpdateReviewById(reviewParams *dtos.UpdateReviewByIdParams, reviewPayload *dtos.UpdateReviewByIdDTO) *utils.AppError {
	reviewService.logger.Info("Update by id review service called...")

	// call the update review by id repository
	repositoryErr := reviewService.ReviewRepository.UpdateReviewById(reviewParams, reviewPayload)
	return repositoryErr
}

func (reviewService *ReviewService) DeleteReviewById(reviewParams *dtos.DeleteReviewByIdDTO) *utils.AppError {
	reviewService.logger.Info("Delete review service called...")

	// call the delete review by id repository
	repositoryErr := reviewService.ReviewRepository.DeleteReviewById(reviewParams)
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
