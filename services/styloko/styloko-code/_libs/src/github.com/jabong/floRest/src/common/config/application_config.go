package config

import (
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/cachestrategy"
	"github.com/jabong/floRest/src/common/monitor"
	"github.com/jabong/floRest/src/common/utils/http"
)

type AppConfig struct {
	AppName             string
	AppVersion          string
	ServerPort          string
	LogConfFile         string
	MonitorConfig       monitor.MonitorConf
	CacheStrategyConfig cachestrategy.Config
	Performance         PerformanceConfigs
	DynamicConfig       DynamicConfigInfo
	HttpConfig          http.Config
	ResponseHeaders     ResponseHeaderFields
}

type PerformanceConfigs struct {
	UseCorePercentage float64
	GCPercentage      float64
}

type ResponseHeaderFields struct {
	CacheControl CacheControlHeaders
}

type CacheControlHeaders struct {
	ResponseType    string
	NoCache         bool
	NoStore         bool
	MaxAgeInSeconds int
}

type Application struct {
	ResponseHeaders ResponseHeaderFields
}

type DynamicConfigInfo struct {
	Active          bool
	RefreshInterval int
	ConfigKey       string
	Cache           cache.Config
}

//Global ApplicationConfig Singleton
var GlobalAppConfig *AppConfig = nil
var ApplicationConfig interface{}
