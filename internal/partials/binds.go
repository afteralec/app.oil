package partials

import fiber "github.com/gofiber/fiber/v2"

var BindLoginErr = fiber.Map{
	"NoticeSectionID": "login-err",
	"NoticeText": []string{
		"The username and password you entered couldn't be verified.",
		"Please try again.",
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
