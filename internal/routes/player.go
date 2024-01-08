package routes

import (
	"fmt"
	"strings"
)

const (
	Players                string = "/players"
	PlayerPasswordParam    string = "players/:id/password"
	Login                  string = "/login"
	Logout                 string = "/logout"
	Register               string = "/player/new"
	Reserved               string = "/player/reserved"
	Profile                string = "/profile"
	Recover                string = "/recover"
	RecoverUsername        string = "/recover/username"
	RecoverUsernameSuccess string = "/recover/username/success"
	RecoverPassword        string = "/recover/password"
	RecoverPasswordSuccess string = "/recover/password/success"
	ResetPassword          string = "/reset/password"
	ResetPasswordSuccess   string = "/reset/password/success"
	SearchPlayer           string = "/player/search"
)

func SearchPlayerPath(dest string) string {
	return fmt.Sprintf("%s/%s", SearchPlayer, dest)
}

func PlayerPasswordPath(pid int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/password", Players, pid)
	return sb.String()
}
