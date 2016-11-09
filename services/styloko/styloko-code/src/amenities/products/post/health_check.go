package post

type ProductCreateHealthCheck struct {
}

func (n ProductCreateHealthCheck) GetName() string {
	return "product Create"
}

func (n ProductCreateHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
