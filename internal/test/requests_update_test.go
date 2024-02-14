package test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
)

func TestUpdateRequestFieldUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	url := MakeTestURL(route.RequestFieldPath(rid, request.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("value", "Test")
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateRequestFieldBadRequestNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid+1, request.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("value", "Test")
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateRequestFieldFatal(t *testing.T) {
	i := service.NewInterfaces()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid, request.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("value", "Test")
	writer.Close()

	i.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateRequestFieldForbiddenUnowned(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField(request.FieldName, "Test")
	writer.Close()

	url := MakeTestURL(route.RequestFieldPath(rid, request.FieldName))

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestUpdateRequestFieldForbiddenNotEditable(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	// TODO: Update this to use a helper that calls the app's API instead of hacking it
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField(request.FieldName, "Test")
	writer.Close()

	url := MakeTestURL(route.RequestFieldPath(rid, request.FieldName))

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestUpdateRequestFieldBadRequestMissingBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid, request.FieldName))

	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateRequestBadRequestFieldMalformedBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid, request.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("notavalue", "Test")
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}
