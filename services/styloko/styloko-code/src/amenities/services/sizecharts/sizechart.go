package sizecharts

import (
	proUtil "amenities/products/common"
	sizeUtils "amenities/sizechart/common"
	"common/ResourceFactory"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strconv"
)

// This function returns sizechart for product while productsizechart migration
func FetchSizeChart(prd proUtil.Product) (interface{}, error) {
	driver, err := ResourceFactory.GetMySqlDriver("ProductMigration")
	if err != nil {
		logger.Error("#FetchSizeChart(): Unable to get mysql driver.", err.Error())
		return nil, err
	}
	query := `SELECT fk_catalog_distinct_sizechart from catalog_config_additional_info ` +
		` where fk_catalog_config = ?`
	var id *int
	rows, getErr := driver.Query(query, prd.SeqId)
	if getErr != nil {
		logger.Error("#FetchSizeChart(): Unable to get sizechart for sku.", getErr.DeveloperMessage)
		return nil, fmt.Errorf(fmt.Sprintf("#FetchSizeChart(): Unable to get sizechart for sku. %s", getErr.DeveloperMessage))
	}
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			logger.Error("#FetchSizeChart(): Unable to get sizechart for sku.", err.Error())
		}
	}
	rows.Close()
	if id == nil {
		logger.Info("#FetchSizeChart(): sizechart doesnot exist for sku:", prd.SKU)
		return nil, fmt.Errorf("#FetchSizeChart()1: sizechart doesnot exist for sku:", prd.SKU)
	}
	sizechart := sizeUtils.GetSizeChartById(*id)
	if sizechart == nil {
		logger.Info("#FetchSizeChart(): Sizechart doesnot exist for sku", prd.SKU)
		return nil, fmt.Errorf("#FetchSizeChart()2: Sizechart doesnot exist for sku", prd.SKU)
	}
	sChartPdp, err := CreateDisplayableSizeChart(prd.SKU, -1, (*sizechart).SizeChartType, *sizechart, prd)
	if sChartPdp == nil {
		return nil, err
	}
	sizeChProd := proUtil.ProSizeChart{Id: (*sizechart).IdCatalogSizeChart, Data: *sChartPdp}
	return sizeChProd, nil
}

// This function returns sizechart for product while migration
func GetSizeChartForProductMigration(prd proUtil.Product) interface{} {
	driver, err := ResourceFactory.GetMySqlDriver("ProductMigration")
	if err != nil {
		logger.Error("#GetSizeChartForProductMigration(): Unable to get mysql driver.", err.Error())
	}
	query := `SELECT fk_catalog_distinct_sizechart from catalog_config_additional_info ` +
		` where fk_catalog_config = ?`
	var id *int
	rows, getErr := driver.Query(query, prd.SeqId)
	if getErr != nil {
		logger.Error("#GetSizeChartForProductMigration(): Unable to get sizechart for sku.", getErr.DeveloperMessage)
	}
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			logger.Error("#GetSizeChartForProductMigration(): Unable to get sizechart for sku.", err.Error())
		}
	}
	rows.Close()
	if id == nil {
		logger.Info("#GetSizeChartForProductMigration(): sizechart doesnot exist for sku:", prd.SKU)
		return nil
	}
	sizechart := sizeUtils.GetSizeChartById(*id)
	if sizechart == nil {
		logger.Info("#GetSizeChartForProductMigration(): Sizechart doesnot exist for sku", prd.SKU)
		return nil
	}
	sChartPdp, _ := CreateDisplayableSizeChart(prd.SKU, -1, (*sizechart).SizeChartType, *sizechart, prd)
	if sChartPdp == nil {
		return nil
	}
	sizeChProd := proUtil.ProSizeChart{Id: (*sizechart).IdCatalogSizeChart, Data: *sChartPdp}
	return sizeChProd
}

