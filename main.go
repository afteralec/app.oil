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
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/handlers/home"
	"petrichormud.com/app/internal/handlers/login"
	"petrichormud.com/app/internal/handlers/logout"
	newemail "petrichormud.com/app/internal/handlers/player/email/new"
	newplayer "petrichormud.com/app/internal/handlers/player/new"
	usernamereserved "petrichormud.com/app/internal/handlers/player/reserved"
	"petrichormud.com/app/internal/handlers/profile"
	"petrichormud.com/app/internal/handlers/verify"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/sessiondata"
	"petrichormud.com/app/internal/queries"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("SET GLOBAL local_infile=true;")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := queries.New(db)

	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
		Protocol: 3,
	})

	views := html.NewFileSystem(http.FS(viewsfs), ".html")
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	readTimeout := time.Second * time.Duration(readTimeoutSecondsCount)
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
	app.Use(sessiondata.New(s, q, r))
	app.Use(bind.New())

	app.Static("/", "./web/static")

	app.Get("/", home.New())

	app.Get("/test", func(c *fiber.Ctx) error {
		b := c.Locals("bind").(fiber.Map)
		b["Username"] = "alec"
		b["Email"] = "after.alec@gmail.com"
		b["VerifyToken"] = c.Params("t")
		return c.Render("web/views/verify", c.Locals("bind"), "web/views/layouts/standalone")
	})

	app.Post("/login", login.New(s, q, r))
	app.Post("/logout", logout.New(s))

	player := app.Group("/player")
	player.Post("/new", newplayer.New(db, s, q, r))
	player.Get("/:id", profile.New(q, r))
	player.Post("/reserved", usernamereserved.New(q))
	email := player.Group("/email")
	email.Post("/new", newemail.New(db, s, q, r))

	app.Get("/verify", verify.New(q, r))
	app.Post("/verify", verify.NewVerify(q, r))

	app.Get("/profile", profile.NewWithoutParams(q, r))

	log.Fatal(app.Listen(":8008"))
}
