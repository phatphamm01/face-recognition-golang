package router

import (
	recognizeRouter "face-recognition-golang/router/recognize"
	socketRouter "face-recognition-golang/router/socket"

	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	recognizeRouter.CreateRouter(app)
	socketRouter.CreateRouter(app)
}
