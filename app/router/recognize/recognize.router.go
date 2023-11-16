package recognizeRouter

import (
	"bytes"
	"encoding/json"
	"face-recognition-golang/db"
	"face-recognition-golang/middleware"
	"face-recognition-golang/validator"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CreateRouter(app fiber.Router) {
	app.Post("/recognize", middleware.ValidateInput(validator.ValidateRecognize{}, true), func(c *fiber.Ctx) error {
		input := c.Locals("input").(*validator.ValidateRecognize)

		isVerifies, err := VerifiedHandler(db.Client, input)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"data":   err.Error(),
			})
		}

		// handle Http request
		type Image struct {
			Url        string `json:"url"`
			IsVerified bool   `json:"isVerified"`
		}
		type Body struct {
			UserId string  `json:"userId"`
			Images []Image `json:"images"`
		}

		body := Body{
			UserId: input.UserID,
			Images: []Image{},
		}

		for _, item := range isVerifies {
			body.Images = append(body.Images, Image{
				Url:        item.Image,
				IsVerified: item.IsAuth,
			})
		}

		bodyToByte, err := json.Marshal(body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"data":   err.Error(),
			})
		}
		fmt.Println("start call:", string(bodyToByte))
		_, err = http.Post("https://finder.sohe.in/api/v1/images/verified", "application/json", bytes.NewBuffer(bodyToByte))
		fmt.Println("end call", err)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"data":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"data":   isVerifies,
		})
	})
}
