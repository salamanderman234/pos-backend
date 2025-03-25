package config

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

func registerValidation() {
	vld.RegisterValidation("after_now", afterNowValidation)
	vld.RegisterValidation("password", passwordValidation)
}

func afterNowValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Int()
	unixTime := time.Unix(val, 0)
	now := time.Now()
	return unixTime.After(now)
}

func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Ensure length >= 8
	if len(password) < 8 {
		return false
	}

	// Regular expressions for validation
	hasUppercase := `[A-Z]`
	hasNumber := `[0-9]`
	hasSymbol := `[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':",\.<>\/?\\|]`

	// Check for uppercase letter
	if match, _ := regexp.MatchString(hasUppercase, password); !match {
		return false
	}

	// Check for number
	if match, _ := regexp.MatchString(hasNumber, password); !match {
		return false
	}

	// Check for special symbol
	if match, _ := regexp.MatchString(hasSymbol, password); !match {
		return false
	}

	return true
}
