package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestDeleteEmailUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

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

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}
	emails, err := i.Queries.ListEmails(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	email := emails[0]

	path := routes.EmailPath(strconv.FormatInt(email.ID, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestDeleteEmailDBError(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

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

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}
	emails, err := i.Queries.ListEmails(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	email := emails[0]

	path := routes.EmailPath(strconv.FormatInt(email.ID, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	i.Close()
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestDeleteEmailUnowned(t *testing.T) {
	// TODO: Figure out why this test isn't working
	t.Skip()
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Handlers(a, &i)
	app.Middleware(a, &i)

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)
	SetupTestDeleteEmail(t, &i, "testify2", TestEmailAddress)

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

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}
	emails, err := i.Queries.ListEmails(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	email := emails[0]

	// Log in as a different user
	CallRegister(t, a, "testify2", TestPassword)
	res = CallLogin(t, a, "testify2", TestPassword)
	cookies = res.Cookies()
	sessionCookie = cookies[0]
	req = AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	path := routes.EmailPath(strconv.FormatInt(email.ID, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestDeleteEmailInvalidID(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	path := routes.EmailPath("invalid")
	// TODO: Make this TestURL prefix into a func
	url := fmt.Sprintf("%s%s", TestURL, path)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestDeleteNonexistantEmail(t *testing.T) {
	// TODO: Figure out why this test isn't working
	t.Skip()
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Handlers(a, &i)
	app.Middleware(a, &i)

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

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

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}
	emails, err := i.Queries.ListEmails(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	email := emails[0]

	path := routes.EmailPath(strconv.FormatInt(email.ID, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	_, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestDeleteEmailSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

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

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}
	emails, err := i.Queries.ListEmails(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	email := emails[0]

	path := routes.EmailPath(strconv.FormatInt(email.ID, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func SetupTestDeleteEmail(t *testing.T, i *shared.Interfaces, u string, e string) {
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
