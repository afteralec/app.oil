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
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestEditEmailUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	req := EditEmailRequest(eid, TestEmailAddressTwo)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditEmailMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	url := MakeTestURL(routes.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPut, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditEmailMalformedInput(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notemail", "malformed")
	writer.Close()

	url := MakeTestURL(routes.EmailPath(strconv.FormatInt(eid, 10)))

	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditEmailFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	req := EditEmailRequest(eid, TestEmailAddressTwo)
	req.AddCookie(sessionCookie)

	i.Close()

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = shared.SetupInterfaces()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestEditEmailUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	req := EditEmailRequest(eid, TestEmailAddressTwo)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditEmailInvalidID(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.EmailPath("invalid"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditEmailNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.EmailPath(strconv.FormatInt(eid, 10)))

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(sessionCookie)

	_, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	req = EditEmailRequest(eid, TestEmailAddress)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestEditEmailForbiddenUnverified(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	req := EditEmailRequest(eid, TestEmailAddress)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditEmailSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	_, err := i.Queries.MarkEmailVerified(context.Background(), eid)
	if err != nil {
		t.Fatal(err)
	}

	req := EditEmailRequest(eid, TestEmailAddress)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
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
