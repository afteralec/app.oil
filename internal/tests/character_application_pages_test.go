package tests

import (
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

func TestCharacterNamePageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterNamePageUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
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

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterGenderPageUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsernameTwo, TestPassword)
	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
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

func TestCharacterSdescPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterShortDescriptionPageUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCharacterShortDescriptionPageSuccess(t *testing.T) {
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

func TestCharacterDescriptionPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterDescriptionPageUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
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

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterBackstoryPageUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsernameTwo, TestPassword)
	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
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

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	url := MakeTestURL(routes.CharacterApplicationReviewPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCharacterReviewPageUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsernameTwo, TestPassword)
	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.CharacterApplicationReviewPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
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
