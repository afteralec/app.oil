package tests

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/rooms"
	"petrichormud.com/app/internal/routes"
)

func TestRoomsPageUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.Rooms)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestRoomsPageForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.Rooms)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestRoomsPageSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerViewAllRoomsName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.Rooms)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestRoomPageUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	rmid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rmid)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerViewAllRoomsName)
	defer DeleteTestPlayerPermission(t, &i, prid)

	url := MakeTestURL(routes.RoomPath(rmid))

	req := httptest.NewRequest(http.MethodGet, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestRoomPageForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	rmid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rmid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomPath(rmid))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestRoomPageNotFound(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	rmid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rmid)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerViewAllRoomsName)
	defer DeleteTestPlayerPermission(t, &i, prid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomPath(rmid + 1))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestRoomPageSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	rmid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rmid)

	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerViewAllRoomsName)
	defer DeleteTestPlayerPermission(t, &i, prid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomPath(rmid))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestNewRoomUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	url := MakeTestURL(routes.Rooms)

	req := httptest.NewRequest(http.MethodPost, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestNewRoomForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.Rooms)

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestNewRoomSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.Rooms)

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestNewRoomWithLinkUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	url := MakeTestURL(routes.Rooms)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(rid, 10))
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestNewRoomWithLinkForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.Rooms)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(rid, 10))
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestNewRoomWithLinkNotFoundInvalidRID(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.Rooms)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d", rid+100)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", sb.String())
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestNewRoomWithLinkBadRequestInvalidDirection(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.Rooms)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(rid, 10))
	writer.WriteField("direction", "weast")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
}

func TestNewRoomWithLinkSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)
	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.Rooms)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(rid, 10))
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestEditRoomPageUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	url := MakeTestURL(routes.EditRoomPath(rid))

	req := httptest.NewRequest(http.MethodGet, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditRoomPageForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.EditRoomPath(rid))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestEditRoomPageSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	permissionID := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, permissionID)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.EditRoomPath(rid))

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestEditRoomExitUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	ridOne := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridOne)
	ridTwo := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridTwo)

	url := MakeTestURL(routes.RoomExitsPath(ridOne))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(ridTwo, 10))
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.WriteField("two-way", "true")
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditRoomExitForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	ridOne := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridOne)
	ridTwo := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomExitsPath(ridOne))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(ridTwo, 10))
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.WriteField("two-way", "true")
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

func TestEditRoomExitRoomNotFound(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomExitsPath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(rid+1000, 10))
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.WriteField("two-way", "true")
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

func TestEditRoomExitBadRequestLinkToSelf(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomExitsPath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(rid, 10))
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.WriteField("two-way", "true")
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

func TestEditRoomExitBadRequestInvalidID(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	ridOne := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridOne)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomExitsPath(ridOne))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", "notanid")
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.WriteField("two-way", "true")
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

func TestEditRoomExitBadRequestEmptyID(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	ridOne := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridOne)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomExitsPath(ridOne))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", "0")
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.WriteField("two-way", "true")
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

func TestEditRoomExitBadRequestInvalidTwoWay(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomExitsPath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(rid+1000, 10))
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.WriteField("two-way", "notaboolean")
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

func TestEditRoomExitSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	ridOne := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridOne)
	ridTwo := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridTwo)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomExitsPath(ridOne))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", strconv.FormatInt(ridTwo, 10))
	writer.WriteField("direction", rooms.DirectionNorth)
	writer.WriteField("two-way", "true")
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

func TestClearRoomExitUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	ridOne := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridOne)
	ridTwo := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridTwo)

	if err := rooms.Link(rooms.LinkParams{
		Queries:   i.Queries,
		ID:        ridOne,
		To:        ridTwo,
		Direction: rooms.DirectionNorth,
		TwoWay:    true,
	}); err != nil {
		t.Fatal(err)
	}

	url := MakeTestURL(routes.RoomExitPath(ridOne, rooms.DirectionNorth))

	req := httptest.NewRequest(http.MethodDelete, url, nil)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestUnmodifiedRooms(t, &i)

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestClearRoomExitForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	ridOne := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridOne)
	ridTwo := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridTwo)

	if err := rooms.Link(rooms.LinkParams{
		Queries:   i.Queries,
		ID:        ridOne,
		To:        ridTwo,
		Direction: rooms.DirectionNorth,
		TwoWay:    true,
	}); err != nil {
		t.Fatal(err)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomExitPath(ridOne, rooms.DirectionNorth))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusForbidden, res.StatusCode)
}

func TestClearRoomExitSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	ridOne := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridOne)
	ridTwo := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, ridTwo)

	if err := rooms.Link(rooms.LinkParams{
		Queries:   i.Queries,
		ID:        ridOne,
		To:        ridTwo,
		Direction: rooms.DirectionNorth,
		TwoWay:    true,
	}); err != nil {
		t.Fatal(err)
	}

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomExitPath(ridOne, rooms.DirectionNorth))

	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.AddCookie(sessionCookie)

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func TestEditRoomTitleUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	url := MakeTestURL(routes.RoomTitlePath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "A bright, calm stretch of ocean")
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditRoomTitleForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomTitlePath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "A bright, calm stretch of ocean")
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

func TestEditRoomTitleNotFound(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomTitlePath(rid + 1000))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "A bright, calm stretch of ocean")
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

func TestEditRoomTitleBadRequestInvalid(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomTitlePath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Invalid room!1234")
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

func TestEditRoomTitleSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomTitlePath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "A bright, calm stretch of ocean")
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

func TestEditRoomDescriptionUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	url := MakeTestURL(routes.RoomDescriptionPath(rid))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s And has changes.", TestRoom.Description)
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

func TestEditRoomDescriptionForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomDescriptionPath(rid))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s And has changes.", TestRoom.Description)
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

func TestEditRoomDescriptionNotFound(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomDescriptionPath(rid + 1000))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s And has changes.", TestRoom.Description)
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

func TestEditRoomDescriptionBadRequestInvalid(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomDescriptionPath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("desc", "Invalid room!1234")
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

func TestEditRoomDescriptionSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomDescriptionPath(rid))

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s And has changes.", TestRoom.Description)
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

func TestEditRoomSizeUnauthorized(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	url := MakeTestURL(routes.RoomSizePath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("size", "3")
	writer.Close()

	req := httptest.NewRequest(http.MethodPatch, url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestEditRoomSizeForbiddenNoPermission(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomSizePath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("size", "3")
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

func TestEditRoomSizeNotFound(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomSizePath(rid + 1000))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("size", "3")
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

func TestEditRoomSizeBadRequestInvalid(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomSizePath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("size", "5")
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

func TestEditRoomSizeSuccess(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	a := fiber.New(config.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	pid := CreateTestPlayer(t, &i, a, TestUsername, TestPassword)
	defer DeleteTestPlayer(t, &i, TestUsername)
	prid := CreateTestPlayerPermission(t, &i, pid, permissions.PlayerCreateRoomName)
	defer DeleteTestPlayerPermission(t, &i, prid)
	rid := CreateTestRoom(t, &i, TestRoom)
	defer DeleteTestRoom(t, &i, rid)

	sessionCookie := LoginTestPlayer(t, a, TestUsername, TestPassword)

	url := MakeTestURL(routes.RoomSizePath(rid))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("size", "3")
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
