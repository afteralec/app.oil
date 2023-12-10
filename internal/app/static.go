package app

import fiber "github.com/gofiber/fiber/v2"

func Static(app *fiber.App) {
	app.Static("/", "./web/static")
	app.Static("/loaders", "./web/svg/loaders")
}
