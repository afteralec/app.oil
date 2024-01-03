package views

const (
	Recover                string = "recover"
	RecoverPassword        string = "recover-password"
	RecoverPasswordSuccess string = "recover-password-success"
	RecoverUsername        string = "recover-username"
	RecoverUsernameSuccess string = "recover-username-success"
)

const (
	PartialRecoverUsernameErrInternal string = "partial-recover-username-err-internal"
	PartialRecoverUsernameErrInvalid  string = "partial-recover-username-err-invalid"
)

const (
	ResetPassword        string = "reset-password"
	ResetPasswordSuccess string = "reset-password-success"
)

const Profile string = "profile"

const (
	PartialProfileEmailNewUnverified string = "partial-profile-email-new-unverified"
	PartialProfileEmailErrInternal   string = "partial-profile-email-err-internal"
	// TODO: Change this to a forbidden error
	PartialProfileEmailErrTooMany  string = "partial-profile-email-err-too-many"
	PartialProfileEmailErrInvalid  string = "partial-profile-email-err-invalid"
	PartialProfileEmailErrConflict string = "partial-profile-email-err-conflict"
)

const (
	PartialProfileEmailDeleteErrUnauthorized string = "partial-profile-email-delete-err-unauthorized"
	PartialProfileEmailDeleteErrInternal     string = "partial-profile-email-delete-err-internal"
	PartialProfileEmailDeleteErrNotFound     string = "partial-profile-email-delete-err-not-found"
	PartialProfileEmailDeleteSuccess         string = "partial-profile-email-delete-success"
)

const Characters string = "characters"

const PlayerPermissionsDetail string = "player-permissions-detail"

const CharacterApplicationQueue string = "character-application-queue"

const (
	PartialPlayerFree        string = "partial-player-free"
	PartialPlayerReserved    string = "partial-player-reserved"
	PartialPlayerReservedErr string = "partial-player-reserved-err"
)

const (
	PartialRegisterErrInvalid  string = "partial-register-err-invalid"
	PartialRegisterErrInternal string = "partial-register-err-internal"
	PartialRegisterErrConflict string = "partial-register-err-conflict"
)

const (
	PartialLoginErr string = "partial-login-err"
)

const (
	PartialPlayerPermissionsSearchResults string = "partial-player-permissions-search-results"
)
