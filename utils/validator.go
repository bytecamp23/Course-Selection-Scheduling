package utils

import (
	"Course-Selection-Scheduling/internal/global"
	"github.com/go-playground/validator/v10"
	"regexp"
	"unicode"
)

var isUpperOrLowerCase = "^[a-zA-Z]{8,20}$"

//var isUpperAndLowerAndDigit = "^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{8,20}$"
var isDigit = "^[0-9]+$"
var isUserType = "^[123]$"

func UserNameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().Interface().(string)
	if ret, _ := regexp.MatchString(isUpperOrLowerCase, username); !ret {
		return false
	}
	return true
}

func PasswordValidator(fl validator.FieldLevel) bool {
	var (
		isUpper  = false
		isLower  = false
		isNumber = false
	)
	password := fl.Field().Interface().(string)
	length := len(password)
	if length < 8 || length > 20 {
		return false
	}
	for _, s := range password {
		switch {
		case unicode.IsUpper(s):
			isUpper = true
		case unicode.IsLower(s):
			isLower = true
		case unicode.IsNumber(s):
			isNumber = true
		default:
		}
	}
	if !isUpper || !isLower || !isNumber {
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
