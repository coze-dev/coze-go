package coze

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoiceConst(t *testing.T) {
	as := assert.New(t)

	t.Run("AudioFormat", func(t *testing.T) {
		as.Equal("mp3", AudioFormatMP3.String())
		as.NotNil(AudioFormatMP3.Ptr())
	})

	t.Run("LanguageCode", func(t *testing.T) {
		as.Equal("zh", LanguageCodeZH.String())
	})
}
