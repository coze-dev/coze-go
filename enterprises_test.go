package coze

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnterprises(t *testing.T) {
	as := assert.New(t)

	t.Run("new enterprises", func(t *testing.T) {
		core := newCore(&clientOption{})
		enterprises := newEnterprises(core)
		as.NotNil(enterprises)
		as.NotNil(enterprises.core)
		as.NotNil(enterprises.Members)
	})
}
