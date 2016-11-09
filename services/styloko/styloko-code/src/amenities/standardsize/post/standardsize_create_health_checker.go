package post

type StandardSizeCreateHealthCheck struct {
}

func (n StandardSizeCreateHealthCheck) GetName() string {
	return "hello world"
}

func (n StandardSizeCreateHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
