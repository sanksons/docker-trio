package attribute

import (
	"time"
)

type AttributeSet struct {
	SeqId      int       `bson:"seqId" json:"seqId"`
	Name       string    `bson:"name" json:"name"`
	Label      string    `bson:"label" json:"label"`
	Identifier string    `bson:"identifier" json:"identifier"`
	CreatedAt  time.Time `bson:"creatdAt,omitempty"`
	UpdatedAt  time.Time `bson:"updatedAt,omitempty"`
}

type AttributeRow struct {
	SeqId                  int        `xorm:"not null pk autoincr INT(10)" bson:"seqId" json:"seqId"`
	Set                    Set        `bson:"set,omitempty" json:"set"`
	Name                   string     `xorm:"not null VARCHAR(255)" bson:"name" json:"name"`
	IsGlobal               int        `bson:"isGlobal" json:"isGlobal"`
	ProductType            string     `xorm:"not null ENUM('config','simple')" bson:"productType" json:"productType"`
	AttributeType          string     `xorm:"not null ENUM('system','option','multi_option','value','custom')" bson:"attributeType" json:"attributeType"`
	MaxLength              *int       `xorm:"MEDIUMINT(8) unsigned" bson:"maxLength,omitempty" json:"maxLength"`
	DecimalPlaces          *int       `xorm:"MEDIUMINT(8) unsigned" bson:"decimalPlaces,omitempty" json:"decimalPlaces"`
	DefaultValue           *string    `xorm:"VARCAHR(255)" bson:"defaultValue,omitempty" json:"defaultValue"`
	UniqueValue            *string    `xorm:"ENUM('global','config')" bson:"uniqueValue,omitempty" json:"uniqueValue"`
	AliceExport            *string    `xorm:"not null default 'no' ENUM('no','meta','attribute')" bson:"aliceExport" json:"aliceExport"`
	DwhFieldName           *string    `xorm:"VARCAHR(255)" bson:"dwhFieldName,omitempty" json:"dwhFieldName"`
	SolrSearchable         *int       `xorm:"default 0 TINYINT(1)" bson:"solrSearchable,omitempty" json:"solrSearchable"`
	SolrFilter             *int       `xorm:"default 0 TINYINT(1)" bson:"solrFilter,omitempty" json:"solrFilter"`
	SolrSuggestions        *int       `xorm:"default 0 TINYINT(1)" bson:"solrSuggestions,omitempty" json:"solrSuggestions"`
	ExportPosition         *int       `xorm:"default 0 TINYINT(1)" bson:"exportPosition,omitempty" json:"exportPosition"`
	ExportName             *string    `xorm:"VARCAHR(255)" bson:"exportName,omitempty" json:"exportName"`
	ImportName             *string    `xorm:"VARCAHR(255)" bson:"importName,omitempty" json:"importName"`
	ImportConfigIdentifier *int       `xorm:"default 0 TINYINT(1)" bson:"importConfigIdentifier,omitempty" json:"importConfigIdentifier"`
	Label                  *string    `xorm:"not null VARCHAR(255)" bson:"label" json:"label"`
	Description            *string    `xorm:"TEXT" bson:"description" json:"description"`
	PetMode                *string    `xorm:"not null default 'edit' ENUM('edit','display','invisible')" bson:"petMode" json:"petMode"`
	PetType                *string    `xorm:"not null default 'textfield' ENUM('textfield','textarea','numberfield','datefield','checkbox','dropdown','multiselect','combo','multicombo')" bson:"petType" json:"petType"`
	PetOverview            *int       `xorm:"TINYINT(1)" bson:"petOverview,omitempty" json:"petOverview"`
	PetOverviewFilter      *int       `xorm:"TINYINT(1)" bson:"petOverviewFilter,omitempty" json:"petOverviewFilter"`
	PetQc                  *int       `xorm:"TINYINT(1)" bson:"petQc,omitempty" json:"petQc"`
	Validation             *string    `xorm:"not null default 'edit' ENUM('decimal','integer','percent','email','url','letters','lettersnumbers')" bson:"validation,omitempty" json:"validation"`
	Mandatory              bool       `xorm:"TINYINT(1)" bson:"mandatory" json:"mandatory"`
	MandatoryImport        *int       `xorm:"TINYINT(1)" bson:"mandatoryImport,omitempty" json:"mandatoryImport"`
	Extra                  *string    `xorm:"TEXT" bson:"extra,omitempty" json:"extra"`
	Visible                *int       `xorm:"not null default 1 TINYINT(1)" bson:"visible" json:"visible"`
	CreatedAt              *time.Time `xorm:"not null TIMESTAMP" bson:"createdAt,omitempty"`
	UpdatedAt              *time.Time `xorm:"not null TIMESTAMP" bson:"updatedAt,omitempty"`
	DataType               *string    `xorm:"VARCHAR(15)" bson:"dataType,omitempty" json:"dataType"`
	IsActive               int        `bson:"isActive" json:"isActive"`
	Options                []Option   `bson:"options,omitempty" json:"options"`
}
type Option struct {
	SeqId     int     `bson:"seqId" json:"seqId"`
	Value     *string `bson:"value" json:"value"`
	Position  *int    `bson:"position" json:"position"`
	IsDefault *int    `bson:"isDefault" json:"isDefault"`
}
type Set struct {
	SeqId int    `bson:"seqId,omitempty" json:"seqId"`
	Name  string `bson:"name,omitempty" json:"name"`
}

// AttributeMap struct for attribute mappings
type AttributeMap struct {
	From    string            `bson:"from"`
	To      string            `bson:"to"`
	Mapping map[string]string `bson:"mapping"`
}
