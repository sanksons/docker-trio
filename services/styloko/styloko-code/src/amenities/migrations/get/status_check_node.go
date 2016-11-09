package get

import (
	"amenities/migrations/common"
	"common/appconstant"
	"common/utils"

	florest_constants "github.com/jabong/floRest/src/common/constants"

	workflow "github.com/jabong/floRest/src/common/orchestrator"
)

// StatusCheck -> struct for node based data
type StatusCheck struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (mg *StatusCheck) SetID(id string) {
	mg.id = id
}

// GetID -> returns current node ID to orchestrator
func (mg StatusCheck) GetID() (id string, err error) {
	return mg.id, nil
}

// Name -> Returns node name to orchestrator
func (mg StatusCheck) Name() string {
	return "StatusCheck"
}

// Execute -> Executes the current workflow
func (mg StatusCheck) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	path, _ := io.IOData.Get(PATH)
	pathStr, _ := path.(string)
	if pathStr == "delete" {
		err := common.ClearAllRedisFlags()
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to delete flag"}
		}
		io.IOData.Set(florest_constants.RESULT, "Deleted all redis flags")
		return io, nil
	}
	if pathStr == "status" {
		key, ok := utils.GetQueryParams(io, "key")
		if !ok || len(key) < 1 {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Missing key", DeveloperMessage: "Please provide a valid key."}
		}

		validKeys := common.GetStatusKeys()
		found := false
		for _, x := range validKeys {
			if key == x {
				found = true
				break
			}
		}
		if !found {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid key", DeveloperMessage: "Please provide a valid key."}
		}
		flag, err := common.GetFlagFromRedis(key)
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to get flag"}
		}
		m := make(map[string]interface{})
		if flag {
			m["key"] = key
			m["status"] = "running"
		} else {
			m["key"] = key
			m["status"] = "finished"
		}

		io.IOData.Set(florest_constants.RESULT, m)
		return io, nil
	}

	return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid path param", DeveloperMessage: "Please provide a valid path param."}

}

// PATH constant for path vairable
const PATH = "PATH"
