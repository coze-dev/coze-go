package coze

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int
	Message string
	LogID   string
}

func NewError(code int, msg, logID string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		LogID:   logID,
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("code=%d, message=%s, logid=%s",
		e.Code,
		e.Message,
		e.LogID)
}

// AsCozeError checks if the error is of type Error
func AsCozeError(err error) (*Error, bool) {
	var cozeErr *Error
	if errors.As(err, &cozeErr) {
		return cozeErr, true
	}
	return nil, false
}

// authErrorFormat represents the error response from Coze API
type authErrorFormat struct {
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
	Error        string `json:"error"`
}

// AuthErrorCode represents authentication error codes
type AuthErrorCode string

const (
	/*
	 * The user has not completed authorization yet, please try again later
	 */
	AuthorizationPending AuthErrorCode = "authorization_pending"
	/*
	 * The request is too frequent, please try again later
	 */
	SlowDown AuthErrorCode = "slow_down"
	/*
	 * The user has denied the authorization
	 */
	AccessDenied AuthErrorCode = "access_denied"
	/*
	 * The token is expired
	 */
	ExpiredToken AuthErrorCode = "expired_token"
)

// String implements the Stringer interface
func (c AuthErrorCode) String() string {
	return string(c)
}

type AuthError struct {
	HttpCode     int
	Code         AuthErrorCode
	ErrorMessage string
	Param        string
	LogID        string
	parent       error
}

func NewCozeAuthExceptionWithoutParent(error *authErrorFormat, statusCode int, logID string) *AuthError {
	return &AuthError{
		HttpCode:     statusCode,
		ErrorMessage: error.ErrorMessage,
		Code:         AuthErrorCode(error.ErrorCode),
		Param:        error.Error,
		LogID:        logID,
	}
}

// Error implements the error interface
func (e *AuthError) Error() string {
	return fmt.Sprintf("HttpCode: %d, Code: %s, Message: %s, Param: %s, LogID: %s",
		e.HttpCode,
		e.Code,
		e.ErrorMessage,
		e.Param,
		e.LogID)
}

// Unwrap returns the parent error
func (e *AuthError) Unwrap() error {
	return e.parent
}

// AsAuthError 判断错误是否为 CozeAuthError 类型
func AsAuthError(err error) (*AuthError, bool) {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr, true
	}
	return nil, false
}
