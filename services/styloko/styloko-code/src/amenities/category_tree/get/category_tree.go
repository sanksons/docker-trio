package get

import (
	"amenities/categories/common"
	"common/appconstant"
	"common/constants"
	"common/utils"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"common/ResourceFactory"

	"github.com/jabong/floRest/src/common/cache"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// CategoryTreeGet -> struct for node based data
type CategoryTreeGet struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CategoryTreeGet) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CategoryTreeGet) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CategoryTreeGet) Name() string {
	return "CategoryTreeGet"
}

// Execute -> Starts node execution.
func (cs CategoryTreeGet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	getCacheFlag := true
	cacheHeader, err := utils.GetRequestHeader(io, "No-Cache")
	if strings.ToLower(cacheHeader) == "true" {
		getCacheFlag = false
	}
	if getCacheFlag {
		tree, err := cs.getFromCache()
		if err == nil {
			io.IOData.Set(florest_constants.RESULT, tree)
			return io, nil
		}
	}
	categories, err := cs.getCategories()
	if err != nil {
		return io, &florest_constants.AppError{
			Code:             appconstant.DataNotFoundErrorCode,
			Message:          "Cannot get active categories from database",
			DeveloperMessage: err.Error(),
		}
	}
	data, _ := cs.preOrderTrav(1, categories)

	// Set category tree to cache
	go cs.setCache(data)

	io.IOData.Set(florest_constants.RESULT, data)
	return io, nil
}

// getCategories returns all categories from Mongo
func (cs CategoryTreeGet) getCategories() ([]common.CategoryGetVerbose, error) {
	var ctgryStruct []common.CategoryGetVerbose
	mgoSession := ResourceFactory.GetMongoSession(constants.CATEGORY_SEARCH)
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(constants.CATEGORY_COLLECTION)
	err := mgoObj.Find(bson.M{"status": "active"}).All(&ctgryStruct)
	if err != nil {
		return ctgryStruct, err
	}
	return ctgryStruct, nil
}

// preOrderTrav traverses the tree in recursive preorder manner
// to start traversal from root, please use starting id as 1
func (cs CategoryTreeGet) preOrderTrav(id int, categories []common.CategoryGetVerbose) (map[string]interface{}, int) {
	def := make(map[string]interface{}, 0)
	if !cs.checkElement(id, categories) {
		return nil, id + 1
	}
	item := cs.getElement(id, categories)
	def["name"] = item.Name
	def["id"] = item.CategoryId
	def["nodes"] = nil
	def["displaySizeConversion"] = item.DisplaySizeConversion
	def["googleTreeMapping"] = item.GoogleTreeMapping
	def["segment"] = item.CategorySeg
	def["parent"] = item.Parent
	def["pdfActive"] = item.PdfActive
	def["sizeChartActive"] = item.SizeChartAcive
	def["sizeChartApplicable"] = item.SizeChartApplicable
	def["status"] = item.Status
	if math.Abs(float64(item.Right-item.Left)) == 1 {
		return def, item.Right
	}
	base := id
	var tmp []map[string]interface{}
	for item.Right-1 > base {
		tmp1, tmp2 := cs.preOrderTrav(base+1, categories)
		base = tmp2
		if tmp1 == nil {
			continue
		}
		tmp = append(tmp, tmp1)
	}
	def["nodes"] = tmp
	return def, item.Right
}

// getElement returns the element from the tree array
func (cs CategoryTreeGet) getElement(id int, categories []common.CategoryGetVerbose) common.CategoryGetVerbose {
	var tmp common.CategoryGetVerbose
	for x := range categories {
		if categories[x].Left == id {
			tmp = categories[x]
		}
	}
	return tmp
}

// checkElement checks for element in case there are missing keys
func (cs CategoryTreeGet) checkElement(id int, categories []common.CategoryGetVerbose) bool {
	for x := range categories {
		if categories[x].Left == id {
			return true
		}
	}
	return false
}

// setCache sets category tree to cache
func (cs CategoryTreeGet) setCache(data interface{}) {
	e, _ := json.Marshal(data)
	i := cache.Item{
		Key:   "CATEGORY_TREE",
		Value: string(e),
	}
	err := cacheObj.Set(i, false, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in setting cache: %s", err.Error()))
	}
}

// getFromCache retrieves from blitz
func (cs CategoryTreeGet) getFromCache() (interface{}, error) {
	item, err := cacheObj.Get("CATEGORY_TREE", false, false)
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}
	var v map[string]interface{}
	json.Unmarshal([]byte(item.Value.(string)), &v)
	return v, nil
}
