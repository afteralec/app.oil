package partial

import (
	"html/template"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
)

type RenderParams struct {
	Bind     fiber.Map
	Template string
	Layout   string
}

func Render(e *html.Engine, p RenderParams) (template.HTML, error) {
	var sb strings.Builder
	if err := e.Render(&sb, p.Template, p.Bind, p.Layout); err != nil {
		return template.HTML(""), err
	}
	return template.HTML(sb.String()), nil
}
