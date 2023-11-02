package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/shared"
)

func TestLogin(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(&i, t)

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

func SetupTestLogin(i *shared.Interfaces, t *testing.T) {
	_, err := i.Database.Exec("DELETE FROM players WHERE username = 'testify';")
	if err != nil {
		t.Fatal(err)
	}
}
