package middleware

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/sessiondata"
	"petrichormud.com/app/internal/shared"
)

func Setup(app *fiber.App, i *shared.Interfaces) {
	app.Use(logger.New())
	app.Use(csrf.New(configs.CSRF(i.Sessions)))
	app.Use(sessiondata.New(i))
	app.Use(bind.New())
}
