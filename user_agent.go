package coze

import (
	"encoding/json"
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
