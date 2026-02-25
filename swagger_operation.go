package coze

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// SwaggerOperationRequest is a generic request payload for swagger-generated API modules.
type SwaggerOperationRequest struct {
	PathParams  map[string]any
	QueryParams map[string]any
	Body        any
}

// SwaggerOperationResponse is a generic response payload for swagger-generated API modules.
type SwaggerOperationResponse struct {
	baseResponse
	Data any `json:"data"`
}

func buildSwaggerOperationURL(path string, pathParams map[string]any, queryParams map[string]any) string {
	urlPath := strings.TrimSpace(path)
	for key, value := range pathParams {
		if key == "" || value == nil {
			continue
		}
		placeholder := "{" + strings.TrimSpace(key) + "}"
		urlPath = strings.ReplaceAll(urlPath, placeholder, url.PathEscape(fmt.Sprint(value)))
	}

	if len(queryParams) == 0 {
		return urlPath
	}

	values := url.Values{}
	for key, value := range queryParams {
		key = strings.TrimSpace(key)
		if key == "" || value == nil {
			continue
		}
		appendSwaggerQueryValue(values, key, value)
	}
	encoded := values.Encode()
	if encoded == "" {
		return urlPath
	}
	if strings.Contains(urlPath, "?") {
		return urlPath + "&" + encoded
	}
	return urlPath + "?" + encoded
}

func appendSwaggerQueryValue(values url.Values, key string, value any) {
	if value == nil {
		return
	}
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			appendSwaggerQueryValue(values, key, v.Index(i).Interface())
		}
	default:
		values.Add(key, fmt.Sprint(v.Interface()))
	}
}
