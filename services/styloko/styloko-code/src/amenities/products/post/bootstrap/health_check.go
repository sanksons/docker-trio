package bootstrap

type BootstrapHealthCheck struct {
}

func (n BootstrapHealthCheck) GetName() string {
	return "product update health check"
}

func (n BootstrapHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
