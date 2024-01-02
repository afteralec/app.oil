package util

import (
	"fmt"
	"regexp"
	"strings"
)

func RegexForExactMatchStrings(ss []string) *regexp.Regexp {
	var sb strings.Builder
	fmt.Fprintf(&sb, "^(%s)$", strings.Join(ss, "|"))
	return regexp.MustCompile(sb.String())
}
