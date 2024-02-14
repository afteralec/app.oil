package test

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
	playerpermission "petrichormud.com/app/internal/player/permission"
	"petrichormud.com/app/internal/route"
)

func TestActorImageNameReservedConflict(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", TestActorImage.Name)
	writer.Close()

	url := MakeTestURL(route.ActorImageReserved)

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusConflict, res.StatusCode)
}

func TestActorImageNameReservedFatal(t *testing.T) {
	i := interfaces.SetupShared()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)

	i.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", TestActorImage.Name)
	writer.Close()

	url := MakeTestURL(route.ActorImageReserved)

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = interfaces.SetupShared()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestActorImageNameReservedOK(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	config := config.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", TestActorImage.Name)
	writer.Close()

	url := MakeTestURL(route.ActorImageReserved)

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestActorImagesPageUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(route.ActorImages)

	req := httptest.NewRequest(http.MethodGet, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestActorImagesPageForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImages)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestActorImagesPageSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerViewAllActorImagesName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImages)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestNewActorImageUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", TestActorImage.Name)
	writer.Close()

	url := MakeTestURL(route.ActorImages)

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestActorImageByName(t, &i, TestActorImage.Name)

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestNewActorImageForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImages)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", TestActorImage.Name)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestActorImageByName(t, &i, TestActorImage.Name)

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestNewActorImageBadRequestMissingBody(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImages)

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestActorImageByName(t, &i, TestActorImage.Name)

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestNewActorImageBadRequestInvalidName(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImages)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "NotAGoodName")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestActorImageByName(t, &i, TestActorImage.Name)

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestNewActorImageSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImages)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", TestActorImage.Name)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestActorImageByName(t, &i, TestActorImage.Name)

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestEditActorImagePageUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	url := MakeTestURL(route.EditActorImagePath(aiid))

	req := httptest.NewRequest(http.MethodGet, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditActorImagePageForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EditActorImagePath(aiid))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditActorImagePageSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EditActorImagePath(aiid))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestActorImagePageUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	url := MakeTestURL(route.ActorImagePath(aiid))

	req := httptest.NewRequest(http.MethodGet, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestActorImagePageForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImagePath(aiid))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestActorImagePageSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerViewAllActorImagesName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImagePath(aiid))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestEditActorImageShortDescriptionUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	url := MakeTestURL(route.ActorImageShortDescriptionPath(aiid))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s, with changes", TestActorImage.ShortDescription)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", sb.String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditActorImageShortDescriptionForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageShortDescriptionPath(aiid))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s, with changes", TestActorImage.ShortDescription)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", sb.String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditActorImageShortDescriptionNotFound(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageShortDescriptionPath(aiid + 1000))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s, with changes", TestActorImage.ShortDescription)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", sb.String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestEditActorImageShortDescriptionBadRequestInvalid(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageShortDescriptionPath(aiid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", "This is an invalid short description - 1234 tell_me that you love me mo4r.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditActorImageShortDescriptionConflictSameAs(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageShortDescriptionPath(aiid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", TestActorImage.ShortDescription)
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusConflict, res.StatusCode)
}

func TestEditActorImageShortDescriptionSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageShortDescriptionPath(aiid))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s, with changes", TestActorImage.ShortDescription)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("sdesc", sb.String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestEditActorImageDescriptionUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	url := MakeTestURL(route.ActorImageDescriptionPath(aiid))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s It is a test actor.", TestActorImage.Description)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("desc", sb.String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditActorImageDescriptionForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageDescriptionPath(aiid))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s It is a test actor.", TestActorImage.Description)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("desc", sb.String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditActorImageDescriptionNotFound(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageDescriptionPath(aiid + 1000))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s It is a test actor.", TestActorImage.Description)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("desc", sb.String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestEditActorImageDescriptionBadRequestInvalid(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageDescriptionPath(aiid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("desc", "This is an invalid description - 1234 tell_me that you love me mo4r.")
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestEditActorImageDescriptionConflictSameAs(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageDescriptionPath(aiid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("desc", TestActorImage.Description)
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusConflict, res.StatusCode)
}

func TestEditActorImageDescriptionSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, playerpermission.PlayerCreateActorImageName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	aiid := CreateTestActorImage(t, &i, TestActorImage)
	defer DeleteTestActorImage(t, &i, aiid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.ActorImageDescriptionPath(aiid))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s It is a test actor.", TestActorImage.Description)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("desc", sb.String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
