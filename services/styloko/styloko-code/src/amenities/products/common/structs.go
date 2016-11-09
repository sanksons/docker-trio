package common

import (
	factory "common/ResourceFactory"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
)

const SOLR_DAYS_LIMIT = 60

var RetailPartners = []int{16503, 16505, 16507, 8845}

type PrdctAttrUpdateCndtn struct {
	ProductSku  string
	ProductType string
	IsGlobal    bool
	PattrMap    map[string]*Attribute
}

type ProductAttrSystemUpdate struct {
	ProConfigId int
	AttrName    string
	AttrValue   interface{}
}

type AttrSearchCondition struct {
	Name        string
	ProductType string
	IsGlobal    bool
	SetId       int
}

type M map[string]interface{}

func (m M) ToString() string {
	return ""
}

func (m M) ToJson() string {
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}
func (m M) ToStringKeys() string {
	var ret string
	for k, _ := range m {
		ret = ret + "; " + k
	}
	return ret
}

type ProSizeChart struct {
	Id   int         `bson:"id" json:"id"`
	Data interface{} `bson:"data" json:"data"`
}

// Store Product data In Cache
type ProductCache struct {
	Data       interface{} `json:"data"`
	Visibility bool        `json:"visibility"`
}

type ProductCacheCollection map[string]ProductCache

func (c ProductCacheCollection) ToStringKeys() string {
	var ret string
	for k, _ := range c {
		ret = ret + "; " + k
	}
	return ret
}

type ProductError struct {
	Code    int
	Message string
}

func (pe *ProductError) Error() string {
	return pe.Message
}

type Category struct {
	Id     int    `bson:"seqId"`
	Left   int    `bson:"lft"`
	Right  int    `bson:"rgt"`
	Parent int    `bson:"parent"`
	Name   string `bson:"name"`
	UrlKey string `bson:"urlKey"`
}

type ProductGroup struct {
	Id   int    `bson:"seqId" json:"seqId"`
	Name string `bson:"name" json:"name"`
}

type TaxClass struct {
	Id         int     `bson:"seqId"`
	Name       string  `bson:"name"`
	Position   int     `bson:"position"`
	IsDefault  int     `bson:"isDefault"`
	TaxPercent float64 `bson:"taxPercent"`
}

type Seller struct {
	Id                int                `json:"seqId"`
	OrgName           string             `json:"orgName"`
	SellerName        string             `json:"slrName"`
	SellerId          string             `json:"slrId"`
	Status            string             `json:"status"`
	Address           string             `json:"addr1"`
	Address2          string             `json:"addr2"`
	Phone             string             `json:"phn"`
	City              string             `json:"city"`
	PostCode          int                `json:"pstcode"`
	CountryCode       string             `json:"cntryCode"`
	Email             string             `json:"ordrEml"`
	Rating            string             `json:"rating"`
	CreatedAt         time.Time          `json:"crtdAt"`
	UpdatedAt         time.Time          `json:"updtdAt"`
	Sync              bool               `json:"sync"`
	UpdatedCommission []SellerCommission `json:"updtComisn"`
	SellerCustInfo    SellerCustomerInfo `json:"slrCustInfo"`
}

type SellerCustomerInfo struct {
	ProcessingTime   string `json:"proctime"`
	DispatchLocation string `json:"dispatchLoc"`
}

type SellerCommission struct {
	CategoryId int     `json:"ctgrId"`
	Percentage float64 `json:"percntge"`
}

type Brand struct {
	SeqId       int    `json:"id"`
	Status      string `json:"status"`
	Name        string `json:"name"`
	Position    int    `json:"position"`
	UrlKey      string `json:"urlKey"`
	ImageName   string `json:"imageName"`
	BrandClass  string `json:"brandClass"`
	IsExclusive int    `json:"isExclusive"`
	BrandInfo   string `json:"brandInfo"`
}

type Score struct {
	Final              float64  `json:"final"`
	FinalWA            float64  `json:"finalWa"`
	MobileA            float64  `json:"mobileA"`
	MobileB            float64  `json:"mobileB"`
	AppA               float64  `json:"appA"`
	Novelty            float64  `json:"novelty"`
	Boost              float64  `json:"boost"`
	SimpleAvailability float64  `json:"simpleAvailability"`
	TopSeller          *float64 `json:"topSeller"`
	Availability       float64  `json:"availability"`
	Random             float64  `json:"random"`
	New                int      `json:"productNew"`
}

