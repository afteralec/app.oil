package tests

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
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestCharactersPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
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

	a := fiber.New(configs.Fiber())
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

	a := fiber.New(configs.Fiber())
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
	url := MakeTestURL(routes.Characters)
	return httptest.NewRequest(http.MethodPost, url, nil)
}

func MakeTestCharacterApplicationNameBody() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "test")
	writer.Close()
	return body, writer
}

func MakeTestCharacterApplicationNameBodyInvalid() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "tes")
	writer.Close()
	return body, writer
}

func MakeTestCharacterApplicationGenderBody() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("gender", character.GenderNonBinary)
	writer.Close()
	return body, writer
}

func MakeTestCharacterApplicationGenderBodyInvalid() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("gender", "NotAGender")
	writer.Close()
	return body, writer
}

func MakeTestCharacterApplicationShortDescriptionBody() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", "test, testing, person")
	writer.Close()
	return body, writer
}

func MakeTestCharacterApplicationShortDescriptionBodyInvalid() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", "test")
	writer.Close()
	return body, writer
}

func MakeTestCharacterApplicationDescriptionBody() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	description := ""
	for len(description) < shared.MinCharacterDescriptionLength {
		description = description + "This is a test actor."
	}
	writer.WriteField("description", description)
	writer.Close()
	return body, writer
}

func MakeTestCharacterApplicationDescriptionBodyInvalid() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("description", "test")
	writer.Close()
	return body, writer
}

func MakeTestCharacterApplicationBackstoryBody() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	backstory := ""
	for len(backstory) < shared.MinCharacterBackstoryLength {
		backstory = backstory + "This is a tragic backstory."
	}
	writer.WriteField("backstory", backstory)
	writer.Close()
	return body, writer
}

func MakeTestCharacterApplicationBackstoryBodyInvalid() (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("backstory", "test")
	writer.Close()
	return body, writer
}
