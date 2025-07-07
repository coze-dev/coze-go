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
	Rooms            *audioRooms
	Speech           *audioSpeech
	Voices           *audioVoices
	Transcriptions   *audioTranscriptions
	VoiceprintGroups *audioVoiceprintGroups
	Live             *audioLive
}

func newAudio(core *core) *audio {
	return &audio{
		Rooms:            newAudioRooms(core),
		Speech:           newAudioSpeech(core),
		Voices:           newAudioVoices(core),
		Transcriptions:   newAudioTranscriptions(core),
		VoiceprintGroups: newAudioVoiceprintGroups(core),
		Live:             newAudioLive(core),
	}
}
