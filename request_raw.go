package coze

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// RawRequestReq ...
type RawRequestReq struct {
	Method   string      // http request method, such as GET, POST
	URL      string      // http request url
	Body     interface{} // http request body, query, path and other parameter information
	IsFile   bool        // send body data as a file
	AuthType int         // 0: token, 1: no-token, 2: jwt(not-jwt-token)
}

func (r *core) rawRequest(ctx context.Context, req *RawRequestReq, resp interface{}) (err error) {
	// 1. parse request
	rawHttpReq, err := r.parseRawHttpRequest(ctx, req)
	if err != nil {
		// 这里日志不需要区分 level, 输出 [error] 日志
		r.Log(ctx, LogLevelError, "[coze] %s %s parse_req failed, err=%s", req.Method, req.URL, err)
		setBaseRespInterface(resp, nil)
		return err
	}

	// 2. request log
	switch r.logLevel {
	case LogLevelDebug:
		r.Log(ctx, LogLevelDebug, "[coze] %s %s start, body=%s", rawHttpReq.Method, rawHttpReq.URL, rawHttpReq.RawBody)
	case LogLevelInfo:
		r.Log(ctx, LogLevelInfo, "[coze] %s %s start", rawHttpReq.Method, rawHttpReq.URL)
	default:
		// error 不需要 req 日志, 合并到 resp 一起
	}

	// 3. do request
	httpResponse, respContent, err := r.doRequest(ctx, rawHttpReq, resp)
	logID, statusCode := getResponseLogID(httpResponse)
	setBaseRespInterface(resp, httpResponse)
	if err != nil {
		switch r.logLevel {
		case LogLevelDebug:
			// [debug]: 详细 error 日志
			r.Log(ctx, LogLevelError, "[coze] %s %s failed, log_id=%s, status=%d, body=%s, err=%s", rawHttpReq.Method, rawHttpReq.URL, logID, statusCode, respContent, err)
		default:
			// [其他]: 简单 error 日志
			r.Log(ctx, LogLevelError, "[coze] %s %s failed, log_id=%s, status=%d, err=%s", rawHttpReq.Method, rawHttpReq.URL, logID, statusCode, err)
		}
		setBaseRespInterface(resp, nil)
		return err
	}
	code, msg := getCodeMsg(resp)

	// 4. response log
	if statusCode >= http.StatusBadRequest || code != 0 {
		r.Log(ctx, LogLevelError, "[coze] %s %s failed, log_id=%s, status=%d, code=%d, msg=%s", rawHttpReq.Method, rawHttpReq.URL, logID, statusCode, code, msg)
	} else {
		switch r.logLevel {
		case LogLevelDebug:
			r.Log(ctx, LogLevelDebug, "[coze] %s %s success, log_id=%s, body=%s", rawHttpReq.Method, rawHttpReq.URL, logID, respContent)
		case LogLevelInfo:
			r.Log(ctx, LogLevelInfo, "[coze] %s %s success, log_id=%s", rawHttpReq.Method, rawHttpReq.URL, logID)
		default:
			// error 不需要 resp 日志
		}
	}

	// 5. response
	if code != 0 {
		return NewError(int(code), msg, logID)
	}

	return nil
}

// 把可读的 RawRequestReq ，解析为 http 请求的参数 rawHttpRequestParam
func (r *core) parseRawHttpRequest(ctx context.Context, req *RawRequestReq) (*rawHttpRequest, error) {
	// 0 init
	rawHttpReq := &rawHttpRequest{
		Method:  strings.ToUpper(req.Method),
		Headers: map[string]string{},
		URL:     r.baseURL + req.URL,
	}

	// 1 headers
	if err := rawHttpReq.parseHeader(ctx, r, req); err != nil {
		return nil, err
	}

	// 2 body
	if err := rawHttpReq.parseRawRequestReqBody(req.Body, req.IsFile); err != nil {
		return nil, err
	}

	// 3 return
	return rawHttpReq, nil
}

