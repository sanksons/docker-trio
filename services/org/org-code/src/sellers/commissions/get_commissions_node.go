package commissions

import (
	"common/ResourceFactory"
	"common/appconstant"
	"common/notification"
	"errors"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"sellers/common"
	"strconv"
	"strings"
)

type GetCommissions struct {
	id string
}

func (s *GetCommissions) SetID(id string) {
	s.id = id
}

func (s GetCommissions) GetID() (id string, err error) {
	return s.id, nil
}

func (s GetCommissions) Name() string {
	return "GET commission"
}

func (s GetCommissions) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SELLER_GET_COMMISSION)
	defer func() {
		logger.EndProfile(profiler, SELLER_GET_COMMISSION)
	}()
	rc, _ := io.ExecContext.Get(florest_constants.REQUEST_CONTEXT)
	logger.Info("entered "+s.Name(), rc)
	io.ExecContext.SetDebugMsg(SELLER_GET_COMMISSION, "Seller get commissions execution started")

	data, err := io.IOData.Get(GET_SEARCH)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting search request data: %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Error while getting search request data", DeveloperMessage: err.Error()}
	}
	// parse query
	dataMap := data.(map[string]interface{})
	if len(dataMap) == 0 {
		logger.Error("Error while getting search params.")
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Improper Format for Search String"}
	}

	strSeqId, _ := dataMap["seqId"]
	stringSeqId := strSeqId.(string)
	intSeqId, _ := strconv.Atoi(stringSeqId)
	resp, err := common.GetById(intSeqId, "Get_Commission")
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting Id Details : %s", err.Error()))
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid Seller Id", DeveloperMessage: "No data found for the passed id."}
	}
	//get update commission value from slave db of SC
	dataVal, _ := dataMap["categories"]
	y := dataVal.(bson.M)
	stringCatIds := y["$in"].(string)
	resp, err = s.GetSellerWithCommissions(resp, stringSeqId, stringCatIds)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in GetUpdateCommission :%s", err.Error()))
		notification.SendNotification("Error while getting commission", err.Error(), nil, "error")
	}

	io.IOData.Set(florest_constants.RESULT, resp)
	return io, nil
}

//takes input schema without update commiission values, gets values from mysql
//appends and returns schema with update commission values
func (s GetCommissions) GetSellerWithCommissions(slrData common.Schema, seqId string, catIds string) (common.Schema, error) {
	comMap, err := s.GetCommissions(seqId, catIds)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting commission: %s", err.Error()))
		return slrData, err
	}
	//logic here
	catArr := common.ParseQuery(catIds)
	for _, v := range catArr {
		if comMap[v] != nil {
			slrData.UpdateCommission = comMap[v]
			return slrData, nil
		}
	}
	return slrData, nil
}

//takes an id and gets commission from the seller centre slave database for the passed id
func (s GetCommissions) GetCommissions(seqid string, catIds string) (map[int][]common.Commission, error) {
	comData := []common.GetCommission{}
	catIds = strings.Replace(catIds, "[", "", -1)
	catIds = strings.Replace(catIds, "]", "", -1)
	sql := `SELECT
                          seller.src_id AS sellerId,
                          catalog_category.src_id AS categoryId,
                          seller_commission.percentage
                        FROM seller_commission
                        INNER JOIN catalog_category
                          ON catalog_category.id_catalog_category = seller_commission.fk_catalog_category
                        INNER JOIN seller
                          ON seller.id_seller=seller_commission.fk_seller
                        WHERE seller.src_id = ` + seqid + ` AND catalog_category.src_id IN (` + catIds + `)`
	driver, errs := ResourceFactory.GetMySqlDriverSC(UPDATE_COMMISSIONS)
	if errs != nil {
		logger.Error(fmt.Sprintf("Couldn't acquire mysql resource. Error: %s", errs.Error()))
		return nil, errs
	}
	rows, err := driver.Query(sql)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while executing query: %v", err))
		return nil, errors.New("Error while executing query to get commission value")
	}
	defer rows.Close()
	for rows.Next() {
		com := common.GetCommission{}
		err := rows.Scan(&com.SeqId,
			&com.CategoryId,
			&com.Percentage)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while scanning commiission values: %v", err))
			return nil, err
		}
		comData = append(comData, com)
	}
	comMap, er := s.GetCommissionMap(comData)
	if er != nil {
		logger.Error(fmt.Sprintf("Error while getting map of commission data: %s", er.Error()))
		return nil, er
	}
	return comMap, nil
}

//populates the recieved commission into struct to be visible in api response
func (s GetCommissions) GetCommissionMap(comData []common.GetCommission) (map[int][]common.Commission, error) {
	comMap := make(map[int][]common.Commission)
	for _, v := range comData {
		comArr := common.Commission{}
		comArr.CategoryId = v.CategoryId
		comArr.Percentage = v.Percentage
		comMap[v.CategoryId] = append(comMap[v.CategoryId], comArr)
	}
	return comMap, nil
}
