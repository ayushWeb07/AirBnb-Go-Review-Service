package models

type ReviewModel struct {
	ID         string `db:"id" json:"id"`
	BookingID  string `db:"booking_id" json:"booking_id" validate:"required"`
	HotelID    string `db:"hotel_id" json:"hotel_id" validate:"required"`
	UserID     string `db:"user_id" json:"user_id" validate:"required"`
	Rating     int    `db:"rating" json:"rating" validate:"required"`
	IsSynced   bool   `db:"is_synced" json:"is_synced" validate:"required"`
	ReviewText string `db:"review_text" json:"review_text" validate:"required"`
	CreatedAt  string `db:"created_at" json:"createdAt"`
	UpdatedAt  string `db:"updated_at" json:"updatedAt"`
}
