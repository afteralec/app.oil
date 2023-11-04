package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/shared"
)

// TODO: To test this, grab the email id from the database - then, pull all the uuid tokens and match the right one

func TestVerifyPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Get(VerifyRoute, VerifyPage(&i))

	url := fmt.Sprintf("%s%s", shared.TestURL, VerifyRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}
