package handlers

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/shared"
)

const (
	TestEmailAddress    = "testify@test.com"
	TestEmailAddressTwo = "testify2@test.com"
)

func TestAddEmailSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(AddEmailRoute, AddEmail(&i))

	SetupTestAddEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	log.Print(res)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestAddEmailWithoutLogin(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()
	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)
	app.Use(session.New(&i))
	app.Post(AddEmailRoute, AddEmail(&i))

	SetupTestAddEmail(t, &i, TestUsername, TestEmailAddress)

	req := AddEmailRequest(TestEmailAddress)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestAddEmailInvalidAddress(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Use(session.New(&i))

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))
	app.Post(AddEmailRoute, AddEmail(&i))

	SetupTestAddEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	// TODO: Add more test cases for possible inputs here
	req := AddEmailRequest("invalid")
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestAddEmailDBDisconnected(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Use(session.New(&i))

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))
	app.Post(AddEmailRoute, AddEmail(&i))

	SetupTestAddEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)

	i.Close()
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestAddEmailMalformedInput(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Use(session.New(&i))

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))
	app.Post(AddEmailRoute, AddEmail(&i))

	SetupTestAddEmail(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notemail", "blahblahblah")
	writer.Close()

	url := fmt.Sprintf("http://petrichormud.com%s", AddEmailRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func SetupTestAddEmail(t *testing.T, i *shared.Interfaces, u string, e string) {
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

func AddEmailRequest(e string) *http.Request {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", e)
	writer.Close()

	url := fmt.Sprintf("%s%s", shared.TestURL, AddEmailRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
