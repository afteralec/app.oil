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

type contentreview struct {
	Inner map[string]string
}

// TODO: Handling for when an invalid fidl status is stored in the system
func (c *contentreview) Status(field string) (string, bool) {
	status, ok := c.Inner[field]
	if !ok {
		return "", false
	}
	return status, true
}
