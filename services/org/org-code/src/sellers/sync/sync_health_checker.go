package sync

type getSyncHealthCheck struct {
}

func (n getSyncHealthCheck) GetName() string {
	return "Seller sync health check successfull"
}

func (n getSyncHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
