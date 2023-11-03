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

	SetupTestLogin(&i, t)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(LoginRoute, Login(&i))

	body, contentType := LoginTestFormData()

	req := httptest.NewRequest(http.MethodPost, "http://petrichormud.com/login", body)
	req.Header.Set("Content-Type", contentType)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestLoginSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(&i, t)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))

	body, contentType := LoginTestFormData()

	// TODO: Extract this test url to a constant?
	url := fmt.Sprintf("http://petrichormud.com%s", RegisterRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)

	body, contentType = LoginTestFormData()

	url = fmt.Sprintf("http://petrichormud.com%s", LoginRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestLoginWithWrongPassword(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(&i, t)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))

	body, contentType := LoginTestFormData()

	// TODO: Extract this test url to a constant?
	url := fmt.Sprintf("http://petrichormud.com%s", RegisterRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)

	body, contentType = LoginTestFormDataWithPW("wrongpassword")

	url = fmt.Sprintf("http://petrichormud.com%s", LoginRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestLoginWithMalformedFormData(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestLogin(&i, t)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))

	body, contentType := LoginTestFormData()

	// TODO: Extract this test url to a constant?
	url := fmt.Sprintf("http://petrichormud.com%s", RegisterRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)

	body, contentType = LoginTestMalformedFormData()

	url = fmt.Sprintf("http://petrichormud.com%s", LoginRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err = app.Test(req)
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

func LoginTestFormData() (io.Reader, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", "testify")
	writer.WriteField("password", "T3sted_tested")
	writer.Close()
	return body, writer.FormDataContentType()
}

func LoginTestMalformedFormData() (io.Reader, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", "testify")
	writer.Close()
	return body, writer.FormDataContentType()
}

func LoginTestFormDataWithPW(pw string) (io.Reader, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", "testify")
	writer.WriteField("password", pw)
	writer.Close()
	return body, writer.FormDataContentType()
}
