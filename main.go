package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
  "time"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}

	var hello string
	if err := db.QueryRow("SELECT 'hello from MySQL!';").Scan(&hello); err != nil {
		log.Fatalf("query: %v", err)
	}

	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("web/views/index", fiber.Map{
      "CopyrightYear": time.Now().Year(),
			"Title": "Hello, World!",
		}, "web/views/layouts/main")
	})

  app.Get("/request/:id", func(c *fiber.Ctx) error {
    return c.Render("web/views/request", fiber.Map{
      "CopyrightYear": time.Now().Year(),
      "ID": c.Params("id"),
      "Status": "Ready",
      "Name": "Test Character",
      "Backstory": "This is a tragic backstory.\nWith a newline.",
      "ShortDescription": "test, testerly man",
      "Description": "This is a test description.",
      "Class": "Crafting",
      "Origin": "LowQuarter",
      "Gender": "Male",
    }, "web/views/layouts/main")
  })

	app.Static("/", "./web/static")

	app.Listen(":8008")
}
