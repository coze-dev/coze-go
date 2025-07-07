# WebSocket Support for coze-go

The coze-go library now includes WebSocket support for real-time audio speech synthesis, audio transcription, and chat interactions.

## Features

- **Audio Speech**: Real-time text-to-speech conversion via WebSocket
- **Audio Transcriptions**: Real-time speech-to-text conversion via WebSocket
- **Chat**: Real-time chat interactions with bots, including tool calls and audio support

## Usage

### Basic Setup

```go
package main

import (
    "github.com/coze-dev/coze-go"
    "github.com/coze-dev/coze-go/websockets"
)

func main() {
    // Create authentication
    auth := coze.NewTokenAuth("your-api-token")
    
    // Create Coze API client
    client := coze.NewCozeAPI(auth, coze.WithBaseURL(coze.CnBaseURL))
    
    // Access WebSocket clients
    speechClient := client.WebSockets.Audio.Speech()
    transcriptionsClient := client.WebSockets.Audio.Transcriptions()
    chatClient := client.WebSockets.Chat.Create(websockets.WithBotID("your-bot-id"))
}
```

### Audio Speech (Text-to-Speech)

```go
// Create speech client
speechClient := client.WebSockets.Audio.Speech()

// Set up event handlers
handler := &websockets.SpeechEventHandler{
    OnSpeechCreated: func(event *websockets.SpeechCreatedEvent) error {
        fmt.Printf("Speech session created: %s\n", event.Data.SessionID)
        return nil
    },
    OnSpeechAudioUpdate: func(event *websockets.SpeechAudioUpdateEvent) error {
        // Decode audio data
        audioData, err := websockets.GetAudioFromDelta(event.Data.Delta)
        if err != nil {
            return err
        }
        fmt.Printf("Received audio chunk: %d bytes\n", len(audioData))
        return nil
    },
    OnSpeechAudioCompleted: func(event *websockets.SpeechAudioCompletedEvent) error {
        fmt.Println("Speech synthesis completed")
        return nil
    },
}

// Register handlers
handler.RegisterHandlers(speechClient)

// Connect
if err := speechClient.Connect(); err != nil {
    log.Fatal(err)
}
defer speechClient.Close()

// Send text for synthesis
if err := speechClient.AppendTextBuffer("Hello, world!"); err != nil {
    log.Fatal(err)
}

if err := speechClient.CompleteTextBuffer(); err != nil {
    log.Fatal(err)
}

// Wait for completion
event, err := speechClient.WaitForSpeechAudioCompleted(30 * time.Second)
if err != nil {
    log.Fatal(err)
}
```

### Audio Transcriptions (Speech-to-Text)

```go
// Create transcriptions client
transcriptionsClient := client.WebSockets.Audio.Transcriptions()

// Set up event handlers
handler := &websockets.TranscriptionsEventHandler{
    OnTranscriptionsMessageUpdate: func(event *websockets.TranscriptionsMessageUpdateEvent) error {
        fmt.Printf("Transcription: %s\n", event.Data.Content)
        return nil
    },
    OnTranscriptionsMessageCompleted: func(event *websockets.WebSocketEvent) error {
        fmt.Println("Transcription completed")
        return nil
    },
}

// Register handlers
handler.RegisterHandlers(transcriptionsClient)

// Connect
if err := transcriptionsClient.Connect(); err != nil {
    log.Fatal(err)
}
defer transcriptionsClient.Close()

// Send audio data
audioData := []byte("your-audio-data-here")
if err := transcriptionsClient.AppendAudioBuffer(audioData); err != nil {
    log.Fatal(err)
}

if err := transcriptionsClient.CompleteAudioBuffer(); err != nil {
    log.Fatal(err)
}

// Wait for completion
event, err := transcriptionsClient.WaitForTranscriptionCompleted(30 * time.Second)
if err != nil {
    log.Fatal(err)
}
```

### Chat

