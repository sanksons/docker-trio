package category

type Category struct {
	SeqId                 int     `xorm:"not null pk autoincr INT(10)" bson:"seqId" mapstructure:"id_catalog_category" json:"seqId"`
	Status                string  `xorm:"not null ENUM('active','inactive','inherited_inactive','deleted')" bson:"status" mapstructure:"status" json:"status"`
	Lft                   int     `xorm:"not null unique INT(11)" bson:"lft" mapstructure:"lft" json:"lft"`
	Rgt                   int     `bson:"rgt" mapstructure:"rgt" json:"rgt"`
	Name                  string  `xorm:"not null VARCHAR(255)" bson:"name" mapstructure:"name" json:"name"`
	NameEn                *string `xorm:"VARCHAR(255)" bson:"-" mapstructure:"name_en" json:"name_en"`
	UrlKey                *string `xorm:"not null VARCHAR(128)" bson:"urlKey,omitempty" mapstructure:"url_key"`
	SizechartActive       *int    `xorm:"default 1 TINYINT(1)" bson:"szchrtActv,omitempty" mapstructure:"sizechart_active,omitempty" json:"sizechartActive"`
	PdfName               *string `xorm:"VARCHAR(255)" bson:"pdfName,omitempty" mapstructure:"pdf_name,omitempty" json:"pdfName"`
	PdfActive             *int    `xorm:"default 0 TINYINT(1)" bson:"pdfActv,omitempty" mapstructure:"pdf_active,omitempty" json:"pdfActive"`
	DisplaySizeConversion *string `xorm:"default '0' VARCHAR(16)" bson:"dispSzConv,omitempty" mapstructure:"display_size_conversion,omitempty" json:"displaySizeConversion"`
	GoogleTreeMapping     *string `xorm:"VARCHAR(255)" bson:"gleTreeMapping,omitempty" mapstructure:"google_tree_mapping,omitempty" json:"googleTreeMapping"`
	SizechartApplicable   *int    `xorm:"not null default 1 TINYINT(1)" bson:"szchrtApp,omitempty" mapstructure:"sizechart_applicable,omitempty" json:"sizechartApplicable"`
}
type CatalogSegment struct {
	IdCatalogSegment int     `xorm:"not null pk autoincr INT(10)" bson:"idCtlgSeg" mapstructure:"id_catalog_segment" json:"seqId"`
	Name             string  `xorm:"VARCHAR(255)" bson:"name" mapstructure:"name" json:"name"`
	NameEn           *string `xorm:"VARCHAR(255)" bson:"-" mapstructure:"name" json:"nameEn"`
	UrlKey           *string `xorm:"VARCHAR(255)" bson:"urlKey,omitempty" mapstructure:"url_key"`
	Genders          *string `xorm:"VARCHAR(64)" bson:"genders,omitempty" mapstructure:"genders"`
	UseForCategories *int    `xorm:"default 0 TINYINT(1)" bson:"useForCtgry" mapstructure:"use_for_categories,omitempty" json:"useForCategories"`
	UseForCapmaigns  *int    `xorm:"default 0 TINYINT(1)" bson:"useForCmpgns" mapstructure:"use_for_campaigns,omitempty" json:"useForCampaigns"`
}
type CatalogSegmentMongo struct {
	CatalogSegment CatalogSegment `bson:"catalogSegment,omitempty" json:"catalogSegment"`
}
type CatalogSegment1 struct {
	SeqId            int     `xorm:"not null pk autoincr INT(10)" bson:"-" mapstructure:"id_catalog_category" json:"-"`
	IdCatalogSegment int     `xorm:"not null pk autoincr INT(10)" bson:"seqId" mapstructure:"id_catalog_segment" json:"seqId"`
	Name             string  `xorm:"VARCHAR(255)" bson:"name" mapstructure:"name" json:"name"`
	NameEn           *string `xorm:"VARCHAR(255)" bson:"-" mapstructure:"name" json:"nameEn"`
	UrlKey           *string `xorm:"VARCHAR(255)" bson:"urlKey,omitempty" mapstructure:"url_key"`
	Genders          *string `xorm:"VARCHAR(64)" bson:"genders,omitempty" mapstructure:"genders"`
	UseForCategories *int    `xorm:"default 0 TINYINT(1)" bson:"-" mapstructure:"use_for_categories,omitempty" json:"useForCategories"`
	UseForCapmaigns  *int    `xorm:"default 0 TINYINT(1)" bson:"-" mapstructure:"use_for_campaigns,omitempty" json:"useForCampaigns"`
}
type CategoryMongo struct {
	SeqId                 int               `xorm:"not null pk autoincr INT(10)" bson:"seqId" mapstructure:"id_catalog_category" json:"seqId"`
	Status                string            `xorm:"not null ENUM('active','inactive','inherited_inactive','deleted')" bson:"status" mapstructure:"status" json:"status"`
	Lft                   int               `bson:"lft" mapstructure:"lft" json:"lft"`
	Rgt                   int               `xorm:"not null unique INT(11)" bson:"rgt" mapstructure:"rgt" json:"rgt"`
	Parent                *int              `bson:"parent" mapstructure:"parent" json:"parent"`
	Name                  string            `xorm:"not null VARCHAR(255)" bson:"name" mapstructure:"name" json:"name"`
	NameEn                *string           `xorm:"VARCHAR(255)" bson:"-" mapstructure:"name_en" json:"name_en"`
	UrlKey                *string           `xorm:"not null VARCHAR(128)" bson:"urlKey,omitempty" mapstructure:"url_key"`
	SizechartActive       *int              `xorm:"default 1 TINYINT(1)" bson:"szchrtActv,omitempty" mapstructure:"sizechart_active,omitempty" json:"sizechartActive"`
	PdfName               *string           `xorm:"VARCHAR(255)" bson:"pdfName,omitempty" mapstructure:"pdf_name,omitempty" json:"pdfName"`
	PdfActive             *int              `xorm:"default 0 TINYINT(1)" bson:"pdfActv,omitempty" mapstructure:"pdf_active,omitempty" json:"pdfActive"`
	DisplaySizeConversion *string           `xorm:"default '0' VARCHAR(16)" bson:"dispSzConv,omitempty" mapstructure:"display_size_conversion,omitempty" json:"displaySizeConversion"`
	GoogleTreeMapping     *string           `xorm:"VARCHAR(255)" bson:"gleTreeMapping,omitempty" mapstructure:"google_tree_mapping,omitempty" json:"googleTreeMapping"`
	SizechartApplicable   *int              `xorm:"not null default 1 TINYINT(1)" bson:"szchrtApp,omitempty" mapstructure:"sizechart_applicable,omitempty" json:"sizechartApplicable"`
	CatalogSegment1       []CatalogSegment1 `bson:"segment,omitempty" json:"segment"`
}
