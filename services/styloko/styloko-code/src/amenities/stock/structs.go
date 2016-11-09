package stock

import (
	"common/utils"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
)

const REDIS_STOCK_KEY = "styloko_stock_"

type StockRequest struct {
	Simple   string `json:"simple"`
	Total    string `json:"totalqty"`
	Reserved string `json:"reservedqty"`
}

func (stock StockRequest) GetStock() int {
	total, err := utils.GetInt(stock.Total)
	if err != nil {
		logger.Error(fmt.Sprintf("(stock StockRequest)#GetStock()1: %s", err.Error()))
		total = 0
	}
	reserved, err := utils.GetInt(stock.Reserved)
	if err != nil {
		logger.Error(fmt.Sprintf("(stock StockRequest)#GetStock()2: %s", err.Error()))
		reserved = 0
	}
	return (total - reserved)
}
