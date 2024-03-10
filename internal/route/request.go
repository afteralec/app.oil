package route

import "fmt"

const CreateRequestCommentPathParam string = "/request/:id/comment/:field"

func CreateRequestCommentPath(rid int64, field string) string {
	return fmt.Sprintf("/request/%d/comment/%s", rid, field)
}
