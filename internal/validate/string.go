package validate

import "regexp"

type StringValidator interface {
	IsValid(s string) bool
}

type StringValidatorGroup struct {
	Validators []StringValidator
}

func NewStringValidatorGroup(validators []StringValidator) StringValidatorGroup {
	return StringValidatorGroup{
		Validators: validators,
	}
}

func (v *StringValidatorGroup) IsValid(s string) bool {
	for _, validator := range v.Validators {
		if !validator.IsValid(s) {
			return false
		}
	}

	return true
}

type StringLengthValidator struct {
	MinLen int
	MaxLen int
}

func NewStringLengthValidator(min, max int) StringLengthValidator {
	return StringLengthValidator{
		MinLen: min,
		MaxLen: max,
	}
}

func (v *StringLengthValidator) IsValid(s string) bool {
	if len(s) < v.MinLen {
		return false
	}
	if len(s) > v.MaxLen {
		return false
	}
	return true
}

type StringRegexMatchValidator struct {
	Regex *regexp.Regexp
}

func NewStringRegexMatchValidator(regex *regexp.Regexp) StringRegexMatchValidator {
	return StringRegexMatchValidator{
		Regex: regex,
	}
}

func (v *StringRegexMatchValidator) IsValid(s string) bool {
	return v.Regex.MatchString(s)
}

type StringRegexNoMatchValidator struct {
	Regex *regexp.Regexp
}

func NewStringRegexNoMatchValidator(regex *regexp.Regexp) StringRegexNoMatchValidator {
	return StringRegexNoMatchValidator{
		Regex: regex,
	}
}

func (v *StringRegexNoMatchValidator) IsValid(s string) bool {
	return !v.Regex.MatchString(s)
}
