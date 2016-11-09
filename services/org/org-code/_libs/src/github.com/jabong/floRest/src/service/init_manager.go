package service

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/jabong/floRest/src/common/cachestrategy"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/env"
	"github.com/jabong/floRest/src/common/monitor"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/jabong/floRest/src/common/utils/responseheaders"
	"github.com/jabong/floRest/src/common/versionmanager"
)

type InitManager struct {
}

var ResourceToThreshold = map[string]int64{}

/*
Entry point for the init module.
Creates the following:
1> Application Configuration Object
2> Logger
3> Application Orchestrator pipelines
4> Data Access Objects
*/
func (im InitManager) Execute() {

	//Intialise OsEnvVariables
	env.GetOsEnviron()

	//Create the Global Application Config
	initConfig()

	//Sets the Application Performance Params
	setAppPerfParams()

	//Initialises Logger
	initLogger()

	// Initilalize Monitor
	initMonitor()

	// Initialize DBAdapterManager
	InitDBAdapterManager()

	//Create the WorkFlows
	InitVersionManager()

	// Initialise http pooling
	InitHttpPool()

	//Initializes custom api init functionality
	InitCustomApiInit()

	//Initialize Apis
	InitApis()

	//Initialise the Health Checks
	InitHealthCheck()
}

//initConfig initialises the Global Application Config
func initConfig() {
	cm := new(ConfigManager)
	cm.InitializeGlobalConfig()
	cm.InitializeAppConfig()
	cm.UpdateConfigFromEnv(config.ApplicationConfig, "application")
	cm.UpdateConfigFromEnv(config.GlobalAppConfig, "global")
}

//initLogger initialises the logger
func initLogger() {
	err := logger.Initialise(config.GlobalAppConfig.LogConfFile)
	if err != nil {
		panic(err)
	}
}

//initVersionManager create the WorkFlows
func InitVersionManager() {
	vmap := versionmanager.VersionMap{
		versionmanager.Version{
			Resource: "SERVICE",
			Version:  "V1",
			Action:   "GET",
			BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
		}: createServiceOrchestrator(),
		versionmanager.Version{
			Resource: constants.HEALTHCHECKAPI,
			Version:  "",
			Action:   "GET",
			BucketId: constants.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
		}: createHealthCheckOrchestrator(),
	}

	addApiVersions(vmap)
	versionmanager.Initialize(vmap)
}

//Calls the custom api init function
func InitCustomApiInit() {
	//If the apiCustomInitFunc is defined then execute it
	if apiCustomInitFunc != nil {
		apiCustomInitFunc()
	}
}

func addApiVersions(vmap versionmanager.VersionMap) {
	for _, apiInstance := range apiList {
		vmap[apiInstance.GetVersion()] = apiInstance.GetOrchestrator()
	}
}

func createServiceOrchestrator() orchestrator.Orchestrator {
	logger.Info("Service Pipeline Creation begin")

	serviceOrchestrator := new(orchestrator.Orchestrator)
	serviceWorkflow := new(orchestrator.WorkFlowDefinition)
	serviceWorkflow.Create()

	//Create and add execution node UriInterpreter
	uriInterpreter := new(UriInterpreter)
	uriInterpreter.SetID("1")
	uerr := serviceWorkflow.AddExecutionNode(uriInterpreter)
	if uerr != nil {
		logger.Error(fmt.Sprintln(uerr))
	}

	//Create and add execution node BusinessLogicExecutor
	businessLogicExecutor := new(BusinessLogicExecutor)
	businessLogicExecutor.SetID("3")
	berr := serviceWorkflow.AddExecutionNode(businessLogicExecutor)
	if berr != nil {
		logger.Error(fmt.Sprintln(berr))
	}
	responseHeaderWriter := new(responseheaders.Writer)
	responseHeaderWriter.SetID("4")
	rhErr := serviceWorkflow.AddExecutionNode(responseHeaderWriter)
	if rhErr != nil {
		logger.Error(fmt.Sprintln(rhErr))
	}

	//Create and add execution node HttpResponseCreator
	httpResponseCreator := new(HttpResponseCreator)
	httpResponseCreator.SetID("5")
	herr := serviceWorkflow.AddExecutionNode(httpResponseCreator)
	if herr != nil {
		logger.Error(fmt.Sprintln(herr))
	}

	//Add the connections between the nodes
	c1err := serviceWorkflow.AddConnection(uriInterpreter, businessLogicExecutor)
	if c1err != nil {
		logger.Error(fmt.Sprintln(c1err))
	}
	c2err := serviceWorkflow.AddConnection(businessLogicExecutor, responseHeaderWriter)
	if c2err != nil {
		logger.Error(fmt.Sprintln(c2err))
	}
	c3err := serviceWorkflow.AddConnection(responseHeaderWriter, httpResponseCreator)
	if c3err != nil {
		logger.Error(fmt.Sprintln(c3err))
	}

	//Set the Workflow Start Node
	serviceWorkflow.SetStartNode(uriInterpreter)

	//Assign the workflow definition to the Orchestrator
	serviceOrchestrator.Create(serviceWorkflow)

	logger.Info(serviceOrchestrator.String())
	logger.Info("Service Pipeline Created")

	return *serviceOrchestrator
}

