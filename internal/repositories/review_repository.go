package repositories

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/utils"
	"go.uber.org/zap"
)

type ReviewRepositoryInterface interface {
	CreateReview(reviewPayload *dtos.CreateReviewDTO) *utils.AppError
	GetAllReviewsByHotelId(reviewPayload *dtos.GetAllReviewsByHotelIdDTO) ([]*models.ReviewModel, *utils.AppError)
	GetReviewById(reviewPayload *dtos.GetReviewByIdDTO) (*models.ReviewModel, *utils.AppError)
	UpdateReviewById(reviewId *dtos.UpdateReviewByIdParams, reviewPayload *dtos.UpdateReviewByIdDTO) *utils.AppError
	DeleteReviewById(reviewPayload *dtos.DeleteReviewByIdDTO) *utils.AppError
}

type ReviewRepository struct {
	db           *sql.DB
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (reviewRepository *ReviewRepository) CreateReview(reviewPayload *dtos.CreateReviewDTO) *utils.AppError {
	// insert into the db
	query := "INSERT INTO reviews (booking_id, hotel_id, user_id, rating, review_text, is_synced) VALUES (?, ?, ?, ?, ?, ?)"
	result, queryExecErr := reviewRepository.db.Exec(query, reviewPayload.BookingID, reviewPayload.HotelID, reviewPayload.UserID, reviewPayload.Rating, reviewPayload.ReviewText, reviewPayload.IsSynced)

	if queryExecErr != nil {
		reviewRepository.logger.Error("Failed to insert review into the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to insert review into the database: " + queryExecErr.Error())
	}

	id, insertErr := result.LastInsertId()

	if insertErr != nil {
		reviewRepository.logger.Error("Failed to insert review into the database",
			zap.String("error", insertErr.Error()))

		return utils.InternalServerError("Failed to insert review into the database: " + insertErr.Error())
	}

	reviewRepository.logger.Info("Successfully inserted review into the database",
		zap.Int64("review_id", id))

	return nil
}

func (reviewRepository *ReviewRepository) GetAllReviewsByHotelId(reviewPayload *dtos.GetAllReviewsByHotelIdDTO) ([]*models.ReviewModel, *utils.AppError) {
	var reviewModels []*models.ReviewModel

	// load the rows
	query := "SELECT * FROM reviews WHERE hotel_id = ?"
	rows, queryErr := reviewRepository.db.Query(query, reviewPayload.HotelID)

	if queryErr != nil {
		reviewRepository.logger.Error("Something went wrong while fetching all the reviews",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching all the reviews: " + queryErr.Error())
	}

	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		reviewModel := &models.ReviewModel{}

		rowScanErr := rows.Scan(&reviewModel.ID, &reviewModel.BookingID, &reviewModel.HotelID, &reviewModel.UserID, &reviewModel.Rating, &reviewModel.ReviewText, &reviewModel.IsSynced, &reviewModel.CreatedAt, &reviewModel.UpdatedAt)

		if rowScanErr != nil {
			reviewRepository.logger.Error("Failed to fetch all the reviews from the database",
				zap.String("error", rowScanErr.Error()))

			return nil, utils.InternalServerError("Something went wrong while fetching all the reviews: " + rowScanErr.Error())
		}

		reviewModels = append(reviewModels, reviewModel)
	}

	rowsErr := rows.Err()

	if rowsErr != nil {
		reviewRepository.logger.Error("Failed to fetch all the reviews from the database",
			zap.String("error", rowsErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching all the reviews: " + rowsErr.Error())
	}

	reviewRepository.logger.Info("Successfully fetched all the reviews from the database",
		zap.Int("count", len(reviewModels)))

	return reviewModels, nil
}

func (reviewRepository *ReviewRepository) GetReviewById(reviewPayload *dtos.GetReviewByIdDTO) (*models.ReviewModel, *utils.AppError) {
	// create the dummy instance
	reviewModel := &models.ReviewModel{}

	// fetch from the db
	query := "SELECT * FROM reviews WHERE id = ?"

	queryErr := reviewRepository.db.QueryRow(query, reviewPayload.ID).Scan(&reviewModel.ID, &reviewModel.BookingID, &reviewModel.HotelID, &reviewModel.UserID, &reviewModel.Rating, &reviewModel.ReviewText, &reviewModel.IsSynced, &reviewModel.CreatedAt, &reviewModel.UpdatedAt)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			reviewRepository.logger.Error("Such review not found",
				zap.Int("review_id", reviewPayload.ID))

			return nil, utils.NotFound("Review with such id not found")
		}

		reviewRepository.logger.Error("Failed to fetch the review from the database",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Failed to fetch the review from the database: " + queryErr.Error())
	}

	reviewRepository.logger.Info("Successfully fetched the review from the database",
		zap.Int("review_id", reviewModel.ID),
	)

	return reviewModel, nil
}

func (reviewRepository *ReviewRepository) UpdateReviewById(reviewId *dtos.UpdateReviewByIdParams, reviewPayload *dtos.UpdateReviewByIdDTO) *utils.AppError {
	// prepare and execute the query
	query := "UPDATE reviews SET rating = ?, review_text = ?, is_synced = ? WHERE id = ?"
	result, queryExecErr := reviewRepository.db.Exec(query, reviewPayload.Rating, reviewPayload.ReviewText, reviewPayload.IsSynced, reviewId.ID)

	if queryExecErr != nil {
		reviewRepository.logger.Error("Failed to update review from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to update review from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		reviewRepository.logger.Error("Failed to update review from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to update review from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		reviewRepository.logger.Error("No review has been updated from the database",
			zap.Int("review_id", reviewId.ID))

		return utils.NotFound("Review with such id not found")
	}

	reviewRepository.logger.Info("Successfully updated the review from the database",
		zap.Int("review_id", reviewId.ID))

	return nil
}

func (reviewRepository *ReviewRepository) DeleteReviewById(reviewPayload *dtos.DeleteReviewByIdDTO) *utils.AppError {
	// prepare and execute the query
	query := "DELETE FROM reviews WHERE id = ?"
	result, queryExecErr := reviewRepository.db.Exec(query, reviewPayload.ID)

	if queryExecErr != nil {
		reviewRepository.logger.Error("Failed to delete review from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to delete review from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		reviewRepository.logger.Error("Failed to delete review from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to delete review from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		reviewRepository.logger.Error("No review has been deleted from the database",
			zap.Int("review_id", reviewPayload.ID))

		return utils.NotFound("Review with such id not found")
	}

	reviewRepository.logger.Info("Successfully deleted the review from the database",
		zap.Int("review_id", reviewPayload.ID))

	return nil
}

func NewReviewRepository(logger *zap.Logger, db *sql.DB, serverConfig *config.ServerConfig) ReviewRepositoryInterface {
	newReviewRepository := &ReviewRepository{
		db:           db,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newReviewRepository
}
