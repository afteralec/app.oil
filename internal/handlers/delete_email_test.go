package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/shared"
)

func TestDeleteEmailUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(AddEmailRoute, AddEmail(&i))
	app.Delete(EmailRoute, DeleteEmail(&i))

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
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

	// TODO: Turn this route into a generator
	url := fmt.Sprintf("%s/player/email/%d", shared.TestURL, email.ID)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestDeleteEmailDBError(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(AddEmailRoute, AddEmail(&i))
	app.Delete(EmailRoute, DeleteEmail(&i))

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
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

	// TODO: Turn this route into a generator
	url := fmt.Sprintf("%s/player/email/%d", shared.TestURL, email.ID)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	i.Close()
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestDeleteEmailUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(LogoutRoute, Logout(&i))
	app.Post(AddEmailRoute, AddEmail(&i))
	app.Delete(EmailRoute, DeleteEmail(&i))

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)
	SetupTestDeleteEmail(t, &i, "testify2", TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
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

	// Log in as a different user
	CallRegister(t, app, "testify2", TestPassword)
	res = CallLogin(t, app, "testify2", TestPassword)
	cookies = res.Cookies()
	sessionCookie = cookies[0]
	req = AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Turn this route into a generator
	url := fmt.Sprintf("%s/player/email/%d", shared.TestURL, email.ID)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestDeleteEmailInvalidID(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(AddEmailRoute, AddEmail(&i))
	app.Delete(EmailRoute, DeleteEmail(&i))

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	url := fmt.Sprintf("%s/player/email/%s", shared.TestURL, "invalid")
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestDeleteNonexistantEmail(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(AddEmailRoute, AddEmail(&i))
	app.Delete(EmailRoute, DeleteEmail(&i))

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
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

	// TODO: Turn this route into a generator
	url := fmt.Sprintf("%s/player/email/%d", shared.TestURL, email.ID)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	_, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestDeleteEmailSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(AddEmailRoute, AddEmail(&i))
	app.Delete(EmailRoute, DeleteEmail(&i))

	SetupTestDeleteEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
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

	// TODO: Turn this route into a generator
	url := fmt.Sprintf("%s/player/email/%d", shared.TestURL, email.ID)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err = app.Test(req)
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
