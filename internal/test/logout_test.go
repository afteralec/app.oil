package test

import (
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

func TestLogoutPageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.Logout)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

// TODO: Add some logic to verify that the user is logged out
func TestLogoutSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.Logout)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
