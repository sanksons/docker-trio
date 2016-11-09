package post

type BrandCreateHealthCheck struct {
}

func (n BrandCreateHealthCheck) GetName() string {
	return "Brand Post API"
}

func (n BrandCreateHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
