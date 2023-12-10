package tests

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestLoginPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, routes.Login)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestLoginNonExistantUser(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestLogin(t, &i, TestUsername)

	res := CallLogin(t, a, TestUsername, TestPassword)
	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestLoginSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(t, &i, TestUsername)

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestLoginWithWrongPassword(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(t, &i, TestUsername)

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, "wrongpassword")
	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestLoginWithMalformedFormData(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(t, &i, TestUsername)

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", "testify")
	writer.Close()

	url := fmt.Sprintf("%s%s", TestURL, routes.Login)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func CallLogin(t *testing.T, app *fiber.App, u string, pw string) *http.Response {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", u)
	writer.WriteField("password", pw)
	writer.Close()

	url := fmt.Sprintf("%s%s", TestURL, routes.Login)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	return res
}

func SetupTestLogin(t *testing.T, i *shared.Interfaces, u string) {
	query := fmt.Sprintf("DELETE FROM players WHERE username = '%s'", u)
	_, err := i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
