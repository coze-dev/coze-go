package coze

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCozeError(t *testing.T) {
	// 测试创建新的 Error
	err := NewError(1001, "test error", "test-log-id")
	assert.NotNil(t, err)
	assert.Equal(t, 1001, err.Code)
	assert.Equal(t, "test error", err.Message)
	assert.Equal(t, "test-log-id", err.LogID)
}

func TestCozeError_Error(t *testing.T) {
	// 测试 Error() 方法
	err := NewError(1001, "test error", "test-log-id")
	expectedMsg := "code=1001, message=test error, logid=test-log-id"
	assert.Equal(t, expectedMsg, err.Error())
}

func TestAsCozeError(t *testing.T) {
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
			assert.Equal(t, tt.wantBool, gotBool)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Code, gotErr.Code)
				assert.Equal(t, tt.wantErr.Message, gotErr.Message)
				assert.Equal(t, tt.wantErr.LogID, gotErr.LogID)
			} else {
				assert.Nil(t, gotErr)
			}
		})
	}
}

func TestAuthErrorCode_String(t *testing.T) {
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
			assert.Equal(t, tt.want, tt.code.String())
		})
	}
}

func TestNewCozeAuthExceptionWithoutParent(t *testing.T) {
	// 测试创建新的认证错误
	errorFormat := &authErrorFormat{
		ErrorMessage: "invalid token",
		ErrorCode:    "invalid_token",
		Error:        "token_error",
	}
	err := NewCozeAuthExceptionWithoutParent(errorFormat, 401, "test-log-id")

	assert.NotNil(t, err)
	assert.Equal(t, 401, err.HttpCode)
	assert.Equal(t, "invalid token", err.ErrorMessage)
	assert.Equal(t, AuthErrorCode("invalid_token"), err.Code)
	assert.Equal(t, "token_error", err.Param)
	assert.Equal(t, "test-log-id", err.LogID)
	assert.Nil(t, err.parent)
}

func TestAuthError_Error(t *testing.T) {
	// 测试 Error() 方法
	err := &AuthError{
		HttpCode:     401,
		Code:         AuthErrorCode("invalid_token"),
		ErrorMessage: "invalid token",
		Param:        "token_error",
		LogID:        "test-log-id",
	}

	expectedMsg := "HttpCode: 401, Code: invalid_token, Message: invalid token, Param: token_error, LogID: test-log-id"
	assert.Equal(t, expectedMsg, err.Error())
}

func TestAuthError_Unwrap(t *testing.T) {
	// 测试无父错误的情况
	t.Run("No Parent", func(t *testing.T) {
		err := &AuthError{}
		assert.Nil(t, err.Unwrap())
	})

	// 测试有父错误的情况
	t.Run("With Parent", func(t *testing.T) {
		parentErr := errors.New("parent error")
		err := &AuthError{
			parent: parentErr,
		}
		assert.Equal(t, parentErr, err.Unwrap())
	})
}

func TestAsAuthError(t *testing.T) {
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
			assert.Equal(t, tt.wantBool, gotBool)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.HttpCode, gotErr.HttpCode)
				assert.Equal(t, tt.wantErr.Code, gotErr.Code)
				assert.Equal(t, tt.wantErr.ErrorMessage, gotErr.ErrorMessage)
				assert.Equal(t, tt.wantErr.Param, gotErr.Param)
				assert.Equal(t, tt.wantErr.LogID, gotErr.LogID)
			} else {
				assert.Nil(t, gotErr)
			}
		})
	}
}
