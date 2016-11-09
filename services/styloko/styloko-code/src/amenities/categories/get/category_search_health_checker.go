package get

// CategorySearchHealthCheck stuct for Healthcheck
type CategorySearchHealthCheck struct {
}

// GetName returns the name of the Healthcheck
func (n CategorySearchHealthCheck) GetName() string {
	return "Category Get Healthcheck"
}

// GetHealth returns Healthcheck status
func (n CategorySearchHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
