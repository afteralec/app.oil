package tests

import (
	"context"
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

func TestSubmitCharacterApplicationUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
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

	a := fiber.New(configs.Fiber())
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

	a := fiber.New(configs.Fiber())
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

	a := fiber.New(configs.Fiber())
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

	a := fiber.New(configs.Fiber())
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

	a := fiber.New(configs.Fiber())
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

	a := fiber.New(configs.Fiber())
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

func TestCancelCharacterApplicationUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	url := MakeTestURL(routes.CharacterApplicationPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCancelCharacterApplicationUnowned(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.CharacterApplicationPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestCancelCharacterApplicationSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)

	url := MakeTestURL(routes.CharacterApplicationPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestCancelCharacterApplicationNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)

	url := MakeTestURL(routes.CharacterApplicationPath(strconv.FormatInt(rid+1, 10)))
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCancelCharacterApplicationFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)

	url := MakeTestURL(routes.CharacterApplicationPath(strconv.FormatInt(rid, 10)))
	i.Close()
	req := httptest.NewRequest(http.MethodDelete, url, nil)
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

func TestPutCharacterApplicationInReviewUnauthorized(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestCharacterApplication(t, &i, rid)
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}

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

	url := MakeTestURL(routes.PutCharacterApplicationInReviewPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestPutCharacterApplicationInReviewOwnApplicationForbidden(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}

	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
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

	res := CallLogin(t, a, TestUsername, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.PutCharacterApplicationInReviewPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestPutCharacterApplicationInReviewNoPermissionForbidden(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.PutCharacterApplicationInReviewPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

// TODO: Write tests for every status except Submitted
func TestPutCharacterApplicationInReviewNotSubmitted(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)

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

	url := MakeTestURL(routes.PutCharacterApplicationInReviewPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestPutCharacterApplicationInReviewSuccess(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	// TODO: This is a hack. Rework this to use valid content and the existing handlers
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}

	CallRegister(t, a, TestUsernameTwo, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsernameTwo)
	if err != nil {
		t.Fatal(err)
	}
	rp, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        p.ID,
		IPID:       p.ID,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
	})
	if err != nil {
		t.Fatal(err)
	}
	permid, err := rp.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestPlayerPermission(t, &i, permid)

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.PutCharacterApplicationInReviewPath(strconv.FormatInt(rid, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestPutCharacterApplicationInReviewNotFound(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	// TODO: This is a hack. Rework this to use valid content and the existing handlers
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}

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

	url := MakeTestURL(routes.PutCharacterApplicationInReviewPath(strconv.FormatInt(rid+1, 10)))
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestPutCharacterApplicationInReviewFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	// TODO: This is a hack. Rework this to use valid content and the existing handlers
	if err := i.Queries.MarkRequestSubmitted(context.Background(), rid); err != nil {
		t.Fatal(err)
	}

	CallRegister(t, a, TestUsernameTwo, TestPassword)
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

	res := CallLogin(t, a, TestUsernameTwo, TestPassword)
	sessionCookie := res.Cookies()[0]

	url := MakeTestURL(routes.PutCharacterApplicationInReviewPath(strconv.FormatInt(rid, 10)))
	i.Close()
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)
	res, err = a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	i = shared.SetupInterfaces()
	defer DeleteTestCharacterApplication(t, &i, rid)
	defer DeleteTestPlayer(t, &i, TestUsername)
	defer DeleteTestPlayer(t, &i, TestUsernameTwo)
	defer DeleteTestPlayerPermission(t, &i, prid)

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}
