package coze

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWebSocketEvent(t *testing.T) {
	as := assert.New(t)
	tests := []struct {
		name      string
		message   []byte
		wantType  reflect.Type
		wantErr   bool
		checkData func(t *testing.T, event IWebSocketEvent)
	}{
		{
			name:     "Parse SpeechCreated Event",
			message:  []byte(`{"event_type": "speech.created", "id": "123"}`),
			wantType: reflect.TypeOf(&WebSocketSpeechCreatedEvent{}),
			wantErr:  false,
			checkData: func(t *testing.T, event IWebSocketEvent) {
				e := event.(*WebSocketSpeechCreatedEvent)
				as.Equal("123", e.ID)
			},
		},
		{
			name:     "Parse ChatMessageDelta Event",
			message:  []byte(`{"event_type": "conversation.message.delta", "id": "456", "data": {"content": "hello"}}`),
			wantType: reflect.TypeOf(&WebSocketConversationMessageDeltaEvent{}),
			wantErr:  false,
			checkData: func(t *testing.T, event IWebSocketEvent) {
				e := event.(*WebSocketConversationMessageDeltaEvent)
				as.Equal("456", e.ID)
				as.Equal("hello", e.Data.Content)
			},
		},
		{
			name:     "Unknown Event Type",
			message:  []byte(`{"event_type": "unknown.event", "id": "789"}`),
			wantType: reflect.TypeOf(&commonWebSocketEvent{}),
			wantErr:  false,
			checkData: func(t *testing.T, event IWebSocketEvent) {
				e := event.(*commonWebSocketEvent)
				as.Equal("unknown.event", string(e.GetEventType()))
			},
		},
		{
			name:    "Invalid JSON",
			message: []byte(`{"event_type": "speech.created"`), // Malformed JSON
			wantErr: true,
		},
		{
			name:    "Mismatched Data Structure",
			message: []byte(`{"event_type": "speech.created", "data": {"invalid_field": "value"}}`),
			wantErr: false, // json.Unmarshal is lenient with extra fields
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseWebSocketEvent(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseWebSocketEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if tt.wantType != nil && reflect.TypeOf(got) != tt.wantType {
				t.Errorf("parseWebSocketEvent() got type = %v, want type %v", reflect.TypeOf(got), tt.wantType)
			}

			if tt.checkData != nil {
				tt.checkData(t, got)
			}
		})
	}
}

func TestBaseWebSocketEvent_Getters(t *testing.T) {
	event := baseWebSocketEvent{
		EventType: WebSocketEventTypeChatCreated,
		ID:        "test-id",
		Detail: &EventDetail{
			LogID: "log-id",
		},
	}

	if event.GetEventType() != WebSocketEventTypeChatCreated {
		t.Errorf("GetEventType() = %v, want %v", event.GetEventType(), WebSocketEventTypeChatCreated)
	}

	if event.GetID() != "test-id" {
		t.Errorf("GetID() = %v, want %v", event.GetID(), "test-id")
	}

	if event.GetDetail().LogID != "log-id" {
		t.Errorf("GetDetail().LogID = %v, want %v", event.GetDetail().LogID, "log-id")
	}
}
