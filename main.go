package main

import (
	"embed"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/shared"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.NewFileSystem(http.FS(viewsfs), ".html")
	a := fiber.New(configs.Fiber(views))

	app.Middleware(a, &i)
	app.Handlers(a, &i)
	app.Static(a)

	log.Fatal(a.Listen(":8008"))
}
