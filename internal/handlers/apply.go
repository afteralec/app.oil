package handlers

import (
	fiber "github.com/gofiber/fiber/v2"
)

func Apply(a *fiber.App) {
	a.Static("/", "./web/static")

	a.Get("/", Home)

	player := a.Group("player")
	player.Post("/new", NewPlayer)

	request := a.Group("request")
	request.Get("/:id", Request)
}
