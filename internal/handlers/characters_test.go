package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/shared"
)

func TestCharactersPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Get(CharactersRoute, CharactersPage(&i))

	url := fmt.Sprintf("%s%s", shared.TestURL, CharactersRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharactersPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Get(CharactersRoute, CharactersPage(&i))

	SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	url := fmt.Sprintf("%s%s", shared.TestURL, CharactersRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharactersPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Get(CharactersRoute, CharactersPage(&i))

	SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	i.Close()
	url := fmt.Sprintf("%s%s", shared.TestURL, CharactersRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestNewCharacter(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(NewCharacterRoute, NewCharacterApplication(&i))

	url := fmt.Sprintf("%s%s", shared.TestURL, NewCharacterRoute)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestNewCharacterSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))

	SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	url := fmt.Sprintf("%s%s", shared.TestURL, NewCharacterRoute)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestNewCharacterFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))

	SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	i.Close()

	url := fmt.Sprintf("%s%s", shared.TestURL, NewCharacterRoute)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplication(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(bind.New())
	app.Use(session.New(&i))

	app.Put(CharacterRoute, UpdateCharacterApplication(&i))

	url := fmt.Sprintf("%s%s", shared.TestURL, CharacterRoute)
	req := httptest.NewRequest(http.MethodPut, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterRoute, UpdateCharacterApplication(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d", shared.TestURL, "/characters/", rid)
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
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterRoute, UpdateCharacterApplication(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d", shared.TestURL, "/characters/", rid+1)
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterRoute, UpdateCharacterApplication(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	url := fmt.Sprintf("%s%s%d", shared.TestURL, "/characters/", rid)
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterRoute, UpdateCharacterApplication(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d", shared.TestURL, "/characters/", rid)
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationName(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(bind.New())
	app.Use(session.New(&i))

	app.Patch(CharacterApplicationNameRoute, UpdateCharacterApplicationName(&i))

	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s", shared.TestURL, CharacterApplicationNameRoute)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationNameSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationNameRoute, UpdateCharacterApplicationName(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/name")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "test")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationNameNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationNameRoute, UpdateCharacterApplicationName(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid+1, "/name")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationNameFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationNameRoute, UpdateCharacterApplicationName(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/name")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationNameMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationNameRoute, UpdateCharacterApplicationName(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/name")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationGender(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(bind.New())
	app.Use(session.New(&i))

	app.Patch(CharacterApplicationGenderRoute, UpdateCharacterApplicationGender(&i))

	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s", shared.TestURL, CharacterApplicationGenderRoute)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationGenderSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationGenderRoute, UpdateCharacterApplicationGender(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/gender")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("gender", "NonBinary")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationGenderNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationGenderRoute, UpdateCharacterApplicationGender(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid+1, "/gender")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationGenderFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationGenderRoute, UpdateCharacterApplicationGender(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/gender")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationGenderMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationGenderRoute, UpdateCharacterApplicationGender(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/gender")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationSdesc(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(bind.New())
	app.Use(session.New(&i))

	app.Patch(CharacterApplicationSdescRoute, UpdateCharacterApplicationSdesc(&i))

	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s", shared.TestURL, CharacterApplicationSdescRoute)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationSdescSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationSdescRoute, UpdateCharacterApplicationSdesc(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/sdesc")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", "testing, testerly person")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationSdescNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationSdescRoute, UpdateCharacterApplicationSdesc(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid+1, "/sdesc")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationSdescFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationSdescRoute, UpdateCharacterApplicationSdesc(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/sdesc")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationSdescMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationSdescRoute, UpdateCharacterApplicationSdesc(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/sdesc")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationDescription(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(bind.New())
	app.Use(session.New(&i))

	app.Patch(CharacterApplicationDescriptionRoute, UpdateCharacterApplicationDescription(&i))

	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s", shared.TestURL, CharacterApplicationDescriptionRoute)
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestUpdateCharacterApplicationDescriptionSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationDescriptionRoute, UpdateCharacterApplicationDescription(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/description")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("description", "This is a test actor.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationDescriptionNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationDescriptionRoute, UpdateCharacterApplicationDescription(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid+1, "/description")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationDescriptionFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationDescriptionRoute, UpdateCharacterApplicationDescription(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/description")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationDescriptionMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationDescriptionRoute, UpdateCharacterApplicationDescription(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/description")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
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

	app.Patch(CharacterApplicationBackstoryRoute, UpdateCharacterApplicationBackstory(&i))

	// TODO: Put this in a generator
	url := fmt.Sprintf("%s%s", shared.TestURL, CharacterApplicationBackstoryRoute)
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
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationBackstoryRoute, UpdateCharacterApplicationBackstory(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/backstory")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("backstory", "This is a tragic backstory.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestUpdateCharacterApplicationBackstoryNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Patch(CharacterApplicationBackstoryRoute, UpdateCharacterApplicationBackstory(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid+1, "/backstory")
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationBackstoryFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationBackstoryRoute, UpdateCharacterApplicationBackstory(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/backstory")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationBackstoryMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Put(CharacterApplicationBackstoryRoute, UpdateCharacterApplicationBackstory(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s%s%d%s", shared.TestURL, "/characters/", rid, "/backstory")
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCharacterNamePage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Get(NewCharacterNameRoute, CharacterNamePage(&i))

	url := fmt.Sprintf("%s%s", shared.TestURL, NewCharacterNameRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterNamePageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterNameRoute, CharacterNamePage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s/characters/new/%d/name", shared.TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterNamePageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterNameRoute, CharacterNamePage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	url := fmt.Sprintf("%s/characters/new/%d/name", shared.TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterNamePageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterNameRoute, CharacterNamePage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	url := fmt.Sprintf("%s/characters/new/%d/name", shared.TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterGenderPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Get(NewCharacterGenderRoute, CharacterGenderPage(&i))

	url := fmt.Sprintf("%s%s", shared.TestURL, NewCharacterGenderRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterGenderPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterGenderRoute, CharacterGenderPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/gender", shared.TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterGenderPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterGenderRoute, CharacterGenderPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/gender", shared.TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterGenderPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterGenderRoute, CharacterGenderPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/gender", shared.TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterSdescPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Get(NewCharacterSdescRoute, CharacterSdescPage(&i))

	url := fmt.Sprintf("%s%s", shared.TestURL, NewCharacterSdescRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterSdescPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterSdescRoute, CharacterSdescPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/sdesc", shared.TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterSdescPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterSdescRoute, CharacterSdescPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/sdesc", shared.TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterSdescPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterSdescRoute, CharacterSdescPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/sdesc", shared.TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterDescriptionPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Get(NewCharacterDescriptionRoute, CharacterDescriptionPage(&i))

	url := fmt.Sprintf("%s%s", shared.TestURL, NewCharacterDescriptionRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterDescriptionPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterDescriptionRoute, CharacterDescriptionPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/description", shared.TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterDescriptionPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterDescriptionRoute, CharacterDescriptionPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/description", shared.TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterDescriptionPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterDescriptionRoute, CharacterDescriptionPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/description", shared.TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterBackstoryPage(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))

	app.Get(NewCharacterBackstoryRoute, CharacterBackstoryPage(&i))

	url := fmt.Sprintf("%s%s", shared.TestURL, NewCharacterBackstoryRoute)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterBackstoryPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterBackstoryRoute, CharacterBackstoryPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/backstory", shared.TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterBackstoryPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterBackstoryRoute, CharacterBackstoryPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/backstory", shared.TestURL, rid+1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterBackstoryPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	app := fiber.New(configs.Fiber(views))
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Post(RegisterRoute, Register(&i))
	app.Post(LoginRoute, Login(&i))
	app.Post(NewCharacterRoute, NewCharacterApplication(&i))
	app.Get(NewCharacterBackstoryRoute, CharacterBackstoryPage(&i))

	rid, sessionCookie := CharacterApplicationRID(t, &i, app)
	i.Close()
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/characters/new/%d/backstory", shared.TestURL, rid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := app.Test(req)
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
	url := fmt.Sprintf("%s%s", shared.TestURL, NewCharacterRoute)
	return httptest.NewRequest(http.MethodPost, url, nil)
}
