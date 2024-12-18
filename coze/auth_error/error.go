package auth_error

import (
	"errors"
	"fmt"
)

// CozeError represents the error response from Coze API
type CozeError struct {
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
func (c *AuthErrorCode) String() string {
	return string(*c)
}

type CozeAuthError struct {
	HttpCode     int
	Code         AuthErrorCode
	ErrorMessage string
	Param        string
	LogID        string
	parent       error
}

func NewCozeAuthExceptionWithoutParent(error *CozeError, statusCode int, logID string) *CozeAuthError {
	return &CozeAuthError{
		HttpCode:     statusCode,
		ErrorMessage: error.ErrorMessage,
		Code:         AuthErrorCode(error.ErrorCode),
		Param:        error.Error,
		LogID:        logID,
	}
}

// Error implements the error interface
func (e *CozeAuthError) Error() string {
	return fmt.Sprintf("HttpCode: %d, Code: %s, ErrorMessage: %s, Param: %s, LogID: %s",
		e.HttpCode,
		e.Code,
		e.ErrorMessage,
		e.Param,
		e.LogID)
}

// Unwrap returns the parent error
func (e *CozeAuthError) Unwrap() error {
	return e.parent
}

// AsCozeAuthError 判断错误是否为 CozeAuthError 类型
func AsCozeAuthError(err error) (*CozeAuthError, bool) {
	var authErr *CozeAuthError
	if errors.As(err, &authErr) {
		return authErr, true
	}
	return nil, false
}