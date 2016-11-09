package post

import (
	"amenities/services/products"
	sizeUtils "amenities/sizechart/common"
	factory "common/ResourceFactory"
	mongodb "common/mongodb"
	"common/notification"
	"common/utils"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
)

// validates the sku level sizechart and convert it to brandwise data, also returns
// array of sku for which validation failed
func validateSkuSizeChart(sizechart sizeUtils.SizeChart) ([]sizeUtils.BrandWiseScData, []string, []string) {
	var failedSku []string
	var successSku []string
	failedSku = []string{}
	successSku = []string{}
	var lastInputSku = ""
	skuWiseSC := [][]string{}
	var brandWiseScArr []sizeUtils.BrandWiseScData
	for _, sizeChartRow := range sizechart.SChartData {
		inputSku := sizeChartRow[0]
		inputColumnHeader := sizeChartRow[1]
		inputRowHeaderName := sizeChartRow[2]
		inputRowHeaderType := sizeChartRow[3]
		inputValue := sizeChartRow[4]

		if inputSku != lastInputSku && lastInputSku != "" {
			// Process SKU sizechart for one sku
			errSkuSc, data := processSkuWiseSC(lastInputSku, sizechart.CategoryId, skuWiseSC)
			if errSkuSc != "" {
				failedSku = append(failedSku, errSkuSc)
			} else {
				brandWiseScArr = append(brandWiseScArr, data)
				successSku = append(successSku, lastInputSku)
			}

			skuWiseSC = [][]string{}
		}
		row := []string{inputColumnHeader, inputRowHeaderName, inputRowHeaderType, inputValue}
		// stores the sizechart data for one sku
		skuWiseSC = append(skuWiseSC, row)

		lastInputSku = inputSku
	}

	// prcess sku sizechart for last sku
	errSkuSc, data := processSkuWiseSC(lastInputSku, sizechart.CategoryId, skuWiseSC)
	if errSkuSc != "" {
		failedSku = append(failedSku, errSkuSc)
	} else {
		brandWiseScArr = append(brandWiseScArr, data)
		successSku = append(successSku, lastInputSku)
	}
	return brandWiseScArr, failedSku, successSku
}

//	This function process the sizechart data for one sku,
//	validates it and convert it to brandwise sizechart data

func processSkuWiseSC(sku string, sizeChCategory int, skuWiseSC [][]string) (string, sizeUtils.BrandWiseScData) {
	var sizeChartDataArr []sizeUtils.SizeChartData
	var brSc sizeUtils.BrandWiseScData
	var errValid string
	brSc = sizeUtils.BrandWiseScData{}
	skuCategory, _, brandId, err := getProductInfo(sku)
	if err != nil {
		errValid = "Sku " + sku + "doesnot exist in the system"
		logger.Error(fmt.Sprintf("#processSkuWiseSC():Sku %s doesnot exist in the system", sku))
		return sku, brSc
	}
	brandName := getBrandNameById(brandId)
	errValid, sizeChartDataArr = checkIfSkuWiseSCValid(skuWiseSC, sizeChCategory, skuCategory, brandName)
	if errValid != "" {
		logger.Error(fmt.Sprintf("#processSkuWiseSC(): Validation failed for sku sizechart with sku %s and Reason is :%s", sku, errValid))
		// Notifying Reason for failed validation
		title := fmt.Sprintf("Validation failed for sizechart.")
		text := fmt.Sprintf("Validation failed for sku level sizechart with sku %s and Reason is :%s", sku, errValid)
		tags := []string{"sizechart", "sku", "validation-failed"}
		notification.SendNotification(title, text, tags, "info")
		return sku, brSc
	}
	brSc.BrandId = brandId
	brSc.ScData = sizeChartDataArr
	return "", brSc

}

//	This function return product Info
func getProductInfo(sku string) (int, int, int, error) {
	productInfo, err := products.BySku(sku, "mongo")
	if err != nil {
		return 0, 0, 0, err
	}
	leafCategory := productInfo.Leaf
	configID := productInfo.SeqId
	brandId := productInfo.BrandId
	return leafCategory[0], configID, brandId, nil
}

