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
	"strings"
	"time"
)

// HTTPClient an interface for making HTTP requests
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
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

	return packInstance(ctx, instance, resp)
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

	if err := c.setCommonHeaders(req); err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	return packInstance(ctx, instance, resp)
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

	if err := c.setCommonHeaders(req); err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err = checkHttpResp(ctx, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *core) StreamRequest(ctx context.Context, method, path string, body any, opts ...RequestOption) (*http.Response, error) {
	resp, err := c.RawRequest(ctx, method, path, body, opts...)
	if err != nil {
		return nil, err
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && strings.Contains(contentType, "application/json") {
		return nil, packInstance(ctx, &baseResponse{}, resp)
	}
	return resp, nil
}

func packInstance(ctx context.Context, instance any, resp *http.Response) error {
	err := checkHttpResp(ctx, resp)
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
		logger.Errorf(ctx, fmt.Sprintf("unmarshal response body: %s", string(bodyBytes)))
		return err
	}
	if baseResp, ok := instance.(baseRespInterface); ok {
		return isResponseSuccess(ctx, baseResp, bodyBytes, httpResponse)
	}
	return nil
}

func isResponseSuccess(ctx context.Context, baseResp baseRespInterface, bodyBytes []byte, httpResponse *httpResponse) error {
	baseResp.SetHTTPResponse(httpResponse)
	if baseResp.GetCode() != 0 {
		logger.Warnf(ctx, "request failed, body=%s, log_id=%s", string(bodyBytes), httpResponse.LogID())
		return NewError(baseResp.GetCode(), baseResp.GetMsg(), httpResponse.LogID())
	}
	return nil
}

func checkHttpResp(ctx context.Context, resp *http.Response) error {
	logID := resp.Header.Get(httpLogIDKey)
	// 鉴权的情况，需要解析
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("coze read response body failed: %w, log_id: %s", err, logID)
		}
		errorInfo := authErrorFormat{}
		err = json.Unmarshal(bodyBytes, &errorInfo)
		if err != nil {
			logger.Errorf(ctx, fmt.Sprintf("unmarshal response body: %s", string(bodyBytes)))
			return errors.New(string(bodyBytes) + " log_id: " + logID)
		}
		return NewAuthError(&errorInfo, resp.StatusCode, logID)
	}
	return nil
}
