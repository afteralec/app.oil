package handlers

import (
	"bytes"
	"database/sql"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	html "github.com/gofiber/template/html/v2"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/shared"
)

func TestLogin(t *testing.T) {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("SET GLOBAL local_infile=true;")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	opts := configs.Redis()
	r := redis.NewClient(&opts)
	defer r.Close()

	// TODO: Update this config to be more secure. Will depend on environment.
	s := session.New()

	i := shared.InterfacesBuilder().Database(db).Redis(r).Sessions(s).Build()

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(LoginRoute, Login(&i))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", "testify")
	writer.WriteField("password", "T3sted_tested")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "http://petrichormud.com/login", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}
