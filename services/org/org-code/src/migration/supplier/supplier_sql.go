package supplier

func getSupplierSql() string {
	sql := `SELECT 
			cs.id_catalog_supplier,
			cs.seller_id,
			cs.name,
			cs.name,
			cs.status,
			cs.order_email,
			cs.contact,
			cs.phone,
			cs.customercare_email,
			cs.customercare_contact,
			cs.customercare_phone,
			cs.created_at,
			cs.updated_at,
			sa.street,
			sa.street_number,
			sa.city,
			sa.postcode,
			country.iso2_code
		   FROM catalog_supplier as cs
		   LEFT JOIN supplier_address as sa
		   ON cs.id_catalog_supplier = sa.fk_id_catalog_supplier
		   LEFT JOIN country
		   ON country.id_country = sa.fk_country`
	return sql
}

func getMaxIdForSupplier() string {
	sql := `SELECT
 			 MAX(id_catalog_supplier)
			FROM catalog_supplier;`
	return sql
}
