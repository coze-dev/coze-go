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
	as.NotNil(ptr(int32(4)))
	as.NotNil(ptr(int16(5)))
	as.NotNil(ptr(int8(6)))
	as.NotNil(ptr(uint(7)))
	as.NotNil(ptr(uint64(8)))
	as.NotNil(ptr(uint32(9)))
	as.NotNil(ptr(uint16(10)))
	as.NotNil(ptr(uint8(11)))
	as.NotNil(ptr(float32(12.1)))
	as.NotNil(ptr(float64(13.1)))
	as.NotNil(ptr(true))
	as.NotNil(ptr(false))

	as.Equal(1, ptrValue(ptr(1)))
	as.Equal("2", ptrValue(ptr("2")))
	as.Equal(int64(3), ptrValue(ptr(int64(3))))
	as.Equal(int32(4), ptrValue(ptr(int32(4))))
	as.Equal(int16(5), ptrValue(ptr(int16(5))))
	as.Equal(int8(6), ptrValue(ptr(int8(6))))
	as.Equal(uint(7), ptrValue(ptr(uint(7))))
	as.Equal(uint64(8), ptrValue(ptr(uint64(8))))
	as.Equal(uint32(9), ptrValue(ptr(uint32(9))))
	as.Equal(uint16(10), ptrValue(ptr(uint16(10))))
	as.Equal(uint8(11), ptrValue(ptr(uint8(11))))
	as.Equal(float32(12.1), ptrValue(ptr(float32(12.1))))
	as.Equal(float64(13.1), ptrValue(ptr(float64(13.1))))
	as.Equal(true, ptrValue(ptr(true)))
	as.Equal(false, ptrValue(ptr(false)))
	var s *string
	as.Equal("", ptrValue(s))
}

func Test_GenerateRandomString(t *testing.T) {
	str1, err := generateRandomString(10)
	assert.Nil(t, err)
	str2, err := generateRandomString(10)
	assert.Nil(t, err)
	assert.NotEqual(t, str1, str2)
}

func Test_MustToJson(t *testing.T) {
	jsonStr := mustToJson(map[string]string{"test": "test"})
	assert.Equal(t, jsonStr, `{"test":"test"}`)
}
