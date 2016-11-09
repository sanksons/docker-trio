package datadog

// Alert types
const (
	INFO    = "info"
	WARNING = "warning"
	ERROR   = "error"
)

// Endpoints for datadog
const (
	baseURL       = "https://app.datadoghq.com"
	eventEndpoint = "/api/v1/events?"
	gaugeEndpoint = "/api/v1/series?"
)
