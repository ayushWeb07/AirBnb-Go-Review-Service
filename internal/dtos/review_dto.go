package dtos

type CreateReviewDTO struct {
	BookingID  int    `db:"booking_id" json:"booking_id" validate:"required,gte=1"`
	HotelID    int    `db:"hotel_id" json:"hotel_id" validate:"required,gte=1"`
	UserID     int    `db:"user_id" json:"user_id" validate:"required,gte=1"`
	Rating     int    `db:"rating" json:"rating" validate:"required,gte=1,lte=5"`
	ReviewText string `db:"review_text" json:"review_text" validate:"required,min=50,max=1000"`
	IsSynced   bool   `db:"is_synced" json:"is_synced"`
}

type GetAllReviewsByHotelIdDTO struct {
	HotelID int `json:"hotel_id" validate:"required,gte=1"`
}

type GetReviewByIdDTO struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type DeleteReviewByIdDTO struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type UpdateReviewByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type UpdateReviewByIdDTO struct {
	Rating     int    `db:"rating" json:"rating" validate:"required,gte=1,lte=5"`
	ReviewText string `db:"review_text" json:"review_text" validate:"required,min=50,max=1000"`
	IsSynced   bool   `db:"is_synced" json:"is_synced"`
}
