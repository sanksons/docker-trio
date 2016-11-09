package product

import (
	"common/utils"
	"common/xorm/mysql"
	db "database/sql"
	"errors"
	"fmt"
	"github.com/go-xorm/core"
	"strconv"
	"sync"
	"time"
)

var r sync.RWMutex

type CatalogAttribute struct {
	Id                    int     `mapstructure:"id_catalog_attribute"`
	fkCatalogAttributeSet *int    `mapstructure:"fk_catalog_attribute_set"`
	name                  string  `mapstructure:"name"`
	productType           string  `mapstructure:"product_type"`
	attributeType         string  `mapstructure:"attribute_type"`
	aliceExport           string  `mapstructure:"alice_export"`
	solrSearchable        *int    `mapstructure:"solr_searchable"`
	solrFilter            *int    `mapstructure:"solr_filter"`
	solrSuggestions       *int    `mapstructure:"solr_suggestions"`
	Validation            *string `mapstructure:"validation"`
	PetType               string  `mapstructure:"-"`
	label                 string  `mapstructure:"label"`
	attrSetName           string  `mapstructure:"setname"`
}

type AttrOption struct {
	Id    int    `bson:"seqId"`
	Value string `bson:"value"`
}

func getAttributesOfSet(attributeSetId int, protype string) ([]*CatalogAttribute, error) {

	var id string = strconv.Itoa(attributeSetId)
	index := id + "#" + protype
	if _, ok := AttributesCache[index]; ok {
		return AttributesCache[index], nil
	}
	var sql string = `SELECT 
		id_catalog_attribute,
		fk_catalog_attribute_set,
		ca.name,
		product_type,
		attribute_type,
		alice_export,
		solr_searchable,
		solr_filter,
		solr_suggestions,
		validation,
		pet_type,
		ca.label,
		IF (fk_catalog_attribute_set IS NULL, "global",cas.name) AS setname
		FROM catalog_attribute AS ca
		LEFT JOIN catalog_attribute_set AS cas
		ON ca.fk_catalog_attribute_set = cas.id_catalog_attribute_set
		WHERE (fk_catalog_attribute_set IS NULL OR fk_catalog_attribute_set = ` + id + `) 
		AND attribute_type!="system"
		AND ca.name NOT IN ("barcode_ean", "ean_code")
		AND product_type="` + protype + `";`
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		return nil, errors.New("getAttributesOfSet(): " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return nil, errors.New("getAttributesOfSet(): Error in Parsing rows")
	}
	var attributes []*CatalogAttribute
	for rows.Next() {
		attribute := &CatalogAttribute{}
		err := rows.Scan(&attribute.Id,
			&attribute.fkCatalogAttributeSet,
			&attribute.name,
			&attribute.productType,
			&attribute.attributeType,
			&attribute.aliceExport,
			&attribute.solrSearchable,
			&attribute.solrFilter,
			&attribute.solrSuggestions,
			&attribute.Validation,
			&attribute.PetType,
			&attribute.label,
			&attribute.attrSetName)
		if err != nil {
			rows.Close()
			return nil, errors.New("getAttributesOfSet(): " + err.Error())
		}
		attributes = append(attributes, attribute)
	}
	rows.Close()
	r.Lock()
	AttributesCache[index] = attributes
	r.Unlock()
	return attributes, nil
}

func processProductAttributes(proId int, attributeSetId int, protype string) ([]*Attribute, error) {
	attributes, err := getAttributesOfSet(attributeSetId, protype)
	if err != nil {
		return nil, errors.New("processProductAttributes(): " + err.Error())
	}
	//var ok bool
	var attrs []*Attribute
	for _, attribute := range attributes {
		var skip = false
		attr := &Attribute{}
		attr.Id = attribute.Id
		attr.AliceExport = attribute.aliceExport
		attr.Label = attribute.label
		attr.Name = attribute.name
		attr.OptionType = attribute.attributeType
		if attribute.solrFilter != nil {
			attr.SolrFilter = *attribute.solrFilter
		}
		if attribute.solrSearchable != nil {
			attr.SolrSearchable = *attribute.solrSearchable
		}
		if attribute.solrSuggestions != nil {
			attr.SolrSuggestions = *attribute.solrSuggestions
		}
		attr.IsGlobal = false
		if attribute.attrSetName == "global" {
			attr.IsGlobal = true
		}
		value, err := getAttributeValue(attribute, proId, protype)
		if err != nil {
			if err == ErrEmpty {
				skip = true
			} else {
				return nil, errors.New("getAttributesOfSet():fail " + err.Error())
			}
		}
		if !skip {
			attr.Value = value
			attrs = append(attrs, attr)
		}
	}
	return attrs, nil
}

