package put

import (
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"sellers/common"
)

//takes input the data that was updated in mongo and
//prepares and sends data to erp
func UpdateOnErp(data interface{}) error {
	res := data.([]common.Schema)
	//preparing data in required format by adapter
	erpData := make([]common.ErpData, 1)
	erpData[0] = common.ErpData{
		Method: "ERP.insertSeller",
		Params: common.ErpSellerData{
			SellerData: res,
		},
	}
	//sending data for updation in erp
	err := common.SendDataToErp(erpData)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while sending data to ERP : %s", err.Error()))
		return err
	}
	return nil
}
