package domain

type MiddleSETERequest struct {
	Data string `json:"data" validate:"required"`
}