func createHealthCheckOrchestrator() orchestrator.Orchestrator {
	logger.Info("Health Check Pipeline Creation begin")

	healthCheckOrchestrator := new(orchestrator.Orchestrator)
	healthCheckOrchestratorWorkflow := new(orchestrator.WorkFlowDefinition)
	healthCheckOrchestratorWorkflow.Create()

	healthCheckExecutor := new(healthcheck.HealthCheckExecutor)
	healthCheckExecutor.SetID("1")
	hcerr := healthCheckOrchestratorWorkflow.AddExecutionNode(healthCheckExecutor)
	if hcerr != nil {
		logger.Error(fmt.Sprintln(hcerr))
	}

	healthCheckOrchestratorWorkflow.SetStartNode(healthCheckExecutor)
	healthCheckOrchestrator.Create(healthCheckOrchestratorWorkflow)

	logger.Info(healthCheckOrchestrator.String())
	logger.Info("Health Check Pipeline Created")

	return *healthCheckOrchestrator
}

//initApis initializes all apis
func InitApis() {
	for _, apiInstance := range apiList {
		apiInstance.Init()
	}
	logger.Info("Initialized apis")
}

//setAppPerfParams sets the Application's performance parameters
func setAppPerfParams() {
	perf := config.GlobalAppConfig.Performance
	setGCPercentage(perf.GCPercentage)
	setNoOfCpuCores(perf.UseCorePercentage)
}

//setGCPercentage sets when to trigger the garbage collection
func setGCPercentage(gcPer float64) {
	debug.SetGCPercent(int(gcPer))
}

//setNoOfCpuCores sets number of CPU cores the app should use
func setNoOfCpuCores(cpuCorePer float64) {
	corePer := config.GlobalAppConfig.Performance.UseCorePercentage
	if corePer <= 0 {
		fmt.Printf("No of Cpu Core to be Used = 1\n")
		return
	}
	totalCpus := float64(runtime.NumCPU())
	cpuCore := totalCpus * (corePer / 100)
	if cpuCore <= 0 {
		fmt.Printf("No of Cpu Core to be Used = 1\n")
		return
	}

	if cpuCore > totalCpus {
		cpuCore = totalCpus
	}
	fmt.Printf("No of Cpu Core to be Used = %d\n", int(cpuCore))
	runtime.GOMAXPROCS(int(cpuCore))
}

func InitHealthCheck() {
	healthCheckArray := make([]healthcheck.HealthCheckInterface, len(apiList))
	//get all Healthcheck instances
	for i, apiInstance := range apiList {
		healthCheckArray[i] = apiInstance.GetHealthCheck()
	}

	healthCheckArray = append(healthCheckArray, new(ServiceHealthCheck))
	healthcheck.Initialise(healthCheckArray)
}

// initMonitor: initlaize monitor
func initMonitor() {
	if config.GlobalAppConfig.MonitorConfig.Enabled {
		if err := monitor.GetInstance().Initialize(&config.GlobalAppConfig.MonitorConfig); err != nil {
			logger.Error(fmt.Sprintln(err))
		}
	}
}

// initDBAdapterManager: initialize dbAdapterManager
func InitDBAdapterManager() {
	cachestrategy.DBAdapterMgr = new(cachestrategy.DBAdapterManager)
	cachestrategy.DBAdapterMgr.Initialize()
}

// InitHttpPool: initialize http pool
func InitHttpPool() {
	http.InitConnPool(&config.GlobalAppConfig.HttpConfig)
}
