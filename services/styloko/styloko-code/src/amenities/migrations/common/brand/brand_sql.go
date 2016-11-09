package brand

//query to get entire Brand flatly
func getBrandSql() string {
	sql := `SELECT *
		   FROM catalog_brand
		   ORDER BY id_catalog_brand ASC`
	return sql
}

//query to get related brands of a brand
func getRelatedBrandSql(brandIds string) string {
	sql := `SELECT
    		cb.id_catalog_brand,
			crb.*
			FROM catalog_brand as cb
			RIGHT JOIN catalog_related_brand as crb
			ON cb.id_catalog_brand = crb.fk_catalog_brand
		   	WHERE cb.id_catalog_brand IN (` + brandIds + `)`
	return sql
}

//query to get brandCertificates associated with a brand
func getBrandCertificateSql(brandIds string) string {
	sql := `SELECT
			cb.id_catalog_brand,
		    bc.*
		    FROM catalog_brand as cb
		    RIGHT JOIN catalog_brand_certificate as bc
		     ON cb.id_catalog_brand = bc.fk_catalog_brand
		    WHERE cb.id_catalog_brand IN (` + brandIds + `)`
	return sql
}
