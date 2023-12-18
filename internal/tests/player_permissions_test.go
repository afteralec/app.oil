package tests

import (
	"context"
	"fmt"
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
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestPlayerPermissionsPageUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.PlayerPermissions)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestPlayerPermissionsPageForbidden(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]
	url := MakeTestURL(routes.PlayerPermissions)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestPlayerPermissionsPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: This is a hack. Update this test to simulate adding permissions via a dummy super-user
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerGrantAllPermissions.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	permid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, permid)

	url := MakeTestURL(routes.PlayerPermissions)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestPlayerPermissionsDetailPageUnauthorized(t *testing.T) {
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

	url := MakeTestURL(routes.PlayerPermissionsDetailPath(strconv.FormatInt(p.ID, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestPlayerPermissionsDetailPageForbidden(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}

	url := MakeTestURL(routes.PlayerPermissionsDetailPath(strconv.FormatInt(p.ID, 10)))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestPlayerPermissionsDetailPageSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CallRegister(t, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: This is a hack. Update this test to simulate adding permissions via a dummy super-user
	pr, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerGrantAllPermissions.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	permid, err := pr.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, permid)

	url := MakeTestURL(routes.PlayerPermissionsDetailPath(TestUsername))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func DeleteTestPlayerPermission(t *testing.T, i *shared.Interfaces, id int64) {
	query := fmt.Sprintf("DELETE FROM player_permissions WHERE id = %d;", id)
	_, err := i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
