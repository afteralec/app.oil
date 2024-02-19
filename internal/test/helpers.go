package test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
)

// TODO: Create queries to delete records for use during these test.
// It makes more sense to use the stable API than to risk breaking the test suite on a table name change
//
// TODO: See if some of these could accept a session cookie and the PID instead of calling up the player itself
// TODO: Turn some of these into transactions, or make a transaction-enabled version?
//
// TODO: Add a function to clean up resources starting with the test usernames -
// pretty much everything can be traced up to a PID - and call it from the CLI

func CreateTestPlayer(t *testing.T, i *service.Interfaces, a *fiber.App, u, pw string) int64 {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", u)
	writer.WriteField("password", pw)
	writer.WriteField("confirmPassword", pw)
	writer.Close()

	url := MakeTestURL(route.Register)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	p, err := i.Queries.GetPlayerByUsername(context.Background(), u)
	if err != nil {
		t.Fatal(err)
	}
	return p.ID
}

func DeleteTestPlayer(t *testing.T, i *service.Interfaces, u string) {
	p, err := i.Queries.GetPlayerByUsername(context.Background(), u)
	if err != nil {
		return
	}

	_, err = i.Database.Exec("DELETE FROM players WHERE username = ?;", u)
	if err != nil {
		t.Fatal(err)
	}

	_, err = i.Database.Exec("DELETE FROM emails WHERE pid = ?;", p.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = i.Database.Exec("DELETE FROM player_settings WHERE pid = ?;", p.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func LoginTestPlayer(t *testing.T, a *fiber.App, u string, pw string) *http.Cookie {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", u)
	writer.WriteField("password", pw)
	writer.Close()

	url := MakeTestURL(route.Login)

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Write a better way to pull the session cookie
	return res.Cookies()[0]
}

func CreateTestEmail(t *testing.T, i *service.Interfaces, a *fiber.App, e, u, pw string) int64 {
	sessionCookie := LoginTestPlayer(t, a, u, pw)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", e)
	writer.Close()

	url := MakeTestURL(route.NewEmailPath())

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	p, err := i.Queries.GetPlayerByUsername(context.Background(), u)
	if err != nil {
		t.Fatal(err)
	}

	email, err := i.Queries.GetEmailByAddressForPlayer(context.Background(), query.GetEmailByAddressForPlayerParams{
		PID:     p.ID,
		Address: e,
	})
	if err != nil {
		t.Fatal(err)
	}

	return email.ID
}

// TODO: Rework this to use endpoints on the app instead of interfaces directly
func CreateTestPlayerPermission(t *testing.T, i *service.Interfaces, pid int64, name string) int64 {
	permissionResult, err := i.Queries.CreatePlayerPermission(context.Background(), query.CreatePlayerPermissionParams{
		PID:  pid,
		IPID: pid,
		Name: name,
	})
	if err != nil {
		t.Fatal(err)
	}
	permissionID, err := permissionResult.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return permissionID
}

func DeleteTestPlayerPermission(t *testing.T, i *service.Interfaces, id int64) {
	query := fmt.Sprintf("DELETE FROM player_permissions WHERE id = %d;", id)
	_, err := i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func CreateTestCharacterApplication(t *testing.T, i *service.Interfaces, a *fiber.App, u, pw string) int64 {
	sessionCookie := LoginTestPlayer(t, a, u, pw)

	url := MakeTestURL(route.Characters)

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.AddCookie(sessionCookie)

	_, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	p, err := i.Queries.GetPlayerByUsername(context.Background(), u)
	if err != nil {
		t.Fatal(err)
	}

	apps, err := i.Queries.ListCharacterApplicationsForPlayer(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}

	return apps[0].Request.ID
}

func DeleteTestCharacterApplication(t *testing.T, i *service.Interfaces, rid int64) {
	DeleteTestRequest(t, i, rid)

	_, err := i.Database.Exec("DELETE FROM character_application_content WHERE id = ?;", rid)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}

	_, err = i.Database.Exec("DELETE FROM character_application_content_history WHERE rid = ?;", rid)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
}

func DeleteTestRequest(t *testing.T, i *service.Interfaces, rid int64) {
	_, err := i.Database.Exec("DELETE FROM requests WHERE id = ?;", rid)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}

	_, err = i.Database.Exec("DELETE FROM request_comments WHERE rid = ?;", rid)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}

	_, err = i.Database.Exec("DELETE FROM request_comment_history WHERE rid = ?;", rid)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}

	_, err = i.Database.Exec("DELETE FROM request_status_change_history WHERE rid = ?;", rid)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
}

type CreateTestRequestCommentParams struct {
	T        *testing.T
	I        *service.Interfaces
	A        *fiber.App
	Username string
	Password string
	Field    string
	RID      int64
}

