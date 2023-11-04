package setup

import (
	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/shared"
)

func App(app *fiber.App, i *shared.Interfaces) {
	Middleware(app, i)
	Static(app)
	Handlers(app, i)
}