//This function Returns the applicable sizechart for the given sku
func GetSizeChartForProduct(prd proUtil.Product) interface{} {
	categoryId := prd.PrimaryCategory
	brandId := prd.BrandId
	ty := prd.TY
	sku := prd.SKU
	SizeChartPriority := []int{1, 2}
	var result *sizeUtils.SizeChartMongo

	// for given product check if sku level sizechart applicable
	result = sizeUtils.GetSkuLevelSizeChartProd(prd.SKU)
	if result == nil {
		// gets the sizechart for product based on priority of SC
		for _, scTy := range SizeChartPriority {
			result = sizeUtils.CheckGivenSizeChartForProduct(categoryId, brandId, ty, scTy)
			if result != nil {
				break
			}
		}
	}

	// case when ty exist for product but sizechart for ty doesnot
	// Look for sizechart without ty field
	if result == nil && ty != 0 {
		for _, scTy := range SizeChartPriority {
			result = sizeUtils.CheckGivenSizeChartForProduct(categoryId, brandId, 0, scTy)
			if result != nil {
				break
			}
		}

	}
	if result == nil {
		logger.Error(fmt.Sprintf("GetSizeChart Service : SizeChart doesnot exist for the sku #%s", sku))
		return nil
	}
	sChartPdp, _ := CreateDisplayableSizeChart(sku, -1, (*result).SizeChartType, *result, prd)
	if sChartPdp == nil {
		return nil
	}
	sizeChProd := proUtil.ProSizeChart{Id: (*result).IdCatalogSizeChart, Data: *sChartPdp}
	return sizeChProd
}

// create a displayable PDP page sizechart for particular SKU
func CreateDisplayableSizeChart(sku string, preTy int, currTy int, doc sizeUtils.SizeChartMongo, prd proUtil.Product) (*sizeUtils.SChart, error) {
	k := 0
	mismatch := true
	upload := false
	chart := sizeUtils.SChart{}
	sizeIndexMapping := make(map[string]int)
	chart.Sizes = make(map[string][]string)

	// Type Mapping
	sizeChartTypeMap := make(map[int]string)
	sizeChartTypeMap = map[int]string{0: "sku", 1: "brand", 2: "brick"}

	// Rules to create sizechart for product,
	// depending upon the previously uploaded sizechart type and currently being uploaded
	switch {
	case preTy == -1, preTy == currTy, currTy < preTy:
		upload = true
	}

	if upload == false {
		logger.Info("#createDisplayableSizeChart() ->Sizechart of higher priority already exists for sku " + sku)
		return nil, fmt.Errorf("Sizechart of higher priority already exists for sku " + sku)
	}

	sizeSku, err := GetSimpleSizesForProd(prd)

	if err != nil {
		logger.Error("#createDisplayableSizeChart() ->Sizes not found for the sku " + sku)
		return nil, fmt.Errorf("Sizes not found for the sku " + sku)
	}
	chart.Headers = getHeader(doc.SizeChartInfo)
	chart.ImageName = doc.SizeChartImagePath
	chart.ScType = sizeChartTypeMap[currTy]
	for _, size := range sizeSku {
		chart.Sizes[strconv.Itoa(k)] = []string{size}
		sizeIndexMapping[size] = k
		k++
	}

	for _, data := range doc.SizeChartInfo {
		if data.RowHeaderType != "" {
			if index, ok := sizeIndexMapping[data.RowHeaderType]; ok {
				chart.Sizes[strconv.Itoa(index)] = append(chart.Sizes[strconv.Itoa(index)], data.Value)
				mismatch = false
			}
		}
	}
	if mismatch {
		logger.Info("#createDisplayableSizeChart() ->Mismatch of sizes for sku " + sku + ".")
		return nil, fmt.Errorf("Mismatch of sizes for sku " + sku + ".")
	}
	return &chart, nil
}

// This function obtains the sizes for all simples of product
func GetSimpleSizesForProd(prd proUtil.Product) ([]string, error) {
	simpleSizesArray := []string{}
	prd.SortSimples()
	prodSimpleArray := prd.Simples
	sizeAttributeName, err := prd.AttributeSet.GetVariationAttributeName()

	if err != nil {
		logger.Error("Unable to get the size attribute Name: ", err)
		return simpleSizesArray, err
	}
	for _, simple := range prodSimpleArray {
		simpleSize := simple.GetSize(sizeAttributeName)
		simpleSizesArray = append(simpleSizesArray, simpleSize)
	}
	return simpleSizesArray, nil
}

// This function returns the header of sizechart
func getHeader(data []sizeUtils.SizeChartData) []string {
	var headers []string
	checkDup := make(map[string]int)

	for _, sizeChartData := range data {
		header := sizeChartData.ColumnHeader

		if _, ok := checkDup[header]; !ok {
			headers = append(headers, header)
			checkDup[header] = 1
		}
	}
	return headers
}
