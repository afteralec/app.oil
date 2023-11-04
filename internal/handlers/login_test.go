package handlers

import (
	"bytes"
	"fmt"
	"io"
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

// TODO: Add operation for revoking and adding permissions

func TestLoginPage(t *testing.T) {
	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Get(LoginRoute, LoginPage())

	url := fmt.Sprintf("http://petrichormud.com%s", LoginRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestLoginNonExistantUser(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(t, &i, TestUsername)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(LoginRoute, Login(&i))

	res := CallLogin(t, app, TestUsername, TestPassword)
	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestLoginSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(t, &i, TestUsername)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))

	CallRegister(t, app, TestUsername, TestPassword)

	res := CallLogin(t, app, TestUsername, TestPassword)
	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestLoginWithWrongPassword(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(t, &i, TestUsername)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))

	CallRegister(t, app, TestUsername, TestPassword)

	res := CallLogin(t, app, TestUsername, "wrongpassword")
	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestLoginWithMalformedFormData(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(t, &i, TestUsername)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))

	CallRegister(t, app, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", "testify")
	writer.Close()

	url := fmt.Sprintf("%s%s", shared.TestURL, LoginRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := app.Test(req)
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

	url := fmt.Sprintf("%s%s", shared.TestURL, LoginRoute)
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

func LoginTestFormDataWithPW(pw string) (io.Reader, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", "testify")
	writer.WriteField("password", pw)
	writer.Close()
	return body, writer.FormDataContentType()
}
