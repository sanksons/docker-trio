package stock

type StockHealthCheck struct {
}

func (n StockHealthCheck) GetName() string {
	return "Stock Get Healthcheck"
}

func (n StockHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
