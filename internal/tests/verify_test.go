package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestVerifyPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Get(handlers.VerifyRoute, handlers.VerifyPage(&i))

	url := fmt.Sprintf("%s%s", TestURL, handlers.VerifyRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestVerifyUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(handlers.VerifyRoute, handlers.Verify(&i))

	url := fmt.Sprintf("%s%s", TestURL, handlers.VerifyRoute)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestVerifyNoToken(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(handlers.RegisterRoute, handlers.Register(&i))
	app.Post(handlers.LoginRoute, handlers.Login(&i))
	app.Post(handlers.VerifyRoute, handlers.Verify(&i))

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := fmt.Sprintf("%s%s", TestURL, handlers.VerifyRoute)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestVerifyExpiredToken(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(handlers.RegisterRoute, handlers.Register(&i))
	app.Post(handlers.LoginRoute, handlers.Login(&i))
	app.Post(routes.NewEmailPath(), handlers.AddEmail(&i))
	app.Post(handlers.VerifyRoute, handlers.Verify(&i))

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]

	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}
	emails, err := i.Queries.ListEmails(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	email := emails[0]

	url := fmt.Sprintf("%s/player/email/%d/resend?t=non-existant-key", TestURL, email.ID)
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func SetupTestVerify(t *testing.T, i *shared.Interfaces, u string, e string) {
	query := fmt.Sprintf("DELETE FROM players WHERE username = '%s'", u)
	_, err := i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
	query = fmt.Sprintf("DELETE FROM emails WHERE address = '%s'", e)
	_, err = i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
