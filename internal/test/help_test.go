package test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
)

func TestHelpPageFatal(t *testing.T) {
	i := service.NewInterfaces()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.Help)

	i.Close()

	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestHelpPageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.Help)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

// TODO: Need to figure out seeding help data during a test run
func TestHelpFilePageNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.HelpFilePath("notahelpfile"))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestHelpFilePageFatal(t *testing.T) {
	i := service.NewInterfaces()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.HelpFilePath("notahelpfile"))

	i.Close()

	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestHelpFilePageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	TestHelpFile.PID = pid
	CreateTestHelpFile(t, &i, TestHelpFile)
	defer DeleteTestHelpFile(t, &i, TestHelpFile.Slug)

	url := MakeTestURL(route.HelpFilePath("test"))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestSearchHelpNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("search", "this string doesn't show up anywhere in help files")
	writer.Close()

	url := MakeTestURL(route.Help)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestSearchHelpSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	TestHelpFile.PID = pid
	CreateTestHelpFile(t, &i, TestHelpFile)
	defer DeleteTestHelpFile(t, &i, TestHelpFile.Slug)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("search", "test")
	writer.WriteField("title", "true")
	writer.Close()

	url := MakeTestURL(route.Help)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
