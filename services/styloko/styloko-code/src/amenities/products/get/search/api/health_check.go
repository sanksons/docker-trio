package api

type ProductSearchHealthCheck struct {
}

func (n ProductSearchHealthCheck) GetName() string {
	return "product search"
}

func (n ProductSearchHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
