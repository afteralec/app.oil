package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

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

func TestSubmitCharacterApplicationUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, "testify2", TestPassword)
	res := CallLogin(t, a, "testify2", TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.SubmitCharacterApplicationPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestSubmitCharacterApplicationSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	// TODO: This is a hack. Rework this to use valid content and the existing handlers
	if err := i.Queries.MarkRequestReady(context.Background(), rid); err != nil {
		t.Fatal(err)
	}
	url := MakeTestURL(routes.SubmitCharacterApplicationPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestSubmitCharacterApplicationNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.SubmitCharacterApplicationPath(strconv.FormatInt(rid+1, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestSubmitCharacterApplicationFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.SubmitCharacterApplicationPath(strconv.FormatInt(rid, 10)))
	i.Close()
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = shared.SetupInterfaces()
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestSubmitCharacterApplicationVersionZero(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	if err := i.Queries.MarkRequestReady(context.Background(), rid); err != nil {
		t.Fatal(err)
	}
	url := MakeTestURL(routes.SubmitCharacterApplicationPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)

	r, err := i.Queries.GetRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, int32(1), r.VID)
}

func TestSubmitCharacterApplicationVersionOne(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	if err := i.Queries.MarkRequestReady(context.Background(), rid); err != nil {
		t.Fatal(err)
	}
	if err := i.Queries.IncrementRequestVersion(context.Background(), rid); err != nil {
		t.Fatal(err)
	}
	url := MakeTestURL(routes.SubmitCharacterApplicationPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)

	r, err := i.Queries.GetRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, int32(1), r.VID)
}

func TestCharacterApplicationSubmittedSuccessPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationSubmittedSuccessPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterApplicationSubmittedSuccessPageUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	// TODO: This is kind of a hack; make this a better representation
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}
	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.CharacterApplicationSubmittedSuccessPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCharacterApplicationSubmittedSuccessPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	// TODO: This is kind of a hack; make this a better representation
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}
	url := MakeTestURL(routes.CharacterApplicationSubmittedSuccessPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterApplicationSubmittedSuccessPageWrongStatus(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationSubmittedSuccessPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCharacterApplicationSubmittedSuccessPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationSubmittedSuccessPath(strconv.FormatInt(rid+1, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterApplicationSubmittedSuccessPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationSubmittedSuccessPath(strconv.FormatInt(rid+1, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = shared.SetupInterfaces()
	defer i.Close()
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCharacterApplicationSubmittedPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationSubmittedPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterApplicationSubmittedPageUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	// TODO: This is kind of a hack; make this a better representation
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}
	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.CharacterApplicationSubmittedPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCharacterApplicationSubmittedPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	// TODO: This is kind of a hack; make this a better representation
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}
	url := MakeTestURL(routes.CharacterApplicationSubmittedPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCharacterApplicationSubmittedPageWrongStatus(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationSubmittedPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestCharacterApplicationSubmittedPageNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationSubmittedPath(strconv.FormatInt(rid+1, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCharacterApplicationSubmittedPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	i.Close()
	url := MakeTestURL(routes.CharacterApplicationSubmittedPath(strconv.FormatInt(rid+1, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = shared.SetupInterfaces()
	defer i.Close()
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}
