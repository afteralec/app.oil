package permissions

const (
	Login      = "Login"
	ViewPlayer = "ViewPlayer"
)

func DefaultSet() [1]string {
	return [1]string{Login}
}

func AdminSet() [2]string {
	return [2]string{Login, ViewPlayer}
}
