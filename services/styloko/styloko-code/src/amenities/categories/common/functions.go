package common

import (
	"database/sql"
	"reflect"
	"strconv"
	"strings"

	"github.com/jabong/floRest/src/common/utils/logger"
)

// InsertSegments -> Inserts segments in array of CategorySegment
func InsertSegments(segIds []int) []CategorySegment {
	one, two, three := false, false, false
	segments := []CategorySegment{}
	for k := range segIds {
		switch segIds[k] {
		case 1:
			if !one {
				tmp := CategorySegment{}
				tmp.IdCategorySegment = 1
				tmp.Genders = "Women|Unisex"
				tmp.Name = "Women"
				tmp.UrlKey = "women"
				segments = append(segments, tmp)
				one = true
			}
			break
		case 2:
			if !two {
				tmp := CategorySegment{}
				tmp.IdCategorySegment = 2
				tmp.Genders = "Men|Unisex"
				tmp.Name = "Men"
				tmp.UrlKey = "men"
				segments = append(segments, tmp)
				two = true
			}
			break
		case 3:
			if !three {
				tmp := CategorySegment{}
				tmp.IdCategorySegment = 3
				tmp.Genders = "Baby|Boys|Girls"
				tmp.Name = "Kids"
				tmp.UrlKey = "kids"
				segments = append(segments, tmp)
				three = true
			}
			break
		default:
			break
		}
	}

	return segments
}

// genQuery -> Generates a query string for values only, i.e. a="value",b="value"
// SQL commands must be added around the returned string.
// Generic function.
func genQuery(container interface{}) string {
	queryValue := ""
	stype := reflect.TypeOf(container)
	sval := reflect.ValueOf(container)
	for i := 0; i < stype.NumField(); i++ {
		field := stype.Field(i)
		val := sval.Field(i)
		jtag := field.Tag.Get("sql")
		if jtag == "" {
			continue
		}
		tagValues := strings.Split(jtag, ",")
		if len(tagValues) > 1 && tagValues[1] == "omitempty" {
			switch val.Kind() {
			case reflect.String:
				if val.String() != "" {
					queryValue += tagValues[0] + "=" + "\"" + val.String() + "\","
					break
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				{
					if val.Int() != 0 {
						queryValue += tagValues[0] + "=" + string(strconv.FormatInt(val.Int(), 10)) + ","
						break
					}
				}
			default:
				break
			}
		} else {
			switch val.Kind() {
			case reflect.String:
				queryValue += tagValues[0] + "=" + val.String() + ","
				break
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				queryValue += tagValues[0] + "=" + string(strconv.FormatInt(val.Int(), 10)) + ","
				break
			default:
				break
			}
		}
	}
	return queryValue[:len(queryValue)-1]
}

// categoryUpdateQuery -> creates update query for the category.
func categoryUpdateQuery(args interface{}) (string, int, []int, error) {
	queryValue := ""
	dataMap, _ := args.(map[string]interface{})
	id, _ := dataMap["id"].(int)
	segIds, _ := dataMap["segIds"].([]int)
	queryValue = genQuery(dataMap["categoryUpdate"])
	updateQuery := "UPDATE catalog_category SET " + queryValue + " WHERE id_catalog_category=" + string(strconv.Itoa(id))
	return updateQuery, id, segIds, nil
}

// getSegments -> returns segment creation and deletion queries
func getSegments(segIds []int, id int) (deleteQuery string, updateQuery []string) {
	deleteQuery = "DELETE from catalog_category_has_catalog_segment WHERE fk_catalog_category=" + strconv.Itoa(id)
	insertQuery := "INSERT into catalog_category_has_catalog_segment (fk_catalog_category, fk_catalog_segment) VALUES "
	for x := range segIds {
		tmp := insertQuery + "(" + strconv.Itoa(id) + "," + strconv.Itoa(segIds[x]) + ")"
		updateQuery = append(updateQuery, tmp)
	}
	return deleteQuery, updateQuery
}

// genQueryCreate -> Generates a query string for columns and values i.e. (prop11, prop2) & ("val1","val2")
// SQL commands must be added around the returned string.
// Generic function.
func genQueryCreate(container interface{}) (string, string) {
	queryValue := "("
	queryColumns := "("
	stype := reflect.TypeOf(container)
	sval := reflect.ValueOf(container)
	for i := 0; i < stype.NumField(); i++ {
		field := stype.Field(i)
		val := sval.Field(i)
		jtag := field.Tag.Get("sql")
		if jtag == "" {
			continue
		}
		tagValues := strings.Split(jtag, ",")
		if len(tagValues) > 1 && tagValues[1] == "omitempty" {
			switch val.Kind() {
			case reflect.String:
				if val.String() != "" {
					queryValue += "\"" + val.String() + "\","
					queryColumns += tagValues[0] + ","
					break
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				{
					if val.Int() != 0 {
						queryValue += string(strconv.FormatInt(val.Int(), 10)) + ","
						queryColumns += tagValues[0] + ","
						break
					}
				}
			default:
				break
			}
		} else {
			switch val.Kind() {
			case reflect.String:
				queryValue += val.String() + ","
				queryColumns += tagValues[0] + ","
				break
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				queryValue += string(strconv.FormatInt(val.Int(), 10)) + ","
				queryColumns += tagValues[0] + ","
				break
			default:
				break
			}
		}
	}
	queryValue = queryValue[:len(queryValue)-1] + ")"
	queryColumns = queryColumns[:len(queryColumns)-1] + ")"
	return queryColumns, queryValue
}

// getSegmentsCreate -> returns array of segment create queries.
// No delete queries in this scenario
func getSegmentsCreate(segIds []int, id int) (updateQuery []string) {
	insertQuery := "INSERT IGNORE into catalog_category_has_catalog_segment (fk_catalog_category, fk_catalog_segment) VALUES "
	for x := range segIds {
		tmp := insertQuery + "(" + strconv.Itoa(id) + "," + strconv.Itoa(segIds[x]) + ")"
		updateQuery = append(updateQuery, tmp)
	}
	return updateQuery
}

// categoryInsertQuery -> return everything i.e category create query, segment queries, rightSubTree queries and
// parentTree queries. One must fire them all and then commit the transaction.
func categoryInsertQuery(args interface{}) (insertCategory string, segmentQueries []string, updateRight string, updateLeft string, id int) {
	dataMap, _ := args.(map[string]interface{})
	id, _ = dataMap["id"].(int)
	segIds, _ := dataMap["segIds"].([]int)
	parentRight, _ := dataMap["parentRight"].(string)
	// parentLeft, _ := dataMap["parentLeft"].(string)
	// Base insertion query values received here.
	queryColumn, queryValue := genQueryCreate(dataMap["categoryCreate"])
	insertCategory = "INSERT into catalog_category " + queryColumn + " VALUES " + queryValue
	// Segment queryies
	segmentQueries = getSegmentsCreate(segIds, id)

	// Right and left update queries
	updateRight = "UPDATE catalog_category SET rgt=rgt+2 WHERE rgt >= " + parentRight + " ORDER BY rgt DESC"
	updateLeft = "UPDATE catalog_category SET lft=lft+2 WHERE lft > " + parentRight + " ORDER BY lft DESC"

	return insertCategory, segmentQueries, updateRight, updateLeft, id
}

// execOrRollback -> executes the query, if it fails, then db rollback is called.
func execOrRollback(txnObj *sql.Tx, query string) (err error) {
	_, err = txnObj.Exec(query)
	if err != nil {
		logger.Error("Exec Failed. Transaction rollback begin.")
		txnObj.Rollback()
	}
	return err
}
