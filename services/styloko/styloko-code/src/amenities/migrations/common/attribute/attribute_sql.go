package attribute

import (
	"fmt"
	"strconv"
)

func getAttributeSetSql() string {
	sql := `SELECT
	       id_catalog_attribute_set as seqId,
	       name,
	       ( IF( label is NULL, label_en, label ) ) AS label,
	       identifier
		   FROM catalog_attribute_set
		   WHERE name != 'global_attribute'
		   ORDER BY id_catalog_attribute_set ASC`
	return sql
}

func getAttributeSql() string {
	sql := `SELECT ca.id_catalog_attribute as seqId,
			( IF( ca.fk_catalog_attribute_set is NULL, 'global', cas.name ) ) AS setName,
			( IF( ca.fk_catalog_attribute_set is NULL, 23, cas.id_catalog_attribute_set ) ) AS setId,
			ca.name,
			( IF( ca.fk_catalog_attribute_set is NULL, 1, 0 ) ) AS isGlobal,
			ca.product_type,
			ca.attribute_type,
			ca.max_length,
			ca.decimal_places,
			ca.default_value,
			ca.unique_value,
			ca.alice_export,
			ca.dwh_field_name,
			ca.solr_searchable,
			ca.solr_filter,
			ca.solr_suggestions,
			ca.export_position,
			ca.export_name,
			ca.import_name,
			ca.import_config_identifier,
			( IF( ca.label is NULL, ca.label_en, ca.label ) ) AS label,
			( IF( ca.description is NULL, ca.description_en, ca.description ) ) AS description,
			ca.pet_mode,
			ca.pet_type,
			ca.pet_overview,
			ca.pet_overview_filter,
			ca.pet_qc,
			ca.validation,
			ca.mandatory,
			ca.mandatory_import,
			ca.extra,
			ca.visible,
			ca.created_at,
			ca.updated_at,
			ca.data_type
	        FROM catalog_attribute AS ca
		    LEFT JOIN catalog_attribute_set AS cas ON id_catalog_attribute_set = fk_catalog_attribute_set
		    ORDER BY ca.id_catalog_attribute ASC`
	return sql
}

func getAttributeIdSql(id string) string {
	sql := `SELECT ca.id_catalog_attribute as seqId,
			( IF( ca.fk_catalog_attribute_set is NULL, 'global', cas.name ) ) AS setName,
			( IF( ca.fk_catalog_attribute_set is NULL, 23, cas.id_catalog_attribute_set ) ) AS setId,
			ca.name,
			( IF( ca.fk_catalog_attribute_set is NULL, 1, 0 ) ) AS isGlobal,
			ca.product_type,
			ca.attribute_type,
			ca.max_length,
			ca.decimal_places,
			ca.default_value,
			ca.unique_value,
			ca.alice_export,
			ca.dwh_field_name,
			ca.solr_searchable,
			ca.solr_filter,
			ca.solr_suggestions,
			ca.export_position,
			ca.export_name,
			ca.import_name,
			ca.import_config_identifier,
			( IF( ca.label is NULL, ca.label_en, ca.label ) ) AS label,
			( IF( ca.description is NULL, ca.description_en, ca.description ) ) AS description,
			ca.pet_mode,
			ca.pet_type,
			ca.pet_overview,
			ca.pet_overview_filter,
			ca.pet_qc,
			ca.validation,
			ca.mandatory,
			ca.mandatory_import,
			ca.extra,
			ca.visible,
			ca.created_at,
			ca.updated_at,
			ca.data_type
	        FROM catalog_attribute AS ca
		    LEFT JOIN catalog_attribute_set AS cas ON id_catalog_attribute_set = fk_catalog_attribute_set
		    WHERE ca.id_catalog_attribute = ` + id
	return sql
}

func getOptionSql(attrName string) string {
	optionTbl := fmt.Sprintf("catalog_attribute_option_global_%s", attrName)
	idOptionTbl := fmt.Sprintf("id_%s", optionTbl)
	sql := `SELECT
                        ` + idOptionTbl + `,
                        ` + idOptionTbl + `,
                        name as value,
                        position,
                        is_default
                   FROM ` + optionTbl + `
                   ORDER BY position ASC, name ASC`
	return sql
}

func getAttributeOptionSql(attrSetName string, attrName string) string {
	optionTbl := fmt.Sprintf("catalog_attribute_option_%s_%s", attrSetName, attrName)
	idOptionTbl := fmt.Sprintf("id_%s", optionTbl)
	sql := `SELECT ` + idOptionTbl + `
				   as id_catalog_attribute_option,
                   IF (name = '', 'default', name) as value,
                   position,
                   is_default
                   FROM ` + optionTbl + `
                   ORDER BY position ASC, value ASC`
	return sql
}

func updateAttributeOption(attrSetName string, attrName string, opt Option) string {
	optionTbl := fmt.Sprintf("catalog_attribute_option_%s_%s", attrSetName, attrName)
	sql := `INSERT INTO ` + optionTbl +
		` VALUES (` +
		strconv.Itoa(opt.SeqId) + `, "` +
		*opt.Value + `", "` +
		*opt.Value + `", ` +
		strconv.Itoa(*opt.Position) + `, ` +
		strconv.Itoa(*opt.IsDefault) + `)`
	return sql
}