```go
// Create chat client
chatClient := client.WebSockets.Chat.Create(
    websockets.WithBotID("your-bot-id"),
)

// Set up event handlers
handler := &websockets.ChatEventHandler{
    OnConversationChatCreated: func(event *websockets.ConversationChatCreatedEvent) error {
        fmt.Printf("Chat created: %s\n", event.Data.ChatID)
        return nil
    },
    OnConversationMessageDelta: func(event *websockets.ConversationMessageDeltaEvent) error {
        fmt.Printf("Message: %s\n", event.Data.Content)
        return nil
    },
    OnConversationChatRequiresAction: func(event *websockets.ConversationChatRequiresActionEvent) error {
        // Handle tool calls
        for _, toolCall := range event.Data.RequiredAction.SubmitToolOutputs.ToolCalls {
            // Execute tool and get result
            result := executeMyTool(toolCall.Function.Name, toolCall.Function.Arguments)
            
            // Submit result
            toolOutputs := []websockets.ToolOutput{
                {
                    ToolCallID: toolCall.ID,
                    Output:     result,
                },
            }
            
            return chatClient.SubmitToolOutputs(event.Data.ChatID, toolOutputs)
        }
        return nil
    },
}

// Register handlers
handler.RegisterHandlers(chatClient)

// Connect
if err := chatClient.Connect(); err != nil {
    log.Fatal(err)
}
defer chatClient.Close()

// Send a message
if err := chatClient.CreateMessage("Hello, how are you?"); err != nil {
    log.Fatal(err)
}

// Or send audio
audioData := []byte("your-audio-data")
if err := chatClient.AppendAudioBuffer(audioData); err != nil {
    log.Fatal(err)
}
if err := chatClient.CompleteAudioBuffer(); err != nil {
    log.Fatal(err)
}

// Wait for completion
event, err := chatClient.WaitForChatCompleted(60 * time.Second)
if err != nil {
    log.Fatal(err)
}
```

## Configuration Options

### Audio Configuration

```go
// Input audio configuration
inputAudio := &websockets.InputAudio{
    Format:     "wav",
    Codec:      "pcm",
    SampleRate: 24000,
    Channel:    1,
    BitDepth:   16,
}

// Output audio configuration
outputAudio := &websockets.OutputAudio{
    Codec:      "pcm",
    SpeechRate: 0,
    VoiceID:    "your-voice-id",
    PCMConfig: &websockets.PCMConfig{
        SampleRate: 24000,
    },
}
```

### Client Options

```go
// Speech client with options
speechClient := client.WebSockets.Audio.Speech(
    websockets.WithOutputAudio(outputAudio),
)

// Transcriptions client with options
transcriptionsClient := client.WebSockets.Audio.Transcriptions(
    websockets.WithInputAudio(inputAudio),
)

// Chat client with options
chatClient := client.WebSockets.Chat.Create(
    websockets.WithBotID("bot-id"),
    websockets.WithChatInputAudio(inputAudio),
    websockets.WithChatOutputAudio(outputAudio),
)
```

## Event Types

### Speech Events
- `speech.created` - Speech session created
- `speech.updated` - Speech configuration updated
- `input_text_buffer.completed` - Text buffer submission completed
- `speech.audio.update` - Audio chunk received
- `speech.audio.completed` - Speech synthesis completed

### Transcription Events
- `transcriptions.created` - Transcription session created
- `transcriptions.updated` - Transcription configuration updated
- `input_audio_buffer.completed` - Audio buffer submission completed
- `transcriptions.message.update` - Transcription result received
- `transcriptions.message.completed` - Transcription completed

### Chat Events
- `chat.created` - Chat session created
- `conversation.chat.created` - Conversation started
- `conversation.message.delta` - Message content received
- `conversation.audio.delta` - Audio content received
- `conversation.chat.completed` - Chat completed
- `conversation.chat.requires_action` - Tool call required

## Examples

See the `examples/` directory for complete working examples:
- `examples/websockets_audio_speech.go` - Text-to-speech example
- `examples/websockets_audio_transcriptions.go` - Speech-to-text example
- `examples/websockets_chat.go` - Chat example with tool calls

## Error Handling

All WebSocket clients support error handling through event handlers:

```go
handler := &websockets.SpeechEventHandler{
    OnError: func(err error) error {
        log.Printf("WebSocket error: %v", err)
        return nil
    },
    OnClosed: func() error {
        log.Println("Connection closed")
        return nil
    },
}
```

## Thread Safety

The WebSocket clients are designed to be thread-safe. You can safely call methods from multiple goroutines.

## Environment Variables

The examples use these environment variables:
- `COZE_API_TOKEN` - Your Coze API token
- `COZE_API_BASE` - API base URL (optional, defaults to CN)
- `COZE_BOT_ID` - Bot ID for chat examples