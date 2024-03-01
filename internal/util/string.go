package util

import (
	"fmt"
	"strings"
)

func RegexForExactMatchStrings(ss []string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "^(%s)$", strings.Join(ss, "|"))
	return sb.String()
}
