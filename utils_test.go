package coze

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHTTP struct {
	Response *http.Response
	Error    error
}

func (m *mockHTTP) Do(*http.Request) (*http.Response, error) {
	return m.Response, m.Error
}

func Test_Ptr(t *testing.T) {
	as := assert.New(t)
	as.NotNil(ptr(1))
	as.NotNil(ptr("2"))
	as.NotNil(ptr(int64(3)))
	as.NotNil(ptr(int8(6)))
	as.NotNil(ptr(uint(7)))
	as.NotNil(ptr(uint32(9)))
	as.NotNil(ptr(float32(12.1)))
	as.NotNil(ptr(true))

	as.Equal(1, ptrValue(ptr(1)))
	as.Equal("2", ptrValue(ptr("2")))
	as.Equal(int64(3), ptrValue(ptr(int64(3))))
	as.Equal(uint(7), ptrValue(ptr(uint(7))))
	as.Equal(uint32(9), ptrValue(ptr(uint32(9))))
	as.Equal(float32(12.1), ptrValue(ptr(float32(12.1))))
	as.Equal(true, ptrValue(ptr(true)))
	var s *string
	as.Equal("", ptrValue(s))

	as.Nil(ptrNotZero(""))
	as.Nil(ptrNotZero(0))
	as.Nil(ptrNotZero(0.0))
	as.Nil(ptrNotZero(false))
	as.NotNil(ptrNotZero("1"))
}

func Test_GenerateRandomString(t *testing.T) {
	as := assert.New(t)
	str1, err := generateRandomString(10)
	as.Nil(err)
	str2, err := generateRandomString(10)
	as.Nil(err)
	as.NotEqual(str1, str2)
}

func Test_mustToJson(t *testing.T) {
	tests := []struct {
		name string
		args any
		want string
	}{
		{"success", map[string]string{"test": "test"}, `{"test":"test"}`},
		{"fail", func() {}, `{}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, mustToJson(tt.args), "mustToJson(%v)", tt.args)
		})
	}
}