func getAttributeValue(attribute *CatalogAttribute, id int, protype string) (interface{}, error) {
	//global
	oldVariations := map[string]string{
		"bags":             "variation",
		"beauty":           "variation",
		"fragrances":       "variation",
		"home":             "variation",
		"sports_equipment": "size",
		"toys":             "variation",
	}
	//old attributes to consider
	type NewAttr struct {
		Attrset []string
		OldAttr string
	}
	oldValueAttrs := map[string]NewAttr{
		"fits":               NewAttr{[]string{"shoes"}, "fit"},
		"jeans_wash_effects": NewAttr{[]string{"app_men", "app_women"}, "jeans_wash_effect"},
		"lens_types":         NewAttr{[]string{"bags"}, "lens_type"},
		"secondary_colors":   NewAttr{[]string{"global"}, "secondary_color"},
		"products_warranty":  NewAttr{[]string{"global"}, "product_warranty"},
		"qualities":          NewAttr{[]string{"home"}, "quality"},
		"materials_code":     NewAttr{[]string{"home"}, "material_code"},
		"pocket":             NewAttr{[]string{"app_men", "app_women"}, "pockets"},
	}

	switch attribute.attributeType {
	case "multi_option":
		data, err := prepareMultiOptionData(
			id, attribute.attrSetName, attribute.name, protype)
		return data, err
	case "option":
		if protype == "simple" && attribute.name == "variations" {
			if v, ok := oldVariations[attribute.attrSetName]; ok {
				data, err := prepareNewVariation(id, v, attribute.attrSetName)
				return data, err
			}
		}
		if protype == "config" {
			if v, ok := oldValueAttrs[attribute.name]; ok {
				if utils.InArrayString(v.Attrset, attribute.attrSetName) {
					data, err := prepareNewOptionAttrs(id, v.OldAttr, attribute.name, attribute.attrSetName)
					return data, err
				}
			}
		}
		data, err := prepareOptionData(
			id, attribute.attrSetName, attribute.name, protype)
		return data, err
	case "value", "custom":
		data, err := prepareValueData(
			id, attribute.attrSetName, attribute.name, protype, attribute)
		return data, err
	default:
		return nil, errors.New("getAttributeValue():Invalid Attribute Type.")
	}
}

func prepareValueData(id int, s1 string, s2 string, protype string, attribute *CatalogAttribute) (interface{}, error) {

	var result []byte
	var resultdate *time.Time

	var validation string
	if attribute.Validation != nil {
		validation = *attribute.Validation
	}

	var idStr string = strconv.Itoa(id)
	var Table = "catalog_" + protype
	var whereField = "id_catalog_" + protype
	if s1 != "global" {
		Table = "catalog_" + protype + "_" + s1
		whereField = "fk_catalog_" + protype
	}
	var sql string = `SELECT ` + SqlSafe(s2) + ` 
    	FROM ` + Table + ` 
    	WHERE ` + whereField + ` = ` + idStr
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return result, nil
		}
		return result, errors.New("prepareValueData(): " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return result, errors.New("prepareValueData(): Assertion failed")
	}
	if attribute.PetType != "datefield" {
		for rows.Next() {
			err = rows.Scan(&result)
			if err != nil {
				continue
			}
		}
		rows.Close()
		strVal := string(result)

		if strVal == "" {
			return "", ErrEmpty
		}
		switch validation {

		case "decimal":
			//convert to decimal
			floatVal, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return 0.00, nil
			}
			return floatVal, nil
		case "integer":
			//convert to interger
			intVal, err := strconv.Atoi(strVal)
			if err != nil {
				return 0, nil
			}
			return intVal, nil
		default:
			if attribute.PetType == "checkbox" {
				intVal, _ := strconv.Atoi(strVal)
				return intVal, nil
			}
			return strVal, nil
		}

	} else {
		for rows.Next() {
			err = rows.Scan(&resultdate)
			if err != nil {
				continue
			}
		}
		//fmt.Println(reflect.TypeOf(resultdate))
		rows.Close()
		if resultdate == nil {
			return "", ErrEmpty
		}
		return utils.ToMySqlTime(resultdate), nil
	}
	return result, nil
}

