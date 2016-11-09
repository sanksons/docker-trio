package set

import (
	"time"
)

type AttributeSet struct {
	SeqId            int     `bson:"seqId" json:"seqId"`
	Name             string  `bson:"name" json:"name"`
	Label            string  `bson:"label" json:"label"`
	GlobalIdentifier *string `bson:"-" json:"globalIdentifier"`
	Identifier       string  `bson:"identifier" json:"identifier"`
}

//struct for seller center
type AttributeSets struct {
	SeqId            int         `bson:"seqId" json:"setId"`
	Name             string      `bson:"name" json:"name"`
	Label            string      `bson:"label" json:"label"`
	GlobalIdentifier *string     `bson:"-" json:"globalIdentifier" `
	Identifier       string      `bson:"identifier" json:"identifier"`
	Attributes       []Attribute `bson:"attributes" json:"attributes,"`
}

//struct for seller center
type Attribute struct {
	SeqId                  int       `bson:"seqId" json:"attributeId"`
	GlobalIdentifier       *string   `bson:"-" json:"globalIdentifier,omitempty"`
	IsGlobal               int       `bson:"isGlobal" json:"isGlobal"`
	AttrSet                Set       `bson:"set,omitempty" json:"attrSet,omitempty"`
	Name                   string    `bson:"name" json:"name"`
	Label                  *string   `bson:"label" json:"label"`
	Description            *string   `bson:"description" json:"description"`
	ProductType            string    `bson:"productType" json:"productType"`
	AttributeType          string    `bson:"attributeType" json:"attributeType"`
	MaxLength              *int      `bson:"maxLength" json:"maxLength,omitempty"`
	DecimalPlaces          *int      `bson:"decimalPlaces" json:"decimalPlaces,omitempty"`
	DefaultValue           *string   `bson:"defaultValue" json:"defaultValue,omitempty"`
	UniqueValue            *string   `bson:"uniqueValue" json:"uniqueValue,omitempty"`
	PetType                *string   `bson:"petType" json:"inputType"`
	PetMode                *string   `bson:"petMode" json:"inputMode"`
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
	Options                []Option  `bson:"options" json:"options,omitempty"`
	CreatedAt              time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt              time.Time `bson:"updatedAt" json:"updatedAt"`
}

type Set struct {
	SeqId *int    `bson:"seqId,omitempty" json:"seqId,omitempty"`
	Name  *string `bson:"name,omitempty" json:"name,omitempty"`
}

type Option struct {
	SeqId            int     `bson:"seqId" json:"optionId"`
	GlobalIdentifier *string `bson:"-" json:"globalIdentifier"`
	Name             *string `bson:"value" json:"name"`
	Position         *int    `bson:"position" json:"position"`
	IsDefault        *bool   `bson:"isDefault" json:"isDefault"`
}

type Parameters struct {
	SetId      *string `bson:"setId"`
	Count      int     `bson:"param"`
	AttrIdName *string `bson:"attrId"`
	EndPoint   *string `bson:"endPoint"`
}
