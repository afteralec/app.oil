package handlers

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

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/shared"
)

const TestEmailAddressTwo = "testify2@test.com"

func TestEditEmailUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(AddEmailRoute, AddEmail(&i))
	app.Put(EmailRoute, EditEmail(&i))

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

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

	req = EditEmailRequest(email.ID, TestEmailAddressTwo)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditEmailDBError(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(AddEmailRoute, AddEmail(&i))
	app.Put(EmailRoute, EditEmail(&i))

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)
	SetupTestEditEmail(t, &i, "testify2", TestEmailAddressTwo)

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

	req = EditEmailRequest(email.ID, TestEmailAddressTwo)
	req.AddCookie(sessionCookie)
	// Close the connection to the DB to simulate a DB error
	i.Close()
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestEditEmailUnowned(t *testing.T) {
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
	app.Put(EmailRoute, EditEmail(&i))

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)
	SetupTestEditEmail(t, &i, "testify2", TestEmailAddressTwo)

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
	CallRegister(t, app, TestUsernameTwo, TestPassword)
	res = CallLogin(t, app, TestUsernameTwo, TestPassword)
	cookies = res.Cookies()
	sessionCookie = cookies[0]
	req = AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	req = EditEmailRequest(email.ID, TestEmailAddressTwo)
	req.AddCookie(sessionCookie)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditEmailInvalidID(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(AddEmailRoute, AddEmail(&i))
	app.Put(EmailRoute, EditEmail(&i))

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	url := fmt.Sprintf("%s/player/email/%s", shared.TestURL, "invalid")
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(sessionCookie)

	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditNonexistantEmail(t *testing.T) {
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
	app.Put(EmailRoute, EditEmail(&i))

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

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
	req, err = http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(sessionCookie)

	_, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	req = EditEmailRequest(email.ID, TestEmailAddress)
	req.AddCookie(sessionCookie)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestEditEmailUnverified(t *testing.T) {
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
	app.Put(EmailRoute, EditEmail(&i))

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

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

	req = EditEmailRequest(email.ID, TestEmailAddress)
	req.AddCookie(sessionCookie)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func SetupTestEditEmail(t *testing.T, i *shared.Interfaces, u string, e string) {
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

func EditEmailRequest(id int64, e string) *http.Request {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", e)
	writer.Close()

	url := fmt.Sprintf("%s/player/email/%d", shared.TestURL, id)
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
