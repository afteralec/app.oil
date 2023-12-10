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

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/shared"
)

const (
	TestUsername    = "testify"
	TestUsernameTwo = "testify2"
	TestPassword    = "T3sted_tested"
)

// TODO: Add failure tests here for bad inputs
func TestRegister(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestRegister(t, &i, TestUsername)

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Post(handlers.RegisterRoute, handlers.Register(&i))

	res := CallRegister(t, app, TestUsername, TestPassword)

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func CallRegister(t *testing.T, app *fiber.App, u string, pw string) *http.Response {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", u)
	writer.WriteField("password", pw)
	writer.WriteField("confirmPassword", pw)
	writer.Close()

	url := fmt.Sprintf("%s%s", TestURL, handlers.RegisterRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	return res
}

func SetupTestRegister(t *testing.T, i *shared.Interfaces, u string) {
	query := fmt.Sprintf("DELETE FROM players WHERE username = '%s'", u)
	_, err := i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
