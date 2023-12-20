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
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestCreateRequestCommentUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
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
	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), request.FieldName))

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
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, prid)

	if err = i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: p.ID,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), "name"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateRequestCommentInvalidText(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, prid)

	if err = i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: p.ID,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), "name"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateRequestCommentBadField(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, prid)

	if err = i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: p.ID,
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

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateRequestCommentNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, prid)

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid+1, 10), "name"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCreateRequestCommentForbiddenOwnRequest(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, prid)

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), "name"))

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
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, prid)

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), "name"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCreateRequestCommentNotReviewer(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, prid)

	CallRegister(t, a, TestUsernameThree, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameThree)
	p, err = i.Queries.GetPlayerByUsername(context.Background(), TestUsernameThree)
	if err != nil {
		t.Fatal(err)
	}
	pr, err = i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err = pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, prid)

	if err = i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: p.ID,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), "name"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCreateRequestCommentNoPermission(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}

	if err = i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: p.ID,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
	writer.Close()

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), "name"))

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.AddCookie(sessionCookie)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCreateRequestCommentSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, prid)

	if err = i.Queries.MarkRequestInReview(context.Background(), queries.MarkRequestInReviewParams{
		RPID: p.ID,
		ID:   rid,
	}); err != nil {
		t.Fatal(t)
	}

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This name is fantastic.")
	writer.Close()

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(rid, 10), "name"))

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.AddCookie(sessionCookie)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
