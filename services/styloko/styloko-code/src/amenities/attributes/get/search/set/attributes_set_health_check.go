package set

type AttributesSetHealthCheck struct {
}

func (n AttributesSetHealthCheck) GetName() string {
	return "hello world"
}

func (n AttributesSetHealthCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
