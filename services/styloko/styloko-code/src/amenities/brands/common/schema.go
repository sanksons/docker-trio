package common

import (
	"time"
)

type BrandCertificate struct {
	IdCatalogBrand            int       `bson:"-" json:"-"`
	IdCatalogBrandCertificate int       `bson:"idCtlgBrndCert,omitempty" json:"idCtlgBrndCert,omitempty"`
	FkCatalogBrand            int       `bson:"-" json:"-"`
	FkCatalogCategory         int       `bson:"fkCtlgCtgry,omitempty" json:"fkCtlgCtgry,omitempty"`
	Filename                  string    `bson:"filename,omitempty" json:"filename,omitempty"`
	ValidFrom                 time.Time `bson:"validFrm,omitempty" json:"validFrm,omitempty"`
	ValidTill                 time.Time `bson:"validTill,omitempty" json:"validTill,omitempty"`
	IsActive                  int       `bson:"isActv,omitempty" json:"isActv,omitempty"`
	FkAclUser                 int       `bson:"fkAclUsr,omitempty" json:"fkAclUsr,omitempty"`
	Ip                        string    `bson:"Ip,omitempty" json:"Ip,omitempty"`
	CreatedAt                 time.Time `bson:"-" json:"crtdAt"`
	UpdatedAt                 time.Time `bson:"-" json:"updtdAt"`
}

type RelatedBrand struct {
	IdRelatedBrand int `bson:"seqId,omitempty" json:"seqId,omitempty" required:"true" sql="id_related_brand,omitempty"`
}

type Brand struct {
	SeqId        int                `bson:"seqId,omitempty" json:"seqId,omitempty"`
	Status       string             `bson:"status,omitempty" json:"status,omitempty"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	Position     int                `bson:"position,omitempty" json:"position,omitempty"`
	UrlKey       string             `bson:"urlKey,omitempty" json:"urlKey,omitempty"`
	ImageName    string             `bson:"imgName,omitempty" json:"imgName,omitempty"`
	BrandClass   string             `bson:"brndClass,omitempty" json:"brndClass,omitempty"`
	IsExclusive  int                `bson:"isExc" json:"isExc,omitempty"`
	BrandInfo    string             `bson:"brandInfo,omitempty" json:"brandInfo,omitempty"`
	RelatedBrand []RelatedBrand     `bson:"relatedBrnd,omitempty" json:"relatedBrnd"`
	Certificate  []BrandCertificate `bson:"brndCert,omitempty" json:"brndCert"`
}

type BrandUpdate struct {
	Status       string             `bson:"status,omitempty" json:"status,omitempty"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	Position     int                `bson:"position,omitempty" json:"position,omitempty"`
	UrlKey       string             `bson:"urlKey,omitempty" json:"urlKey,omitempty"`
	ImageName    string             `bson:"imgName,omitempty" json:"imgName,omitempty"`
	BrandClass   string             `bson:"brndClass,omitempty" json:"brndClass,omitempty"`
	IsExclusive  int                `bson:"isExc,omitempty" json:"isExc,omitempty"`
	BrandInfo    string             `bson:"brandInfo,omitempty" json:"brandInfo"`
	RelatedBrand []RelatedBrand     `bson:"relatedBrnd,omitempty" json:"relatedBrnd,omitempty"`
	Certificate  []BrandCertificate `bson:"brndCert,omitempty" json:"brndCert,omitempty"`
}

type BrandData struct {
	Branddata []Brand `json:"brandDataList"`
}
