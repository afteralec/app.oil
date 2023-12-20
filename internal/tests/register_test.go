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
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

// TODO: Add failure tests here for bad inputs
func TestRegister(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestRegister(t, &i, TestUsername)

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	res := CallRegister(t, a, TestUsername, TestPassword)

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func CallRegister(t *testing.T, app *fiber.App, u string, pw string) *http.Response {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", u)
	writer.WriteField("password", pw)
	writer.WriteField("confirmPassword", pw)
	writer.Close()

	url := MakeTestURL(routes.Register)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	return res
}

func SetupTestRegister(t *testing.T, i *shared.Interfaces, u string) {
	DeleteTestPlayer(t, i, u)
}
