package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

func pcmWriteToWavFile(file string, audioPCMData []byte) error {
	outFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// WAV 文件头信息
	var (
		chunkID       = []byte{'R', 'I', 'F', 'F'}
		format        = []byte{'W', 'A', 'V', 'E'}
		subchunk1ID   = []byte{'f', 'm', 't', ' '}
		subchunk1Size = uint32(16) // PCM
		audioFormat   = uint16(1)  // PCM = 1 (线性量化)
		numChannels   = uint16(1)  // Mono = 1, Stereo = 2
		sampleRate    = uint32(24000)
		byteRate      = sampleRate * uint32(numChannels) * uint32(audioFormat) // SampleRate * NumChannels * BitsPerSample/8
		blockAlign    = numChannels * uint16(audioFormat)
		bitsPerSample = uint16(16)
		subchunk2ID   = []byte{'d', 'a', 't', 'a'}
	)

	// 预留空间写入 ChunkSize 和 Subchunk2Size
	if _, err := outFile.Seek(44, 0); err != nil {
		return err
	}

	// 模拟音频数据
	for i := 0; i < len(audioPCMData)-1; i += 2 {
		err := binary.Write(outFile, binary.LittleEndian, audioPCMData[i:i+2])
		if err != nil {
			return err
		}
	}

	// 获取文件大小
	fileInfo, err := outFile.Stat()
	if err != nil {
		return err
	}

	// 计算 ChunkSize 和 Subchunk2Size
	fileSize := fileInfo.Size()
	chunkSize := uint32(fileSize - 8)
	subchunk2Size := uint32(fileSize - 44)

	// 回写 WAV 文件头
	if _, err := outFile.Seek(0, 0); err != nil {
		return err
	}

	headers := [][]interface{}{
		{chunkID, chunkSize, format},
		{subchunk1ID, subchunk1Size, audioFormat, numChannels, sampleRate, byteRate, blockAlign, bitsPerSample},
		{subchunk2ID, subchunk2Size},
	}

	for _, headerSection := range headers {
		for _, headerField := range headerSection {
			if err := binary.Write(outFile, binary.LittleEndian, headerField); err != nil {
				return err
			}
		}
	}
	return nil
}

type handler struct {
	coze.BaseWebSocketAudioSpeechHandler
	data []byte
}

func (r *handler) OnClientError(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketClientErrorEvent) error {
	log.Printf("speech client_error: %v", event)
	return nil
}

func (r *handler) OnError(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketErrorEvent) error {
	log.Printf("speech error: %v", event)
	return nil
}

func (r *handler) OnSpeechAudioUpdate(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketSpeechAudioUpdateEvent) error {
	r.data = append(r.data, event.Data.Delta...)
	return nil
}

func (r *handler) OnSpeechAudioCompleted(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketSpeechAudioCompletedEvent) error {
	err := pcmWriteToWavFile("output.wav", r.data)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return err
	}

	log.Printf("speech completed, audio write to %s", "output.wav")
	return nil
}

func main() {
	cozeAPIToken := os.Getenv("COZE_API_TOKEN")
	cozeAPIBase := os.Getenv("COZE_API_BASE")
	if cozeAPIBase == "" {
		cozeAPIBase = coze.CnBaseURL
	}

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(cozeAPIBase), coze.WithLogLevel(coze.LogLevelDebug))

	// Create speech WebSocket client
	speechClient := client.WebSockets.Audio.Speech.Create(context.Background(), &coze.CreateWebsocketAudioSpeechReq{})
	speechClient.RegisterHandler(&handler{})

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := speechClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer speechClient.Close()

	// Send text to be converted to speech
	text := "今天天气不错"
	fmt.Printf("Sending text: %s\n", text)

	if err := speechClient.InputTextBufferAppend(&coze.WebSocketInputTextBufferAppendEventData{
		Delta: text,
	}); err != nil {
		log.Fatalf("Failed to append text: %v", err)
	}

	if err := speechClient.InputTextBufferComplete(nil); err != nil {
		log.Fatalf("Failed to complete text buffer: %v", err)
	}

	// time.Sleep(time.Hour)
	// Wait for speech completion
	fmt.Println("Waiting for speech completion...")
	event, err := speechClient.Wait(300000 * time.Second)
	if err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}
	fmt.Printf("Speech completed! Event: %+v\n", event)
}
