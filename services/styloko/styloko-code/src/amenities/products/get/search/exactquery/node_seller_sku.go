package exactquery

import (
	proUtil "amenities/products/common"
	searchUtil "amenities/products/get/search"
	"common/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

const (
	SELLER_ID_SKU = "sellerIdSku"
	SUCCESS       = "success"
	FAILED        = "failed"
)

type SellerIdSku struct {
	SellerSku string `json:"sellerSku"`
	SellerId  int    `json:"sellerId"`
	SimpleSku string `json:"simpleSku"`
}

// Checks if the query we have is of SellerSkuApi type.
type SellerSkuApi struct {
	id string
}

func (ss *SellerSkuApi) SetID(id string) {
	ss.id = id
}

func (ss SellerSkuApi) GetID() (id string, err error) {
	return ss.id, nil
}

func (ss SellerSkuApi) Name() string {
	return "SellerSkuApi"
}

func (ss SellerSkuApi) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Debug("Enter SellerSkuApi Execute node")
	sellerIdSkuMap, err := ss.validateSellerSkuApi(io)
	if err != nil {
		return io, &constants.AppError{
			Code:             constants.ParamsInValidErrorCode,
			Message:          "Validation Error",
			DeveloperMessage: err.Error(),
		}
	}
	response, err := ss.getSellerIdSku(sellerIdSkuMap)
	if err != nil {
		return io, &constants.AppError{
			Code:             constants.DbErrorCode,
			Message:          "Database Error",
			DeveloperMessage: err.Error(),
		}
	}
	io.IOData.Set(constants.RESULT, response)
	logger.Debug("Exit sellerSkuApi Execute node")
	return io, nil
}

//Check if it is a valid seller id sku query
func (ss SellerSkuApi) validateSellerSkuApi(io workflow.WorkFlowData) (map[int][]string, error) {
	valStr, ok := utils.GetQueryParams(io, SELLER_ID_SKU)
	if !ok {
		return nil, fmt.Errorf("Key: %v missing in request.", SELLER_ID_SKU)
	}
	valStr = strings.Replace(valStr, "[", "", -1)
	valStr = strings.Replace(valStr, "]", "", -1)
	if valStr == "" {
		return nil, fmt.Errorf("Array cannot be empty")
	}
	valStrArr := strings.Split(valStr, ",")
	if len(valStrArr) > searchUtil.SellerSkuLimit {
		return nil, fmt.Errorf("Array length exceeded limit")
	}
	retMap := make(map[int][]string, 0)
	for _, val := range valStrArr {
		valFinal := strings.Split(val, ":")
		if len(valFinal) != 2 {
			return nil, fmt.Errorf("Seller Id, Seller Sku pair incorrect: %v", val)
		}
		sellerId, err := strconv.Atoi(valFinal[0])
		if err != nil {
			return nil, fmt.Errorf("SellerId: %v invalid", valFinal[0])
		}
		retMap[sellerId] = append(retMap[sellerId], valFinal[1])
	}
	return retMap, nil
}

func (ss SellerSkuApi) getSellerIdSku(sellerIdSkuMap map[int][]string) (map[string][]SellerIdSku, error) {
	success := make([]SellerIdSku, 0)
	failed := make([]SellerIdSku, 0)

	for sellerId, sellerSkuArr := range sellerIdSkuMap {
		result := []proUtil.ProductSmallSimples{}
		result, err := proUtil.GetAdapter(searchUtil.DbAdapterName).
			GetProductBySellerIdSku(sellerId, sellerSkuArr)
		if err != nil {
			return nil, err
		}
		sellerSkuMap := make(map[string]bool, 0)
		for _, sellerSku := range sellerSkuArr {
			sellerSkuMap[sellerSku] = false
		}
		for _, productSmall := range result {
			for _, simple := range productSmall.Simples {
				if _, ok := sellerSkuMap[simple.SellerSku]; ok {
					success = append(success,
						SellerIdSku{SellerSku: simple.SellerSku,
							SellerId: sellerId, SimpleSku: simple.Sku})
					sellerSkuMap[simple.SellerSku] = true
				}
			}
		}
		for sellerSku, found := range sellerSkuMap {
			if !found {
				failed = append(failed,
					SellerIdSku{SellerSku: sellerSku, SellerId: sellerId})
			}
		}
	}
	response := make(map[string][]SellerIdSku, 0)
	response[SUCCESS] = success
	response[FAILED] = failed
	return response, nil
}