func (sc *Score) GetBoost(configId int) float64 {
	sql := `SELECT
			  boost
			FROM catalog_product_boost
			WHERE fk_catalog_config = ` + strconv.Itoa(configId)
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return 0
	}
	result, sqlerr := driver.Query(sql)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return 0
	}
	defer result.Close()
	var boost int
	for result.Next() {
		result.Scan(&boost)
	}
	return Round(float64(boost)/float64(200)+float64(0.5), .5, 4)
}

func (sc *Score) GetSimpleAvailabilityScore(availableSimples int, totalSimples int) (float64, error) {
	return strconv.ParseFloat(
		fmt.Sprintf("%.4f", float64(availableSimples)/float64(totalSimples)),
		64,
	)
}

func (sc *Score) GetRandomScore() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Round(r.Float64(), .5, 4)
}

func (sc *Score) IsNewProduct(activatedAt *time.Time) int {
	if activatedAt == nil {
		return 0
	}
	now := time.Now()
	diff := now.Sub(*activatedAt)
	days := int(diff.Hours() / 24)
	if days <= SOLR_DAYS_LIMIT {
		return 1
	}
	return 0
}

type PriceMap struct {
	MaxPrice            float64  `json:"maxPrice"`
	Price               float64  `json:"price"`
	MaxOriginalPrice    float64  `json:"maxOriginalPrice"`
	OriginalPrice       float64  `json:"originalPrice"`
	SpecialPrice        *float64 `json:"specialPrice"`
	SpecialPriceFrom    *string  `json:"specialPriceFrom"`
	SpecialPriceTo      *string  `json:"specialPriceTo"`
	MaxSavingPercentage float64  `json:"maxSavingPercentage"`
	DiscountedPrice     *float64 `json:"discountedPrice"`
}

type ProductSmall struct {
	Id  int    `bson:"seqId"`
	Sku string `bson:"sku"`
}

type ProductSmallSimple struct {
	Sku       string `bson:"sku"`
	SellerSku string `bson:"sellerSku"`
}

type ProductSmallSimples struct {
	Simples []ProductSmallSimple `bson:"simples"`
}

type ProUpdateCriteria struct {
	ActivatedAt TimeNull
	PetApproved IntNull
	Status      StringNull
}

type StringNull struct {
	Isset bool
	Value *string
}

func (s StringNull) GetStringValue() string {
	if s.Value == nil {
		return ""
	}
	return *s.Value
}

type IntNull struct {
	Isset bool
	Value *int
}

func (i IntNull) GetStringValue() string {
	if i.Value == nil {
		return ""
	}
	return strconv.Itoa(*i.Value)
}

type FloatNull struct {
	Isset bool
	Value *float64
}

type TimeNull struct {
	Isset bool
	Value *time.Time
}

func (t TimeNull) GetStringValue() string {
	if t.Value == nil {
		return ""
	}
	return ToMySqlTime(t.Value)
}

//pet status handler
type PetStatus struct {
	Created bool
	Edited  bool
	Image   bool
}

type PriceGetter struct {
	DoNotPickPriceFrmMysql bool `bson:"-" json:"-"`
	UseMaster              bool `bson:"-" json:"-"`
}

//Update Pattterns
type PriceUpdate struct {
	SimpleId        int
	Price           float64
	UpdateSP        bool //if we need to change sp
	SpecialPrice    *float64
	SpecialFromDate *time.Time
	SpecialToDate   *time.Time
}

type JabongDiscount struct {
	SimpleId int
	Discount float64
	FromDate *time.Time
	ToDate   *time.Time
}

type PriceDetails struct {
	Price               *float64
	SpecialPrice        *float64
	SpecialFromDate     *string
	SpecialToDate       *string
	DiscountedPrice     float64
	MaxSavingPercentage float64
}

type SellerData struct {
	SellerIds []string `json:"sellerIds"`
}
