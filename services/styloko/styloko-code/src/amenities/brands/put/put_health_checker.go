package put

type BrandsHealthCheck struct {
}

func (n BrandsHealthCheck) GetName() string {
	return "Brand Put API"
}

func (n BrandsHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
