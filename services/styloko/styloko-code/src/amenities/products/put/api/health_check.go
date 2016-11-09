package api

type ProductUpdateHealthCheck struct {
}

func (n ProductUpdateHealthCheck) GetName() string {
	return "product update health check"
}

func (n ProductUpdateHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
