package socketRouter

import (
	"github.com/gofiber/fiber/v2"
)

func CreateRouter(app fiber.Router) {
	app.Get("/ws-recognize-record", WsRecognizeRecord)
}
