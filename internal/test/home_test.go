package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	route "petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
)

func TestHomePage(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	req := httptest.NewRequest(http.MethodGet, MakeTestURL(route.Home), nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
