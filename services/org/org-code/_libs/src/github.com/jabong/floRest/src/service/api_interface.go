package service

import (
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/versionmanager"
)

type ApiInterface interface {
	GetVersion() versionmanager.Version

	GetOrchestrator() orchestrator.Orchestrator

	GetHealthCheck() healthcheck.HealthCheckInterface

	Init()
}
