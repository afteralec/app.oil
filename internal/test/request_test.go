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
	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/request/definition"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
)

func TestCreateRequestUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("type", request.TypeCharacterApplication)
	writer.Close()

	url := MakeTestURL(route.Requests)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCreateRequestFatal(t *testing.T) {
	i := service.NewInterfaces()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	i.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("type", request.TypeCharacterApplication)
	writer.Close()

	url := MakeTestURL(route.Requests)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCreateRequestBadRequestMissingBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.Requests)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateRequestBadRequestInvalidType(t *testing.T) {
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
	writer.WriteField("type", "not-a-type")
	writer.Close()

	url := MakeTestURL(route.Requests)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateRequestSuccess(t *testing.T) {
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
	writer.WriteField("type", request.TypeCharacterApplication)
	writer.Close()

	url := MakeTestURL(route.Requests)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestCreateCharacterApplicationUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("type", request.TypeCharacterApplication)
	writer.Close()

	url := MakeTestURL(route.Characters)
	req := httptest.NewRequest(http.MethodPost, url, body)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCreateCharacterApplicationFatal(t *testing.T) {
	i := service.NewInterfaces()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	i.Close()

	url := MakeTestURL(route.Characters)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCreateCharacterApplicationSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.Characters)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestRequestFieldPageUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestRequestFieldPageUnowned(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestRequestFieldPageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	// TODO: Hack
	if err := i.Queries.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     rid,
		Status: request.StatusReady,
	}); err != nil {
		t.Fatal(err)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestRequestFieldPageNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid+1, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestRequestFieldPageFatal(t *testing.T) {
	i := service.NewInterfaces()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	i.Close()

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestRequestPageUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	url := MakeTestURL(route.RequestPath(rid))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestRequestFieldPageForbiddenUnowned(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(route.RequestPath(rid))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestRequestPageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestPath(rid))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestRequestPageNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestPath(rid + 1))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestRequestPageFatal(t *testing.T) {
	i := service.NewInterfaces()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	i.Close()

	url := MakeTestURL(route.RequestPath(rid))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharactersPageUnauthorized(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	url := MakeTestURL(route.Characters)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharactersPageSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.Characters)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharactersPageFatal(t *testing.T) {
	i := service.NewInterfaces()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	i.Close()

	url := MakeTestURL(route.Characters)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateRequestFieldUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

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

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid+1, definition.FieldCharacterApplicationName.Type))

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

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

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
	defer DeleteTestRequest(t, &i, rid)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateRequestFieldForbiddenUnowned(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField(definition.FieldCharacterApplicationName.Type, "Test")
	writer.Close()

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

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

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	// TODO: Update this to use a helper that calls the app's API instead of hacking it
	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusSubmitted,
	}); err != nil {
		t.Fatal(t)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField(definition.FieldCharacterApplicationName.Type, "Test")
	writer.Close()

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

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

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateRequestFieldBadRequestMalformedBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

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

func TestUpdateRequestStatusUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	url := MakeTestURL(route.RequestStatusPath(rid))

	req := httptest.NewRequest(http.MethodPost, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateRequestStatusBadRequestNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestStatusPath(rid + 1000))

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateRequestStatusFatal(t *testing.T) {
	i := service.NewInterfaces()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestStatusPath(rid))

	i.Close()

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateRequestFieldStatusUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	prid := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestPlayerPermission(t, &i, prid)
	defer DeleteTestRequest(t, &i, rid)
	if err := i.Queries.UpdateRequestReviewer(context.Background(), query.UpdateRequestReviewerParams{
		ID:   rid,
		RPID: pid,
	}); err != nil {
		t.Fatal(err)
	}
	if err := i.Queries.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     rid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(err)
	}

	url := MakeTestURL(route.RequestFieldStatusPath(rid, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodPost, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateRequestFieldStatusBadRequestNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	prid := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestPlayerPermission(t, &i, prid)
	defer DeleteTestRequest(t, &i, rid)
	if err := i.Queries.UpdateRequestReviewer(context.Background(), query.UpdateRequestReviewerParams{
		ID:   rid,
		RPID: pid,
	}); err != nil {
		t.Fatal(err)
	}
	if err := i.Queries.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     rid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(err)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldStatusPath(rid+1000, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateRequestFieldStatusFatal(t *testing.T) {
	i := service.NewInterfaces()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	prid := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	if err := i.Queries.UpdateRequestReviewer(context.Background(), query.UpdateRequestReviewerParams{
		ID:   rid,
		RPID: pid,
	}); err != nil {
		t.Fatal(err)
	}
	if err := i.Queries.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     rid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(err)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.RequestFieldStatusPath(rid, definition.FieldCharacterApplicationName.Type))

	i.Close()

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestPlayerPermission(t, &i, prid)
	defer DeleteTestRequest(t, &i, rid)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCreateRequestChangeRequestUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)

	url := MakeTestURL(route.RequestChangeRequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is terrible.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCreateRequestChangeRequestBadRequestMissingBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	url := MakeTestURL(route.RequestChangeRequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateRequestChangeRequestBadRequestInvalidText(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	url := MakeTestURL(route.RequestChangeRequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateRequestChangeRequestBadRequestInvalidField(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	url := MakeTestURL(route.RequestChangeRequestFieldPath(rid, "notafield"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This whatever is also fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateRequestChangeRequestNotFoundNoRequest(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	url := MakeTestURL(route.RequestChangeRequestFieldPath(rid+1000, definition.FieldCharacterApplicationName.Type))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCreateRequestChangeRequestForbiddenNotInReview(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(route.RequestChangeRequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCreateRequestChangeRequestForbiddenNotReviewer(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionID)

	pid = CreateTestPlayer(t, &i, a, TestUsernameThree, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameThree)
	permissionID = CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionID)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(route.RequestChangeRequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCreateRequestChangeRequestForbiddenNoPermission(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is fantastic.")
	writer.Close()

	url := MakeTestURL(route.RequestChangeRequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.AddCookie(sessionCookie)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCreateRequestChangeRequestSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionID := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionID)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is fantastic.")
	writer.Close()

	url := MakeTestURL(route.RequestChangeRequestFieldPath(rid, definition.FieldCharacterApplicationName.Type))

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestDeleteRequestChangeRequestUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id))

	req := httptest.NewRequest(http.MethodDelete, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestDeleteRequestChangeRequestNotFoundNoChangeRequest(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id + 1000))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestDeleteRequestChangeRequestNotFoundNoRequest(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id + 1000))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestDeleteRequestChangeRequestSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionID := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionID)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestEditRequestChangeRequestUnauthorizedNotLoggedIn(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is not fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditRequestChangeRequestBadRequestMissingBody(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id + 1000))

	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditRequestChangeRequestBadRequestInvalidText(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id + 1000))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "tooshort")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditRequestChangeRequestNotFoundNoChangeRequest(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id + 1000))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is not fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestEditRequestChangeRequestNotFoundNoRequest(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionId)
	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id + 1000))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is not fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestEditRequestChangeRequestSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber(i.Templates))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestRequest(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionID := CreateTestPlayerPermission(t, &i, pid, player.PermissionReviewCharacterApplications.Name)
	defer DeleteTestPlayerPermission(t, &i, permissionID)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	if err := request.UpdateStatus(i.Queries, request.UpdateStatusParams{
		RID:    rid,
		PID:    pid,
		Status: request.StatusInReview,
	}); err != nil {
		t.Fatal(t)
	}

	id := CreateTestRequestChangeRequest(CreateTestRequestChangeRequestParams{
		T:        t,
		I:        &i,
		A:        a,
		Username: TestUsernameTwo,
		Password: TestPassword,
		Field:    definition.FieldCharacterApplicationName.Type,
		RID:      rid,
	})

	url := MakeTestURL(route.RequestChangeRequestPath(id))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "This name is not fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)

	field, err := i.Queries.GetRequestFieldByType(context.Background(), query.GetRequestFieldByTypeParams{
		RID:  rid,
		Type: definition.FieldCharacterApplicationName.Type,
	})
	if err != nil {
		t.Fatal(err)
	}
	change, err := i.Queries.GetOpenRequestChangeRequestForRequestField(context.Background(), field.ID)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, change.Text, "This name is not fantastic.")
}
