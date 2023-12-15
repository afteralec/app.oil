package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/permission"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestReviewCharacterApplicationsPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	_ = CallLogin(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.ReviewCharacterApplicationsPath())
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestReviewCharacterApplicationssPageNoPermission(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.ReviewCharacterApplicationsPath())
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestReviewCharacterApplicationssPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}

	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permission.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Rework this to delete all permissions by username
	defer DeleteTestPlayerPermission(t, &i, prid)

	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.ReviewCharacterApplicationsPath())
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestReviewCharacterApplicationsPageFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestCharacters(t, &i, TestUsername)

	CallRegister(t, a, TestUsername, TestPassword)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permission.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	prid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Rework this to delete all permissions by username

	res := CallLogin(t, a, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
	i.Close()

	url := MakeTestURL(routes.ReviewCharacterApplicationsPath())
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = shared.SetupInterfaces()
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestPlayerPermission(t, &i, prid)

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}
