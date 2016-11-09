package get

import (
	proUtil "amenities/products/common"
	factory "common/ResourceFactory"
	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type DataAccessorNode struct {
	id string
}

func (cs *DataAccessorNode) SetID(id string) {
	cs.id = id
}

func (cs DataAccessorNode) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs DataAccessorNode) Name() string {
	return "DataAccessorNode"
}

func (cs DataAccessorNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Info("Enter DataAccessorNode node")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, GET_TAXCLASS_FETCH_NODE)
	defer logger.EndProfile(profiler, GET_TAXCLASS_FETCH_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(DEBUG_KEY_NODE, "DataAccessorNode")
	tClasses, err := cs.GetAllTaxClasses()
	if err != nil {
		return io, &constants.AppError{
			Code:             constants.IncorrectDataErrorCode,
			Message:          "cannot Load data",
			DeveloperMessage: err.Error(),
		}
	}
	io.IOData.Set(IODATA, tClasses)
	logger.Info("Exit DataAccessorNode node")
	return io, nil
}

func (cs DataAccessorNode) GetAllTaxClasses() ([]proUtil.TaxClass, error) {
	session := factory.GetMongoSession("TaxClass")
	defer session.Close()

	tClass := []proUtil.TaxClass{}
	err := session.SetCollection(proUtil.TAXCLASS_COLLECTION).
		Find(nil).Limit(100).All(&tClass)
	if err != nil {
		return tClass, err
	}
	return tClass, nil
}
