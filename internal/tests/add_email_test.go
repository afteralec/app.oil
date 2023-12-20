package tests

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

const (
	TestEmailAddress    = "testify@test.com"
	TestEmailAddressTwo = "testify2@test.com"
)

func TestAddEmailSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestAddEmailUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	config := configs.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	req := AddEmailRequest(TestEmailAddress)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestAddEmailInvalidAddress(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	config := configs.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	// TODO: Add more test cases for possible inputs here
	req := AddEmailRequest("invalid")
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestAddEmailFatal(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	config := configs.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)

	i.Close()

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = shared.SetupInterfaces()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestAddEmailMalformedInput(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	config := configs.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notemail", "blahblahblah")
	writer.Close()

	url := fmt.Sprintf("%s%s", TestURL, routes.NewEmailPath())
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func SetupTestAddEmail(t *testing.T, i *shared.Interfaces, u string, e string) {
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

func AddEmailRequest(e string) *http.Request {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", e)
	writer.Close()

	url := MakeTestURL(routes.NewEmailPath())
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
