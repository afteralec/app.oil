package request

type content struct {
	Inner map[string]string
}

func (c *content) Value(field string) (string, bool) {
	value, ok := c.Inner[field]
	if !ok {
		return "", false
	}
	return value, true
}
