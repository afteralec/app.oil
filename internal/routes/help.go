package routes

import (
	"fmt"
	"strings"
)

const (
	Help          string = "/help"
	HelpPathParam string = "/help/:slug"
)

func HelpPath(slug string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%s", Help, slug)
	return sb.String()
}
