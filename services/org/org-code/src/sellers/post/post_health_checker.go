package post

type SellerCreateHealthCheck struct {
}

func (n SellerCreateHealthCheck) GetName() string {
	return "Seller get api health check successfull"
}

func (n SellerCreateHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
