package validator

type ValidateRecognize struct {
	UserID string   `json:"userId" validate:"required"`
	Images []string `json:"images" validate:"required"`
}
