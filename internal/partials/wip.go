package partials

const (
	ResendVerificationEmailErrInternal        string = "partial-resend-verification-email-err-internal"
	ResendVerificationEmailErrNotFound        string = "partial-resend-verification-email-err-not-found"
	ResendVerificationEmailErrConflict        string = "partial-resend-verification-email-err-conflict"
	ResendVerificationEmailErrConflictUnowned string = "partial-resend-verification-email-err-conflict-unowned"
	ResendVerificationEmailSuccess            string = "partial-resend-verification-email-success"
)

const ProfileEmailUnverified string = "partial-profile-email-unverified"

const (
	ProfileEmailNewUnverified string = "partial-profile-email-new-unverified"
	ProfileEmailErrInternal   string = "partial-profile-email-err-internal"
	// TODO: Change this to a forbidden error
	ProfileEmailErrTooMany  string = "partial-profile-email-err-too-many"
	ProfileEmailErrInvalid  string = "partial-profile-email-err-invalid"
	ProfileEmailErrConflict string = "partial-profile-email-err-conflict"
)

const (
	ProfileEmailDeleteErrUnauthorized string = "partial-profile-email-delete-err-unauthorized"
	ProfileEmailDeleteErrInternal     string = "partial-profile-email-delete-err-internal"
	ProfileEmailDeleteErrNotFound     string = "partial-profile-email-delete-err-not-found"
	ProfileEmailDeleteSuccess         string = "partial-profile-email-delete-success"
)

const (
	ProfileEmailEditErrUnauthorized string = "partial-profile-email-edit-err-unauthorized"
	ProfileEmailEditErrInternal     string = "partial-profile-email-edit-err-internal"
	ProfileEmailEditErrNotFound     string = "partial-profile-email-edit-err-not-found"
	ProfileEmailEditSuccess         string = "partial-profile-email-edit-success"
)

const (
	RecoverUsernameErrInternal string = "partial-recover-username-err-internal"
	RecoverUsernameErrInvalid  string = "partial-recover-username-err-invalid"
)

const (
	VerifyEmailSuccess string = "partial-verify-email-success"
)

const RequestCommentCurrent string = "partial-request-comment-current"

const (
	PlayerPermissionsSearchResults string = "partial-player-permissions-search-results"
)

const (
	LoginErr string = "partial-login-err"
)

const (
	PlayerFree        string = "partial-player-free"
	PlayerReserved    string = "partial-player-reserved"
	PlayerReservedErr string = "partial-player-reserved-err"
)

const (
	RegisterErrInvalid  string = "partial-register-err-invalid"
	RegisterErrInternal string = "partial-register-err-internal"
	RegisterErrConflict string = "partial-register-err-conflict"
)
