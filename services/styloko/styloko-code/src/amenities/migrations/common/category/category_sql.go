package category

import (
	"strconv"
)

//query to get info of all categories
func getCatalogCategorySql() string {
	sql := `SELECT cc.*
			FROM catalog_category as cc
			ORDER BY id_catalog_category ASC`
	return sql
}

//query to get all catalog segments
func getCatalogSegmentSql() string {
	sql := `SELECT cs.*
			FROM catalog_segment as cs
			ORDER BY id_catalog_segment ASC`
	return sql
}

//query to get the parent of each node
func getParentSql(SeqId int) string {
	SeqId1 := strconv.Itoa(SeqId)
	sql := `SELECT
			(SELECT id_catalog_category
			FROM catalog_category t2
			WHERE t2.lft < t1.lft
			AND t2.rgt > t1.rgt
			ORDER BY t2.rgt-t1.rgt ASC limit 1)
			AS parent FROM catalog_category t1
			WHERE t1.id_catalog_category IN (` + SeqId1 + `)
			ORDER BY rgt-lft DESC`
	return sql
}

//query to find all segments associated with each category
func categorySegmentJoinSql(SeqId string) string {
	sql := `SELECT
			cc.id_catalog_category,
			cs.*
			FROM catalog_category as cc
			INNER JOIN catalog_category_has_catalog_segment as cchcs
			ON cc.id_catalog_category=cchcs.fk_catalog_category
			INNER JOIN catalog_segment as cs
			ON cchcs.fk_catalog_segment=cs.id_catalog_segment
			WHERE cc.id_catalog_category IN (` + SeqId + `)`
	return sql
}
