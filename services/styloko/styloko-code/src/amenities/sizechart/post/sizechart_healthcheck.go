package post

type SizeChartHealthCheck struct {
}

func (n SizeChartHealthCheck) GetName() string {
	return "Sizechart health check"
}

func (n SizeChartHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
