package routes

import (
	"fmt"
	"strings"
)

const (
	Help              string = "/help"
	HelpFilePathParam string = "/help/:slug"
)

func HelpFilePath(slug string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%s", Help, slug)
	return sb.String()
}
