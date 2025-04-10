package coze

import "net/http"

type Responser interface {
	Response() HTTPResponse
}

type HTTPResponse interface {
	LogID() string
}

type httpResponse struct {
	Status        int
	Header        http.Header
	ContentLength int64

	logid string
}

func (r *httpResponse) LogID() string {
	if r.logid == "" {
		r.logid = r.Header.Get(httpLogIDKey)
	}
	return r.logid
}

type baseResponse struct {
	Code         int           `json:"code"`
	Msg          string        `json:"msg"`
	HTTPResponse *httpResponse `json:"http_response"`
}

func (r *baseResponse) SetHTTPResponse(httpResponse *httpResponse) {
	r.HTTPResponse = httpResponse
}

func (r *baseResponse) SetCode(code int) {
	r.Code = code
}

func (r *baseResponse) SetMsg(msg string) {
	r.Msg = msg
}

func (r *baseResponse) GetCode() int {
	return r.Code
}

func (r *baseResponse) GetMsg() string {
	return r.Msg
}

func (r *baseResponse) LogID() string {
	return r.HTTPResponse.LogID()
}

type baseRespInterface interface {
	SetHTTPResponse(httpResponse *httpResponse)
	SetCode(code int)
	SetMsg(msg string)
	GetMsg() string
	GetCode() int
}

type baseModel struct {
	httpResponse *httpResponse
}

func (r *baseModel) setHTTPResponse(httpResponse *httpResponse) {
	r.httpResponse = httpResponse
}

func (r *baseModel) Response() HTTPResponse {
	return r.httpResponse
}

func (r *baseModel) LogID() string {
	return r.httpResponse.LogID()
}

func newHTTPResponse(resp *http.Response) *httpResponse {
	return &httpResponse{
		Status:        resp.StatusCode,
		Header:        resp.Header,
		ContentLength: resp.ContentLength,
	}
}
