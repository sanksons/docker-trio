package monitor

type MonitorConf struct {
	AppName       string
	Platform      string
	AgentServer   string
	Verbose       bool
	Enabled       bool
	MetricsServer string
}