func (r *core) doRequest(ctx context.Context, rawHttpReq *rawHttpRequest, realResponse interface{}) (*http.Response, string, error) {
	if rawHttpReq.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, rawHttpReq.Timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, rawHttpReq.Method, rawHttpReq.URL, rawHttpReq.Body)
	if err != nil {
		return nil, "", err
	}
	for k, v := range rawHttpReq.Headers {
		req.Header.Set(k, v)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return resp, "", err
	}

	_, media, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	respFilename := ""
	if media != nil {
		respFilename = media["filename"]
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, "", err
	}

	var respContent string
	if respFilename == "" {
		respContent = string(bs)
	} else {
		respContent = fmt.Sprintf("<FILE: %d>", len(bs))
	}

	if realResponse != nil {
		if resp != nil && resp.StatusCode == http.StatusOK {
			isSpecResp := false
			if setter, ok := realResponse.(readerSetter); ok {
				isSpecResp = true
				setter.SetReader(bytes.NewReader(bs))
			}
			if setter, ok := realResponse.(filenameSetter); ok {
				isSpecResp = true
				setter.SetFilename(respFilename)
			}
			if isSpecResp {
				return resp, respContent, nil
			}
		}

		if len(bs) == 0 && resp.StatusCode >= http.StatusBadRequest {
			return resp, respContent, fmt.Errorf("request fail: %s", resp.Status)
		}

		if err = json.Unmarshal(bs, realResponse); err != nil {
			return resp, respContent, fmt.Errorf("invalid json: %s, err: %s", bs, err)
		}
	}

	return resp, respContent, nil
}

func (r *rawHttpRequest) parseHeader(ctx context.Context, ins *core, req *RawRequestReq) error {
	// agent
	r.Headers["User-Agent"] = userAgent
	r.Headers["X-Coze-Client-User-Agent"] = clientUserAgent

	// logid
	if ins.enableLogID {
		logID, ok := getStringFromContext(ctx, ctxLogIDKey)
		if ok {
			r.Headers[httpLogIDKey] = logID
		}
	}

	// auth
	switch req.AuthType {
	case 0:
		token, err := ins.auth.Token(ctx)
		if err != nil {
			return err
		}
		r.Headers[authorizeHeader] = "Bearer " + token
	}

	return nil
}

