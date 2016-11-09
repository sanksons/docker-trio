package attributes

type AttributesHealthCheck struct {
}

func (n AttributesHealthCheck) GetName() string {
	return "hello world"
}

func (n AttributesHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
