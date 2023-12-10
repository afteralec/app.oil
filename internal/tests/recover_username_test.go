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
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestRecoverUsernamePage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.RecoverUsername)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestRecoverUsernameSuccessPageRedirectsWithoutToken(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.RecoverUsernameSuccess)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusFound, res.StatusCode)
}

func TestRecoverUsernameMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestRecoverUsername(t, &i, TestUsername, TestEmailAddress)

	url := MakeTestURL(routes.RecoverUsername)
	req := httptest.NewRequest(http.MethodPost, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestRecoverUsernameMalformedBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestRecoverUsername(t, &i, TestUsername, TestEmailAddress)

	url := MakeTestURL(routes.RecoverUsername)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notemail", "notanemail")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestRecoverUsernameSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestRecoverUsername(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	emails := ListEmailsForPlayer(t, &i, TestUsername)
	email := emails[0]
	_, err = i.Queries.MarkEmailVerified(context.Background(), email.ID)
	if err != nil {
		t.Fatal(err)
	}

	url := MakeTestURL(routes.RecoverUsername)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddress)
	writer.Close()
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func SetupTestRecoverUsername(t *testing.T, i *shared.Interfaces, u string, e string) {
	query := fmt.Sprintf("DELETE FROM players WHERE username = '%s';", u)
	_, err := i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
	query = fmt.Sprintf("DELETE FROM emails WHERE address = '%s';", e)
	_, err = i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
