package routers

import (
	"net/http"
	"strconv"

	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/controllers"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/middlewares"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ReviewRouter struct {
	ReviewController controllers.ReviewControllerInterface
	logger           *zap.Logger
	serverConfig     *config.ServerConfig
}

func (reviewRouter *ReviewRouter) Register(r *chi.Mux) {
	r.Route("/api/v1/reviews", func(r chi.Router) {
		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CreateReviewDTO]).Post("/", reviewRouter.ReviewController.CreateReview)

		r.With(middlewares.DecodeAndValidateParams[dtos.GetAllReviewsByHotelIdDTO](
			func(req *http.Request) (*dtos.GetAllReviewsByHotelIdDTO, *utils.AppError) {
				hotelId, err := strconv.Atoi(chi.URLParam(req, "hotel_id"))

				if err != nil {
					return nil, utils.InternalServerError("Hotel id must be provided in integer: " + err.Error())
				}

				return &dtos.GetAllReviewsByHotelIdDTO{
					HotelID: hotelId,
				}, nil
			},
		)).Get("/hotel/{hotel_id}", reviewRouter.ReviewController.GetAllReviewsByHotelId)

		r.With(middlewares.DecodeAndValidateParams[dtos.GetReviewByIdDTO](
			func(req *http.Request) (*dtos.GetReviewByIdDTO, *utils.AppError) {
				reviewId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.InternalServerError("Hotel id must be provided in integer: " + err.Error())
				}

				return &dtos.GetReviewByIdDTO{
					ID: reviewId,
				}, nil
			},
		)).Get("/{id}", reviewRouter.ReviewController.GetReviewById)

		r.With(middlewares.DecodeAndValidateParams[dtos.UpdateReviewByIdParams](
			func(req *http.Request) (*dtos.UpdateReviewByIdParams, *utils.AppError) {
				reviewId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.InternalServerError("Hotel id must be provided in integer: " + err.Error())
				}

				return &dtos.UpdateReviewByIdParams{
					ID: reviewId,
				}, nil
			},
		)).With(middlewares.DecodeAndValidateRequestBody[dtos.UpdateReviewByIdDTO]).Put("/{id}", reviewRouter.ReviewController.UpdateReviewById)

		r.With(middlewares.DecodeAndValidateParams[dtos.DeleteReviewByIdDTO](
			func(req *http.Request) (*dtos.DeleteReviewByIdDTO, *utils.AppError) {
				reviewId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.InternalServerError("Hotel id must be provided in integer: " + err.Error())
				}

				return &dtos.DeleteReviewByIdDTO{
					ID: reviewId,
				}, nil
			},
		)).Delete("/{id}", reviewRouter.ReviewController.DeleteReviewById)
	})
}

func NewReviewRouter(controller controllers.ReviewControllerInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RouterInterface {
	newReviewRouter := &ReviewRouter{
		ReviewController: controller,
		logger:           logger,
		serverConfig:     serverConfig,
	}

	return newReviewRouter
}
