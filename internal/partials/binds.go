package partials

import fiber "github.com/gofiber/fiber/v2"

var BindLoginErr = fiber.Map{
	"NoticeSectionID": "register-err",
	"NoticeText": []string{
		"The username and password you entered couldn't be verified.",
		"Please try again.",
	},
}

var BindRegisterErrInternal = fiber.Map{
	"NoticeSectionID": "register-err",
	"NoticeText": []string{
		"Something's gone terribly wrong.",
	},
	"RefreshButton": true,
}

var BindRegisterErrInvalidUsername = fiber.Map{
	"NoticeSectionID": "register-err",
	"NoticeText": []string{
		"What you entered isn't a valid username.",
		"Please follow the prompts and try again.",
	},
}

var BindRegisterErrInvalidPassword = fiber.Map{
	"NoticeSectionID": "register-err",
	"NoticeText": []string{
		"What you entered isn't a valid password.",
		"Please follow the prompts and try again.",
	},
}

var BindRegisterErrInvalidConfirmPassword = fiber.Map{
	"NoticeSectionID": "register-err",
	"NoticeText": []string{
		"That password and password confirmation don't match.",
		"Please re-enter your password confirmation.",
	},
}

var BindRegisterErrConflict = fiber.Map{
	"NoticeSectionID": "register-err",
	"NoticeText": []string{
		"Sorry! That username is already taken.",
		"Please try a different username.",
	},
}

var BindRecoverUsernameErrInvalid = fiber.Map{
	"NoticeSectionID": "recover-username-err",
	"NoticeText": []string{
		"What you entered isn't a valid email address.",
		"Please try again.",
	},
}

var BindRecoverUsernameErrInternal = fiber.Map{
	"NoticeSectionID": "recover-username-err",
	"NoticeText": []string{
		"Something's gone horribly wrong.",
	},
	"RefreshButton": true,
}

var BindRecoverPasswordErrInvalidEmail = fiber.Map{
	"NoticeSectionID": "recover-password-err",
	"NoticeText": []string{
		"What you entered isn't a valid email address.",
		"Please try again.",
	},
}

var BindRecoverPasswordErrInvalidUsername = fiber.Map{
	"NoticeSectionID": "recover-password-err",
	"NoticeText": []string{
		"What you entered isn't a valid username.",
		"Please try again.",
	},
}

var BindRecoverPasswordErrInternal = fiber.Map{
	"NoticeSectionID": "recover-password-err",
	"NoticeText": []string{
		"Something's gone horribly wrong.",
	},
	"RefreshButton": true,
}

var BindResetPasswordErr = fiber.Map{
	"NoticeSectionID": "login-err",
	"NoticeText": []string{
		"The username and password you entered couldn't be verified.",
		"Please try again.",
	},
}
