package coze

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCozeError(t *testing.T) {
	as := assert.New(t)

	t.Run("new", func(t *testing.T) {
		err := NewError(1001, "test error", "test-log-id")
		as.NotNil(err)
		as.Equal(1001, err.Code)
		as.Equal("test error", err.Message)
		as.Equal("test-log-id", err.LogID)
	})

	t.Run("Error()", func(t *testing.T) {
		err := NewError(1001, "test error", "test-log-id")
		expectedMsg := "code=1001, message=test error, logid=test-log-id"
		as.Equal(expectedMsg, err.Error())
	})
}

func TestAsCozeError(t *testing.T) {
	as := assert.New(t)
	tests := []struct {
		name     string
		err      error
		wantErr  *Error
		wantBool bool
	}{
		{
			name:     "nil error",
			err:      nil,
			wantErr:  nil,
			wantBool: false,
		},
		{
			name:     "non-Error",
			err:      errors.New("standard error"),
			wantErr:  nil,
			wantBool: false,
		},
		{
			name:     "Error",
			err:      NewError(1001, "test error", "test-log-id"),
			wantErr:  NewError(1001, "test error", "test-log-id"),
			wantBool: true,
		},
		{
			name: "wrapped Error",
			err: fmt.Errorf("wrapped: %w",
				NewError(1001, "test error", "test-log-id")),
			wantErr:  NewError(1001, "test error", "test-log-id"),
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, gotBool := AsCozeError(tt.err)
			as.Equal(tt.wantBool, gotBool)
			if tt.wantErr != nil {
				as.Equal(tt.wantErr.Code, gotErr.Code)
				as.Equal(tt.wantErr.Message, gotErr.Message)
				as.Equal(tt.wantErr.LogID, gotErr.LogID)
			} else {
				as.Nil(gotErr)
			}
		})
	}
}

func TestAuthErrorCode_String(t *testing.T) {
	as := assert.New(t)
	tests := []struct {
		name string
		code AuthErrorCode
		want string
	}{
		{
			name: "AuthorizationPending",
			code: AuthorizationPending,
			want: "authorization_pending",
		},
		{
			name: "SlowDown",
			code: SlowDown,
			want: "slow_down",
		},
		{
			name: "AccessDenied",
			code: AccessDenied,
			want: "access_denied",
		},
		{
			name: "ExpiredToken",
			code: ExpiredToken,
			want: "expired_token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as.Equal(tt.want, tt.code.String())
		})
	}
}

func TestNewCozeAuthExceptionWithoutParent(t *testing.T) {
	as := assert.New(t)
	// 测试创建新的认证错误
	errorFormat := &authErrorFormat{
		ErrorMessage: "invalid token",
		ErrorCode:    "invalid_token",
		Error:        "token_error",
	}
	err := NewAuthError(errorFormat, 401, "test-log-id")

	as.NotNil(err)
	as.Equal(401, err.HttpCode)
	as.Equal("invalid token", err.ErrorMessage)
	as.Equal(AuthErrorCode("invalid_token"), err.Code)
	as.Equal("token_error", err.Param)
	as.Equal("test-log-id", err.LogID)
	as.Nil(err.parent)
}

func TestAuthError_Error(t *testing.T) {
	as := assert.New(t)
	err := &AuthError{
		HttpCode:     401,
		Code:         AuthErrorCode("invalid_token"),
		ErrorMessage: "invalid token",
		Param:        "token_error",
		LogID:        "test-log-id",
	}

	expectedMsg := "HttpCode: 401, Code: invalid_token, Message: invalid token, Param: token_error, LogID: test-log-id"
	as.Equal(expectedMsg, err.Error())
}

func TestAuthError_Unwrap(t *testing.T) {
	as := assert.New(t)
	t.Run("No Parent", func(t *testing.T) {
		err := &AuthError{}
		as.Nil(err.Unwrap())
	})

	// 测试有父错误的情况
	t.Run("With Parent", func(t *testing.T) {
		parentErr := errors.New("parent error")
		err := &AuthError{
			parent: parentErr,
		}
		as.Equal(parentErr, err.Unwrap())
	})
}

func TestAsAuthError(t *testing.T) {
	as := assert.New(t)
	tests := []struct {
		name     string
		err      error
		wantErr  *AuthError
		wantBool bool
	}{
		{
			name:     "nil error",
			err:      nil,
			wantErr:  nil,
			wantBool: false,
		},
		{
			name:     "non-AuthError",
			err:      errors.New("standard error"),
			wantErr:  nil,
			wantBool: false,
		},
		{
			name: "AuthError",
			err: &AuthError{
				HttpCode:     401,
				Code:         AuthErrorCode("invalid_token"),
				ErrorMessage: "invalid token",
				Param:        "token_error",
				LogID:        "test-log-id",
			},
			wantErr: &AuthError{
				HttpCode:     401,
				Code:         AuthErrorCode("invalid_token"),
				ErrorMessage: "invalid token",
				Param:        "token_error",
				LogID:        "test-log-id",
			},
			wantBool: true,
		},
		{
			name: "wrapped AuthError",
			err: fmt.Errorf("wrapped: %w", &AuthError{
				HttpCode:     401,
				Code:         AuthErrorCode("invalid_token"),
				ErrorMessage: "invalid token",
				Param:        "token_error",
				LogID:        "test-log-id",
			}),
			wantErr: &AuthError{
				HttpCode:     401,
				Code:         AuthErrorCode("invalid_token"),
				ErrorMessage: "invalid token",
				Param:        "token_error",
				LogID:        "test-log-id",
			},
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, gotBool := AsAuthError(tt.err)
			as.Equal(tt.wantBool, gotBool)
			if tt.wantErr != nil {
				as.Equal(tt.wantErr.HttpCode, gotErr.HttpCode)
				as.Equal(tt.wantErr.Code, gotErr.Code)
				as.Equal(tt.wantErr.ErrorMessage, gotErr.ErrorMessage)
				as.Equal(tt.wantErr.Param, gotErr.Param)
				as.Equal(tt.wantErr.LogID, gotErr.LogID)
			} else {
				as.Nil(gotErr)
			}
		})
	}
}
