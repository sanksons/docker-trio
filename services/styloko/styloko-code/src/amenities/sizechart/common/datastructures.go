package common

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

type SkuSizeChart struct {
	CategoryId int        `json:"categoryID"`
	CsvName    string     `json:"sizeChartCsvName"`
	ImageName  string     `json:"sizeChartImagefile"`
	AclUserId  int        `json:"aclUserId"`
	SChartData [][]string `json:"sizeChartCsvData"`
}

type SizeChart struct {
	CategoryId  int        `json:"categoryID"`
	CatalogType string     `json:"categoryType"`
	CsvName     string     `json:"sizeChartCsvName"`
	ImageName   string     `json:"sizeChartImagefile"`
	AclUserId   int        `json:"aclUserId"`
	SChartData  [][]string `json:"sizeChartCsvData"`
}

type SizeChartData struct {
	Brand         string
	ColumnHeader  string
	RowHeaderName string
	RowHeaderType string
	Value         string
}

type BrandWiseScData struct {
	BrandId int
	ScData  []SizeChartData
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

type SizeChartForMysql struct {
	Skus        []string
	SizeChartId int
	SizeChartTy int
	CategoryId  int
	CreatedAt   time.Time
}

type SizeChartSkuMapping struct {
	Sku         string `bson:"sku"`
	SizeChartId int    `bson:"sizechartId"`
}
