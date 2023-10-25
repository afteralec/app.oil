package handlers

import (
	fiber "github.com/gofiber/fiber/v2"
)

func Apply(a *fiber.App) {
	a.Static("/", "./web/static")

	a.Get("/", Home)

	a.Post("/login", Login)
	a.Post("/logout", Logout)

	player := a.Group("player")
	player.Post("/new", NewPlayer)
	player.Post("/reserved", PlayerReserved)

	request := a.Group("request")
	request.Get("/:id", Request)
}
