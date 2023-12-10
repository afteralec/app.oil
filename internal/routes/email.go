package routes

import "fmt"

const (
	Email = "/player/email"
)

func NewEmailPath() string {
	return fmt.Sprintf("%s/new", Email)
}

func EmailPath() string {
	return fmt.Sprintf("%s/%s", Email, ":id")
}
