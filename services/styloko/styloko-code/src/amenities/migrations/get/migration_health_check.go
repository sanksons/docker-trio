package get

// MigrationHealthCheck stuct for Healthcheck
type MigrationHealthCheck struct {
}

// GetName returns the name of the Healthcheck
func (n MigrationHealthCheck) GetName() string {
	return "Category Get Healthcheck"
}

// GetHealth returns Healthcheck status
func (n MigrationHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
