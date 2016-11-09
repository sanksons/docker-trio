package sizechart

import (
	"time"
)

type ProdSChart struct {
	Id        int    `json:"id" bson:"id"`
	SizeChart SChart `json:"data" bson:"data"`
}

type SChart struct {
	Headers   []string            `json:"headers" bson:"headers"`
	ImageName string              `json:"image_name" bson:"imageName"`
	Sizes     map[string][]string `json:"sizes" bson:"sizes"`
	ScType    string              `json:"sizechart_type" bson:"sizeChartType"`
}
type SizeChartData struct {
	//IdCatalogSizeChart int     `xorm:"not null pk autoincr index INT(10) " mapstructure "id_catalog_sizechart"`
	Brand         string  `xorm:"varchar(255) not null default null" mapstructure "brand"`
	ColumnHeader  *string `xorm:"varchar(255) default null" mapstructure "column_header"`
	RowHeaderType *string `xorm:"varchar(255) default null" mapstructure "row_header_type"`
	RowHeaderName *string `xorm:"varchar(255) not null default null" mapstructure "row_header_name"`
	Value         string  `xorm:"varchar(255) not null default null" mapstructure "value"`
}

type SizechartMapping struct {
	Sku         string `mapstructure "sku" bson:"sku"`
	SizeChartId int    `mapstructure "id" bson:"sizechartId"`
}

type SizeChart struct {
	IdCatalogDistinctSizeChart int       `mapstructure "id_catalog_distinct_sizechart"`
	FkCatalogCategory          int       `mapstructure "fk_catalog_category"`
	FkCatalogBrand             *int      `mapstructure "fk_catalog_brand"`
	FkCatalogTy                *int      `mapstructure "fk_catalog_ty"`
	SizeChartName              string    `mapstructure "sizechart_name"`
	SizeChartType              *int      `mapstructure "sizechart_type"`
	FkAclUser                  *int      `mapstructure "fk_acl_user"`
	SizeChartImage             string    `mapstructure "image_path"`
	CreatedAt                  time.Time `mapstructure "created_at"`
	UpdatedAt                  time.Time `mapstructure "updated_at"`
}

type SizeChartMongo struct {
	IdCatalogSizeChart int             `bson:"seqId"`
	FkCatalogCategory  int             `bson:"categoryId"`
	FkCatalogBrand     int             `bson:"brandId,omitempty"`
	FkCatalogTy        int             `bson:"tyId,omitempty"`
	SizeChartName      string          `bson:"sizeChartName,omitempty"`
	SizeChartType      int             `bson:"sizeChartType"`
	FkAclUser          int             `bson:"aclUser,omitempty"`
	SizeChartImagePath string          `bson:"sizeChartImagePath"`
	SizeChartInfo      []SizeChartData `bson:"data"`
	CreatedAt          time.Time       `bson:"createdAt"`
	UpdatedAt          time.Time       `bson:"updatedAt"`
}
