package brand

import (
	"time"
)

type Brand struct {
	SeqId       int     `xorm:"not null pk autoincr INT(10)" bson:"seqId" mapstructure:"id_catalog_brand" json: "seqId"`
	Status      string  `xorm:"not null ENUM('active','inactive','deleted')" bson:"status" mapstructure:"status" json: "status"`
	Name        string  `xorm:"unique VARCHAR(255)" bson:"name" mapstructure:"name" json: "name"`
	NameEn      *string `bson:"-" mapstructure:"name_en" json: "nameEn"`
	Position    *int    `xorm:"INT(5)" bson:"pos,omitempty" mapstructure:"position" json: "position"`
	UrlKey      *string `xorm:"VARCHAR(255)" bson:"urlKey" mapstructure:"url_key,omitempty" json: "urlKey"`
	ImageName   *string `xorm:"VARCHAR(255)" bson:"imgName,omitempty" mapstructure:"image_name,omitempty" json: "imageName"`
	BrandClass  string  `xorm:"VARCHAR(255) default 'regular' ENUM('regular','premium','regular_and_premium','boutique','designer')" bson:"brndClass" mapstructure:"brand_class" json: "brandClass"`
	IsExclusive int     `xorm:"default 0 TINYINT(1)" bson:"isExc" mapstructure:"override_order" json: "isExclusive"`
	BrandInfo   *string `xorm:"VARCHAR(255)" bson:"brndInfo,omitempty" mapstructure:"brand_info" json: "brandInfo"`
}

type BrandCertificate struct {
	SeqId                     int       `xorm:"not null autoincr INT(10)" bson:"seqId" mapstructure:"id_catalog_brand" json: "seqId"`
	IdCatalogBrandCertificate int       `xorm:"not null pk autoincr INT(10)" bson:"idCtlgBrndCert" mapstructure:"id_catalog_brand_certificate" json: "idCatalogBrandCertificate"`
	FkCatalogBrand            int       `xorm:"INT(10)" bson:"-" mapstructure:"fk_catalog_brand" json: "fkCatalogBrand"`
	FkCatalogCategory         int       `xorm:"not null INT(10)" bson:"fkCtlgCtgry" mapstructure:"fk_catalog_category" json: "fkCatalogCategory"`
	Filename                  string    `xorm:"not null VARCHAR(128)" bson:"filename" mapstructure:"filename" json: "filename"`
	ValidFrom                 time.Time `xorm:"not null" bson:"validFrm" mapstructure:"valid_from" json: "validFrom"`
	ValidTill                 time.Time `xorm:"not null" bson:"validTill" mapstructure:"valid_till" json: "validTill"`
	IsActive                  *int      `xorm:"default 0 TINYINT(1)" bson:"isActv,omitempty" mapstructure:"is_active" json: "isActive"`
	FkAclUser                 *int      `xorm:"default 0 INT(10)" bson:"fkAclUsr,omitempty" mapstructure:"fk_acl_user" json: "fkAclUser"`
	Ip                        *string   `xorm:"VARCHAR(32)" bson:"Ip,omitempty" mapstructure:"ip" json: "ip"`
	CreatedAt                 time.Time `bson:"-" mapstructure:"created_at" json: "createdAt"`
	UpdatedAt                 time.Time `bson:"-" mapstructure:"updated_at" json: "updatedAt"`
}

type RelatedBrand struct {
	SeqId                 int `xorm:"not null autoincr INT(10)" mapstructure:"id_catalog_brand" bson:"-"`
	IdCatalogRelatedBrand int `xorm:"not null pk autoincr INT(10)" mapstructure:"id_catalog_related_brand" bson:"-"`
	FkCatalogBrand        int `xorm:"not null INT(10)" mapstructure:"fk_catalog_brand" bson:"-"`
	IdRelatedBrand        int `xorm:"not null INT(10)" bson:"seqId" mapstructure:"id_related_brand"`
}

type BrandMongo struct {
	SeqId        int                `xorm:"not null pk autoincr INT(10)" xorm:"not null pk autoincr INT(10)" bson:"seqId" mapstructure:"id_catalog_brand" json : "seqId"`
	Status       string             `xorm:"not null ENUM('active','inactive','deleted')" bson:"status" mapstructure:"status" json: "status"`
	Name         string             `xorm:"unique VARCHAR(255)" bson:"name" mapstructure:"name" json : "name"`
	Position     *int               `xorm:"INT(5)" bson:"pos,omitempty" mapstructure:"position" json:"position"`
	UrlKey       *string            `xorm:"VARCHAR(255)" bson:"urlKey" mapstructure:"url_key,omitempty" json: "urlKey"`
	ImageName    *string            `xorm:"VARCHAR(255)" bson:"imgName,omitempty" mapstructure:"image_name,omitempty" json : "imageName"`
	BrandClass   string             `xorm:"VARCHAR(255) default 'regular' ENUM('regular','premium','regular_and_premium','boutique','designer')" bson:"brndClass" mapstructure:"brand_class" json : "brandClass"`
	IsExclusive  int                `xorm:"default 0 TINYINT(1)" bson:"isExc" mapstructure:"is_exclusive" json: "isExclusive"`
	BrandInfo    *string            `xorm:"VARCHAR(255)" bson:"brndInfo,omitempty" mapstructure:"brand_info" json: "brandInfo"`
	RelatedBrand []RelatedBrand     `bson:"relatedBrnd,omitempty" json:"relatedBrand"`
	Certificate  []BrandCertificate `bson:"brndCert,omitempty" json:"brandCertificate"`
}
