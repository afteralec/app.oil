package web

import (
	"embed"
	"net/http"

	"github.com/gofiber/template/html/v2"
)

//go:embed views/*
var viewsfs embed.FS

func ViewsEngine() *html.Engine {
	return html.NewFileSystem(http.FS(viewsfs), ".html")
}
