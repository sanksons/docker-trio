package put

// CategoryUpdateHealthCheck struct
type CategoryUpdateHealthCheck struct {
}

// GetName returns healthcheck API name
func (n CategoryUpdateHealthCheck) GetName() string {
	return "Category update healthcheck"
}

// GetHealth returns health of API
func (n CategoryUpdateHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
