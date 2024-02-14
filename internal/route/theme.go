package route

import (
	"fmt"
	"strings"
)

const Theme string = "/theme"

const ThemePathParam string = "/theme/:theme"

func ThemePath(theme string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%s", Theme, theme)
	return sb.String()
}
