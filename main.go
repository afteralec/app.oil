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

	"petrichormud.com/app/internal/handlers"
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

	app.Get("/", handlers.Home())

	app.Post("/login", handlers.Login(s, q, r))
	app.Post("/logout", handlers.Logout(s))

	player := app.Group("/player")
	player.Post("/new", handlers.CreatePlayer(db, s, q, r))
	player.Post("/reserved", handlers.UsernameReserved(q))
	email := player.Group("/email")
	email.Post("/new", handlers.AddEmail(db, s, q, r))

	app.Get("/verify", handlers.Verify(q, r))
	app.Post("/verify", handlers.VerifyEmail(q, r))

	app.Get("/profile", handlers.Profile(q, r))

	log.Fatal(app.Listen(":8008"))
}

type Shared struct {
	Database *sql.DB
	Redis    *redis.Client
	Queries  *queries.Queries
	Sessions *session.Store
}

type SharedBuilder struct {
	Shared Shared
}

func (builder *SharedBuilder) Database(db *sql.DB) {
	builder.Shared.Database = db
	builder.Shared.Queries = queries.New(db)
}

func (builder *SharedBuilder) Redis(r *redis.Client) {
	builder.Shared.Redis = r
}

func (builder *SharedBuilder) Sessions(s *session.Store) {
	builder.Shared.Sessions = s
}

func (builder *SharedBuilder) Build() Shared {
	return builder.Shared
}
