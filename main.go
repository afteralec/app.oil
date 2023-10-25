package main

import (
	"embed"
	"log"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/middleware"
	"petrichormud.com/app/internal/queries"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	queries.Connect()
	queries.Build()
	defer queries.Disconnect()
	config := configs.Fiber(viewsfs)
	app := fiber.New(config)

	middleware.Apply(app)
	handlers.Apply(app)

	log.Fatal(app.Listen(":8008"))
}
