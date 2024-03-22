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

func (c *contentreview) AllAre(status string) bool {
	// TODO: Report this as an error to error tracking
	if !IsFieldStatusValid(status) {
		return false
	}

	for _, fs := range c.Inner {
		if fs != status {
			return false
		}
	}

	return true
}
