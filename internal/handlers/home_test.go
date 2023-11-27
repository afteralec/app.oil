package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"
	"petrichormud.com/app/internal/configs"
)

func TestHomePage(t *testing.T) {
	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Get(HomeRoute, HomePage())

	req := httptest.NewRequest(http.MethodGet, "http://petrichormud.com", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
