package tests

import (
	"bytes"
	"context"
	"database/sql"
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
	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestCharactersPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, handlers.CharactersRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharactersPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	url := fmt.Sprintf("%s%s", TestURL, handlers.CharactersRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharactersPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	i.Close()
	url := fmt.Sprintf("%s%s", TestURL, handlers.CharactersRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestNewCharacter(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, handlers.NewCharacterApplicationRoute)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestNewCharacterSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	url := fmt.Sprintf("%s%s", TestURL, handlers.NewCharacterApplicationRoute)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestNewCharacterFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	i.Close()

	url := fmt.Sprintf("%s%s", TestURL, handlers.NewCharacterApplicationRoute)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplication(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, handlers.CharacterApplicationRoute)
	req := httptest.NewRequest(http.MethodPut, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s%d", TestURL, "/character/application/", rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "test")
	writer.WriteField("gender", "NonBinary")
	writer.WriteField("sdesc", "test, testerly person")
	writer.WriteField("description", "This is a test character application.")
	writer.WriteField("backstory", "This is a tragic backtory.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s%d", TestURL, "/character/application/", rid+1)
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	url := fmt.Sprintf("%s%s%d", TestURL, "/character/application/", rid)
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s%s%d", TestURL, "/character/application/", rid)
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationName(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s", TestURL, handlers.CharacterApplicationNameRoute)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationNameSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s%s%d%s", TestURL, "/character/application/", rid, "/name")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "test")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationNameNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s%s%d%s", TestURL, "/character/application/", rid+1, "/name")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationNameFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	path := routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationNameMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	path := routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationGender(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s", TestURL, handlers.CharacterApplicationGenderRoute)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationGenderSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s%s%d%s", TestURL, "/character/application/", rid, "/gender")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("gender", "NonBinary")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationGenderNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s%s%d%s", TestURL, "/character/application/", rid+1, "/gender")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationGenderFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	path := routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationGenderMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s%s%d%s", TestURL, "/character/application/", rid, "/gender")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationShortDescription(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, routes.CharacterApplicationShortDescriptionPath(routes.ID))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationShortDescriptionSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	path := routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", "testing, testerly person")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationShortDescriptionNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	path := routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid+1, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	// TODO: Move this into a fixture
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", "testing, testerly person")
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

func TestUpdateCharacterApplicationShortDescriptionFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	path := routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", "testing, testerly person")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationShortDescriptionMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	path := routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationDescription(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, routes.CharacterApplicationDescriptionPath(routes.ID))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationDescriptionSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s%s%d%s", TestURL, "/character/application/", rid, "/description")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("description", "This is a test actor.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationDescriptionNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s%s%d%s", TestURL, "/character/application/", rid+1, "/description")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationDescriptionFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationDescriptionMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationBackstory(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(bind.New())
	app.Use(session.New(&i))

	app.Patch(handlers.CharacterApplicationBackstoryRoute, handlers.UpdateCharacterApplicationBackstory(&i))

	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s", TestURL, handlers.CharacterApplicationBackstoryRoute)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationBackstorySuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s%s%d%s", TestURL, "/character/application/", rid, "/backstory")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("backstory", "This is a tragic backstory.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationBackstoryNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s%s%d%s", TestURL, "/character/application/", rid+1, "/backstory")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationBackstoryFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationBackstoryMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCharacterNamePage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, handlers.CharacterApplicationNameRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterNamePageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s/character/application/%d/name", TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterNamePageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	url := fmt.Sprintf("%s/character/application/%d/name", TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterNamePageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	url := fmt.Sprintf("%s/character/application/%d/name", TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterGenderPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, handlers.CharacterApplicationGenderRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterGenderPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/character/application/%d/gender", TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterGenderPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/gender", TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterGenderPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/character/application/%d/gender", TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterSdescPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, routes.CharacterApplicationShortDescriptionPath(routes.ID))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterSdescPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	path := routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterSdescPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	path := routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid+1, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterShortDescriptionPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	path := routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
	url := fmt.Sprintf("%s%s", TestURL, path)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterDescriptionPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, routes.CharacterApplicationDescriptionPath(routes.ID))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterDescriptionPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/character/application/%d/description", TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterDescriptionPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/character/application/%d/description", TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterDescriptionPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/character/application/%d/description", TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterBackstoryPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := fmt.Sprintf("%s%s", TestURL, handlers.CharacterApplicationBackstoryRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterBackstoryPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/character/application/%d/backstory", TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterBackstoryPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/character/application/%d/backstory", TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterBackstoryPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CharacterApplicationRID(t, &i, a)
	i.Close()
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/character/application/%d/backstory", TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func SetupTestCharacters(t *testing.T, i *shared.Interfaces, u string) {
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}

	pid := p.ID
	reqs, err := i.Queries.ListCharacterApplicationsForPlayer(context.Background(), pid)
	if err != nil {
		t.Fatal(err)
	}

	for _, req := range reqs {
		query := fmt.Sprintf("DELETE FROM character_application_content WHERE rid = %d;", req.ID)
		_, err := i.Database.Exec(query)
		if err != nil {
			t.Fatal(err)
		}

		query = fmt.Sprintf("DELETE FROM character_application_content_history WHERE rid = %d;", req.ID)
		_, err = i.Database.Exec(query)
		if err != nil {
			t.Fatal(err)
		}
	}

	query := fmt.Sprintf("DELETE FROM requests WHERE pid = %d;", pid)
	_, err = i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}

	query = fmt.Sprintf("DELETE FROM players WHERE username = '%s';", u)
	_, err = i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func CharacterApplicationRID(t *testing.T, i *shared.Interfaces, app *fiber.App) (int64, *http.Cookie) {
	SetupTestCharacters(t, i, TestUsername)
	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
	req := NewCharacterRequest()
	req.AddCookie(sessionCookie)
	_, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}
	reqs, err := i.Queries.ListCharacterApplicationsForPlayer(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	r := reqs[0]
	return r.ID, sessionCookie
}

func NewCharacterRequest() *http.Request {
	url := fmt.Sprintf("%s%s", TestURL, handlers.NewCharacterApplicationRoute)
	return httptest.NewRequest(http.MethodPost, url, nil)
}
