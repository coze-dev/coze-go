package coze

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const version = "0.1.0"

var (
	userAgentSDK         = "cozego"
	userAgentLang        = "go"
	userAgentLangVersion = strings.TrimPrefix(runtime.Version(), "go")
	userAgentOsName      = runtime.GOOS
	userAgentOsVersion   = os.Getenv("OSVERSION")
	userAgent            = userAgentSDK + "/" + version + " " + userAgentLang + "/" + userAgentLangVersion + " " + userAgentOsName + "/" + userAgentOsVersion
	clientUserAgent      string
)

func (c *core) setCommonHeaders(req *http.Request) error {
	// agent
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Coze-Client-User-Agent", clientUserAgent)

	// logid
	if c.enableLogID {
		v := req.Context().Value(ctxLogIDKey)
		if logid, ok := v.(string); ok && logid != "" {
			req.Header.Set(httpLogIDKey, logid)
		}
	}

	// auth
	if c.auth != nil {
		// auth 相关请求, c.auth 为 nil
		accessToken, err := c.auth.Token(req.Context())
		if err != nil {
			logger.Errorf(req.Context(), "failed to get access_token: %s", err)
			return err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	return nil
}

func init() {
	clientUserAgent = getCozeClientUserAgent()
}

type userAgentInfo struct {
	Version     string `json:"version"`
	Lang        string `json:"lang"`
	LangVersion string `json:"lang_version"`
	OsName      string `json:"os_name"`
	OsVersion   string `json:"os_version"`
}

func getCozeClientUserAgent() string {
	data, _ := json.Marshal(userAgentInfo{
		Version:     version,
		Lang:        userAgentSDK,
		LangVersion: userAgentLangVersion,
		OsName:      userAgentOsName,
		OsVersion:   userAgentOsVersion,
	})
	return string(data)
}
