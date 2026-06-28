package dtos

type FetchBookingDTO struct {
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
}
