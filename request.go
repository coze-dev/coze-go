package coze

import (
	"context"
	"net/http"
)

// HTTPClient an interface for making HTTP requests
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func isResponseSuccess(ctx context.Context, baseResp baseRespInterface, bodyBytes []byte, httpResponse *httpResponse) error {
	baseResp.SetHTTPResponse(httpResponse)
	if baseResp.GetCode() != 0 {
		logger.Warnf(ctx, "request failed, body=%s, log_id=%s", string(bodyBytes), httpResponse.LogID())
		return NewError(baseResp.GetCode(), baseResp.GetMsg(), httpResponse.LogID())
	}
	return nil
}
