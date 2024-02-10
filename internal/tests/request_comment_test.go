package tests

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/requests"
	"petrichormud.com/app/internal/routes"
)

func TestCreateRequestCommentUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerReviewCharacterApplicationsName)
	defer DeleteTestPlayerPermission(t, &i, permissionId)

	// TODO: Make a map of valid Character Application fields
	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), requests.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCreateRequestCommentMissingBody(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerReviewCharacterApplicationsName)
	defer DeleteTestPlayerPermission(t, &i, permissionId)

	if err := i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: pid,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), requests.FieldName))

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateRequestCommentInvalidText(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerReviewCharacterApplicationsName)
	defer DeleteTestPlayerPermission(t, &i, permissionId)

	if err := i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: pid,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), requests.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "")
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

func TestCreateRequestCommentBadField(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerReviewCharacterApplicationsName)
	defer DeleteTestPlayerPermission(t, &i, permissionId)

	if err := i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: pid,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), "notafield"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This whatever is also fantastic.")
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

func TestCreateRequestCommentNotFound(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerReviewCharacterApplicationsName)
	defer DeleteTestPlayerPermission(t, &i, permissionId)

	if err := i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: pid,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid+1, 10), requests.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
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

func TestCreateRequestCommentForbiddenOwnRequest(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	if err := i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: pid,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), requests.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
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

func TestCreateRequestCommentNotInReview(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionId := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerReviewCharacterApplicationsName)
	defer DeleteTestPlayerPermission(t, &i, permissionId)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), requests.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
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

func TestCreateRequestCommentNotReviewer(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerReviewCharacterApplicationsName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)

	pid = CreateTestPlayer(t, &i, a, TestUsernameThree, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameThree)
	permissionID = CreateTestPlayerPermission(t, &i, pid, permissions.PlayerReviewCharacterApplicationsName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)

	if err := i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: pid,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), requests.FieldName))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
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

func TestCreateRequestCommentNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	if err := i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: pid,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
	writer.Close()

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), requests.FieldName))

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.AddCookie(sessionCookie)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCreateRequestCommentSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	rid := CreateTestCharacterApplication(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	permissionID := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerReviewCharacterApplicationsName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)

	if err := i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: pid,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
	writer.Close()

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), requests.FieldName))

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