func (r *rawHttpRequest) parseRawRequestReqBody(body interface{}, isFile bool) error {
	var reader io.Reader
	fileKey := ""
	query := url.Values{}
	isNeedBody := false
	fileData := map[string]string{}

	if err := rangeStruct(body, func(fieldVV reflect.Value, fieldVT reflect.StructField) error {
		if path := fieldVT.Tag.Get("path"); path != "" {
			if strings.Contains(r.URL, ":"+path) {
				r.URL = strings.ReplaceAll(r.URL, ":"+path, reflectToString(fieldVV))
			} else {
				r.URL = strings.ReplaceAll(r.URL, "{"+path+"}", reflectToString(fieldVV))
			}
		} else if queryKey := fieldVT.Tag.Get("query"); queryKey != "" {
			value := reflectToQueryString(fieldVV)
			sep := fieldVT.Tag.Get("join_sep")
			if sep != "" {
				query.Add(queryKey, strings.Join(value, sep))
			} else {
				for _, v := range value {
					query.Add(queryKey, v)
				}
			}
		} else if j := fieldVT.Tag.Get("json"); j != "" {
			j = strings.TrimSuffix(j, ",omitempty")
			if isFile {
				fileKey = j
				if r, ok := fieldVV.Interface().(io.Reader); ok {
					reader = r
				} else {
					fileData[j] = reflectToString(fieldVV)
				}
			} else {
				isNeedBody = true
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if isNeedBody {
		bs, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r.Body = bytes.NewBuffer(bs)
		r.RawBody = bs
		r.Headers["Content-Type"] = "application/json; charset=utf-8"
	}

	if isFile {
		contentType, bod, err := newFileUploadRequest(fileData, fileKey, reader)
		if err != nil {
			return err
		}
		r.Headers["Content-Type"] = contentType
		r.Body = bod
		r.RawBody = []byte("<FILE>")
	}

	if len(query) > 0 {
		r.URL = r.URL + "?" + query.Encode()
	}

	return nil
}

type rawHttpRequest struct {
	Method  string
	URL     string
	Body    io.Reader
	RawBody []byte
	Headers map[string]string
	Timeout time.Duration
}

func newFileUploadRequest(params map[string]string, filekey string, reader io.Reader) (string, io.Reader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(filekey, "unknown-file")
	if err != nil {
		return "", nil, err
	}
	if reader != nil {
		if _, err = io.Copy(part, reader); err != nil {
			return "", nil, err
		}
	}
	for key, val := range params {
		if err = writer.WriteField(key, val); err != nil {
			return "", nil, err
		}
	}
	if err = writer.Close(); err != nil {
		return "", nil, err
	}

	return writer.FormDataContentType(), body, nil
}

type readerSetter interface {
	SetReader(file io.Reader)
}

type filenameSetter interface {
	SetFilename(filename string)
}

func getCodeMsg(v interface{}) (code int64, msg string) {
	if v == nil {
		return 0, ""
	}
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		return 0, ""
	}
	codeField := vv.FieldByName("Code")
	if codeField.IsValid() {
		if isInReflectKind(codeField.Kind(), []reflect.Kind{
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
		}) {
			code = int64(codeField.Int())
		} else if isInReflectKind(codeField.Kind(), []reflect.Kind{
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
		}) {
			code = int64(codeField.Uint())
		}
	}

	msgField := vv.FieldByName("Msg")
	if msgField.IsValid() {
		if msgField.Kind() == reflect.String {
			msg = msgField.String()
		}
	}

	return
}

func getResponseLogID(response *http.Response) (logID string, statusCode int) {
	if response == nil {
		return
	}
	logID = response.Header.Get(httpLogIDKey)
	statusCode = response.StatusCode
	return
}

func getStringFromContext(ctx context.Context, key string) (string, bool) {
	if ctx == nil {
		return "", false
	}

	v := ctx.Value(key)
	if v == nil {
		return "", false
	}

	switch v := v.(type) {
	case string:
		return v, true
	case *string:
		if v == nil {
			return "", false
		}
		return *v, true
	}
	return "", false
}

func rangeStruct(v interface{}, f func(fieldVV reflect.Value, fieldVT reflect.StructField) error) error {
	vv := reflect.ValueOf(v)
	vt := reflect.TypeOf(v)
	if vt.Kind() == reflect.Ptr {
		vv = vv.Elem()
		vt = vt.Elem()
	}

	for i := 0; i < vt.NumField(); i++ {
		fieldVV := vv.Field(i)
		fieldVT := vt.Field(i)

		if fieldVV.Kind() == reflect.Ptr && fieldVV.IsNil() {
			continue
		}
		if fieldVV.Kind() == reflect.Slice && fieldVV.Len() == 0 {
			continue
		}

		err := f(fieldVV, fieldVT)
		if err != nil {
			return err
		}
	}

	return nil
}

// 从任意类型中提取 baseModel 的指针（广度优先）
func getBaseModelPointer(input any) *baseModel {
	if input == nil {
		return nil
	}

	return findBaseModelInValueBFS(reflect.ValueOf(input))
}

func findBaseModelInValueBFS(v reflect.Value) *baseModel {
	queue := []reflect.Value{v} // queue

	for len(queue) > 0 {
		current := queue[0] // pop
		queue = queue[1:]

		for current.Kind() == reflect.Ptr {
			if current.IsNil() {
				break
			}
			current = current.Elem()
		}
		if current.Kind() != reflect.Struct {
			continue
		}

		t := current.Type()
		for i := 0; i < current.NumField(); i++ {
			field := current.Field(i)
			fieldType := t.Field(i)

			if fieldType.Type.Name() == "baseModel" && fieldType.Anonymous {
				if field.CanAddr() {
					// MUST SUCCESS
					ptr := unsafe.Pointer(field.UnsafeAddr())
					return (*baseModel)(ptr)
				}
			}
		}

		for i := 0; i < current.NumField(); i++ {
			field := current.Field(i)
			if field.Kind() == reflect.Ptr && !field.IsNil() {
				queue = append(queue, field)
			}
			if field.Kind() == reflect.Struct {
				queue = append(queue, field)
			}
		}
	}

	return nil
}

func setBaseRespInterface(resp any, httpResponse *http.Response) {
	if resp == nil {
		return
	}
	h := newHTTPResponse(httpResponse)
	if v, ok := resp.(baseRespInterface); ok {
		v.SetHTTPResponse(h)
	}
	if v := getBaseModelPointer(resp); v != nil {
		v.setHTTPResponse(h)
	}
}

func reflectToString(v reflect.Value) (s string) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	default:
		return v.String()
	}
}

func reflectToQueryString(v reflect.Value) (s []string) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		return []string{v.String()}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []string{strconv.FormatInt(v.Int(), 10)}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []string{strconv.FormatUint(v.Uint(), 10)}
	case reflect.Bool:
		return []string{strconv.FormatBool(v.Bool())}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			s = append(s, reflectToString(v.Index(i)))
		}
		return s
	default:
		return []string{v.String()}
	}
}

func isInReflectKind(v reflect.Kind, list []reflect.Kind) bool {
	for _, vv := range list {
		if vv == v {
			return true
		}
	}
	return false
}
