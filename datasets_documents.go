package coze

import (
	"context"
	"encoding/base64"
	"net/http"
)

func (r *datasetsDocuments) Create(ctx context.Context, req *CreateDatasetsDocumentsReq) (*CreateDatasetsDocumentsResp, error) {
	method := http.MethodPost
	uri := "/open_api/knowledge/document/create"
	resp := &createDatasetsDocumentsResp{}
	err := r.client.Request(ctx, method, uri, req, resp, r.commonHeaderOpt...)
	if err != nil {
		return nil, err
	}
	resp.CreateDatasetsDocumentsResp.setHTTPResponse(resp.HTTPResponse)
	return resp.CreateDatasetsDocumentsResp, nil
}

func (r *datasetsDocuments) Update(ctx context.Context, req *UpdateDatasetsDocumentsReq) (*UpdateDatasetsDocumentsResp, error) {
	method := http.MethodPost
	uri := "/open_api/knowledge/document/update"
	resp := &updateDatasetsDocumentsResp{}
	err := r.client.Request(ctx, method, uri, req, resp, r.commonHeaderOpt...)
	if err != nil {
		return nil, err
	}
	result := &UpdateDatasetsDocumentsResp{}
	result.setHTTPResponse(resp.HTTPResponse)
	return result, nil
}

func (r *datasetsDocuments) Delete(ctx context.Context, req *DeleteDatasetsDocumentsReq) (*DeleteDatasetsDocumentsResp, error) {
	method := http.MethodPost
	uri := "/open_api/knowledge/document/delete"
	resp := &deleteDatasetsDocumentsResp{}
	err := r.client.Request(ctx, method, uri, req, resp, r.commonHeaderOpt...)
	if err != nil {
		return nil, err
	}
	result := &DeleteDatasetsDocumentsResp{}
	result.setHTTPResponse(resp.HTTPResponse)
	return result, nil
}

func (r *datasetsDocuments) List(ctx context.Context, req *ListDatasetsDocumentsReq) (*NumberPaged[Document], error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}
	return NewNumberPaged[Document](
		func(request *PageRequest) (*PageResponse[Document], error) {
			uri := "/open_api/knowledge/document/list"
			resp := &listDatasetsDocumentsResp{}
			doReq := &ListDatasetsDocumentsReq{
				DatasetID: req.DatasetID,
				Size:      request.PageSize,
				Page:      request.PageNum,
			}
			err := r.client.Request(ctx, http.MethodPost, uri, doReq, resp, r.commonHeaderOpt...)
			if err != nil {
				return nil, err
			}
			return &PageResponse[Document]{
				Total:   int(resp.Total),
				HasMore: request.PageSize <= len(resp.DocumentInfos),
				Data:    resp.DocumentInfos,
				LogID:   resp.HTTPResponse.LogID(),
			}, nil
		}, req.Size, req.Page)
}

type datasetsDocuments struct {
	client          *core
	commonHeaderOpt []RequestOption
}

func newDocuments(core *core) *datasetsDocuments {
	return &datasetsDocuments{client: core, commonHeaderOpt: []RequestOption{
		withHTTPHeader("Agw-Js-Conv", "str"),
	}}
}

// Document represents a document in the datasets
type Document struct {
	// The ID of the file.
	DocumentID string `json:"document_id"`

	// The total character count of the file content.
	CharCount int `json:"char_count"`

	// The chunking rules. For detailed instructions, refer to the ChunkStrategy object.
	ChunkStrategy *DocumentChunkStrategy `json:"chunk_strategy"`

	// The upload time of the file, in the format of a 10-digit Unix timestamp.
	CreateTime int `json:"create_time"`

	// The last modified time of the file, in the format of a 10-digit Unix timestamp.
	UpdateTime int `json:"update_time"`

	// The type of file format. Values include:
	// 0: Document type, such as txt, pdf, online web pages, etc.
	// 1: Spreadsheet type, such as xls spreadsheets, etc.
	// 2: Photo type, such as png images, etc.
	FormatType DocumentFormatType `json:"format_type"`

	// The number of times the file has been hit in conversations.
	HitCount int `json:"hit_count"`

	// The name of the file.
	Name string `json:"name"`

	// The size of the file in bytes.
	Size int `json:"size"`

	// The number of slices the file has been divided into.
	SliceCount int `json:"slice_count"`

	// The method of uploading the file. Values include:
	// 0: Upload local files.
	// 1: Upload online web pages.
	SourceType DocumentSourceType `json:"source_type"`

	// The processing status of the file. Values include:
	// 0: Processing
	// 1: Completed
	// 9: Processing failed, it is recommended to re-upload
	Status DocumentStatus `json:"status"`

	// The format of the local file, i.e., the file extension, such as "txt".
	// Supported formats include PDF, TXT, DOC, DOCX.
	Type string `json:"type"`

	// The frequency of automatic updates for online web pages, in hours.
	UpdateInterval int `json:"update_interval"`

	// Whether the online web page is automatically updated. Values include:
	// 0: Do not automatically update
	// 1: Automatically update
	UpdateType DocumentUpdateType `json:"update_type"`
}

