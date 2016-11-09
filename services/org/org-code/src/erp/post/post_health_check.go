package post

type ErpUdateHealthCheck struct {
}

func (n ErpUdateHealthCheck) GetName() string {
	return "Update Erp api health check successful"
}

func (n ErpUdateHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
