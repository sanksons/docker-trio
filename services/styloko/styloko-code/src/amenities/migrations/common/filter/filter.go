package filter

import (
	"common/xorm/mysql"

	"github.com/jabong/floRest/src/common/utils/logger"
)

func StartFilterMigration() error {
	err := filterTable()
	if err != nil {
		return err
	}
	return nil
}

//function that integrates filter details and writes it to mongo
func filterTable() error {
	logger.Info("Started Migrating Filter Table")
	filterInfo, err := getFilterInfo()
	if err != nil {
		return err
	}
	filter := TransformFilter(filterInfo)
	writeFilterToMongo(filter)
	return nil
}

//function that gets filter information from catalog_filter
func getFilterInfo() (map[int]*CatalogFilter, error) {
	sql := getFilterSql()
	rows, err := mysql.GetInstance().Query(sql, false)
	if err != nil {
		return nil, err
	}
	filterInfo, err := processFilterRows(rows)
	if err != nil {
		return nil, err
	}
	return filterInfo, nil
}
