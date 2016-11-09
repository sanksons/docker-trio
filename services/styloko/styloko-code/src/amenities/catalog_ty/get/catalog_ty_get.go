package get

import (
	"common/ResourceFactory"
	"common/appconstant"
	"common/constants"
	"common/utils"
	"log"
	"strconv"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
)

// CatalogTyGet -> struct for node based data
type CatalogTyGet struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *CatalogTyGet) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs CatalogTyGet) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs CatalogTyGet) Name() string {
	return "CatalogTyGet"
}

// Execute -> Starts node execution.
func (cs CatalogTyGet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	qParams, _ := utils.GetQueryParams(io, "id")
	if len(qParams) < 1 {
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid ID prodided", DeveloperMessage: "Please provide a valid ID"}
	}

	_, err := strconv.Atoi(qParams)
	if err != nil {
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid ID prodided", DeveloperMessage: "Please provide a valid ID"}
	}

	driver, err := ResourceFactory.GetMySqlDriver(constants.CATALOG_TY_API)
	if err != nil {
		return io, &florest_constants.AppError{Code: appconstant.ServiceFailureCode, Message: "Cannot connect to mysql", DeveloperMessage: "Cannot acquire mysql resource from resource factory."}
	}
	query := "select * from catalog_ty where id_catalog_ty=ANY (select fk_catalog_ty from catalog_category_ty where fk_catalog_category=" + qParams + ");"
	rows, serr := driver.Query(query)
	defer rows.Close()
	if serr != nil {
		return io, &florest_constants.AppError{Code: appconstant.ServiceFailureCode, Message: "Some error occured with mysql. Code: " + serr.ErrCode, DeveloperMessage: serr.DeveloperMessage}
	}
	var results []interface{}
	for rows.Next() {
		data := make([]string, 5)
		if err := rows.Scan(&data[0], &data[1], &data[2], &data[3], &data[4]); err != nil {
			log.Fatal(err)
		}
		dataMap := make(map[string]interface{}, 0)
		id, _ := strconv.Atoi(data[0])
		dataMap["seqId"] = id
		dataMap["name"] = data[2]
		dataMap["url_key"] = data[4]
		results = append(results, dataMap)
	}
	rows.Close()
	io.IOData.Set(florest_constants.RESULT, results)
	return io, nil
}
