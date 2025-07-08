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
	Folders       *folders
	Templates     *templates
	Users         *users
	Variables     *variables
	Enterprises   *enterprises
	baseURL       string
}

type clientOption struct {
	baseURL     string
	client      HTTPClient
	logLevel    LogLevel
	logger      Logger
	auth        Auth
	enableLogID bool
	headers     http.Header
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
		opt.logger = logger
		setLogger(logger)
	}
}

func WithEnableLogID(enableLogID bool) CozeAPIOption {
	return func(opt *clientOption) {
		opt.enableLogID = enableLogID
	}
}

func WithHeaders(headers http.Header) CozeAPIOption {
	return func(opt *clientOption) {
		opt.headers = headers
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
		Folders:       newFolders(core),
		Templates:     newTemplates(core),
		Users:         newUsers(core),
		Variables:     newVariables(core),
		Enterprises:   newEnterprises(core),
		baseURL:       opt.baseURL,
	}
	return cozeClient
}

type core struct {
	*clientOption
}

func newCore(opt *clientOption) *core {
	if opt.client == nil {
		opt.client = &http.Client{
			Timeout: time.Second * 5,
		}
	}
	return &core{
		clientOption: opt,
	}
}
