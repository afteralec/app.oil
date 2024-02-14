package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/route"
)

func TestVerifyEmailPageUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestVerifyEmailPageSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestVerifyEmailPageForbiddenUnowned(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	CreateTestPlayer(t, &i, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsernameTwo, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestVerifyPageExpiredToken(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	key := keys[0]
	if err := i.Redis.Expire(context.Background(), key, 0).Err(); err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(key, ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestVerifyEmailUnauthorizedNotLoggedIn(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[0]))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestVerifyEmailBadRequestMissingToken(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(route.VerifyEmail)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestVerifyEmailNotFoundExpiredToken(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	key := keys[0]
	if err := i.Redis.Expire(context.Background(), key, 0).Err(); err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(key, ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestVerifyEmailSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	CreateTestEmail(t, &i, a, TestEmailAddress, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer FlushTestRedis(t, &i)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	keys, err := i.Redis.Keys(context.Background(), email.VerificationKey("*")).Result()
	if err != nil {
		t.Fatal(err)
	}
	keyParts := strings.Split(keys[0], ":")

	url := MakeTestURL(route.VerifyEmailWithToken(keyParts[1]))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}
