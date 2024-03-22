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

// TODO: Add failure tests here for bad inputs
func TestRegister(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.WriteField("password", TestPassword)
	writer.WriteField("confirmPassword", TestPassword)
	writer.Close()

	url := MakeTestURL(route.Register)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}
