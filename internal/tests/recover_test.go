package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/middleware/bind"
)

func TestRecoverPage(t *testing.T) {
	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(bind.New())

	app.Get(handlers.RecoverRoute, handlers.RecoverPage())

	url := fmt.Sprintf("%s%s", TestURL, handlers.RecoverRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
