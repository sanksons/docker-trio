package logmonitor

import (
	"testing"
)

func getTestConfig() *Config {
	c := new(Config)
	c.ApiKey = ""
	c.AppKey = ""
	c.AppName = "TestJadeGO"
	c.Platform = DatadogAgent
	c.ServerAddr = "127.0.0.1:8125"
	c.Verbose = true
	c.Enabled = true
	return c
}

func getTestDatadogAgentClient() (MonitorInterface, error) {
	c := getTestConfig()
	d, err := Get(c)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func Test_Get(t *testing.T) {
	_, err := getTestDatadogAgentClient()
	if err != nil {
		t.Errorf("Get DatadogAgentClient Failed %v", err)
	}
}

func Test_Info(t *testing.T) {
	d, err := getTestDatadogAgentClient()
	if err != nil {
		t.Errorf("Failed to get DatadogAgentClient for event infog %v", err)
	}

	data := new(AppLogData)
	data.Title = "Test Info"
	data.Body = "Hello World Info"
	data.Tags = map[string]string{"env": "local"}
	if err := d.Info(data); err != nil {
		t.Errorf("Error event failed %v", err)
	}
}

func Test_Gauges(t *testing.T) {
	d, err := getTestDatadogAgentClient()
	if err != nil {
		t.Errorf("Failed to get DatadogAgentClient for event Gauges %v", err)
	}
	if err := d.Gauge("Test_Gauge", 121.5, []string{"jadehol838"}, 1); err != nil {
		t.Errorf("Error event failed %v", err)
	}
}

func Test_Set(t *testing.T) {
	d, err := getTestDatadogAgentClient()
	if err != nil {
		t.Errorf("Failed to get DatadogAgentClient for event Set %v", err)
	}
	if err := d.Set("Test_Set", "Hello World", []string{"jadehol838"}, 1); err != nil {
		t.Errorf("Error event failed %v", err)
	}
}
