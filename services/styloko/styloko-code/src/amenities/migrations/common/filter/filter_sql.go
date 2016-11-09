package filter

func getFilterSql() string {
	sql := `SELECT *
			FROM catalog_filter
			ORDER BY id_catalog_filter ASC`
	return sql
}