func CreateTestRequestComment(params CreateTestRequestCommentParams) {
	sessionCookie := LoginTestPlayer(params.T, params.A, params.Username, params.Password)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("comment", "This is a test comment, for sure.")
	writer.Close()

	url := MakeTestURL(route.CreateRequestCommentPath(strconv.FormatInt(params.RID, 10), params.Field))

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	_, err := params.A.Test(req)
	if err != nil {
		params.T.Fatal(err)
	}
}

func FlushTestRedis(t *testing.T, i *service.Interfaces) {
	if err := i.Redis.FlushAll(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
}

func ListEmailsForPlayer(t *testing.T, i *service.Interfaces, username string) []query.Email {
	p, err := i.Queries.GetPlayerByUsername(context.Background(), username)
	if err != nil {
		t.Fatal(err)
	}
	emails, err := i.Queries.ListEmails(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	return emails
}

type CreateTestRoomParams struct {
	Title       string
	Description string
	Size        int32
}

func CreateTestRoom(t *testing.T, i *service.Interfaces, p CreateTestRoomParams) int64 {
	result, err := i.Queries.CreateRoom(context.Background(), query.CreateRoomParams{
		Title:       p.Title,
		Description: p.Description,
		Size:        p.Size,
	})
	if err != nil {
		t.Fatal(err)
	}

	rid, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	return rid
}

func DeleteTestRoom(t *testing.T, i *service.Interfaces, id int64) {
	_, err := i.Database.Exec("DELETE FROM rooms WHERE id = ?;", id)
	if err != nil {
		t.Fatal(err)
	}
}

func DeleteTestUnmodifiedRooms(t *testing.T, i *service.Interfaces) {
	// TODO: Get a helper to delete rooms that are orphaned, off-grid, closely resemble the base room, etc
	_, err := i.Database.Exec("DELETE FROM rooms WHERE unmodified = true;")
	if err != nil {
		t.Fatal(err)
	}
}

type CreateTestActorImageParams struct {
	Gender           string
	Name             string
	ShortDescription string
	Description      string
}

func CreateTestActorImage(t *testing.T, i *service.Interfaces, p CreateTestActorImageParams) int64 {
	result, err := i.Queries.CreateActorImage(context.Background(), query.CreateActorImageParams{
		Gender:           p.Gender,
		Name:             p.Name,
		ShortDescription: p.ShortDescription,
		Description:      p.Description,
	})
	if err != nil {
		t.Fatal(err)
	}
	rid, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return rid
}

func DeleteTestActorImage(t *testing.T, i *service.Interfaces, aiid int64) {
	_, err := i.Database.Exec("DELETE FROM actor_images WHERE id = ?;", aiid)
	if err != nil {
		t.Fatal(err)
	}
}

func DeleteTestActorImageByName(t *testing.T, i *service.Interfaces, name string) {
	_, err := i.Database.Exec("DELETE FROM actor_images WHERE name = ?;", name)
	if err != nil {
		t.Fatal(err)
	}
}

type CreateTestHelpFileRelatedParams struct {
	Title string
	Sub   string
	Slug  string
}

type CreateTestHelpFileTagParams struct {
	Tag string
}

type CreateTestHelpFileParams struct {
	Slug     string
	Title    string
	Sub      string
	Category string
	Raw      string
	HTML     string
	Related  []CreateTestHelpFileRelatedParams
	Tags     []CreateTestHelpFileTagParams
	PID      int64
}

func CreateTestHelpFile(t *testing.T, i *service.Interfaces, p CreateTestHelpFileParams) {
	_, err := i.Database.Exec(
		"INSERT INTO help (slug, title, sub, category, pid, raw, html) VALUES (?, ?, ?, ?, ?, ?, ?);",
		p.Slug,
		p.Title,
		p.Sub,
		p.Category,
		p.PID,
		p.Raw,
		p.HTML,
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, related := range p.Related {
		_, err := i.Database.Exec(
			"INSERT INTO help_related (slug, related_title, related_sub, related_slug) VALUES (?, ?, ?, ?);",
			p.Slug,
			related.Title,
			related.Sub,
			related.Slug,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, tag := range p.Tags {
		_, err := i.Database.Exec(
			"INSERT INTO help_tags (slug, tag) VALUES (?, ?);",
			p.Slug,
			tag.Tag,
		)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func DeleteTestHelpFile(t *testing.T, i *service.Interfaces, slug string) {
	_, err := i.Database.Exec("DELETE FROM help WHERE slug = ?;", slug)
	if err != nil {
		t.Fatal(err)
	}

	_, err = i.Database.Exec("DELETE FROM help_related WHERE slug = ?;", slug)
	if err != nil {
		t.Fatal(err)
	}

	_, err = i.Database.Exec("DELETE FROM help_tags WHERE slug = ?;", slug)
	if err != nil {
		t.Fatal(err)
	}
}
