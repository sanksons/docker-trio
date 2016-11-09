package get

type StandardSizeGetHealthCheck struct {
}

func (n StandardSizeGetHealthCheck) GetName() string {
	return "hello world"
}

func (n StandardSizeGetHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
