package commissions

type getCommissionsHealthCheck struct {
}

func (n getCommissionsHealthCheck) GetName() string {
	return "Seller get commissions api health check successfull"
}

func (n getCommissionsHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
