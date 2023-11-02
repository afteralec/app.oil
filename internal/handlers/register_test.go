package handlers

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
	"petrichormud.com/app/internal/shared"
)

func TestRegister(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestRegister(&i, t)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Post(RegisterRoute, Register(&i))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", "testify")
	writer.WriteField("password", "T3sted_tested")
	writer.Close()

	// TODO: Extract this test url to a constant?
	url := fmt.Sprintf("http://petrichormud.com%s", RegisterRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func SetupTestRegister(i *shared.Interfaces, t *testing.T) {
	_, err := i.Database.Exec("DELETE FROM players WHERE username = 'testify';")
	if err != nil {
		t.Fatal(err)
	}
}
