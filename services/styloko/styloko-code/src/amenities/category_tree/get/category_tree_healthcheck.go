package get

// CategoryTreeHealthCheck stuct for Healthcheck
type CategoryTreeHealthCheck struct {
}

// GetName returns the name of the Healthcheck
func (n CategoryTreeHealthCheck) GetName() string {
	return "Category Get Healthcheck"
}

// GetHealth returns Healthcheck status
func (n CategoryTreeHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
