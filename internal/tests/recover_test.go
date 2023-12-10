package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestRecoverPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, routes.Recover)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
