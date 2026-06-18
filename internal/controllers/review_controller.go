package controllers

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/services"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/utils"
	renderPkg "github.com/unrolled/render"
	"go.uber.org/zap"
)

var render *renderPkg.Render

func init() {
	render = renderPkg.New()
}

type ReviewControllerInterface interface {
	CreateReview(resWriter http.ResponseWriter, req *http.Request)
	GetAllReviewsByHotelId(resWriter http.ResponseWriter, req *http.Request)
	GetReviewById(resWriter http.ResponseWriter, req *http.Request)
	UpdateReviewById(resWriter http.ResponseWriter, req *http.Request)
	DeleteReviewById(resWriter http.ResponseWriter, req *http.Request)
}

type ReviewController struct {
	ReviewService services.ReviewServiceInterface
	logger        *zap.Logger
	serverConfig  *config.ServerConfig
}

func (reviewController *ReviewController) CreateReview(resWriter http.ResponseWriter, req *http.Request) {
	reviewPayload := req.Context().Value("payload").(*dtos.CreateReviewDTO)

	// call the create user service
	serviceErr := reviewController.ReviewService.CreateReview(reviewPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while creating the review",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusCreated, resWriter, map[string]any{
		"success": true,
		"message": "Successfully created the review",
	})
}

func (reviewController *ReviewController) GetAllReviewsByHotelId(resWriter http.ResponseWriter, req *http.Request) {
	reviewParams := req.Context().Value("params").(*dtos.GetAllReviewsByHotelIdDTO)

	// call the fetch user by id service
	reviewModels, serviceErr := reviewController.ReviewService.GetAllReviewsByHotelId(reviewParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the reviews of the hotel",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully fetched the reviews of the hotel",
		"reviews": reviewModels,
	})
}

func (reviewController *ReviewController) GetReviewById(resWriter http.ResponseWriter, req *http.Request) {
	reviewParams := req.Context().Value("params").(*dtos.GetReviewByIdDTO)

	// call the fetch user by id service
	reviewModel, serviceErr := reviewController.ReviewService.GetReviewById(reviewParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the review by id",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully fetched the review by id",
		"review":  reviewModel,
	})
}

func (reviewController *ReviewController) UpdateReviewById(resWriter http.ResponseWriter, req *http.Request) {
	reviewParams := req.Context().Value("params").(*dtos.UpdateReviewByIdParams)
	reviewPayload := req.Context().Value("payload").(*dtos.UpdateReviewByIdDTO)

	// call the delete user service
	serviceErr := reviewController.ReviewService.UpdateReviewById(reviewParams, reviewPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while updating the review",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully updated the review",
	})
}

func (reviewController *ReviewController) DeleteReviewById(resWriter http.ResponseWriter, req *http.Request) {
	reviewParams := req.Context().Value("params").(*dtos.DeleteReviewByIdDTO)

	// call the delete user service
	serviceErr := reviewController.ReviewService.DeleteReviewById(reviewParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while deleting the review",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully deleted the review",
	})
}

func NewReviewController(service services.ReviewServiceInterface, logger *zap.Logger, serverConfig *config.ServerConfig) ReviewControllerInterface {
	newReviewController := &ReviewController{
		ReviewService: service,
		logger:        logger,
		serverConfig:  serverConfig,
	}

	return newReviewController
}
