package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestVerifyPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.VerifyEmail)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestVerifyPageUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	if err := i.Redis.FlushAll(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]
	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	rv, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	rvParts := strings.Split(rv[0], ":")

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	res = CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie = res.Cookies()[0]
	url := MakeTestURL(routes.VerifyEmailWithToken(rvParts[1]))
	req = httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestVerifyPageExpiredToken(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	if err := i.Redis.FlushAll(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]
	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	key := keys[0]
	if err := i.Redis.Expire(context.Background(), key, 0).Err(); err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(key, ":")

	url := MakeTestURL(routes.VerifyEmailWithToken(keyParts[1]))
	req = httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestVerifyEmailUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	if err := i.Redis.FlushAll(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]
	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	rv, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	rvParts := strings.Split(rv[0], ":")

	url := MakeTestURL(routes.VerifyEmailWithToken(rvParts[0]))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestVerifyEmailNoToken(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	if err := i.Redis.FlushAll(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]
	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	url := MakeTestURL(routes.VerifyEmail)
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestVerifyEmailExpiredToken(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	if err := i.Redis.FlushAll(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]
	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	rv, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	key := rv[0]
	if err := i.Redis.Expire(context.Background(), key, 0).Err(); err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(key, ":")

	url := MakeTestURL(routes.VerifyEmailWithToken(keyParts[1]))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestVerifyEmailSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.New("../..", ".html")
	a := fiber.New(configs.Fiber(views))
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	if err := i.Redis.FlushAll(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	DeleteTestPlayer(t, &i, TestUsername)
	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]
	req := AddEmailRequest(TestEmailAddress)
	req.AddCookie(sessionCookie)
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	rv, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	key := rv[0]
	keyParts := strings.Split(key, ":")

	url := MakeTestURL(routes.VerifyEmailWithToken(keyParts[1]))
	req = httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func SetupTestVerify(t *testing.T, i *shared.Interfaces, u string, e string) {
	query := fmt.Sprintf("DELETE FROM players WHERE username = '%s';", u)
	_, err := i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
	query = fmt.Sprintf("DELETE FROM emails WHERE address = '%s';", e)
	_, err = i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
