package healthcheck

var healthCheckApiList []HealthCheckInterface = nil

//Initialise initialises an app monitor
func Initialise(apiList []HealthCheckInterface) {
	if healthCheckApiList == nil {
		healthCheckApiList = apiList
	}
}
