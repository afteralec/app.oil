package web

import (
	"embed"
	"net/http"

	"github.com/gofiber/template/html/v2"
)

//go:embed templates/*
var templates embed.FS

func ViewsEngine() *html.Engine {
	return html.NewFileSystem(http.FS(templates), ".html")
}
