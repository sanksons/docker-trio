package sizechart

import (
	"strconv"
)

func getSkuSizechartMappingInfo() string {
	sql := `SELECT cc.sku AS sku , ca.fk_catalog_distinct_sizechart AS id FROM catalog_config_additional_info ca 
	JOIN catalog_config cc ON cc.id_catalog_config = ca.fk_catalog_config AND ca.sizechart_type = 0`
	return sql
}

func getCatalogAdditionalInfoScSku(sizeChartId interface{}) string {
	sql := `select fk_catalog_config from catalog_config_additional_info where ` +
		` fk_catalog_distinct_sizechart = ` + strconv.Itoa(sizeChartId.(int))
	//	sql1 := `select cc.sku from  catalog_config_additional_info ccai INNER JOIN catalog_config cc ` +
	//		`ON ccai.fk_catalog_config = cc.id_catalog_config where fk_catalog_distinct_sizechart = ` + strconv.Itoa(sizeChartId.(int))
	return sql
}

func getDistinctSizeChartAfterId(id int) string {
	sql := `Select sch.id_catalog_distinct_sizechart, sch.fk_catalog_category,
			sch.fk_catalog_brand, sch.fk_catalog_ty, sch.sizechart_name,
			sch.sizechart_type, sch.fk_acl_user,sim.image_path, sch.created_at, sch.updated_at
	 		from catalog_distinct_sizechart sch INNER JOIN
	 		catalog_category_has_sizechart_image sim ON
	 		sch.id_catalog_distinct_sizechart = sim.fk_catalog_distinct_sizechart where sch.id_catalog_distinct_sizechart >` + strconv.Itoa(id)
	return sql
}

func getDistinctSizeChartSql() string {
	sql := `Select sch.id_catalog_distinct_sizechart, sch.fk_catalog_category,
			sch.fk_catalog_brand, sch.fk_catalog_ty, sch.sizechart_name,
			sch.sizechart_type, sch.fk_acl_user,sim.image_path, sch.created_at, sch.updated_at
	 		from catalog_distinct_sizechart sch INNER JOIN
	 		catalog_category_has_sizechart_image sim ON
	 		sch.id_catalog_distinct_sizechart = sim.fk_catalog_distinct_sizechart`
	return sql
}

func getSizeChartDataSql(sizeChartId int) string {

	sql := `select csz.brand, csz.column_header, 
			csz.row_header_type, csz.row_header_name, csz.value from catalog_sizechart csz where csz.fk_catalog_distinct_sizechart = ` + strconv.Itoa(sizeChartId)
	return sql
}

func getCatalogConfigSql(category interface{}, brand interface{}, typeId interface{}) string {

	var where string = ``
	if brand != nil {
		where = ` AND cc.fk_catalog_brand = ` + strconv.Itoa(brand.(int))
	}

	if typeId != nil {
		where = where + ` AND cc.fk_catalog_ty = ` + strconv.Itoa(typeId.(int))
	}

	sql := `select cc.id_catalog_config from catalog_config cc INNER JOIN catalog_config_has_catalog_category ccc
			ON cc.id_catalog_config = ccc.fk_catalog_config	where ccc.fk_catalog_category = ` + strconv.Itoa(category.(int)) + where

	//return "select cc.id_catalog_config from catalog_config cc INNER JOIN catalog_config_has_catalog_category ccc ON cc.id_catalog_config = ccc.fk_catalog_config  where ccc.fk_catalog_category = 3270 AND cc.fk_catalog_brand = 131 AND cc.fk_catalog_ty is NULL"
	return sql
}
