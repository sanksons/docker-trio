package category

import (
	"common/xorm/mysql"

	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/utils/logger"
)

func StartCategoryMigration() error {
	//to insert category in mongo
	err := getCategoryTree()
	if err != nil {
		return err
	}

	//to insert catalogSegments in mongo
	er := getCatalogSegment()
	if er != nil {
		return er
	}
	return nil
}

//function to get the category tree and catalog segments
func getCategoryTree() error {
	logger.Info("Started Migrating Category Table")
	category, cIds, err := getCategories()
	if err != nil {
		return err
	}

	catalogSegment, err := getJoin(*cIds)
	if err != nil {
		return err
	}

	cat := TransformCategory(category, catalogSegment)
	writeCategoryToMongo(cat)
	return nil
}

//function to produce join between categories and catalog segments
func getJoin(cIds string) (map[int][]CatalogSegment1, error) {
	logger.Info("Preparing CatalogSegments")
	sql := categorySegmentJoinSql(cIds)
	rows, err := mysql.GetInstance().Query(sql, false)
	if err == nil {
		cs, er := processJoinedRows(rows)
		if er == nil {
			return cs, nil
		}
		return nil, er
	}
	return nil, err
}

//function to get info pertaining to all catalog segments from mysql and insert it in mongo
func getCatalogSegment() error {
	logger.Info("Started Migrating CatalogSegment Table")
	catalogSegment, er := getCatalogSegments()
	if er != nil {
		return er
	}
	writeCatalogSegmentToMongo(catalogSegment)
	return nil
}

//function to get immediate parentId corresponding to every node
func getParent(mySqlId int) (*int, error) {
	logger.Info("Getting parent of each id Category")
	sql := getParentSql(mySqlId)
	rows, err := mysql.GetInstance().Query(sql, false)
	var parentId *int
	if err != nil {
		logger.Info(err.Error())
		return nil, err
	}
	row := rows.(*core.Rows)
	for row.Next() {
		e := row.Scan(&parentId)
		if e != nil {
			logger.Info(err.Error())
			return nil, e
		}
	}
	row.Close()
	return parentId, nil
}

//function to get info pertaining to all categories from mysql and insert it in mongo
func getCategories() (map[int]*Category, *string, error) {
	logger.Info("Preparing Category")
	sql := getCatalogCategorySql()
	rows, err := mysql.GetInstance().Query(sql, false)
	if err == nil {
		category, cIds, err := processAllCategories(rows)
		if err == nil {
			return category, cIds, nil
		}
		return nil, nil, err
	}
	return nil, nil, err
}

func getCatalogSegments() (map[int]CatalogSegment, error) {
	logger.Info("Preparing CatalogSegments")
	sql := getCatalogSegmentSql()
	rows, err := mysql.GetInstance().Query(sql, false)
	if err == nil {
		catSeg, err := processAllCategorySegments(rows)
		if err == nil {
			return catSeg, nil
		}
		return nil, err
	}
	return nil, err
}