// DocumentBase represents base information for creating a document
type DocumentBase struct {
	// The name of the file.
	Name string `json:"name"`

	// The metadata information of the file.
	SourceInfo *DocumentSourceInfo `json:"source_info"`

	// The update strategy for online web pages. Defaults to no automatic update.
	UpdateRule *DocumentUpdateRule `json:"update_rule,omitempty"`
}

// DocumentChunkStrategy represents chunking strategy for datasetsDocuments
type DocumentChunkStrategy struct {
	// The chunking settings. Values include:
	// 0: Automatic chunking and cleaning. Uses preset rules for data chunking and processing.
	// 1: Custom. In this case, details need to be specified through separator, max_tokens,
	// remove_extra_spaces, and remove_urls_emails.
	ChunkType int `json:"chunk_type"`

	// Maximum chunk length, ranging from 100 to 2000.
	// Required when chunk_type=1.
	MaxTokens int `json:"max_tokens,omitempty"`

	// Whether to automatically filter consecutive spaces, line breaks, and tabs.
	// Values include:
	// true: Automatically filter
	// false: (Default) Do not automatically filter
	// Takes effect when chunk_type=1.
	RemoveExtraSpaces bool `json:"remove_extra_spaces,omitempty"`

	// Whether to automatically filter all URLs and email addresses.
	// Values include:
	// true: Automatically filter
	// false: (Default) Do not automatically filter
	// Takes effect when chunk_type=1.
	RemoveUrlsEmails bool `json:"remove_urls_emails,omitempty"`

	// The chunk identifier.
	// Required when chunk_type=1.
	Separator string `json:"separator,omitempty"`
}

// DocumentSourceInfo represents source information for a document
type DocumentSourceInfo struct {
	// Base64 encoding of the local file.
	// Required when uploading local files.
	FileBase64 string `json:"file_base64,omitempty"`

	// The format of the local file, i.e., the file extension, such as "txt".
	// Supported formats include PDF, TXT, DOC, DOCX.
	// The uploaded file type should match the knowledge base type.
	// Required when uploading local files.
	FileType string `json:"file_type,omitempty"`

	// The URL of the webpage.
	// Required when uploading webpages.
	WebUrl string `json:"web_url,omitempty"`

	// The upload method of the file.
	// Set to 1 to indicate uploading online webpages.
	// Required when uploading online webpages.
	DocumentSource int `json:"document_source,omitempty"`
}

// DocumentUpdateRule represents update rules for datasetsDocuments
type DocumentUpdateRule struct {
	// Whether the online webpage is automatically updated.
	// Values include:
	// 0: Do not automatically update
	// 1: Automatically update
	UpdateType DocumentUpdateType `json:"update_type"`

	// The frequency of automatic updates for online webpages, in hours.
	// Minimum value is 24.
	UpdateInterval int `json:"update_interval"`
}

// DocumentFormatType represents the format type of a document
type DocumentFormatType int

const (
	// Document type, such as txt, pdf, online web pages, etc.
	DocumentFormatTypeDocument DocumentFormatType = 0
	// Spreadsheet type, such as xls spreadsheets, etc.
	DocumentFormatTypeSpreadsheet DocumentFormatType = 1
	// Photo type, such as png images, etc.
	DocumentFormatTypeImage DocumentFormatType = 2
)

// DocumentSourceType represents the source type of a document
type DocumentSourceType int

const (
	// Upload local files.
	DocumentSourceTypeLocalFile DocumentSourceType = 0
	// Upload online web pages.
	DocumentSourceTypeOnlineWeb DocumentSourceType = 1
)

// DocumentStatus represents the status of a document
type DocumentStatus int

const (
	// Processing
	DocumentStatusProcessing DocumentStatus = 0
	// Completed
	DocumentStatusCompleted DocumentStatus = 1
	// Processing failed, it is recommended to re-upload
	DocumentStatusFailed DocumentStatus = 9
)

// DocumentUpdateType represents the update type of a document
type DocumentUpdateType int

const (
	// Do not automatically update
	DocumentUpdateTypeNoAutoUpdate DocumentUpdateType = 0
	// Automatically update
	DocumentUpdateTypeAutoUpdate DocumentUpdateType = 1
)

// CreateDatasetsDocumentsReq represents request for creating document
type CreateDatasetsDocumentsReq struct {
	// The ID of the knowledge base.
	DatasetID int64 `json:"dataset_id"`

	// The metadata information of the files awaiting upload. The array has a maximum length of 10,
	// meaning up to 10 files can be uploaded at a time. For detailed instructions, refer to the
	// DocumentBase object.
	DocumentBases []*DocumentBase `json:"document_bases"`

	// Chunk strategy. These rules must be set only when uploading a file to new knowledge for the
	// first time. For subsequent file uploads to this knowledge, it is not necessary to pass these
	// rules; the default is to continue using the initial settings, and modifications are not
	// supported. For detailed instructions, refer to the ChunkStrategy object.
	ChunkStrategy *DocumentChunkStrategy `json:"chunk_strategy,omitempty"`
}

