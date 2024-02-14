package route

import "fmt"

func CreateRequestCommentPath(rid, field string) string {
	return fmt.Sprintf("/request/%s/comment/%s", rid, field)
}
