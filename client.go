package coze

import (
	"net/http"
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
	baseURL       string
}

type newCozeAPIOpt struct {
	baseURL  string
	client   *http.Client
	logLevel LogLevel
}

type CozeAPIOption func(*newCozeAPIOpt)

// WithBaseURL adds the base URL for the API
func WithBaseURL(baseURL string) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		opt.baseURL = baseURL
	}
}

// WithHttpClient sets a custom HTTP core
func WithHttpClient(client *http.Client) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		opt.client = client
	}
}

// WithLogLevel sets the logging level
func WithLogLevel(level LogLevel) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		opt.logLevel = level
	}
}

func WithLogger(logger Logger) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		setLogger(logger)
	}
}

func NewCozeAPI(auth Auth, opts ...CozeAPIOption) CozeAPI {
	opt := &newCozeAPIOpt{
		baseURL:  ComBaseURL,
		logLevel: LogLevelInfo, // Default log level is Info
	}
	for _, option := range opts {
		option(opt)
	}
	if opt.client == nil {
		opt.client = http.DefaultClient
	}
	saveTransport := opt.client.Transport
	if saveTransport == nil {
		saveTransport = http.DefaultTransport
	}
	opt.client.Transport = &authTransport{
		auth: auth,
		next: saveTransport,
	}
	core := newCore(opt.client, opt.baseURL)
	setLevel(opt.logLevel)
	// Set log level
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
		baseURL:       opt.baseURL,
	}
	return cozeClient
}

type authTransport struct {
	auth Auth
	next http.RoundTripper
}

func (h *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	accessToken, err := h.auth.Token(req.Context())
	if err != nil {
		logger.Errorf(req.Context(), "Failed to get access token: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return h.next.RoundTrip(req)
}
