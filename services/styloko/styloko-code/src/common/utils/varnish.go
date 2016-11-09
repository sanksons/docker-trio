package utils

import (
	"net/http"

	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
)

var port string
var appName string
var version string
var client *http.Client

// InvalidateCache invalidates cache for provided url
// q must be a proper query
// path also must be of the form /path/string
func InvalidateCache(endpoint string, path string, q string) {
	baseURL := "127.0.0.1"
	url := baseURL + ":" + port + "/" + appName + "/" + version + "/" + endpoint
	if path != "" {
		url += path
	}
	if q != "" {
		url += q
	}
	req, _ := http.NewRequest("PURGE", url, nil)
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		logger.Error(err)
	}
}

// InitVarnishCache initializes varnish settings
func InitVarnishCache() {
	global := config.GlobalAppConfig
	port = string(global.ServerPort)
	appName = string(global.AppName)
	version = string(global.AppVersion)
	client = &http.Client{}
}
