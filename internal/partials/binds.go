package partials

import (
	"fmt"
	"html/template"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/shared"
)

var BindLoginErr = fiber.Map{
	"NoticeSectionID": "register-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"The username and password you entered couldn't be verified.",
		"Please try again.",
	},
}

var BindRegisterErrInternal = fiber.Map{
	"NoticeSectionID": "register-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"Something's gone terribly wrong.",
	},
	"RefreshButton": true,
}

var BindRegisterErrInvalidUsername = fiber.Map{
	"NoticeSectionID": "register-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid username.",
		"Please follow the prompts and try again.",
	},
}

var BindRegisterErrInvalidPassword = fiber.Map{
	"NoticeSectionID": "register-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid password.",
		"Please follow the prompts and try again.",
	},
}

var BindRegisterErrInvalidConfirmPassword = fiber.Map{
	"NoticeSectionID": "register-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"That password and password confirmation don't match.",
		"Please re-enter your password confirmation.",
	},
}

var BindRegisterErrConflict = fiber.Map{
	"NoticeSectionID": "register-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"Sorry! That username is already taken.",
		"Please try a different username.",
	},
}

var BindRecoverUsernameErrInvalid = fiber.Map{
	"NoticeSectionID": "recover-username-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid email address.",
		"Please try again.",
	},
}

var BindRecoverUsernameErrInternal = fiber.Map{
	"NoticeSectionID": "recover-username-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"Something's gone horribly wrong.",
	},
	"RefreshButton": true,
}

var BindRecoverPasswordErrInvalidEmail = fiber.Map{
	"NoticeSectionID": "recover-password-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid email address.",
		"Please try again.",
	},
}

var BindRecoverPasswordErrInvalidUsername = fiber.Map{
	"NoticeSectionID": "recover-password-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid username.",
		"Please try again.",
	},
}

var BindRecoverPasswordErrInternal = fiber.Map{
	"NoticeSectionID": "recover-password-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"Something's gone horribly wrong.",
	},
	"RefreshButton": true,
}

var BindResetPasswordErr = fiber.Map{
	"NoticeSectionID": "login-err",
	"Class":           "pt-4",
	"NoticeText": []string{
		"The username and password you entered couldn't be verified.",
		"Please try again.",
	},
}

var BindProfileAddEmailErrUnauthorized = fiber.Map{
	"NoticeSectionID": "add-email-error",
	"Class":           "pt-2 w-[60%]",
	"NoticeText": []string{
		"Your session has expired.",
	},
	"RefreshButton": true,
}

var BindProfileAddEmailErrInternal = fiber.Map{
	"NoticeSectionID": "add-email-error",
	"Class":           "pt-2 w-[60%]",
	"NoticeText": []string{
		"Something's gone horribly wrong.",
	},
	"RefreshButton": true,
}

var BindProfileAddEmailErrInvalid = fiber.Map{
	"NoticeSectionID": "add-email-error",
	"Class":           "pt-2 w-[60%]",
	"NoticeText": []string{
		"What you entered isn't a valid email address.",
		"Please try again.",
	},
}

func BindProfileAddEmailErrTooMany() fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "You've already added the maximum %d emails.", shared.MaxEmailCount)
	return fiber.Map{
		"NoticeSectionID": "add-email-error",
		"Class":           "pt-2 w-[60%]",
		"NoticeText": []string{
			sb.String(),
		},
	}
}

func BindProfileAddEmailErrConflict(email string) fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "<span class=\"font-semibold\">%s</span> is already in use.", email)

	return fiber.Map{
		"NoticeSectionID": "add-email-error",
		"Class":           "pt-2 w-[60%]",
		"NoticeText": []template.HTML{
			template.HTML(sb.String()),
			template.HTML("Please try a different address."),
		},
	}
}
