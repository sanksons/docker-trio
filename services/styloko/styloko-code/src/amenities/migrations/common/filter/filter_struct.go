package filter

import (
	"time"
)

type FilterMongo struct {
	SeqId              int        `xorm:"not null pk autoincr INT(10)" bson:"seqId" mapstructure:"id_catalog_filter" json:"seqId"`
	FkCatalogAttribute *int       `xorm:"unique INT(10)" bson:"fkCtlgAttr" mapstructure:"fk_catalog_attribute" json: "fkCatalogAttribute"`
	Name               *string    `xorm:"not null VARCHAR(128)" bson:"name" mapstructure:"name"`
	Param              *string    `xorm:"not null VARCHAR(64)" bson:"param,omitempty" mapstructure:"param"`
	Description        *string    `xorm:"VARCHAR(255)" bson:"descrp,omitempty" mapstructure:"description"`
	View               *string    `xorm:"not null default 'deflt' VARCHAR(64)" bson:"view,omitempty" mapstructure:"view"`
	ShowOne            *int       `xorm:"default 0 TINYINT(1)" bson:"showOne,omitempty" mapstructure:"show_one"`
	SolrFacetSearch    *string    `xorm:"not null VARCHAR(128)" bson:"solrFacetSrch,omitempty" mapstructure:"solr_facet_search"`
	SolrFacetValue     *string    `xorm:"VARCHAR(64)" bson:"solrFacetVal,omitempty" mapstructure:"solr_facet_value"`
	SolrQueryOperator  *string    `xorm:"not null default 'OR' ENUM('OR','AND')" bson:"solrQryOp,omitempty" mapstructure:"solr_query_operator"`
	SortBy             *string    `xorm:"not null default 'sku' ENUM('alphabetic','sku')" bson:"sortBy,omitempty" mapstructure:"sort_by"`
	SortOrder          *string    `xorm:"not null default 'desc' ENUM('asc','desc')" bson:"sortOrdr,omitempty" mapstructure:"sort_order"`
	OverrideOrder      *int       `xorm:"default 0 TINYINT(1)" bson:"overrideOrdr,omitempty" mapstructure:"override_order"`
	DefaultOrder       *int       `xorm:"TINYINT(2)" bson:"dfltOrdr,omitempty" mapstructure:"default_order"`
	ExtraOptions       *string    `xorm:"VARCHAR(8192)" bson:"extraOpt,omitempty" mapstructure:"extra_options"`
	Status             *string    `xorm:"not null default 'active' ENUM('active','inactive')" bson:"status,omitempty" mapstructure:"status"`
	CreatedAt          *time.Time `xorm:"not null DATETIME" bson:"crtdAt,omitempty" mapstructure:"created_at"`
	UpdatedAt          *time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP" bson:"updtdAt,omitempty" mapstructure:"updated_at"`
}

type CatalogFilter struct {
	SeqId              int        `xorm:"not null pk autoincr INT(10)" bson:"seqId" mapstructure:"id_catalog_filter"`
	FkCatalogAttribute *int       `xorm:"unique INT(10)" bson:"fkCtlgAttr" mapstructure:"fk_catalog_attribute"`
	Name               *string    `xorm:"not null VARCHAR(128)" bson:"name" mapstructure:"name"`
	Param              *string    `xorm:"not null VARCHAR(64)" bson:"param" mapstructure:"param"`
	Description        *string    `xorm:"VARCHAR(255)" bson:"descrp" mapstructure:"description"`
	View               *string    `xorm:"not null default 'deflt' VARCHAR(64)" bson:"view" mapstructure:"view"`
	ShowOne            *int       `xorm:"default 0 TINYINT(1)" bson:"showOne" mapstructure:"show_one"`
	SolrFacetSearch    *string    `xorm:"not null VARCHAR(128)" bson:"solrFacetSrch" mapstructure:"solr_facet_search"`
	SolrFacetValue     *string    `xorm:"VARCHAR(64)" bson:"solrFacetVal" mapstructure:"solr_facet_value"`
	SolrQueryOperator  *string    `xorm:"not null default 'OR' ENUM('OR','AND')" bson:"solrQryOp" mapstructure:"solr_query_operator"`
	SortBy             *string    `xorm:"not null default 'sku' ENUM('alphabetic','sku')" bson:"sortBy" mapstructure:"sort_by"`
	SortOrder          *string    `xorm:"not null default 'desc' ENUM('asc','desc')" bson:"sortOrdr" mapstructure:"sort_order"`
	OverrideOrder      *int       `xorm:"default 0 TINYINT(1)" bson:"overrideOrdr" mapstructure:"override_order"`
	DefaultOrder       *int       `xorm:"TINYINT(2)" bson:"dfltOrdr" mapstructure:"default_order"`
	ExtraOptions       *string    `xorm:"VARCHAR(8192)" bson:"extraOpt" mapstructure:"extra_options"`
	Status             *string    `xorm:"not null default 'active' ENUM('active','inactive')" bson:"status" mapstructure:"status"`
	CreatedAt          *time.Time `xorm:"not null DATETIME" bson:"crtdAt" mapstructure:"created_at"`
	UpdatedAt          *time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP" bson:"updtdAt" mapstructure:"updated_at"`
}