func prepareNewVariation(
	simpleId int,
	oldattrname string,
	attrsetName string,
) (AttrOption, error) {
	result := AttrOption{}
	//get old value
	simpletableName := fmt.Sprintf("catalog_simple_%s", attrsetName)
	sql := "SELECT " + oldattrname + " FROM " + simpletableName + " WHERE fk_catalog_simple=" + strconv.Itoa(simpleId)
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return result, nil
		}
		return result, errors.New("prepareNewVariation()1: " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return result, errors.New("prepareNewVariation(): Assertion failed")
	}
	var variation string
	for rows.Next() {
		err = rows.Scan(&variation)
		if err != nil {
			rows.Close()
			return result, errors.New("prepareNewVariation()2:" + err.Error())
		}
	}
	rows.Close()
	fieldName1 := fmt.Sprintf("id_catalog_attribute_option_%s_variations", attrsetName)
	tablename1 := fmt.Sprintf("catalog_attribute_option_%s_variations", attrsetName)
	sql2 := "SELECT " + fieldName1 + ",name FROM " + tablename1 + " WHERE name='" + variation + "';"
	response, err = mysql.GetInstance().Query(sql2, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return result, nil
		}
		return result, errors.New("prepareNewVariation()1: " + err.Error())
	}
	rows, ok = response.(*core.Rows)
	if !ok {
		return result, errors.New("prepareNewVariation(): Assertion failed")
	}
	for rows.Next() {
		err = rows.Scan(&result.Id, &result.Value)
		if err != nil {
			rows.Close()
			return result, errors.New("prepareNewVariation()2:" + err.Error())
		}
	}
	rows.Close()
	if result.Id == 0 {
		return result, ErrEmpty
	}
	return result, nil
}

func prepareOptionData(id int, s1 string, s2 string, protype string) (AttrOption, error) {
	result := AttrOption{}
	var idStr string = strconv.Itoa(id)
	var Table = "catalog_" + protype
	var whereField = "id_catalog_" + protype
	if s1 != "global" {
		Table = "catalog_" + protype + "_" + s1
		whereField = "fk_catalog_" + protype
	}
	var optionTable = "catalog_attribute_option_" + s1 + "_" + s2
	var sql string = `SELECT 
    	fk_` + optionTable + ` AS id, 
    	s2.name AS val 
    	FROM ` + Table + ` AS s1
		INNER JOIN ` + optionTable + ` AS s2
		ON s1.fk_` + optionTable + ` = s2.id_` + optionTable + `
		WHERE ` + whereField + ` =` + idStr + `
		;`
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return result, nil
		}
		return result, errors.New("prepareOptionData(): " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return result, errors.New("prepareOptionData(): Assertion failed")
	}
	for rows.Next() {
		err = rows.Scan(&result.Id, &result.Value)
		if err != nil {
			rows.Close()
			return result, errors.New("prepareOptionData():" + err.Error())
		}
	}
	rows.Close()
	if result.Id == 0 && s2 == "is_returnable" {
		result.Id = 4
		return result, nil
	}
	if result.Id == 0 {
		return result, ErrEmpty
	}
	return result, nil
}

