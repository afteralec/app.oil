package actors

import (
	"fmt"
	"strings"
)

func ImageTitleWithID(name string, id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "[%d] %s", id, name)
	return sb.String()
}
