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

func TestResendUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestResend(t, &i, TestUsername, TestEmailAddress)

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

	url := MakeTestURL(routes.ResendEmailVerificationPath(strconv.FormatInt(email.ID, 10)))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestResendDBError(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestResend(t, &i, TestUsername, TestEmailAddress)

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

	url := MakeTestURL(routes.ResendEmailVerificationPath(strconv.FormatInt(email.ID, 10)))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	i.Close()
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestResendUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestResend(t, &i, TestUsername, TestEmailAddress)
	defer SetupTestResend(t, &i, TestUsername, TestEmailAddress)
	SetupTestResend(t, &i, TestUsernameTwo, TestEmailAddressTwo)
	defer SetupTestResend(t, &i, TestUsernameTwo, TestEmailAddressTwo)

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

	url := MakeTestURL(routes.ResendEmailVerificationPath(strconv.FormatInt(email.ID, 10)))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestResendInvalidID(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestResend(t, &i, TestUsername, TestEmailAddress)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	url := MakeTestURL(routes.ResendEmailVerificationPath("invalid"))
	req, err := http.NewRequest(http.MethodPost, url, body)
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

func TestResendNonexistantEmail(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestResend(t, &i, TestUsername, TestEmailAddress)

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

	url := MakeTestURL(routes.EmailPath(strconv.FormatInt(email.ID, 10)))
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)
	_, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	url = MakeTestURL(routes.ResendEmailVerificationPath(strconv.FormatInt(email.ID, 10)))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestEditEmailVerified(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
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

	emails := ListEmailsForPlayer(t, &i, TestUsername)
	email := emails[0]
	_, err = i.Queries.MarkEmailVerified(context.Background(), email.ID)
	if err != nil {
		t.Fatal(err)
	}

	url := MakeTestURL(routes.ResendEmailVerificationPath(strconv.FormatInt(email.ID, 10)))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusConflict, res.StatusCode)
}

func TestResendSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
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

	emails := ListEmailsForPlayer(t, &i, TestUsername)
	email := emails[0]

	url := MakeTestURL(routes.ResendEmailVerificationPath(strconv.FormatInt(email.ID, 10)))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func SetupTestResend(t *testing.T, i *shared.Interfaces, u string, e string) {
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
