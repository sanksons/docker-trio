package healthcheck

type HealthCheckInterface interface {
	GetName() string
	GetHealth() map[string]interface{}
}
