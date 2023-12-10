package tests

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestResetPasswordPage(t *testing.T) {
	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(bind.New())

	app.Get(handlers.ResetPasswordRoute, handlers.ResetPasswordPage())

	id := uuid.NewString()
	url := fmt.Sprintf("%s%s?t=%s", TestURL, handlers.ResetPasswordRoute, id)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestResetPasswordSuccessPage(t *testing.T) {
	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(bind.New())

	app.Get(handlers.ResetPasswordSuccessRoute, handlers.ResetPasswordSuccessPage())

	url := fmt.Sprintf("%s%s", TestURL, handlers.ResetPasswordSuccessRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestResetPasswordMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Post(handlers.ResetPasswordRoute, handlers.ResetPassword(&i))

	SetupTestResetPassword(t, &i, TestUsername, TestEmailAddress)

	url := fmt.Sprintf("%s%s", TestURL, handlers.ResetPasswordRoute)
	req := httptest.NewRequest(http.MethodPost, url, nil)

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestResetPasswordMalformedBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Post(handlers.ResetPasswordRoute, handlers.ResetPassword(&i))

	SetupTestResetPassword(t, &i, TestUsername, TestEmailAddress)

	url := fmt.Sprintf("%s%s", TestURL, handlers.ResetPasswordRoute)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notusername", "notausername")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestResetPasswordSuccess(t *testing.T) {
	t.Skip()
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(handlers.RegisterRoute, handlers.Register(&i))
	app.Post(handlers.LoginRoute, handlers.Login(&i))
	app.Post(routes.NewEmailPath(), handlers.AddEmail(&i))
	app.Post(handlers.RecoverPasswordRoute, handlers.RecoverPassword(&i))
	app.Post(handlers.ResetPasswordRoute, handlers.ResetPassword(&i))

	SetupTestRecoverPassword(t, &i, TestUsername, TestEmailAddress)

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
	e := emails[0]

	_, err = i.Queries.MarkEmailVerified(context.Background(), e.ID)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("%s%s", TestURL, handlers.RecoverPasswordRoute)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.WriteField("password", TestPassword)
	writer.WriteField("confirm", TestPassword)
	writer.Close()
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func SetupTestResetPassword(t *testing.T, i *shared.Interfaces, u string, e string) {
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
