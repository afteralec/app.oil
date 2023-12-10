package tests

import (
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

func TestLogout(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Post(handlers.LogoutRoute, handlers.Logout(&i))

	req := httptest.NewRequest(http.MethodPost, "http://petrichormud.com/logout", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestLogoutPage(t *testing.T) {
	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Get(handlers.LogoutRoute, handlers.LogoutPage())

	req := httptest.NewRequest(http.MethodGet, "http://petrichormud.com/logout", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
