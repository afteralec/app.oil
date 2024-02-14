package test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
)

func TestDeleteEmailUnauthorized(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))

	req := httptest.NewRequest(http.MethodDelete, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestDeleteEmailFatal(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	i.Close()

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = service.NewInterfaces()
	defer i.Close()
	defer DeleteTestPlayer(t, &i, TestUsername)

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestDeleteEmailUnowned(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestDeleteEmailInvalidID(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath("invalid"))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestDeleteEmailNotFound(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid+1, 10)))
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestDeleteEmailSuccess(t *testing.T) {
	i := service.NewInterfaces()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	eid := CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.EmailPath(strconv.FormatInt(eid, 10)))
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
