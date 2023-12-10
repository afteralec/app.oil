package tests

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
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

func TestEditEmailUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

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

	req = EditEmailRequest(email.ID, TestEmailAddressTwo)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditEmailMissingInput(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

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
	req = httptest.NewRequest(http.MethodPut, url, nil)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditEmailMalformedInput(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

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

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notemail", "malformed")
	writer.Close()

	path := routes.EmailPath(strconv.FormatInt(email.ID, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req = httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditEmailDBError(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)
	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)
	SetupTestEditEmail(t, &i, "testify2", TestEmailAddressTwo)

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

	req = EditEmailRequest(email.ID, TestEmailAddressTwo)
	req.AddCookie(sessionCookie)
	// Close the connection to the DB to simulate a DB error
	i.Close()
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestEditEmailUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)
	SetupTestEditEmail(t, &i, "testify2", TestEmailAddressTwo)

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
	CallRegister(t, a, TestUsernameTwo, TestPassword)
	res = CallLogin(t, a, TestUsernameTwo, TestPassword)
	cookies = res.Cookies()
	sessionCookie = cookies[0]
	req = AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	req = EditEmailRequest(email.ID, TestEmailAddressTwo)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditEmailInvalidID(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	path := routes.EmailPath("invalid")
	url := fmt.Sprintf("%s%s", TestURL, path)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditNonexistantEmail(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

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
	req, err = http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(sessionCookie)

	_, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	req = EditEmailRequest(email.ID, TestEmailAddress)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestEditEmailUnverified(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

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

	req = EditEmailRequest(email.ID, TestEmailAddress)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditEmailSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestEditEmail(t, &i, TestUsername, TestEmailAddress)

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

	// TODO: Extract this block of functionality to a helper
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}
	emails, err := i.Queries.ListEmails(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	email := emails[0]
	_, err = i.Queries.MarkEmailVerified(context.Background(), email.ID)
	if err != nil {
		t.Fatal(err)
	}

	req = EditEmailRequest(email.ID, TestEmailAddress)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
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

	path := routes.EmailPath(strconv.FormatInt(id, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
