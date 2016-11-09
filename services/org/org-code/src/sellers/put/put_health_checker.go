package put

type sellerUpdateHealthCheck struct {
}

func (n sellerUpdateHealthCheck) GetName() string {
	return "Seller get api health check successfull"
}

func (n sellerUpdateHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
