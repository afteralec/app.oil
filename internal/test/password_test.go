package test

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/player/password"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
)

func TestRecoverPasswordPageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.RecoverPassword)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestRecoverPasswordSuccessPageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	kid, err := password.SetupRecoverySuccess(&i, TestEmailAddress)
	if err != nil {
		t.Fatal(err)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s?t=%s", route.RecoverPasswordSuccess, kid)
	url := MakeTestURL(sb.String())
	req := httptest.NewRequest(http.MethodGet, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestRecoverPasswordBadRequestMissingBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.RecoverPassword)
	req := httptest.NewRequest(http.MethodPost, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestRecoverPasswordBadRequestMalformedBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notemail", "notanemail")
	writer.Close()

	url := MakeTestURL(route.RecoverPassword)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestRecoverPasswordSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	if err := i.Queries.MarkEmailVerified(context.Background(), eid); err != nil {
		t.Fatal(err)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.WriteField("email", TestEmailAddress)
	writer.Close()

	url := MakeTestURL(route.RecoverPassword)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestResetPasswordPage(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	id := uuid.NewString()
	url := MakeTestURL(fmt.Sprintf("%s?t=%s", route.ResetPassword, id))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestResetPasswordSuccessPageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.ResetPasswordSuccess)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestResetPasswordMissingBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.ResetPassword)
	req := httptest.NewRequest(http.MethodPost, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestResetPasswordMalformedBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notusername", "notausername")
	writer.Close()

	url := MakeTestURL(route.ResetPassword)

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestResetPasswordSuccess(t *testing.T) {
	t.Skip()
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	if err := i.Queries.MarkEmailVerified(context.Background(), eid); err != nil {
		t.Fatal(err)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.WriteField("password", TestPassword)
	writer.WriteField("confirm", TestPassword)
	writer.Close()

	url := MakeTestURL(route.RecoverPassword)

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestChangePasswordMissingBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.PlayerPasswordPath(pid))
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestChangePasswordMalformedBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("password", TestPassword)
	writer.WriteField("confirmPassword", TestPassword)
	writer.Close()

	url := MakeTestURL(route.PlayerPasswordPath(pid))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestChangePasswordNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s%s", TestPassword, "!!")
	newPassword := sb.String()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("current", "notTestPassword")
	writer.WriteField("password", newPassword)
	writer.WriteField("confirm", newPassword)
	writer.Close()

	url := MakeTestURL(route.PlayerPasswordPath(pid))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestChangePasswordWrongCurrentPassword(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s%s", TestPassword, "!!")
	newPassword := sb.String()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("current", "notTestPassword")
	writer.WriteField("password", newPassword)
	writer.WriteField("confirm", newPassword)
	writer.Close()

	url := MakeTestURL(route.PlayerPasswordPath(pid))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestChangePasswordInvalidPassword(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("current", TestPassword)
	writer.WriteField("password", "bad")
	writer.WriteField("confirm", "bad")
	writer.Close()

	url := MakeTestURL(route.PlayerPasswordPath(pid))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestChangePasswordInvalidConfirmPassword(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s%s", TestPassword, "!!")
	newPassword := sb.String()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("current", TestPassword)
	writer.WriteField("password", newPassword)
	writer.WriteField("confirm", "notnewpassword")
	writer.Close()

	url := MakeTestURL(route.PlayerPasswordPath(pid))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestChangePasswordSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s%s", TestPassword, "!!")
	newPassword := sb.String()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("current", TestPassword)
	writer.WriteField("password", newPassword)
	writer.WriteField("confirm", newPassword)
	writer.Close()

	url := MakeTestURL(route.PlayerPasswordPath(pid))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
