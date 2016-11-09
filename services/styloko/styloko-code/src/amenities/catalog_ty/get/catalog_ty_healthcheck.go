package get

// CatalogTyHeathCheck stuct for Healthcheck
type CatalogTyHealthCheck struct {
}

// GetName returns the name of the Healthcheck
func (n CatalogTyHealthCheck) GetName() string {
	return "Category Get Healthcheck"
}

// GetHealth returns Healthcheck status
func (n CatalogTyHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
