package delete

type ProductDeleteHealthCheck struct {
}

func (n ProductDeleteHealthCheck) GetName() string {
	return "product Create"
}

func (n ProductDeleteHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
