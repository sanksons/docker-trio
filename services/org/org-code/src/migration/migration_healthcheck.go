package migration

type MigrationHealthCheck struct {
}

func (m MigrationHealthCheck) GetName() string {
	return "Migration api health check successful"
}

func (m MigrationHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
