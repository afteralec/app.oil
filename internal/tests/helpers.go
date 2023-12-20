package tests

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
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

// TODO: Create queries to delete records for use during these tests.
// It makes more sense to use the stable API than to risk breaking the test suite on a table name change
//
// TODO: See if some of these could accept a session cookie and the PID instead of calling up the player itself
// TODO: Turn some of these into transactions, or make a transaction-enabled version?

func CreateTestPlayer(t *testing.T, i *shared.Interfaces, a *fiber.App, u, pw string) int64 {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", u)
	writer.WriteField("password", pw)
	writer.WriteField("confirmPassword", pw)
	writer.Close()

	url := MakeTestURL(routes.Register)
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

func DeleteTestPlayer(t *testing.T, i *shared.Interfaces, u string) {
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
}

func LoginTestPlayer(t *testing.T, a *fiber.App, u string, pw string) *http.Cookie {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", u)
	writer.WriteField("password", pw)
	writer.Close()

	url := MakeTestURL(routes.Login)

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Write a better way to pull the session cookie
	return res.Cookies()[0]
}

func CreateTestEmail(t *testing.T, i *shared.Interfaces, a *fiber.App, e, u, pw string) int64 {
	sessionCookie := LoginTestPlayer(t, a, u, pw)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", e)
	writer.Close()

	url := MakeTestURL(routes.NewEmailPath())

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

	email, err := i.Queries.GetEmailByAddressForPlayer(context.Background(), queries.GetEmailByAddressForPlayerParams{
		PID:     p.ID,
		Address: e,
	})
	if err != nil {
		t.Fatal(err)
	}

	return email.ID
}

func CreateTestPlayerPermission(t *testing.T, i *shared.Interfaces, pid int64) int64 {
	permissionResult, err := i.Queries.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
		PID:        pid,
		IPID:       pid,
		Permission: permissions.PlayerReviewCharacterApplicationsName,
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

func DeleteTestPlayerPermission(t *testing.T, i *shared.Interfaces, id int64) {
	query := fmt.Sprintf("DELETE FROM player_permissions WHERE id = %d;", id)
	_, err := i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func CreateTestCharacterApplication(t *testing.T, i *shared.Interfaces, a *fiber.App, u, pw string) int64 {
	sessionCookie := LoginTestPlayer(t, a, u, pw)

	url := MakeTestURL(routes.NewCharacterApplicationPath())

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

func DeleteTestCharacterApplication(t *testing.T, i *shared.Interfaces, rid int64) {
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

func DeleteTestRequest(t *testing.T, i *shared.Interfaces, rid int64) {
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
	I        *shared.Interfaces
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

	url := MakeTestURL(routes.CreateRequestCommentPath(strconv.FormatInt(params.RID, 10), params.Field))

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)

	_, err := params.A.Test(req)
	if err != nil {
		params.T.Fatal(err)
	}
}

func ListEmailsForPlayer(t *testing.T, i *shared.Interfaces, username string) []queries.Email {
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