// DeleteDatasetsDocumentsReq represents request for deleting datasetsDocuments
type DeleteDatasetsDocumentsReq struct {
	DocumentIDs []int64 `json:"document_ids"`
}

// ListDatasetsDocumentsReq represents request for listing datasetsDocuments
type ListDatasetsDocumentsReq struct {
	// The ID of the knowledge base.
	DatasetID int64 `json:"dataset_id"`

	// The page number for paginated queries. Default is 1, meaning the data return starts from the
	// first page.
	Page int `json:"page,omitempty"`

	// The size of pagination. Default is 10, meaning that 10 data entries are returned per page.
	Size int `json:"size,omitempty"`
}

// UpdateDatasetsDocumentsReq represents request for updating document
type UpdateDatasetsDocumentsReq struct {
	// The ID of the knowledge base file.
	DocumentID int64 `json:"document_id"`

	// The new name of the knowledge base file.
	DocumentName string `json:"document_name,omitempty"`

	// The update strategy for online web pages. Defaults to no automatic updates.
	// For detailed information, refer to the UpdateRule object.
	UpdateRule *DocumentUpdateRule `json:"update_rule,omitempty"`
}

// createDatasetsDocumentsResp represents response for creating document
type createDatasetsDocumentsResp struct {
	baseResponse
	*CreateDatasetsDocumentsResp
}

// CreateDatasetsDocumentsResp represents response for creating document
type CreateDatasetsDocumentsResp struct {
	baseModel
	DocumentInfos []*Document `json:"document_infos"`
}

// listDatasetsDocumentsResp represents response for listing datasetsDocuments
type listDatasetsDocumentsResp struct {
	baseResponse
	*ListDatasetsDocumentsResp
}

// ListDatasetsDocumentsResp represents response for listing datasetsDocuments
type ListDatasetsDocumentsResp struct {
	baseModel
	Total         int64       `json:"total"`
	DocumentInfos []*Document `json:"document_infos"`
}

// deleteDatasetsDocumentsResp represents response for deleting datasetsDocuments
type deleteDatasetsDocumentsResp struct {
	baseResponse
}

// DeleteDatasetsDocumentsResp represents response for deleting datasetsDocuments
type DeleteDatasetsDocumentsResp struct {
	baseModel
}

// updateDatasetsDocumentsResp represents response for updating document
type updateDatasetsDocumentsResp struct {
	baseResponse
}

// UpdateDatasetsDocumentsResp represents response for updating document
type UpdateDatasetsDocumentsResp struct {
	baseModel
}

// BuildWebPage creates basic document information for webpage type
func BuildWebPage(name string, url string) *DocumentBase {
	return &DocumentBase{
		Name:       name,
		SourceInfo: BuildWebPageSourceInfo(url),
		UpdateRule: BuildNoAutoUpdateRule(),
	}
}

// BuildWebPageWithInterval creates webpage type document information with auto-update interval
func BuildWebPageWithInterval(name string, url string, interval int) *DocumentBase {
	return &DocumentBase{
		Name:       name,
		SourceInfo: BuildWebPageSourceInfo(url),
		UpdateRule: BuildAutoUpdateRule(interval),
	}
}

// BuildLocalFile creates basic document information for local file type
func BuildLocalFile(name string, content string, fileType string) *DocumentBase {
	return &DocumentBase{
		Name:       name,
		SourceInfo: BuildLocalFileSourceInfo(content, fileType),
	}
}

// BuildWebPageSourceInfo creates document source information for webpage type
func BuildWebPageSourceInfo(url string) *DocumentSourceInfo {
	return &DocumentSourceInfo{
		WebUrl:         url,
		DocumentSource: 1,
	}
}

// BuildLocalFileSourceInfo creates document source information for local file type
func BuildLocalFileSourceInfo(content string, fileType string) *DocumentSourceInfo {
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))
	return &DocumentSourceInfo{
		FileBase64: encodedContent,
		FileType:   fileType,
	}
}

// BuildNoAutoUpdateRule creates a rule for no automatic updates
func BuildNoAutoUpdateRule() *DocumentUpdateRule {
	return &DocumentUpdateRule{
		UpdateType: DocumentUpdateTypeNoAutoUpdate,
	}
}

// BuildAutoUpdateRule creates a rule for automatic updates with specified interval
func BuildAutoUpdateRule(interval int) *DocumentUpdateRule {
	return &DocumentUpdateRule{
		UpdateType:     DocumentUpdateTypeAutoUpdate,
		UpdateInterval: interval,
	}
}
