package simplifier

type CsvErrorHeathCheck struct {
}

func (n CsvErrorHeathCheck) GetName() string {
	return "CsvErrorHeathCheck"
}

func (n CsvErrorHeathCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
