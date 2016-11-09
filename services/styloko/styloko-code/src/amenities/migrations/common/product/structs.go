package product

import "time"

type ProductMigrationStatus struct {
	Id    int
	State bool
	Msg   string
}

type ProAttributeSet struct {
	Id    int    `bson:"seqId"`
	Name  string `bson:"name"`
	Label string `bson:"label"`
}

type Attribute struct {
	Id              int         `bson:"seqId"`
	IsGlobal        bool        `bson:"isGlobal"`
	Label           string      `bson:"label"`
	Name            string      `bson:"name"`
	Value           interface{} `bson:"value"`
	AliceExport     string      `bson:"aliceExport"`
	OptionType      string      `bson:"optionType"`
	SolrSearchable  int         `bson:"solrSearchable"`
	SolrFilter      int         `bson:"solrFilter"`
	SolrSuggestions int         `bson:"solrSuggestions"`
}

type Sizes map[string][]string

type SizeChart struct {
	Headers   []string
	ImageName string
	Sizes     Sizes
}

type Rating struct {
}

type PriceChart struct {
	MaxPrice            string
	Price               string
	MaxOriginalPrice    string
	OriginalPrice       string
	SpecialPrice        string
	MaxSavingPercentage string
	SpecialPriceFrom    string
	SpecialPriceTo      string
}

type ProductImage struct {
	Id               int        `bson:"seqId"`
	ImageNo          int        `bson:"imageNo"`
	Orientation      string     `bson:"orientation" mapstructure:"orientation"`
	Main             int        `bson:"main"`
	OriginalFileName string     `bson:"originalFilename"`
	ImageName        string     `bson:"imageName"`
	UpdatedAt        *time.Time `bson:"updatedAt"`
}

type ProductVideo struct {
	Id        int       `bson:"seqId" xorm:"'id_video'"`
	FileName  string    `bson:"fileName" xorm:"'file_name'"`
	Thumbnail string    `bson:"thumbnail" xorm:"'thumbnail'"`
	Size      int       `bson:"size" xorm:"'size'"`
	Duration  int       `bson:"duration" xorm:"'duration'"`
	Status    string    `bson:"status" xorm:"'status'"`
	CreatedAt time.Time `bson:"createdAt" xorm:"'created_at'"`
	UpdatedAt time.Time `bson:"updatedAt" xorm:"'updated_at'"`
}

type ProductFetchCriteria struct {
	SellerId    int
	Status      string
	MinId       int
	MaxId       int
	Type        string
	BrandId     int
	PromotionId int
}

func (c ProductFetchCriteria) GetCriteriaType() string {
	/*if c.SellerId > 0 {
		return "seller"
	}
	if c.Status == "active" {
		return "active"
	}*/
	return c.Type
}
