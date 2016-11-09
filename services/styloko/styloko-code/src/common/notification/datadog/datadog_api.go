package datadog

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/jabong/floRest/src/common/utils/logger"
)

// Datadog -> Struct for the datadog API
type Datadog struct {
	apiKey string
	appKey string
	host   string
}

type event struct {
	Title     string
	Text      string
	Tags      []string
	AlertType string
	Host      string
}

type series struct {
	Series []gauge
}

type gauge struct {
	Metric string
	Points float64
	Type   string
	Tags   []string
	Host   string
}

// Init -> This is the constructor function.
func Init(apiKey, appKey string) (dd Datadog) {
	host, _ := os.Hostname()
	d := Datadog{
		apiKey: apiKey,
		appKey: appKey,
		host:   host,
	}
	return d
}

func (dd *Datadog) queryBuilder(baseURL, apiKey, appKey string) string {
	params := url.Values{}
	params.Add("api_key", apiKey)
	params.Add("application_key", appKey)
	finalURL := baseURL + params.Encode()
	return finalURL
}

func (dd *Datadog) buildEvent(title, text string, tags []string, alertType string) event {
	if alertType != INFO && alertType != WARNING && alertType != ERROR {
		logger.Error("Invalid alert type.")
	}

	evt := event{
		Title:     title,
		Text:      text,
		Tags:      tags,
		AlertType: alertType,
		Host:      dd.host,
	}
	return evt
}

func (dd *Datadog) buildMetric(metric string, points float64, ty string, tags []string) series {
	g := gauge{
		Metric: metric,
		Points: points,
		Type:   ty,
		Tags:   tags,
		Host:   dd.host,
	}
	arr := []gauge{g}
	s := series{
		Series: arr,
	}
	return s
}

// PostEvent -> Post event to datadog using API endpoint.
func (dd *Datadog) PostEvent(title, text string, tags []string, alertType string) bool {
	evt := dd.buildEvent(title, text, tags, alertType)
	jsonStr, err := json.Marshal(evt)
	if err != nil {
		logger.Error(err)
		return false
	}
	url := dd.queryBuilder(baseURL+eventEndpoint, dd.apiKey, dd.appKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Error(err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return false
	}
	defer resp.Body.Close()
	if (resp.StatusCode)/10 == 20 {
		return true
	}
	return false
}

// SendMetric -> Send out a metric to datadog
func (dd *Datadog) SendMetric(metric string, points float64, ty string, tags []string) bool {
	series := dd.buildMetric(metric, points, ty, tags)
	jsonStr, err := json.Marshal(series)
	if err != nil {
		logger.Error(err)
		return false
	}
	url := dd.queryBuilder(baseURL+gaugeEndpoint, dd.apiKey, dd.appKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Error(err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return false
	}
	defer resp.Body.Close()
	if (resp.StatusCode)/10 == 20 {
		return true
	}
	return false
}
