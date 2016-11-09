package category

import (
	"amenities/migrations/common/util"
	"common/ResourceFactory"
	"fmt"
	"strings"

	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/utils/logger"
)

//funtion takes input *core.Rows, gets data from mysql and
//gives []Category as output
func processAllCategories(response interface{}) (map[int]*Category, *string, error) {
	rows, ok := response.(*core.Rows)
	categoryArr := make(map[int]*Category)
	var cIds string
	if ok {
		for rows.Next() {
			c := Category{}
			err := rows.ScanStructByIndex(&c)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in ProcessAllCategories() %v", err.Error()))
				rows.Close()
				return nil, nil, err
			}
			categoryArr[c.SeqId] = &c
			cIds += fmt.Sprintf("%d,", c.SeqId)
		}
		rows.Close()
	}
	cIds = strings.TrimRight(cIds, ",")
	return categoryArr, &cIds, nil
}

//funtion takes input *core.Rows, gets data from mysql and
//gives []CatalogSegment as output
func processAllCategorySegments(response interface{}) (map[int]CatalogSegment, error) {
	rows, ok := response.(*core.Rows)
	csArr := make(map[int]CatalogSegment)
	if ok {
		for rows.Next() {
			cs := CatalogSegment{}
			err := rows.ScanStructByIndex(&cs)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in processAllCategorySegments() %v", err.Error()))
				rows.Close()
				return nil, err
			}
			csArr[cs.IdCatalogSegment] = cs
		}
		rows.Close()
	}
	return csArr, nil
}

//funtion takes input *core.Rows, gets data from mysql and
//gives join between segments and category as output
func processJoinedRows(response interface{}) (map[int][]CatalogSegment1, error) {
	rows, ok := response.(*core.Rows)
	cs := make(map[int][]CatalogSegment1)
	if ok {
		for rows.Next() {
			c := CatalogSegment1{}
			err := rows.ScanStructByIndex(&c)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in processJoinedRows() %v", err.Error()))
				rows.Close()
				return nil, err
			}
			cs[c.SeqId] = append(cs[c.SeqId], c)
		}
	}
	rows.Close()
	return cs, nil
}

//This function takes input map[int]*Category as input,
//calls function to get the parent id of each node and finally maps into
//map[int]CategoryMongo and returns it
func TransformCategory(categories map[int]*Category, catalogSegment map[int][]CatalogSegment1) map[int]CategoryMongo {
	cat := make(map[int]CategoryMongo)
	for k, _ := range categories {
		tmp := CategoryMongo{}
		tmp.SeqId = categories[k].SeqId
		tmp.Status = categories[k].Status
		tmp.Lft = categories[k].Lft
		tmp.Rgt = categories[k].Rgt
		tmp.Name = categories[k].Name
		tmp.UrlKey = categories[k].UrlKey
		tmp.SizechartActive = categories[k].SizechartActive
		tmp.PdfName = categories[k].PdfName
		tmp.PdfActive = categories[k].PdfActive
		tmp.DisplaySizeConversion = categories[k].DisplaySizeConversion
		tmp.GoogleTreeMapping = categories[k].GoogleTreeMapping
		tmp.SizechartApplicable = categories[k].SizechartApplicable
		parentId, _ := getParent(categories[k].SeqId)
		tmp.Parent = parentId
		tmp.CatalogSegment1 = catalogSegment[k]
		cat[k] = tmp
	}
	return cat
}

//This function writes []Category into mongo
func writeCategoryToMongo(category map[int]CategoryMongo) {
	logger.Info("Starting to write Category into Mongo")
	mgoSession := ResourceFactory.GetMongoSession("CategoryMigration")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.Categories)
	_ = mongodb.DropCollection()
	for _, v := range category {
		err := mongodb.Insert(v)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting Category into Mongo %v", err.Error()))
		}
	}
}

//This function writes []CatalogSegment into mongo
func writeCatalogSegmentToMongo(catalogSegment map[int]CatalogSegment) {
	logger.Info("Starting to write Catalog Segment into Mongo")
	mgoSession := ResourceFactory.GetMongoSession("CategoryMigration")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.CategorySegments)
	_ = mongodb.DropCollection()
	for _, v := range catalogSegment {
		err := mongodb.Insert(v)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting Catalog Segment into Mongo %v", err.Error()))
		}
	}
}
