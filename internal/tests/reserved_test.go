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
	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/shared"
)

func TestReserved(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestReserved(&i, t)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.WriteField("password", TestPassword)
	writer.WriteField("confirmPassword", TestPassword)
	writer.Close()

	// TODO: Clean this up using the existing test utilities
	url := fmt.Sprintf("%s%s", TestURL, handlers.RegisterRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	body = new(bytes.Buffer)
	writer = multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.Close()

	url = fmt.Sprintf("%s%s", TestURL, handlers.ReservedRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusConflict, res.StatusCode)
}

func SetupTestReserved(i *shared.Interfaces, t *testing.T) {
	_, err := i.Database.Exec("DELETE FROM players WHERE username = 'testify';")
	if err != nil {
		t.Fatal(err)
	}
}
