package coze

import (
	"net/http"
	"time"
)

type CozeAPI struct {
	Audio         *audio
	Bots          *bots
	Chat          *chat
	Conversations *conversations
	Workflows     *workflows
	Workspaces    *workspace
	Datasets      *datasets
	Files         *files
	Templates     *templates
	Users         *users
	baseURL       string
}

type clientOption struct {
	baseURL     string
	client      HTTPClient
	logLevel    LogLevel
	auth        Auth
	enableLogID bool
}

type CozeAPIOption func(*clientOption)

// WithBaseURL adds the base URL for the API
func WithBaseURL(baseURL string) CozeAPIOption {
	return func(opt *clientOption) {
		opt.baseURL = baseURL
	}
}

// WithHttpClient sets a custom HTTP core
func WithHttpClient(client HTTPClient) CozeAPIOption {
	return func(opt *clientOption) {
		opt.client = client
	}
}

// WithLogLevel sets the logging level
func WithLogLevel(level LogLevel) CozeAPIOption {
	return func(opt *clientOption) {
		opt.logLevel = level
	}
}

func WithLogger(logger Logger) CozeAPIOption {
	return func(opt *clientOption) {
		setLogger(logger)
	}
}

func WithEnableLogID(enableLogID bool) CozeAPIOption {
	return func(opt *clientOption) {
		opt.enableLogID = enableLogID
	}
}

func NewCozeAPI(auth Auth, opts ...CozeAPIOption) CozeAPI {
	opt := &clientOption{
		baseURL:  ComBaseURL,
		client:   nil,
		logLevel: LogLevelInfo, // Default log level is Info
		auth:     auth,
	}
	for _, option := range opts {
		option(opt)
	}
	if opt.client == nil {
		opt.client = &http.Client{
			Timeout: time.Second * 5,
		}
	}

	core := newCore(opt)
	setLevel(opt.logLevel)

	cozeClient := CozeAPI{
		Audio:         newAudio(core),
		Bots:          newBots(core),
		Chat:          newChats(core),
		Conversations: newConversations(core),
		Workflows:     newWorkflows(core),
		Workspaces:    newWorkspace(core),
		Datasets:      newDatasets(core),
		Files:         newFiles(core),
		Templates:     newTemplates(core),
		Users:         newUsers(core),
		baseURL:       opt.baseURL,
	}
	return cozeClient
}

type authTransport struct {
	auth Auth
	next http.RoundTripper
}

func (h *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if isAuthContext(req.Context()) {
		return h.next.RoundTrip(req)
	}
	accessToken, err := h.auth.Token(req.Context())
	if err != nil {
		logger.Errorf(req.Context(), "Failed to get access token: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return h.next.RoundTrip(req)
}
