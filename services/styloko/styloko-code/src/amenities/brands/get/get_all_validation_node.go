package get

import (
	"amenities/brands/common"
	"common/utils"
	"errors"
	"fmt"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

//Struct for GetAll Decision Node
type GetAllDecision struct {
	id string
}

//Function to SetID for current node from orchestrator
func (d *GetAllDecision) SetID(id string) {
	d.id = id
}

//Function that returns current node ID to orchestrator
func (d GetAllDecision) GetID() (id string, err error) {
	return d.id, nil
}

//Function that returns node name to orchestrator
func (d GetAllDecision) Name() string {
	return "GetAllDecision node for GET"
}

//Function that returns bool value to trigger GetAll or Search Nodes
func (d GetAllDecision) GetDecision(io workflow.WorkFlowData) (bool, error) {

	//Enable profiling to track performance
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.BRAND_GET)
	defer func() {
		logger.EndProfile(profiler, common.BRAND_GET)
	}()
	io.ExecContext.Set(florest_constants.MONITOR_CUSTOM_METRIC_PREFIX, common.CUSTOM_BRAND_GET_ALL)
	//Reading path params
	searchMap, _, err := utils.GetSearchQueries(io)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting search map from search queries :%s", err.Error()))
		return false, err
	}
	if len(searchMap) == 0 {
		return true, nil
	}
	bsonMap, err := d.GetBsonMapFromSearchMap(searchMap)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting bson map from search map :%s", err.Error()))
		return false, err
	}
	io.IOData.Set(common.BRAND_SEARCH, bsonMap)
	return false, nil
}

//Generic function to parse queryMap into bsonMap that can be used to search in mongo
func (d GetAllDecision) GetBsonMapFromSearchMap(queryMap []map[string]interface{}) (map[string]interface{}, error) {
	bsonMap := make(map[string]interface{})
	for _, v := range queryMap {
		for x, y := range v {
			xNew := d.GetBrandMapping(x)
			val := y.(map[string]interface{})
			if x == "" {
				return nil, errors.New("Key field missing")
			}
			if val["value"] == nil {
				return nil, errors.New("Value field missing")
			}
			switch val["operator"] {
			case "eq":
				bsonMap[xNew] = val["value"]
				break
			case "in":
				bsonMap[xNew] = bson.M{"$in": val["value"]}
				break
			default:
				return nil, errors.New("Incorrect value for operator")
			}
		}
	}
	return bsonMap, nil
}

func (d GetAllDecision) GetBrandMapping(key string) string {
	brandMapping := make(map[string]string)
	brandMapping["id"] = "seqId"
	brandMapping["name"] = "name"
	brandMapping["status"] = "status"
	brandMapping["url"] = "urlKey"
	brandMapping["class"] = "brandClass"
	return brandMapping[key]
}
