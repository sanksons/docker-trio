package get

type SellerGetHealthCheck struct {
}

func (n SellerGetHealthCheck) GetName() string {
	return "Seller get api health check successfull"
}

func (n SellerGetHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
