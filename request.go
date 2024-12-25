package coze

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/coze-dev/coze-go/log"
)

// HTTPClient an interface for making HTTP requests
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type core struct {
	httpClient HTTPClient
	baseURL    string
}

func newCore(httpClient HTTPClient, baseURL string) *core {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 5,
		}
	}
	return &core{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// RequestOption 请求选项函数类型
type RequestOption func(*http.Request) error

// withHTTPHeader add http header
func withHTTPHeader(key, value string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

// withHTTPQuery add http query
func withHTTPQuery(key, value string) RequestOption {
	return func(req *http.Request) error {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

// Request send http request
func (c *core) Request(ctx context.Context, method, path string, body any, instance any, opts ...RequestOption) error {
	resp, err := c.RawRequest(ctx, method, path, body, opts...)
	if err != nil {
		return err
	}

	return packInstance(instance, resp)
}

// UploadFile 上传文件
func (c *core) UploadFile(ctx context.Context, path string, reader io.Reader, fileName string, fields map[string]string, instance any, opts ...RequestOption) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}

	if _, err = io.Copy(part, reader); err != nil {
		return fmt.Errorf("copy file content: %w", err)
	}

	// 添加其他字段
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return fmt.Errorf("write field %s: %w", key, err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.baseURL, path), body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 应用请求选项
	for _, opt := range opts {
		if err := opt(req); err != nil {
			return fmt.Errorf("apply option: %w", err)
		}
	}
	setUserAgent(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	return packInstance(instance, resp)
}

func (c *core) RawRequest(ctx context.Context, method, path string, body any, opts ...RequestOption) (*http.Response, error) {
	urlInfo := fmt.Sprintf("%s%s", c.baseURL, path)

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, urlInfo, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 设置默认请求头
	req.Header.Set("Content-Type", "application/json")

	// 应用请求选项
	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	setUserAgent(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = checkHttpResp(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func packInstance(instance any, resp *http.Response) error {
	err := checkHttpResp(resp)
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	httpResponse := newHTTPResponse(resp)
	err = json.Unmarshal(bodyBytes, instance)
	if err != nil {
		log.Errorf(fmt.Sprintf("unmarshal response body: %s", string(bodyBytes)))
		return err
	}
	if baseResp, ok := instance.(baseRespInterface); ok {
		return isResponseSuccess(baseResp, bodyBytes, httpResponse)
	}
	return nil
}

func isResponseSuccess(baseResp baseRespInterface, bodyBytes []byte, httpResponse *httpResponse) error {
	baseResp.SetHTTPResponse(httpResponse)
	if baseResp.GetCode() != 0 {
		log.Warnf("request unsuccessful: %s, log_id:%s", string(bodyBytes), httpResponse.GetLogID())
		return NewCozeError(baseResp.GetCode(), baseResp.GetMsg(), httpResponse.GetLogID())
	}
	return nil
}

func checkHttpResp(resp *http.Response) error {
	logID := getLogID(resp.Header)
	// 鉴权的情况，需要解析
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		errorInfo := authErrorFormat{}
		err = json.Unmarshal(bodyBytes, &errorInfo)
		if err != nil {
			log.Errorf(fmt.Sprintf("unmarshal response body: %s", string(bodyBytes)))
			return errors.New(string(bodyBytes) + "log_id:%s" + logID)
		}
		return NewCozeAuthExceptionWithoutParent(&errorInfo, resp.StatusCode, logID)
	}
	return nil
}