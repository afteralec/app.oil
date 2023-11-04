package main

import (
	"embed"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/middleware"
	"petrichormud.com/app/internal/setup"
	"petrichormud.com/app/internal/shared"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.NewFileSystem(http.FS(viewsfs), ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	middleware.Setup(app, &i)
	app.Static("/", "./web/static")
	app.Static("/loaders", "./web/svg/loaders")
	setup.Handlers(app, &i)

	log.Fatal(app.Listen(":8008"))
}
