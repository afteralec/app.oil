package tests

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/route"
)

func TestChangePasswordMissingBody(t *testing.T) {
	i := interfaces.SetupShared()
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
	i := interfaces.SetupShared()
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
	i := interfaces.SetupShared()
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
	i := interfaces.SetupShared()
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
	i := interfaces.SetupShared()
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
	i := interfaces.SetupShared()
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
	i := interfaces.SetupShared()
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
