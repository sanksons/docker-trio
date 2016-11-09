package attributes

import (
	"time"
)

type Attribute struct {
	IsGlobal               int      `bson:"isGlobal" json:"isGlobal"`
	Set                    string   `bson:"set,omitempty" json:"set"`
	Name                   *string  `bson:"name" json:"name"`
	Label                  *string  `bson:"label" json:"label"`
	Description            *string  `bson:"description" json:"description"`
	ProductType            string   `bson:"productType" json:"productType"`
	AttributeType          string   `bson:"attributeType" json:"attributeType"`
	MaxLength              *int     `bson:"maxLength" json:"maxLength,omitempty"`
	DecimalPlaces          *int     `bson:"decimalPlaces" json:"decimalPlaces,omitempty"`
	DefaultValue           *string  `bson:"defaultValue" json:"defaultValue,omitempty"`
	UniqueValue            *string  `bson:"uniqueValue" json:"uniqueValue,omitempty"`
	PetType                string   `bson:"petType" json:"inputType"`
	PetMode                string   `bson:"petMode" json:"inputMode"`
	Validation             *string  `bson:"validation" json:"validation,omitempty"`
	Mandatory              *bool    `bson:"mandatory" json:"mandatory"`
	MandatoryImport        *bool    `bson:"mandatoryImport" json:"mandatoryImport"`
	AliceExport            string   `bson:"aliceExport" json:"aliceExport"`
	PetQc                  *int     `bson:"petQc" json:"petQc"`
	ImportConfigIdentifier *int     `bson:"importConfigIdentifier" json:"importConfigIdentifier"`
	SolrSearchable         *int     `bson:"solrSearchable" json:"solrSearchable"`
	SolrFilter             *int     `bson:"solrFilter" json:"solrFilter"`
	SolrSuggestions        *int     `bson:"solrSuggestions" json:"solrSuggestions"`
	Visible                int      `bson:"visible" json:"visible"`
	IsActive               int      `bson:"isActive" json:"isActive"`
	FilterType             *string  `bson:"filterType" json:"filterType"`
	Options                []Option `bson:"options,omitempty" json:"options,omitempty"`
}

type Parameters struct {
	AttrId   int     `bson:"attrId"`
	Count    int     `bson:"param"`
	EndPoint *string `bson:"endPoint"`
}

type Option struct {
	SeqId          int    `bson:"seqId" json:"optionId"`
	Value          string `bson:"value" json:"name" validate:"required"`
	Position       int    `bson:"position" json:"position"`
	IsDefault      int    `bson:"isDefault" json:"isDefault"`
	FilterAttrName string `bson:"-" json:"filterAttrName,omitempty"`
	FilterAttrVal  string `bson:"-" json:"filterAttrVal,omitempty"`
}

type CheckOptions struct {
	AttributeType string   `bson:"attributeType,omitempty" json:"attributeType"`
	Options       []Option `bson:"options" json:"options"`
}

type WAttribute struct {
	SeqId                  int       `bson:"seqId" json:"seqId"`
	IsGlobal               int       `bson:"isGlobal" json:"isGlobal"`
	Set                    Set       `bson:"set,omitempty" json:"set,omitempty"`
	Name                   *string   `bson:"name" json:"name"`
	Label                  *string   `bson:"label" json:"label"`
	Description            *string   `bson:"description" json:"description"`
	ProductType            string    `bson:"productType" json:"productType"`
	AttributeType          string    `bson:"attributeType" json:"attributeType"`
	MaxLength              *int      `bson:"maxLength" json:"maxLength,omitempty"`
	DecimalPlaces          *int      `bson:"decimalPlaces" json:"decimalPlaces,omitempty"`
	DefaultValue           *string   `bson:"defaultValue" json:"defaultValue,omitempty"`
	UniqueValue            *string   `bson:"uniqueValue" json:"uniqueValue,omitempty"`
	PetType                string    `bson:"petType" json:"inputType"`
	PetMode                string    `bson:"petMode" json:"inputMode"`
	Validation             *string   `bson:"validation" json:"validation,omitempty"`
	Mandatory              *bool     `bson:"mandatory" json:"mandatory"`
	MandatoryImport        *bool     `bson:"mandatoryImport" json:"mandatoryImport"`
	AliceExport            string    `bson:"aliceExport" json:"aliceExport"`
	PetQc                  *int      `bson:"petQc" json:"petQc"`
	ImportConfigIdentifier *int      `bson:"importConfigIdentifier" json:"importConfigIdentifier"`
	SolrSearchable         *int      `bson:"solrSearchable" json:"solrSearchable"`
	SolrFilter             *int      `bson:"solrFilter" json:"solrFilter"`
	SolrSuggestions        *int      `bson:"solrSuggestions" json:"solrSuggestions"`
	Visible                int       `bson:"visible" json:"visible"`
	IsActive               int       `bson:"isActive" json:"isActive"`
	FilterType             *string   `bson:"filterType" json:"filterType"`
	UpdatedAt              time.Time `bson:"updatedAt,omitempty" json:"updatedAt"`
	CreatedAt              time.Time `bson:"createdAt,omitempty"`
	Options                []Option  `bson:"options,omitempty" json:"options,omitempty"`
}

type Set struct {
	SeqId *int    `bson:"seqId,omitempty" json:"seqId,omitempty"`
	Name  *string `bson:"name,omitempty" json:"name,omitempty"`
}

type OptionCreateResponse struct {
	Success   bool   `json:"success"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
