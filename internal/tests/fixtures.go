package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"

	"petrichormud.com/app/internal/routes"
)

// TODO: Rename this to include Fixture?
func EditEmailRequest(id int64, e string) *http.Request {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("email", e)
	writer.Close()

	url := MakeTestURL(routes.EmailPath(strconv.FormatInt(id, 10)))
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
