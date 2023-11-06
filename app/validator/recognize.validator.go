package validator

type ValidateRecognize struct {
	UserID string `json:"userId" validate:"required"`
	Image  string `json:"image" validate:"required"`
}
