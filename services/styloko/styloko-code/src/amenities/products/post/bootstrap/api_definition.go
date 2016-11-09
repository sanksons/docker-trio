package bootstrap

import (
	"common/appconfig"
	"common/pool"

	"strconv"

	"github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/versionmanager"
)

//Define Global variables
var (
	conf            *appconfig.AppConfig
	dbAdapterName   string
	workerPool      pool.Safe
	processedCount  *counters
	poolWorkerCount int
	poolQueueSize   int
)

type counters struct {
	counter int64
}

//
// API definitionn for Product Update
//
type BootstrapAPI struct {
}

func (a *BootstrapAPI) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: "BOOTSTRAP",
		Version:  "V1",
		Action:   "POST",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *BootstrapAPI) GetOrchestrator() orchestrator.Orchestrator {

	bootstrapNode := new(BootstrapNode)
	bootstrapNode.SetID("bootstrapNode")

	workflow := new(orchestrator.WorkFlowDefinition)
	workflow.Create()
	workflow.AddExecutionNode(bootstrapNode)
	//set start node
	workflow.SetStartNode(bootstrapNode)

	Orchestrator := new(orchestrator.Orchestrator)
	Orchestrator.Create(workflow)
	return *Orchestrator
}

func (a *BootstrapAPI) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(BootstrapHealthCheck)
}

func (a *BootstrapAPI) Init() {
	//Store config into global object
	conf = config.ApplicationConfig.(*appconfig.AppConfig)
	dbAdapterName = conf.DbAdapter
	var err error
	poolQueueSize, err = strconv.Atoi(conf.Bootstrap.QueueSize)
	if err != nil {
		poolQueueSize = 100
	}
	poolWorkerCount, err = strconv.Atoi(conf.Bootstrap.WorkerCount)
	if err != nil {
		poolWorkerCount = 2
	}
	//Initialize bootstrap pool
	workerPool = pool.NewWorkerSafe(
		"bootstrap", poolWorkerCount, poolQueueSize, 0, 0,
	)
	processedCount = &counters{}
	workerPool.StartWorkers(PublishProduct)
}