//	This function receives sizechart data for one sku, and it checks
//	if data is valid or not
func checkIfSkuWiseSCValid(skuWiseSC [][]string, categoryId int, skuCategory int, brandN string) (string, []sizeUtils.SizeChartData) {
	sizeChartDataArr := []sizeUtils.SizeChartData{}
	var err []string
	var firstColumnHeader string
	var row []string
	columnHeaderBunch := -1
	arrBaseValues := []string{}
	checkDuplicate := []string{}
	headerWiseType := make(map[string][]string)
	lastColumnHeader := ""
	// check if sku category and sizechart category uploaded are same
	if categoryId != skuCategory {
		err = append(err, "Sku Doesnot belong to given category")
	}
	if len(skuWiseSC) == 0 {
		err = append(err, "Empty/Invalid Sizechart csv data")
	}
	// process sku sizechart data , one row at a time
	for k, sizeChartRow := range skuWiseSC {
		inputColumnHeader := sizeChartRow[0]
		inputRowHeaderType := sizeChartRow[2]
		inputValue := sizeChartRow[3]

		// create data structure to store sc data for sku
		sizech := sizeUtils.SizeChartData{}
		sizech.Brand = brandN
		sizech.ColumnHeader = inputColumnHeader
		sizech.RowHeaderName = sizeChartRow[1]
		sizech.RowHeaderType = inputRowHeaderType
		sizech.Value = inputValue

		sizeChartDataArr = append(sizeChartDataArr, sizech)
		sizech = sizeUtils.SizeChartData{}
		// validation code

		str := inputColumnHeader + "--" + inputRowHeaderType + "--" + inputValue

		if stringInSlice(str, checkDuplicate) {
			err = append(err, "Duplicate rows found for sizes")
		}

		checkDuplicate = append(checkDuplicate, str)
		// gets the first columnheader for sku
		if k == 0 {
			firstColumnHeader = inputColumnHeader
		}
		// Check If Base size ie first column header has blank row header type
		if firstColumnHeader == inputColumnHeader {
			// stores value field for first column header ie base size
			arrBaseValues = append(arrBaseValues, inputValue)
			if inputRowHeaderType != "" {
				err = append(err, "Base Size must have blank row header type")
			}
		}
		// process when column header changes and store RowHeaderTypes for each columnheader
		if lastColumnHeader != inputColumnHeader && lastColumnHeader != "" && columnHeaderBunch > 0 {
			headerWiseType[lastColumnHeader] = row
			row = []string{}
		} else if lastColumnHeader != inputColumnHeader {
			columnHeaderBunch++
		}
		// check if other sizes have empty rowheadertype
		if firstColumnHeader != inputColumnHeader {
			row = append(row, inputRowHeaderType)
			if inputRowHeaderType == "" {
				err = append(err, "Others size must have row header type")
			}
		}
		lastColumnHeader = inputColumnHeader

	}
	headerWiseType[lastColumnHeader] = row

	if len(headerWiseType) == 0 {
		err = append(err, "Others size must be in CSV")
	}

	countFirstHeaderValue := len(arrBaseValues)
	// checks if rowheadertype and values of base size matches or not
	if countFirstHeaderValue > 0 {
		for _, typeArr := range headerWiseType {
			if len(typeArr) != countFirstHeaderValue {
				err = append(err, "Base size count and other size count must be same")
			}
			if reflect.DeepEqual(typeArr, arrBaseValues) == false {
				err = append(err, "Base size values and other size row header type values must be same")
			}
		}
	}
	return strings.Join(err, ","), sizeChartDataArr
}

