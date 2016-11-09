package post

import (
	proUtil "amenities/products/common"
	"common/appconfig"
	"fmt"
	"github.com/jabong/floRest/src/common/config"
	florestConsts "github.com/jabong/floRest/src/common/constants"
	"github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/healthcheck"
	"github.com/jabong/floRest/src/common/versionmanager"
	validator "gopkg.in/go-playground/validator.v8"
)

//Global Variables to be used throughout this API
var (
	storageAdapter string
	validate       *validator.Validate
	conf           *appconfig.AppConfig
)

//API declaration for create Products
type CreateProductsApi struct {
}

//Set Resource Specification
func (a *CreateProductsApi) GetVersion() versionmanager.Version {
	return versionmanager.Version{
		Resource: proUtil.PRODUCT_API,
		Version:  "V1",
		Action:   "POST",
		BucketId: florestConsts.ORCHESTRATOR_BUCKET_DEFAULT_VALUE,
	}
}

func (a *CreateProductsApi) GetOrchestrator() orchestrator.Orchestrator {
	//Initiate required nodes
	ValidateNode := new(ValidateNode)
	ValidateNode.SetID("ValidateRequest")

	InsertNode := new(InsertNode)
	InsertNode.SetID("InsertNode")

	ResponseNode := new(ResponseNode)
	ResponseNode.SetID("ResponseNode")

	Workflow := new(orchestrator.WorkFlowDefinition)
	Workflow.Create()

	//
	//        Validate
	//           |
	//         Insert
	//           |
	//        Response

	//Set up workflow

	//add execution nodes
	Workflow.AddExecutionNode(ValidateNode)
	Workflow.AddExecutionNode(InsertNode)
	Workflow.AddExecutionNode(ResponseNode)

	//add connections
	err := Workflow.AddConnection(ValidateNode, InsertNode)
	if err != nil {
		fmt.Sprintln(err)
	}
	err = Workflow.AddConnection(InsertNode, ResponseNode)
	if err != nil {
		fmt.Sprintln(err)
	}
	//set start node
	Workflow.SetStartNode(ValidateNode)

	Orchestrator := new(orchestrator.Orchestrator)
	Orchestrator.Create(Workflow)
	return *Orchestrator
}

//HealthCheck API
func (a *CreateProductsApi) GetHealthCheck() healthcheck.HealthCheckInterface {
	return new(ProductCreateHealthCheck)
}

//API init() declarations
func (a *CreateProductsApi) Init() {
	//store conf
	conf = config.ApplicationConfig.(*appconfig.AppConfig)
	//store storage adapter
	storageAdapter = conf.DbAdapter
	//set validator config
	vconfig := &validator.Config{TagName: "validate"}
	validate = validator.New(vconfig)
	proUtil.IntializeGlobalVariables(conf.Sellers, conf.BrandsProcTime)
}