func prepareMultiOptionData(id int, s1 string, s2 string, protype string) ([]AttrOption, error) {
	var result []AttrOption
	var idStr string = strconv.Itoa(id)
	var optionTable string = "catalog_attribute_option_" + s1 + "_" + s2
	var linkTable string = "catalog_attribute_link_" + s1 + "_" + s2
	var whereField = "fk_catalog_" + protype
	var extraJoin string
	if s1 != "global" {
		extraJoin = "INNER JOIN catalog_" + protype + "_" + s1 + " AS s3 " +
			" ON s1.fk_catalog_" + protype + "_" + s1 + " = s3.id_catalog_" + protype + "_" + s1
	}

	var sql string = `SELECT 
    		fk_` + optionTable + ` AS id, 
            s2.name AS val 
    	FROM ` + linkTable + ` AS s1 
		INNER JOIN ` + optionTable + ` AS s2
		ON s1.fk_` + optionTable + ` = s2.id_` + optionTable + `
		` + extraJoin + ` 
		WHERE ` + whereField + `=` + idStr + `;`
	//fmt.Println(sql)

	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return result, nil
		}
		return result, errors.New("prepareMultiOptionData(): " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return result, errors.New("prepareMultiOptionData(): Assertion failed")
	}

	for rows.Next() {
		attrOpt := AttrOption{}
		err = rows.Scan(&attrOpt.Id, &attrOpt.Value)
		if err != nil {
			printErr(errors.New("prepareMultiOptionData(): " + err.Error()))
			continue
		}
		result = append(result, attrOpt)
	}
	rows.Close()
	//check if empty
	if len(result) == 0 {
		return result, ErrEmpty
	}
	return result, nil
}

func prepareNewOptionAttrs(
	configId int,
	oldattrname string,
	newattrname string,
	attrsetName string,
) (AttrOption, error) {

	result := AttrOption{}
	//get old value
	var tableName string
	var whereField string
	if attrsetName == "global" {
		//global
		tableName = "catalog_config"
		whereField = "id_catalog_config"
	} else {
		tableName = fmt.Sprintf("catalog_config_%s", attrsetName)
		whereField = "fk_catalog_config"
	}
	sql := "SELECT " + oldattrname + " FROM " + tableName + " WHERE " + whereField + "=" + strconv.Itoa(configId)
	response, err := mysql.GetInstance().Query(sql, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return result, nil
		}
		return result, errors.New("prepareNewOptionAttrs()1: " + err.Error())
	}
	rows, ok := response.(*core.Rows)
	defer rows.Close()
	if !ok {
		return result, errors.New("prepareNewOptionAttrs()2: Assertion failed")
	}
	var oldValue *string
	for rows.Next() {
		err = rows.Scan(&oldValue)
		if err != nil {
			rows.Close()
			return result, errors.New("prepareNewOptionAttrs()3:" + err.Error())
		}
	}
	rows.Close()
	if oldValue == nil {
		return result, ErrEmpty
	}
	//fetch new attribute value, based on old

	fieldName1 := fmt.Sprintf("id_catalog_attribute_option_%s_%s", attrsetName, newattrname)
	tablename1 := fmt.Sprintf("catalog_attribute_option_%s_%s", attrsetName, newattrname)
	sql2 := "SELECT " + fieldName1 + ",name FROM " + tablename1 + " WHERE name='" + *oldValue + "';"
	response, err = mysql.GetInstance().Query(sql2, QUERY_MYSQL_MASTER)
	if err != nil {
		if err == db.ErrNoRows {
			return result, nil
		}
		return result, errors.New("prepareNewOptionAttrs()4: " + err.Error())
	}
	rows1, ok := response.(*core.Rows)
	defer rows1.Close()
	if !ok {
		return result, errors.New("prepareNewOptionAttrs()5: Assertion failed")
	}
	for rows1.Next() {
		err = rows1.Scan(&result.Id, &result.Value)
		if err != nil {
			rows1.Close()
			return result, errors.New("prepareNewOptionAttrs()6:" + err.Error())
		}
	}
	rows1.Close()
	if result.Id == 0 {
		return result, ErrEmpty
	}
	return result, nil
}
