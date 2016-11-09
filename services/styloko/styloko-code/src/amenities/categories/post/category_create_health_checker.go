package post

// CategoryCreateHealthCheck -> Basic struct
type CategoryCreateHealthCheck struct {
}

// GetName -> Returns name
func (n CategoryCreateHealthCheck) GetName() string {
	return "Category Create Healthcheck"
}

// GetHealth -> Returns health status
func (n CategoryCreateHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
