package get

type TaxClassHealthCheck struct {
}

func (n TaxClassHealthCheck) GetName() string {
	return "product search"
}

func (n TaxClassHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
