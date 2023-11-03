package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/middleware/sessiondata"
	"petrichormud.com/app/internal/shared"
)

const TestAddress = "testify@test.com"

func TestAddEmailSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestAddEmail(t, &i)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Use(sessiondata.New(&i))

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))
	app.Post(AddEmailRoute, AddEmail(&i))

	body, contentType := RegisterTestFormData()

	// TODO: Extract this test url to a constant?
	url := fmt.Sprintf("http://petrichormud.com%s", RegisterRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	_, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	body, contentType = LoginTestFormData()

	url = fmt.Sprintf("http://petrichormud.com%s", LoginRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	cookies := res.Cookies()
	sessionCookie := cookies[0]

	body, contentType = AddEmailTestFormData(TestAddress)

	url = fmt.Sprintf("http://petrichormud.com%s", AddEmailRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(sessionCookie)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestAddEmailWithoutLogin(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestAddEmail(t, &i)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	app.Use(sessiondata.New(&i))

	app.Post(AddEmailRoute, AddEmail(&i))

	body, contentType := AddEmailTestFormData(TestAddress)

	url := fmt.Sprintf("http://petrichormud.com%s", AddEmailRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestAddEmailInvalidAddress(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestAddEmail(t, &i)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	// TODO: Extract all of this setup to its own functions?
	app.Use(sessiondata.New(&i))

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))
	app.Post(AddEmailRoute, AddEmail(&i))

	body, contentType := RegisterTestFormData()

	// TODO: Extract this test url to a constant?
	url := fmt.Sprintf("http://petrichormud.com%s", RegisterRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	_, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	body, contentType = LoginTestFormData()

	url = fmt.Sprintf("http://petrichormud.com%s", LoginRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	cookies := res.Cookies()
	sessionCookie := cookies[0]

	// TODO: Add more test cases for possible inputs here
	body, contentType = AddEmailTestFormData("invalid")

	url = fmt.Sprintf("http://petrichormud.com%s", AddEmailRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(sessionCookie)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestAddEmailDBDisconnected(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	SetupTestAddEmail(t, &i)

	views := html.New("../..", ".html")
	config := configs.Fiber(views)
	app := fiber.New(config)

	// TODO: Extract all of this setup to its own functions?
	app.Use(sessiondata.New(&i))

	app.Post(LoginRoute, Login(&i))
	app.Post(RegisterRoute, Register(&i))
	app.Post(AddEmailRoute, AddEmail(&i))

	body, contentType := RegisterTestFormData()

	// TODO: Extract this test url to a constant?
	url := fmt.Sprintf("http://petrichormud.com%s", RegisterRoute)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	_, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	body, contentType = LoginTestFormData()

	url = fmt.Sprintf("http://petrichormud.com%s", LoginRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	cookies := res.Cookies()
	sessionCookie := cookies[0]

	// TODO: Add more test cases for possible inputs here
	body, contentType = AddEmailTestFormData(TestAddress)

	i.Close()
	url = fmt.Sprintf("http://petrichormud.com%s", AddEmailRoute)
	req = httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(sessionCookie)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func SetupTestAddEmail(t *testing.T, i *shared.Interfaces) {
	_, err := i.Database.Exec("DELETE FROM players WHERE username = 'testify';")
	if err != nil {
		t.Fatal(err)
	}
	_, err = i.Database.Exec("DELETE FROM emails WHERE address = 'testify@test.com';")
	if err != nil {
		t.Fatal(err)
	}
}

func AddEmailTestFormData(address string) (io.Reader, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", address)
	writer.Close()
	return body, writer.FormDataContentType()
}
