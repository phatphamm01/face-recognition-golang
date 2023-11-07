package recognizeRouter

import (
	"face-recognition-golang/db"
	"face-recognition-golang/middleware"
	"face-recognition-golang/validator"
	"image"
	"os"

	"github.com/gofiber/fiber/v2"
)

func CreateRouter(app fiber.Router) {
	app.Post("/recognize", middleware.ValidateInput(validator.ValidateRecognize{}, true), func(c *fiber.Ctx) error {
		input := c.Locals("input").(*validator.ValidateRecognize)

		base64Image := VerifiedHandler(db.Client, input)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Thành công",
			"data":    base64Image,
		})
	})
}

func GetImageFromPath(path string) (*image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	return &img, err
}
