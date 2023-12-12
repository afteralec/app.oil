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

	url := MakeTestURL(routes.Characters)
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

	url := MakeTestURL(routes.Characters)
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
	url := MakeTestURL(routes.Characters)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestNewCharacterApplicationUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	// TODO: Make this kind of "new" route a utility func on the routes package
	url := MakeTestURL(routes.NewCharacterApplicationPath())
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestNewCharacterApplicationSuccess(t *testing.T) {
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

	url := MakeTestURL(routes.NewCharacterApplicationPath())
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestNewCharacterApplicationFatal(t *testing.T) {
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

	url := MakeTestURL(routes.NewCharacterApplicationPath())
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestNewCharacterApplicationMaxOpen(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestCharacters(t, &i, TestUsername)
	defer SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, a, TestUsername, TestPassword)
	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]

	for i := 1; i <= shared.MaxOpenCharacterApplications; i++ {
		req := NewCharacterApplicationRequest()
		req.AddCookie(sessionCookie)
		_, err := a.Test(req)
		if err != nil {
			t.Fatal(err)
		}
	}

	url := MakeTestURL(routes.NewCharacterApplicationPath())
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestUpdateCharacterApplicationNameUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.CharacterApplicationNamePath(routes.ID))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationNamePath(strconv.FormatInt(rid+1, 10)))
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

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestUpdateCharacterApplicationNameFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10)))
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

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestUpdateCharacterApplicationNameMissingBody(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationGenderUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.CharacterApplicationGenderPath(routes.ID))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationGenderPath(strconv.FormatInt(rid+1, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10)))
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

	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(routes.ID))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid+1, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10)))
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

	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(routes.ID))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid+1, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestUpdateCharacterApplicationBackstoryUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(routes.ID))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	res, err := a.Test(req)
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid+1, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCharacterNamePageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.CharacterApplicationNamePath(routes.ID))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationNamePath(strconv.FormatInt(rid+1, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterGenderPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.CharacterApplicationGenderPath(routes.ID))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationGenderPath(strconv.FormatInt(rid+1, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10)))
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

	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(routes.ID))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid+1, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10)))
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

	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(routes.ID))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid+1, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterBackstoryPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(routes.ID))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid+1, 10)))
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

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterReviewPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.CharacterApplicationReviewPath(routes.ID))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterReviewPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationReviewPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterReviewPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationReviewPath(strconv.FormatInt(rid+1, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterReviewPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationReviewPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestSubmitCharacterApplicationUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.SubmitCharacterApplicationPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
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
		query := fmt.Sprintf("DELETE FROM character_application_content WHERE rid = %d;", req.Request.ID)
		_, err := i.Database.Exec(query)
		if err != nil {
			t.Fatal(err)
		}

		query = fmt.Sprintf("DELETE FROM character_application_content_history WHERE rid = %d;", req.Request.ID)
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

func CreateTestPlayerAndCharacterApplication(t *testing.T, i *shared.Interfaces, app *fiber.App) (int64, *http.Cookie) {
	SetupTestCharacters(t, i, TestUsername)
	CallRegister(t, app, TestUsername, TestPassword)
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
	req := NewCharacterApplicationRequest()
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
	return r.Request.ID, sessionCookie
}

func NewCharacterApplicationRequest() *http.Request {
	url := MakeTestURL(routes.NewCharacterApplicationPath())
	return httptest.NewRequest(http.MethodPost, url, nil)
}
