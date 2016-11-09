package get

type BrandsGetHealthCheck struct {
}

func (n BrandsGetHealthCheck) GetName() string {
	return "Brand Get API"
}

func (n BrandsGetHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
