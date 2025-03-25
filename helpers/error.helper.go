package helpers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/salamanderman234/pos-backend/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// map library error to system error (check error entity)
var errorMap = map[error]config.Response{
	// gorm
	gorm.ErrRecordNotFound: config.ErrNotFound,
	gorm.ErrDuplicatedKey:  config.ErrConflict,
	// end of gorm
	// bcrypt
	bcrypt.ErrMismatchedHashAndPassword: config.ErrInvalidCredentials,
	// end of brcrypt
	// jwt
	jwt.ErrTokenExpired:          config.ErrInvalidToken,
	jwt.ErrTokenSignatureInvalid: config.ErrInvalidToken,
	// end of jwt
}

// map validation error message
var errorMessagesMap = map[string]string{
	"required":   "%s is required",
	"password":   "The password must be at least 8 characters long, 1 uppercase, 1 number and 1 special character",
	"email":      "Invalid email format",
	"after_now":  "Must be today or a future date",
	"oneof":      "Valid value : %e",
	"numeric":    "Only numeric value",
	"max":        "Maximum %e character",
	"min":        "Minimum %e character",
	"kyc_upload": "no_kk and nik is required if etc not provided",
}

func HandleError(c echo.Context, err error) error {
	// return if err is nil
	if err == nil {
		return nil
	}
	// return if error is system error (no need to translate)
	convert, ok := err.(config.Response)
	if ok {
		return c.JSON(convert.Status, convert)
	}
	// check if error is validation error
	convertErrs, ok := err.(validator.ValidationErrors)
	// create validation error (with field error)
	if ok {
		errs := []config.FieldError{}
		for _, convertErr := range convertErrs {
			field := convertErr.StructField()
			tag := convertErr.Tag()
			msg, ok := errorMessagesMap[tag]
			if !ok {
				errs = append(errs, config.FieldError{
					Field: field,
					Error: fmt.Sprintf("Failed on %s", convertErr.ActualTag()),
				})
				break
			}
			if strings.Contains(msg, "%s") {
				msg = fmt.Sprintf(msg, field)
			}
			if strings.Contains(msg, "%e") {
				msg = strings.Replace(msg, "%e", "%s", 1)
				param := strings.Join(strings.Split(convertErr.Param(), " "), ", ")
				msg = fmt.Sprintf(msg, param)
			}

			errs = append(errs, config.FieldError{
				Field: field,
				Error: msg,
			})
		}
		return c.JSON(http.StatusBadRequest, config.Response{
			Status:  http.StatusBadRequest,
			Message: "Validation Error",
			Data:    errs,
		})
	}
	// grab status and message from map if exists
	con, ok := errorMap[err]
	// if not exists in map then return 500
	if !ok {
		var debug *string = nil
		if config.ApplicationDebugStatus() {
			errMsg := err.Error()
			debug = &errMsg
		}
		return c.JSON(http.StatusInternalServerError, config.Response{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
			Debug:   debug,
		})
	}
	return c.JSON(con.Status, con)
}
