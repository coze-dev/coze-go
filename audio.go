package coze

// AudioFormat represents the audio format type
type AudioFormat string

const (
	AudioFormatWAV     AudioFormat = "wav"
	AudioFormatPCM     AudioFormat = "pcm"
	AudioFormatOGGOPUS AudioFormat = "ogg_opus"
	AudioFormatM4A     AudioFormat = "m4a"
	AudioFormatAAC     AudioFormat = "aac"
	AudioFormatMP3     AudioFormat = "mp3"
)

func (f AudioFormat) String() string {
	return string(f)
}

func (f AudioFormat) Ptr() *AudioFormat {
	return &f
}

// LanguageCode represents the language code
type LanguageCode string

const (
	LanguageCodeZH LanguageCode = "zh"
	LanguageCodeEN LanguageCode = "en"
	LanguageCodeJA LanguageCode = "ja"
	LanguageCodeES LanguageCode = "es"
	LanguageCodeID LanguageCode = "id"
	LanguageCodePT LanguageCode = "pt"
)

func (l LanguageCode) String() string {
	return string(l)
}

type audio struct {
	Rooms          *audioRooms
	Speech         *audioSpeech
	Voices         *audioVoices
	Transcriptions *audioTranscriptions
}

func newAudio(core *core) *audio {
	return &audio{
		Rooms:          newRooms(core),
		Speech:         newSpeech(core),
		Voices:         newVoice(core),
		Transcriptions: newTranscriptions(core),
	}
}
