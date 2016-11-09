package sizechart

import (
	"amenities/migrations/common/util"
	"common/ResourceFactory"
	"common/xorm/mysql"
	"fmt"
	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

func MigrateSizechMapping() error {
	data, err := GetSKUSizechartMappingData()
	if err != nil {
		return err
	}
	mgoSession := ResourceFactory.GetMongoSession("SizechartMigration")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.SizeChartMappingCollec)
	for _, doc := range data {
		err := mongodb.Insert(doc)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting sizechart mapping into Mongo %v", err.Error()))
			return err
		}
	}
	index := mgo.Index{
		Key:        []string{"sku"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	indErr := mongodb.EnsureIndex(index)
	if indErr != nil {
		logger.Error("#MigrateSizechMapping(): Unable to create indexes on sizecharts-mapping.", indErr)
		return indErr
	}
	return nil
}

func WriteSizeChartToProd(chart *SChart, currScTy int, sku int, sizechart map[string]interface{}) bool {

	sizeChartTypeMap := make(map[int]string)
	sizeChartTypeMap = map[int]string{0: "sku", 1: "brand", 2: "brick"}
	chart.ScType = sizeChartTypeMap[currScTy]

	sizes, _ := getPrimarySizesForSku(sku)
	fmt.Println("Primary sizes for sku are ", sizes)

	// create standard size and  index mapping for sku
	chart.Sizes = make(map[string][]string)
	sizeIndexMapping := make(map[string]string)
	for i, size := range sizes {
		chart.Sizes[strconv.Itoa(i)] = []string{size}
		sizeIndexMapping[size] = strconv.Itoa(i)
	}

	sizechartdata := sizechart["data"].([]interface{})
	mismatch := true

	for _, sc := range sizechartdata {
		sch := sc.(map[string]interface{})
		if sch["rowheadertype"] != nil {
			if pos, ok := sizeIndexMapping[sch["rowheadertype"].(string)]; ok {
				mismatch = false
				chart.Sizes[pos] = append(chart.Sizes[pos], sch["value"].(string))
			}
		}
	}
	// If sizechart uploaded and standard sizes doesnot match.
	if mismatch {
		return false
	}
	// Create and upload sizechart to product
	sizeChProd := ProdSChart{sizechart["seqId"].(int), *chart}
	//sChProd := WrapSChart{sizeChProd}
	mgoSession := ResourceFactory.GetMongoSession("SizechartMigration")
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(util.Products)
	criteria := bson.M{"seqId": sku}
	updateCri := bson.M{"$set": bson.M{"sizeChart": sizeChProd}}
	err := mgoObj.Update(criteria, updateCri)

	if err != nil {
		fmt.Println("Error while updating sizechart to product# ", sku, err.Error())
	}
	fmt.Println("The sizechart is updated to sku# ", sku)
	return true
}

func UpdateSizeChartToProducts(configIds []int, sizechart map[string]interface{}) {

	var upload bool
	// Get type of current sizechart
	currTy := sizechart["sizeChartType"].(int)
	chart := SChart{}

	chart.Headers = GetHeaderForSizeChart(sizechart)

	if len(chart.Headers) == 0 {
		fmt.Println("Headers doesnot Exist for Sizechart with Id ", sizechart["seqId"])
		return
	}
	chart.ImageName = sizechart["sizeChartImagePath"].(string)

	for _, sku := range configIds {

		// find the type of Previously uploaded sizechart
		fmt.Println("Processing for SKU# ", sku)

		preTy := GetPreviousScTypeForSku(sku)

		if preTy < -1 {
			fmt.Println("SKU ", sku, " doesnot exist in our system")
			continue
		}
		fmt.Println("Previous and current type of sizechart for sku ", preTy, currTy)
		// Rules to create sizechart for product
		switch {
		case preTy == -1, preTy == currTy, currTy < preTy:
			upload = true
		}

		if upload == false {
			fmt.Println("Dont upload for sku#  ", sku)
			continue
		}

		result := WriteSizeChartToProd(&chart, currTy, sku, sizechart)

		if !result {
			fmt.Println("Total mismatch of sizes for product# ", sku)
			continue
		}
	}

}

//Takes input a prodId
//and returns the sizes from the attribute option table
//or from the attribute simple table
func getPrimarySizesForSku(configId int) ([]string, error) {
	option, attrSetName := GetSizeOptionName(configId)
	if option == "" {
		variation, err := getVariation(attrSetName, configId)
		if err != nil {
			logger.Error(fmt.Sprintf("Error occured while getting variation: %s", err.Error()))
			return nil, err
		}
		return variation, nil
	}

	// Get the sizes from the attribute table.

	simplesAttrTable := `catalog_simple_` + attrSetName
	optionsTable := `catalog_attribute_option_` + option
	sizesSql := `SELECT ` + optionsTable + `.name FROM ` + optionsTable +
		` INNER JOIN ` + simplesAttrTable + ` ON ` + simplesAttrTable + `.fk_` + optionsTable + ` = ` +
		optionsTable + `.id_` + optionsTable + ` INNER JOIN catalog_simple ON ` + simplesAttrTable +
		`.fk_catalog_simple = catalog_simple.id_catalog_simple INNER JOIN ` +
		`catalog_config ON catalog_simple.fk_catalog_config = catalog_config.id_catalog_config WHERE ` +
		`catalog_config.id_catalog_config = ` + strconv.Itoa(configId) + ` ORDER BY ` + optionsTable + `.position`

	data, err := mysql.GetInstance().Query(sizesSql, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error occured while getting option sizes: %s", err.Error()))
		return nil, err
	}

	sizeRows := data.(*core.Rows)
	sizes := make([]string, 0)
	for sizeRows.Next() {
		var size string
		err := sizeRows.Scan(&size)
		if err != nil {
			logger.Error(fmt.Sprintf("Error occured while getting option sizes: %s", err.Error()))
			sizeRows.Close()
			return nil, err
		}
		sizes = append(sizes, size)
	}
	sizeRows.Close()
	return sizes, err
}

//gets the sizes for products of bags,
//beauty,fragrances,home sports and toys
//given the prodId and attribute set name
func getVariation(attrSetName string, configId int) ([]string, error) {
	var size string
	var ok bool
	mapping := make(map[string]string)
	mapping["bags"] = "variation"
	mapping["beauty"] = "variation"
	mapping["fragrances"] = "variation"
	mapping["home"] = "variation"
	mapping["sports_equipment"] = "size"
	mapping["toys"] = "variation"
	if size, ok = mapping[attrSetName]; !ok {
		return nil, nil
	}
	sql := `SELECT csa.` + size + ` FROM catalog_simple_` + attrSetName + ` as csa
			 INNER JOIN catalog_simple as cs
			 ON cs.id_catalog_simple = csa.fk_catalog_simple
			INNER JOIN catalog_config as cc
			 ON cs.fk_catalog_config = cc.id_catalog_config
			WHERE cc.id_catalog_config = ` + strconv.Itoa(configId)

	data, err := mysql.GetInstance().Query(sql, false)

	if err != nil {
		logger.Error(fmt.Sprintf("Error occured while getting variation for attribute set name: %s", err.Error()))
		return nil, nil
	}
	rows := data.(*core.Rows)

	var variations []string
	for rows.Next() {
		var variation string
		err := rows.Scan(&variation)
		if err != nil {
			logger.Error(fmt.Sprintf("Error occured while reading rows: %s", err.Error()))
			continue
		}
		variations = append(variations, variation)
	}
	rows.Close()
	return variations, err
}

// Get the attribute set size option name for
// shoes,men apparel,women apparel,kids apparel and jewellery
// configId: The product ID.
// Returns:
// string: The name of the attribute set size option.
// string: The attribute set name.

func GetSizeOptionName(skuId int) (string, string) {

	var attributeSet string
	sql := `SELECT
                name
            FROM catalog_attribute_set
            WHERE id_catalog_attribute_set = (
                SELECT
                    fk_catalog_attribute_set
                FROM catalog_config
                WHERE id_catalog_config = ` + strconv.Itoa(skuId) + `)`

	data, err := mysql.GetInstance().Query(sql, false)

	if err != nil {
		fmt.Println("Error occured in getting attribute set of sku", err.Error())
	}

	rows := data.(*core.Rows)

	for rows.Next() {
		err := rows.Scan(&attributeSet)
		if err != nil {
			fmt.Println("Unable to fetch attribute set", err)
		}
	}
	rows.Close()
	mapping := make(map[string]string)
	mapping["shoes"] = "shoes_sh_size"
	mapping["app_men"] = "app_men_apm_size"
	mapping["app_women"] = "app_women_apw_size"
	mapping["app_kids"] = "app_kids_apk_size"
	mapping["jewellery"] = "jewellery_variation"

	if name, ok := mapping[attributeSet]; ok {
		return name, attributeSet
	}
	return "", attributeSet

}

func GetHeaderForSizeChart(sizechart map[string]interface{}) []string {

	headers := []string{}
	if sizechart["data"] == nil {
		return headers
	}
	infoArray := sizechart["data"].([]interface{})
	fmt.Println("The sizechart is ", sizechart)
	var check map[string]bool
	check = make(map[string]bool)

	for _, info := range infoArray {
		mapInfo := info.(map[string]interface{})

		if mapInfo["columnheader"] == nil {
			continue
		}
		v := mapInfo["columnheader"].(string)

		if check[v] == false {
			check[v] = true
			headers = append(headers, v)
		}

	}
	return headers
}
