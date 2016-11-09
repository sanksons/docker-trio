package post

import (
	"amenities/categories/common"
	mongo "common/ResourceFactory"
	"common/appconstant"
	"common/constants"
	"fmt"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/jabong/floRest/src/common/utils/logger"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
)

// CategoriesCreate -> struct for node based data
type CategoriesCreate struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CategoriesCreate) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CategoriesCreate) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CategoriesCreate) Name() string {
	return "CategoriesCreate"
}

// Execute -> Starts node execution.
// TODO: Caching strategy.
func (cs CategoriesCreate) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, common.POST)
	defer func() {
		logger.EndProfile(profiler, common.POST)
	}()
	logger.Info("Create categories started")
	tmp, _ := io.IOData.Get(common.CATEGORY_VALID_DATA)
	ctgryCreate, _ := tmp.(*common.CategoryCreate)

	// Segment insertion into struct
	segIds := common.InsertSegments(ctgryCreate.SegIds)
	ctgryCreate.CategorySeg = segIds

	res, ok, err := cs.create(*ctgryCreate)
	if err != nil {
		return io, &florest_constants.AppError{Code: appconstant.ServiceFailureCode, Message: "Mongo Error.", DeveloperMessage: err.Error()}
	}
	if !ok {
		return io, &florest_constants.AppError{Code: appconstant.FailedToCreateErrorCode, Message: "Create failed.", DeveloperMessage: "Create in MongoDB failed."}
	}
	io.IOData.Set(florest_constants.RESULT, res)
	// TODO Rebuild Tree Has to be done. Without that insert is as good as useless.
	logger.Info("Create categories finished")
	go cs.deleteTree()
	return io, nil
}

func (cs CategoriesCreate) create(ctgryStruct common.CategoryCreate) (interface{}, bool, error) {
	logger.Info("Get category by ID started")

	// Mongo session and ID get
	mgoSession := mongo.GetMongoSession(constants.CATEGORY_CREATE)
	defer mgoSession.Close()
	seqID := mgoSession.GetNextSequence(constants.CATEGORY_COLLECTION)
	mgoObj := mgoSession.SetCollection(constants.CATEGORY_COLLECTION)

	// Parent Node
	parentID := ctgryStruct.Parent
	var parent common.CategoryCreate
	err := mgoObj.Find(bson.M{"seqId": parentID}).One(&parent)
	if err != nil {
		logger.Info("Invalid parent ID provided.")
		logger.Debug(err)
		return nil, false, err
	}

	// First inserting new category into mongo with parent ID.
	ctgryStruct.Id = seqID
	ctgryStruct.Left = parent.Right
	ctgryStruct.Right = parent.Right + 1
	updateQuery := bson.M{"$set": ctgryStruct}
	query := bson.M{"seqId": seqID}
	res, err := mgoSession.FindAndModify(constants.CATEGORY_COLLECTION, updateQuery, query, true)
	if err != nil {
		logger.Info("Mongo insert failure")
		logger.Debug(err)
		return nil, false, err
	}
	logger.Info("Mongo insert success. New category created.")
	logger.Info("Tree rebuild starts.")

	// Tree rebuild has two parts. Parent subtree rebuild and parent right subtree rebuild.
	// PART 1
	// Rebuild parents below
	// Algorithm -> if lft<=parent.lft & rgt>=parent.rgt -> increase right by 2

	var parentTree []common.CategoryCreate
	err = mgoObj.Find(bson.M{"lft": bson.M{"$lte": parent.Left}, "rgt": bson.M{"$gte": parent.Right}}).All(&parentTree)
	for index := range parentTree {
		parentTree[index].Right += 2
		err = mgoObj.Update(bson.M{"seqId": parentTree[index].Id}, bson.M{"$set": parentTree[index]})
		if err != nil {
			logger.Info("Category rebuild process failed at parent rebuild.")
			logger.Debug(err)
			return nil, false, err
		}
	}
	logger.Info("Parents modified.")

	// PART 2
	// Rebuild right subtree below.
	// Algorithm -> if lft>parent.rgt -> increase right & left by 2

	var rightSubtree []common.CategoryCreate
	err = mgoObj.Find(bson.M{"lft": bson.M{"$gt": parent.Right}}).All(&rightSubtree)
	for index := range rightSubtree {
		rightSubtree[index].Right += 2
		rightSubtree[index].Left += 2
		err = mgoObj.Update(bson.M{"seqId": rightSubtree[index].Id}, bson.M{"$set": rightSubtree[index]})
		if err != nil {
			logger.Info("Category rebuild process failed at right subtree rebuild.")
			logger.Debug(err)
			return nil, false, err
		}
	}
	logger.Info("Right sub tree modified.")

	// SQL data preperation for worker
	sqlData := cs.prepData(seqID, ctgryStruct, parent.Right, parent.Left)

	// Worker start for MySQL data push in background.
	categoryCreatePool.StartJob(sqlData)

	return res, true, nil
}

// prepData -> prepares data in a map[string]interface{} to send to the parallel worker
func (cs CategoriesCreate) prepData(seqID int, ctgryStruct common.CategoryCreate, right int, left int) map[string]interface{} {
	sqlData := make(map[string]interface{})
	sqlData["id"] = seqID
	sqlData["segIds"] = ctgryStruct.SegIds
	sqlData["categoryCreate"] = ctgryStruct
	sqlData["parentRight"] = strconv.Itoa(right)
	sqlData["parentLeft"] = strconv.Itoa(left)
	return sqlData
}

func (cs CategoriesCreate) deleteTree() {
	err := cacheObj.Delete("CATEGORY_TREE")
	if err != nil {
		logger.Error(fmt.Sprintf("Error while deleting category tree from cache: %v", err.Error()))
	}
}
