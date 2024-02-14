package handler

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/view"
)

func DesignDictionaryPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(view.DesignDictionary, view.Bind(c), layout.Main)
	}
}
