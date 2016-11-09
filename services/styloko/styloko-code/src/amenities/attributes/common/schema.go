package common

import "time"

//Attribute as in Attribute collection
type AttributeMongo struct {
	SeqId                  int                    `bson:"seqId"`
	Set                    AttributeSet           `bson:"set"`
	GlobalIdentifier       *string                `bson:"-"`
	IsGlobal               int                    `bson:"isGlobal"`
	Name                   string                 `bson:"name"`
	Label                  string                 `bson:"label"`
	Description            *string                `bson:"description"`
	ProductType            string                 `bson:"productType"`
	AttributeType          string                 `bson:"attributeType"`
	MaxLength              *int                   `bson:"maxLength"`
	DecimalPlaces          *int                   `bson:"decimalPlaces"`
	DefaultValue           *string                `bson:"defaultValue"`
	UniqueValue            *string                `bson:"uniqueValue"`
	PetType                *string                `bson:"petType"`
	PetMode                *string                `bson:"petMode"`
	Validation             *string                `bson:"validation"`
	Mandatory              *int                   `bson:"mandatory"`
	MandatoryImport        *int                   `bson:"mandatoryImport"`
	AliceExport            string                 `bson:"aliceExport"`
	PetQc                  *int                   `bson:"petQc"`
	ImportConfigIdentifier *int                   `bson:"importConfigIdentifier"`
	SolrSearchable         int                    `bson:"solrSearchable"`
	SolrFilter             int                    `bson:"solrFilter"`
	SolrSuggestions        int                    `bson:"solrSuggestions"`
	Visible                int                    `bson:"visible"`
	CreatedAt              time.Time              `bson:"creatdAt"`
	UpdatedAt              time.Time              `bson:"updatedAt"`
	IsActive               int                    `bson:"isActive"`
	Options                []AttributeMongoOption `bson:"options"`
}

//Attribute Options as in Attribute collection.
type AttributeMongoOption struct {
	SeqId            int     `bson:"seqId"`
	GlobalIdentifier *string `bson:"-"`
	Name             string  `bson:"value"`
	Position         int     `bson:"position"`
	IsDefault        int     `bson:"isDefault"`
}

type AttributeSet struct {
	SeqId int    `bson:"seqId"`
	Name  string `bson:"name"`
}