// This function validates the complete sizechart data and
// returns the sizechart data brandwise ie csv contains sizechart
// data for many brands, it will process and check if data is valid and
// return the sizechart data brandwise
func validateSizeChart(sizechart sizeUtils.SizeChart) ([]sizeUtils.BrandWiseScData, interface{}) {
	var lastBrandName = ""
	var firstColumnHeader, lastColHeader string

	var sizeChartData []sizeUtils.SizeChartData
	var row, rowBase []string
	var arrOtherType, arrBaseValues map[string][]string
	arrBaseValues = make(map[string][]string)
	arrOtherType = make(map[string][]string)
	var brandWiseOtherType map[string]map[string][]string
	brandWiseOtherType = make(map[string]map[string][]string)

	var brandWiseSC sizeUtils.BrandWiseScData
	var brandWiseArray []sizeUtils.BrandWiseScData

	if len(sizechart.SChartData) == 0 {
		return nil, "Invalid/Blank CSV"
	}
	// process the sizechart data ie. process each row of csv at a time
	for _, sizeChartRow := range sizechart.SChartData {
		// check if brand is not given, name it generic
		inputBrand := sizeChartRow[0]
		if strings.TrimSpace(sizeChartRow[0]) == "" {
			inputBrand = "_generic"
		}
		inputColumnHeader := sizeChartRow[1]
		inputRowHeaderName := sizeChartRow[2]
		inputRowHeaderType := sizeChartRow[3]
		inputValue := sizeChartRow[4]
		// check for first brand in csv/ every time brand changes in csv
		if lastBrandName == "" || lastBrandName != inputBrand {
			firstColumnHeader = inputColumnHeader

			// if it is the first brand in CSV
			if lastBrandName == "" {
				lastBrandName = inputBrand
			} else {
				arrOtherType[lastColHeader] = row
				row = []string{}
				brandWiseOtherType[lastBrandName] = arrOtherType
				arrBaseValues[lastBrandName] = rowBase
				rowBase = []string{}
				arrOtherType = map[string][]string{}
				// process when change of brand in csv
				brandId, msg := prepareBrand(lastBrandName)
				if brandId < 0 {
					return nil, msg
				}
				brandWiseSC.BrandId = brandId
				brandWiseSC.ScData = sizeChartData
				sizeChartData = []sizeUtils.SizeChartData{}
				brandWiseArray = append(brandWiseArray, brandWiseSC)
				lastBrandName = inputBrand
			}
			lastColHeader = inputColumnHeader
		}
		// check if first column header has empty rowheadertype
		if firstColumnHeader == inputColumnHeader {
			rowBase = append(rowBase, inputValue)
			if inputRowHeaderType != "" {
				return nil, "Base Size must have blank row header type"
			}
		}
		// check for other column headers of brand, rowheadertype should be non-empty
		if firstColumnHeader != inputColumnHeader {
			if lastColHeader != inputColumnHeader && lastColHeader != firstColumnHeader {
				arrOtherType[lastColHeader] = row
				row = []string{}
			}
			lastColHeader = inputColumnHeader
			row = append(row, inputRowHeaderType)
			if inputRowHeaderType == "" {
				return nil, "Other than base size must have row header type"
			}
		}
		// create size chart data struct
		scData := sizeUtils.SizeChartData{}
		scData.Brand = inputBrand
		scData.ColumnHeader = inputColumnHeader
		scData.RowHeaderName = inputRowHeaderName
		scData.RowHeaderType = inputRowHeaderType
		scData.Value = inputValue
		sizeChartData = append(sizeChartData, scData)
	}
	// Process for last brand data
	arrOtherType[lastColHeader] = row
	arrBaseValues[lastBrandName] = rowBase
	brandWiseOtherType[lastBrandName] = arrOtherType
	// Remaning validation check for data ie base size values and rowheadertype for rest columnheder should match

	for brandName, _ := range arrBaseValues {
		if len(brandWiseOtherType[brandName]) == 0 {
			return nil, "Others size must be in CSV for brand " + brandName
		}
		for colHeader, typeArr := range brandWiseOtherType[brandName] {
			if len(typeArr) != len(arrBaseValues[brandName]) {
				return nil, "Base size count and count of RowHeaderType for " + colHeader + " for brand " + brandName + " must be same"
			}
			if reflect.DeepEqual(typeArr, arrBaseValues[brandName]) == false {
				return nil, "Base size values and values of RowHeaderType for " + colHeader + " for brand " + brandName + " must be same"
			}
		}
	}
	// add the entry for last size chart
	brandId, msg := prepareBrand(lastBrandName)
	if brandId < 0 {
		return nil, msg
	}
	brandWiseSC.BrandId = brandId
	brandWiseSC.ScData = sizeChartData
	brandWiseArray = append(brandWiseArray, brandWiseSC)

	return brandWiseArray, nil
}

// check if brand exist in system and return its id
func prepareBrand(name string) (int, string) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, sizeUtils.CheckBrandInSystem)
	defer func() {
		logger.EndProfile(profiler, sizeUtils.CheckBrandInSystem)
	}()
	var brd interface{}
	if name == "_generic_" {
		return 0, ""
	}
	var mongoDriver *mongodb.MongoDriver
	mongoDriver = factory.GetMongoSession(sizeUtils.SizeChartAPI)
	defer mongoDriver.Close()
	mgoObj := mongoDriver.SetCollection(sizeUtils.BrandCollection)
	err := mgoObj.Find(bson.M{"name": name}).Select(bson.M{"_id": 0, "seqId": 1, "name": 1, "urlKey": 1, "imgName": 1}).One(&brd)

	if err != nil {
		return -1, "Brand " + name + " doesnot exist in the system"
	}
	brandMap := brd.(bson.M)
	brandId, bErr := utils.GetInt(brandMap["seqId"])
	if bErr != nil {
		logger.Error(fmt.Sprintf("#prepareBrand: Unable to get brand Id. Error is: %s", bErr))
		return -1, "Unable to verify the brand in the system."
	}
	return brandId, ""

}
