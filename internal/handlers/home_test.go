package handlers

import (
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"
)

func TestHome(t *testing.T) {
	views := html.New("../..", ".html")
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	readTimeout := time.Second * time.Duration(readTimeoutSecondsCount)
	config := fiber.Config{
		Views:       views,
		ViewsLayout: "web/views/layouts/main",
		ReadTimeout: readTimeout,
	}
	app := fiber.New(config)

	app.Get(HomeRoute, Home())

	req := httptest.NewRequest("GET", "http://petrichormud.com", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
