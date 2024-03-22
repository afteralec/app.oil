package test

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
)

func TestAddEmailSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddress)
	writer.Close()

	url := MakeTestURL(route.NewEmailPath())
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestAddEmailUnauthorized(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddress)
	writer.Close()

	url := MakeTestURL(route.NewEmailPath())
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestAddEmailInvalidAddress(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	// TODO: Add more test cases for possible inputs here
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", "invalid")
	writer.Close()

	url := MakeTestURL(route.NewEmailPath())
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestAddEmailFatal(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddress)
	writer.Close()

	url := MakeTestURL(route.NewEmailPath())
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	i.Close()

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestAddEmailMalformedInput(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notemail", "blahblahblah")
	writer.Close()

	url := fmt.Sprintf("%s%s", TestURL, route.NewEmailPath())
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditEmailUnauthorized(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditEmailMissingBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPut, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditEmailMalformedInput(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
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

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))

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
	i := service.NewInterfaces()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	i.Close()

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestEditEmailUnowned(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditEmailInvalidID(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath("invalid"))

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
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(sessionCookie)

	_, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddress)
	writer.Close()

	url = MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req = httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestEditEmailForbiddenUnverified(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddress)
	writer.Close()

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditEmailConflictAlreadyVerified(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	eidTwo := CreateTestEmail(t, &i, a, TestEmailAddressTwo, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	if err := i.Queries.MarkEmailVerified(context.Background(), eid); err != nil {
		t.Fatal(err)
	}

	if err := i.Queries.MarkEmailVerified(context.Background(), eidTwo); err != nil {
		t.Fatal(err)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusConflict, res.StatusCode)
}

func TestEditEmailSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	if err := i.Queries.MarkEmailVerified(context.Background(), eid); err != nil {
		t.Fatal(err)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestDeleteEmailUnauthorized(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))

	req := httptest.NewRequest(http.MethodDelete, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestDeleteEmailFatal(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	i.Close()

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestDeleteEmailUnowned(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestDeleteEmailInvalidID(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath("invalid"))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestDeleteEmailNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid+1, 10)))
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestDeleteEmailSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestVerifyEmailPageUnauthorized(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestVerifyEmailPageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestVerifyEmailPageForbiddenUnowned(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestVerifyPageExpiredToken(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	key := keys[0]
	if err := i.Redis.Expire(context.Background(), key, 0).Err(); err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(key, ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestVerifyEmailUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[0]))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestVerifyEmailBadRequestMissingToken(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.VerifyEmail)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestVerifyEmailNotFoundExpiredToken(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	key := keys[0]
	if err := i.Redis.Expire(context.Background(), key, 0).Err(); err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(key, ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestVerifyEmailSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestResendVerificationEmailUnauthorized(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	url := MakeTestURL(route.ResendEmailVerificationPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestResendVerificationEmailFatal(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ResendEmailVerificationPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	i.Close()

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestResendVerificationEmailBadRequestUnowned(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(route.ResendEmailVerificationPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestResendEmailVerificationBadRequestInvalidID(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddressTwo)
	writer.Close()

	url := MakeTestURL(route.ResendEmailVerificationPath("invalid"))

	req, err := http.NewRequest(http.MethodPost, url, body)
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

func TestResendEmailVerificationNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", TestEmailAddress)
	writer.Close()

	url := MakeTestURL(route.NewEmailPath())
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	emails := ListEmailsForPlayer(t, &i, TestUsername)
	email := emails[0]

	url = MakeTestURL(route.EmailPath(strconv.FormatInt(email.ID, 10)))
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)
	_, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	url = MakeTestURL(route.ResendEmailVerificationPath(strconv.FormatInt(email.ID, 10)))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestResendEmailVerificationSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ResendEmailVerificationPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
