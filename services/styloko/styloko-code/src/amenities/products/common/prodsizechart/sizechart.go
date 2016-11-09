package prodsizechart

import (
	proUtil "amenities/products/common"
	sizeChService "amenities/services/sizecharts"
	"common/utils"
	"fmt"
	"reflect"

	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
)

// This function updates the sizechrt for product on creation/updation
func UpdateProductWithSizeChart(p proUtil.Product) (proUtil.Product, error) {
	sizeChart := sizeChService.GetSizeChartForProduct(p)
	var ok bool

	if sizeChart == nil {
		return p, fmt.Errorf(
			"UpdateProductWithSizeChart()[productId:%d]#: SizeChart doesnot exist for this product",
			p.SeqId,
		)
	}
	p.SizeChart, ok = sizeChart.(proUtil.ProSizeChart)
	if !ok {
		logger.Error(fmt.Sprintf(
			"#UpdateProductWithSizeChart(): Sizechart assertion failed, looking for proUtil.ProSizeChart but found %v",
			reflect.TypeOf(sizeChart),
		))
	}
	return p, nil
}

// This function recalculate the sizechart considering deleted/inactive simples
func CalculateSizeChart(p proUtil.Product) error {
	if p.SizeChart.Data == nil {
		return nil
	}
	simpleSizes, err := p.GetInactiveDeletedSimpleSizes()
	if err != nil {
		return fmt.Errorf("#CalculateSizeChart: Unable to get simpleSizes %s", err.Error())
	}
	if len(simpleSizes) == 0 {
		return nil
	}
	// process the sizechart
	var count int
	var result map[string]interface{}
	result = make(map[string]interface{})

	sizechartMap := p.SizeChart.Data.(bson.M)
	sizechartData := sizechartMap["sizes"]
	sizeDataMap := sizechartData.(bson.M)
	sortedData := make([][]interface{}, len(sizeDataMap))

	for mapIndex, mapData := range sizeDataMap {
		intIndex, err := utils.GetInt(mapIndex)
		if err != nil {
			return fmt.Errorf("CalculateSizeChart: err in get int %s", err.Error())
		}
		sortedData[intIndex] = mapData.([]interface{})
	}
	for _, rowData := range sortedData {
		// check if this size is to be retained/remove from sizechart
		if !utils.InArrayString(simpleSizes, rowData[0].(string)) {
			result[utils.ToString(count)] = rowData
			count++
		}
	}
	sizechartMap["sizes"] = result
	return nil
}
