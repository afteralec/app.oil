package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/handlers/home"
	"petrichormud.com/app/internal/handlers/login"
	"petrichormud.com/app/internal/handlers/logout"
	newplayer "petrichormud.com/app/internal/handlers/player/new"
	usernamereserved "petrichormud.com/app/internal/handlers/player/reserved"
	"petrichormud.com/app/internal/middleware/sessiondata"
	"petrichormud.com/app/internal/queries"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}
	defer db.Close()

	q := queries.New(db)

	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	readTimeout := time.Second * time.Duration(readTimeoutSecondsCount)

	views := html.NewFileSystem(http.FS(viewsfs), ".html")

	config := fiber.Config{
		Views:       views,
		ViewsLayout: "web/views/layouts/main",
		ReadTimeout: readTimeout,
	}
	app := fiber.New(config)

	// TODO: Update this config to be more secure. Will depend on environment.
	s := session.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(sessiondata.New(s))

	app.Static("/", "./web/static")

	app.Get("/", home.New())

	app.Post("/login", login.New(s, q))
	app.Post("/logout", logout.New(s))

	player := app.Group("/player")
	player.Post("/new", newplayer.New(s, q))
	player.Post("/reserved", usernamereserved.New(q))

	log.Fatal(app.Listen(":8008"))
}
