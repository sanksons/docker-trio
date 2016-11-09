package filter

import (
	"amenities/migrations/common/util"
	"common/ResourceFactory"
	"fmt"
	"time"

	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/utils/logger"
)

//This function writes []Filters into mongo
func writeFilterToMongo(filter map[int]FilterMongo) {
	logger.Info("Starting to write Filter into Mongo")
	mgoSession := ResourceFactory.GetMongoSession("FilterMigration")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.Filters)
	_ = mongodb.DropCollection()
	for _, v := range filter {
		err := mongodb.Insert(v)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting Filters into Mongo %v", err.Error()))
		}
	}
}

//funtion takes input *core.Rows, gets data from mysql and
//gives map of CatalogFilter and filter ids as output
func processFilterRows(response interface{}) (map[int]*CatalogFilter, error) {
	rows, ok := response.(*core.Rows)
	filterInfo := make(map[int]*CatalogFilter)
	//var fids string
	if ok {
		for rows.Next() {
			f := CatalogFilter{}
			err := rows.ScanStructByIndex(&f)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in processFilterRows() %v", err.Error()))
				rows.Close()
				return nil, err
			}
			filterInfo[f.SeqId] = &f
		}
		rows.Close()
	}
	return filterInfo, nil
}

//This function takes input map[int]CatalogFilter as input and places
//these into to final map[int]FilterMongo and returns it
func TransformFilter(filterInfo map[int]*CatalogFilter) map[int]FilterMongo {
	filter := make(map[int]FilterMongo)
	for k, _ := range filterInfo {
		var tmp FilterMongo
		tmp.SeqId = filterInfo[k].SeqId
		tmp.FkCatalogAttribute = filterInfo[k].FkCatalogAttribute
		tmp.Name = filterInfo[k].Name
		tmp.Param = filterInfo[k].Param
		tmp.Description = filterInfo[k].Description
		tmp.View = filterInfo[k].View
		tmp.ShowOne = filterInfo[k].ShowOne
		tmp.SolrFacetSearch = filterInfo[k].SolrFacetSearch
		tmp.SolrFacetValue = filterInfo[k].SolrFacetValue
		tmp.SolrQueryOperator = filterInfo[k].SolrQueryOperator
		tmp.SortBy = filterInfo[k].SortBy
		tmp.SortOrder = filterInfo[k].SortOrder
		tmp.OverrideOrder = filterInfo[k].OverrideOrder
		tmp.DefaultOrder = filterInfo[k].DefaultOrder
		tmp.ExtraOptions = filterInfo[k].ExtraOptions
		tmp.Status = filterInfo[k].Status
		if filterInfo[k].CreatedAt == nil {
			time := time.Now()
			filterInfo[k].CreatedAt = &time
		}
		tmp.CreatedAt = filterInfo[k].CreatedAt
		if filterInfo[k].UpdatedAt == nil {
			time := time.Now()
			filterInfo[k].UpdatedAt = &time
		}
		tmp.UpdatedAt = filterInfo[k].UpdatedAt
		filter[k] = tmp
	}
	return filter
}
