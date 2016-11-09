package attributes

import (
	"time"
)

type Attribute struct {
	IsGlobal               int       `bson:"isGlobal,omitempty" json:"isGlobal"`
	Set                    string    `bson:"set,omitempty" json:"-"`
	Name                   string    `bson:"name,omitempty" json:"name"`
	Label                  *string   `bson:"label,omitempty" json:"label"`
	Description            *string   `bson:"description,omitempty" json:"description"`
	ProductType            string    `bson:"productType,omitempty" json:"productType"`
	AttributeType          string    `bson:"attributeType,omitempty" json:"attributeType"`
	MaxLength              *int      `bson:"maxLength,omitempty" json:"maxLength,omitempty"`
	DecimalPlaces          *int      `bson:"decimalPlaces,omitempty" json:"decimalPlaces,omitempty"`
	DefaultValue           *string   `bson:"defaultValue,omitempty" json:"defaultValue,omitempty"`
	UniqueValue            *string   `bson:"uniqueValue,omitempty" json:"uniqueValue,omitempty"`
	PetType                *string   `bson:"petType,omitempty" json:"inputType"`
	PetMode                string    `bson:"petMode,omitempty" json:"inputMode"`
	Validation             *string   `bson:"validation,omitempty" json:"validation,omitempty"`
	Mandatory              *bool     `bson:"mandatory" json:"mandatory"`
	MandatoryImport        *bool     `bson:"mandatoryImport" json:"mandatoryImport"`
	AliceExport            string    `bson:"aliceExport,omitempty" json:"aliceExport"`
	PetQc                  *int      `bson:"petQc,omitempty" json:"petQc"`
	ImportConfigIdentifier *int      `bson:"importConfigIdentifier,omitempty" json:"importConfigIdentifier"`
	SolrSearchable         *int      `bson:"solrSearchable,omitempty" json:"solrSearchable"`
	SolrFilter             *int      `bson:"solrFilter,omitempty" json:"solrFilter"`
	SolrSuggestions        *int      `bson:"solrSuggestions,omitempty" json:"solrSuggestions"`
	Visible                int       `bson:"visible,omitempty" json:"visible"`
	IsActive               int       `bson:"isActive,omitempty" json:"isActive"`
	FilterType             *string   `bson:"filterType,omitempty" json:"filterType"`
	UpdatedAt              time.Time `bson:"updatedAt,omitempty" json:"updatedAt"`
}

type Parameters struct {
	AttrId   int     `bson:"attrId"`
	Count    int     `bson:"param"`
	OptionId int     `bson:"optionId"`
	EndPoint *string `bson:"endPoint"`
}

type CheckOptions struct {
	AttributeType string   `bson:"attributeType,omitempty" json:"attributeType"`
	Options       []Option `bson:"options" json:"options"`
}

type Option struct {
	Value     string `bson:"value" json:"name" validate:"required"`
	Position  int    `bson:"position" json:"position"  validate:"omitempty"`
	IsDefault int    `bson:"isDefault" json:"isDefault"  validate:"omitempty"`
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

type OptionUpdateResponse struct {
	Success   bool   `json:"success"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
