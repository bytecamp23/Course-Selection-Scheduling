package utils

import (
	"Course-Selection-Scheduling/internal/global"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var isUpperOrLowerCase = "^[a-zA-Z]+$"
var isUpperOrLowerOrDigit = "^[a-zA-Z0-9]+$"
var isDigit = "^[0-9]+$"
var isUserType = "^[123]$"

func UserNameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().Interface().(string)
	if len(username) < 8 || len(username) > 20 {
		return false
	}
	if ret, _ := regexp.MatchString(isUpperOrLowerCase, username); !ret {
		return false
	}
	return true
}

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().Interface().(string)
	if len(password) < 8 || len(password) > 20 {
		return false
	}
	if ret, _ := regexp.MatchString(isUpperOrLowerOrDigit, password); !ret {
		return false
	}
	return true
}

func UserTypeValidator(fl validator.FieldLevel) bool {
	userType := fl.Field().Interface().(global.UserType)
	if userType != global.Admin && userType != global.Student && userType != global.Teacher {
		return false
	}
	return true
}

func IsDigitValidator(fl validator.FieldLevel) bool {
	userID := fl.Field().Interface().(string)
	if ret, _ := regexp.MatchString(isDigit, userID); !ret {
		return false
	}
	return true
}

func IsUpperOrLowerOrDigitValidator(fl validator.FieldLevel) bool {
	userID := fl.Field().Interface().(string)
	if ret, _ := regexp.MatchString(isDigit, userID); !ret {
		return false
	}
	return true
}
