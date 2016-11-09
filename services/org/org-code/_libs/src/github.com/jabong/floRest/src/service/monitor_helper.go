package service

import (
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
)

func getCustomMetricPrefix(data workflow.WorkFlowData) string {
	var monitorMetricPrefix string = ""
	monitorCustomMetricPrefix, mcmpError := data.ExecContext.Get(constants.MONITOR_CUSTOM_METRIC_PREFIX)
	if mcmpError == nil {
		monitorMetricPrefix, _ = monitorCustomMetricPrefix.(string)
	}
	return monitorMetricPrefix
}
