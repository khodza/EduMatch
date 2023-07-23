package custom_errors

import (
	"errors"
	"net/http"
	"strings"
)

// Map to represent the set of errors
var CustomErrors = map[string]int{
	// utils errors
	ErrHandleJSONBinding.Error(): http.StatusBadRequest,
	// eduCenter errors
	ErrEduCenterExist.Error():    http.StatusBadRequest,
	ErrEduCenterNotFound.Error(): http.StatusNotFound,
	// users errors
	ErrUserExist.Error():     http.StatusBadRequest,
	ErrUserNotFound.Error():  http.StatusNotFound,
	ErrWrongPassword.Error(): http.StatusBadRequest,
	// course errors
	ErrCourseNotFound.Error(): http.StatusNotFound,
	ErrCourseExists.Error():   http.StatusBadRequest,
	ErrInvalidToken.Error():   http.StatusUnauthorized,
	//auth
	ErrUnauthorized.Error(): http.StatusUnauthorized,
}

// utils errors
var (
	ErrHandleJSONBinding = errors.New("invalid request payload")
)

// eduCenter errors
var ErrEduCenterExist = errors.New("education center already exist")
var ErrEduCenterNotFound = errors.New("education center not found")

// users errors
var ErrUserExist = errors.New("user with this username or email already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrWrongPassword = errors.New("wrong password provided")

// course errors
var ErrCourseNotFound = errors.New("course not found")
var ErrCourseExists = errors.New("course is exists")

// auth errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrUnauthorized = errors.New("you are not allowed to this endpoint")
)

// validation(not handles as usual errors)
var ErrValidation = errors.New("validation failed")

func IsValidationErr(err string) bool {
	return strings.HasPrefix(err, ErrValidation.Error())
}
