package routes

import "fmt"

const (
	Login                  = "/login"
	Logout                 = "/logout"
	Register               = "/player/new"
	Reserved               = "/player/reserved"
	Profile                = "/profile"
	Recover                = "/recover"
	RecoverUsername        = "/recover/username"
	RecoverUsernameSuccess = "/recover/username/success"
	RecoverPassword        = "/recover/password"
	RecoverPasswordSuccess = "/recover/password/success"
	ResetPassword          = "/reset/password"
	ResetPasswordSuccess   = "/reset/password/success"
	SearchPlayer           = "/player/search"
)

func SearchPlayerPath(dest string) string {
	return fmt.Sprintf("%s/%s", SearchPlayer, dest)
}
