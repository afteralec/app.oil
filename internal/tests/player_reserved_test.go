package tests

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
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/route"
)

func TestReservedConflict(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.Close()

	url := MakeTestURL(route.Reserved)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusConflict, res.StatusCode)
}

func TestReservedFatal(t *testing.T) {
	i := interfaces.SetupShared()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)

	i.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.Close()

	url := MakeTestURL(route.Reserved)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = interfaces.SetupShared()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestReservedOK(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.Close()

	url := MakeTestURL(route.Reserved)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
