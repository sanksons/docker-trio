package common

import (
	"amenities/services/standardsize"
	factory "common/ResourceFactory"
	"common/utils"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jabong/floRest/src/common/sqldb"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type ProductSimple struct {
	Id                     int                   `bson:"seqId" mapstructure:"id_catalog_simple"`
	SKU                    string                `bson:"sku" mapstructure:"sku"`
	SellerSKU              string                `bson:"sellerSku"`
	BarcodeEan             string                `bson:"barcodeEan"`
	EanCode                string                `bson:"eanCode"`
	SupplierSKU            *string               `bson:"supplierSku" mapstructure:"sku_supplier_simple"`
	Price                  *float64              `bson:"price" mapstructure:"price"`
	OriginalPrice          *float64              `bson:"originalPrice" mapstructure:"original_price"`
	SpecialPrice           *float64              `bson:"specialPrice" mapstructure:"special_price"`
	SpecialFromDate        *time.Time            `bson:"specialFromDate" mapstructure:"special_from_date"`
	SpecialToDate          *time.Time            `bson:"specialToDate" mapstructure:"special_to_date"`
	JabongDiscount         float64               `bson:"jabongDiscount" mapstructure:"jabong_discount"`
	JabongDiscountFromDate *time.Time            `bson:"jabongDiscountFromDate" mapstructure:"jabong_discount_from_date"`
	JabongDiscountToDate   *time.Time            `bson:"jabongDiscountToDate" mapstructure:"jabong_discount_to_date"`
	TaxClass               *int                  `bson:"taxClass" mapstructure:"fk_catalog_tax_class"`
	Attributes             map[string]*Attribute `bson:"attributes"`
	Global                 map[string]*Attribute `bson:"global"`
	Quantity               *int                  `bson:"-"`
	Status                 string                `bson:"status" mapstructure:"status"`
	CreationSource         *string               `bson:"creationSource" mapstructure:"creation_source_simple"`
	CreatedAt              time.Time             `bson:"createdAt" mapstructure:"created_at"`
	UpdatedAt              time.Time             `bson:"updatedAt" mapstructure:"updated_at"`
	Position               int                   `bson:-`
}

type ProdSimples []*ProductSimple

func (slice ProdSimples) Len() int {
	return len(slice)
}

func (slice ProdSimples) Less(i, j int) bool {
	return slice[i].Position < slice[j].Position
}

func (slice ProdSimples) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (simple *ProductSimple) GetSize(attrName string) string {
	var variation string
	if val, ok := simple.Attributes[utils.SnakeToCamel(attrName)]; ok {
		variaI, err := val.GetValue("value")
		if err != nil {
			logger.Error(err)
			return variation
		}
		var ok bool
		if variation, ok = variaI.(string); ok {
			return variation
		}
	}
	return variation
}

func (simple *ProductSimple) GetSizePosition(attrName string) string {
	if val, ok := simple.Attributes[utils.SnakeToCamel(attrName)]; ok {
		if val.OptionType != OPTION_TYPE_SINGLE {
			return ""
		}
		attrMongo, err := GetAdapter(DB_READ_ADAPTER).GetAttributeMongoById(val.Id)
		if err != nil {
			logger.Error("(simple *ProductSimple) GetSizePosition: " + err.Error())
			return ""
		}
		variaI, err := val.GetValue("value")
		if err != nil {
			logger.Error("(simple *ProductSimple) GetSizePosition: " + err.Error())
			return ""
		}
		for _, option := range attrMongo.Options {
			if option.Name == variaI {
				return strconv.Itoa(option.Position)
			}
		}
	}
	return ""
}

//
// Set Attributes of Product
//
func (simple *ProductSimple) SetAttributes(
	attrSetId int,
	attrs AttributeMapSet,
	isAddDefault bool,
	adapter string,
	inclusionList bool,
) error {

	var (
		global map[string]*Attribute
		normal map[string]*Attribute
	)

	if inclusionList {
		global = simple.Global
		normal = simple.Attributes
	} else {
		global = make(map[string]*Attribute)
		normal = make(map[string]*Attribute)
	}

	//set default value attributes
	if isAddDefault {
		defattrs := GetMappedAttributesSimple(attrSetId)
		for k, s := range defattrs {
			if _, ok := attrs[k]; !ok {
				//key does not exist, add default val
				attrs[k] = s
			}
		}
	}
	attrSlice, err := attrs.ProcessAtributes(PRODUCT_TYPE_SIMPLE, true, adapter)
	if err != nil {
		return err
	}

	if inclusionList {
		attrAllwd := AttributesInfo[ATTRIBUTES_ALLOWED].([]interface{})
		attrbts := make([]string, 0)
		for _, attrName := range attrAllwd {
			attrbts = append(attrbts, attrName.(string))
		}
		for _, v := range attrSlice {
			if utils.InArrayString(attrbts, v.Name) {
				if v.IsGlobal {
					global[utils.SnakeToCamel(v.Name)] = v
				} else {
					normal[utils.SnakeToCamel(v.Name)] = v
				}
			}
		}
	} else {
		for _, v := range attrSlice {
			if v.IsGlobal {
				global[utils.SnakeToCamel(v.Name)] = v
			} else {
				normal[utils.SnakeToCamel(v.Name)] = v
			}
		}
	}

	simple.Global = global
	simple.Attributes = normal
	return nil
}

//
// Prepare simple sku
//
func (s *ProductSimple) SetSKU(sku string) error {
	if s.Id <= 0 {
		return errors.New("Set Simple ID first")
	}
	ids := strconv.Itoa(s.Id)
	s.SKU = sku + "-" + ids
	return nil
}

func (s *ProductSimple) GetQuantityPck(packQty int) int {
	q := s.GetQuantity()
	if q == 0 || packQty == 1 {
		return q
	}
	return (q / packQty)
}

//
// Get Product quantity
//
func (s *ProductSimple) GetQuantity() int {
	//fetch quantity form Redis
	if s.Quantity != nil {
		return *s.Quantity
	}
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, STOCK_CALL)
	defer logger.EndProfile(profiler, STOCK_CALL)

	key := "{stock}_" + strconv.Itoa(s.Id)
	errPrefix := fmt.Sprintf("(s *ProductSimple)#GetQuantity(%s):", key)
	client, rerr := factory.GetStockRedisDriver("Product")
	if rerr != nil {
		logger.Error(errPrefix + "cannot initiate redis" + rerr.Error())
		return 0
	}

	content := client.HGetAllMap(key)
	resultMap, err := content.Result()
	if err != nil {
		logger.Error(errPrefix + err.Error())
		return 0
	}
	var quantity int
	var reserved int
	if v, ok := resultMap["quantity"]; ok {
		quantity, _ = strconv.Atoi(v)
	}
	if v, ok := resultMap["reserved"]; ok {
		reserved, _ = strconv.Atoi(v)
	}
	total := (quantity - reserved)
	if total > 0 {
		s.Quantity = &total
	} else {
		var q int
		s.Quantity = &(q)
	}
	return *s.Quantity
}

//
// Remove special price from product
//
func (s *ProductSimple) RemoveSpecialPrice() {
	s.SpecialPrice = nil
	s.SpecialFromDate = nil
	s.SpecialToDate = nil
}

//
// Add special price to product
//
func (s *ProductSimple) AddSpecialPrice(
	sprice *float64,
	from *string,
	to *string,
) error {
	if sprice == nil || from == nil || to == nil {
		return errors.New("Cannot set Special Price")
	}
	s.SpecialPrice = sprice

	f, err := FromMysqlTime(*from, true)
	if err != nil {
		return err
	}
	s.SpecialFromDate = f

	t, err := FromMysqlTime(*to, true)
	if err != nil {
		return err
	}
	*t = (*t).Add(time.Duration(TO_DATE_DIFF) * time.Second)
	s.SpecialToDate = t

	return nil
}

//
// Set updated At
//
func (s *ProductSimple) SetUpdatedAt() {
	s.UpdatedAt = time.Now()
	return
}

//
// Create New Simple to be inserted
// -> set sequence
// -> set CreatedAt
// -> set UpdatedAt
//
func CreateNewSimple(adapter string) (ProductSimple, error) {
	s := ProductSimple{}
	seqId, err := GetAdapter(adapter).GenerateNextSequence(SIMPLE_COLLECTION)
	if err != nil {
		logger.Error("CreateNewSimple()" + err.Error())
	}
	if seqId <= 0 {
		return s, errors.New("Unable to Generate Sequence")
	}
	s.Id = seqId
	t := time.Now()
	s.CreatedAt = t
	s.UpdatedAt = t
	return s, nil
}

func (s *ProductSimple) SetStandardSize(p *Product,
	dbAdapterName string) {

	varAttrName, err := p.AttributeSet.GetVariationAttributeName()
	if err != nil {
		logger.Error(fmt.Sprintf("Cannot get variation attribute name for %v", p.AttributeSet))
		return
	}
	brndSize := s.GetSize(varAttrName)
	stndrdSizeVal, ok := standardsize.GetStandardSize(p.AttributeSet.Id, p.BrandId, p.Leaf, brndSize)
	if !ok {
		logger.Error("No standard size found")
		return
	}
	attrSr := AttrSearchCondition{ATTR_STANDARD_SIZE,
		PRODUCT_TYPE_SIMPLE, true, 0}
	attrMongo, err := GetAdapter(dbAdapterName).GetAtrributeByCriteria(attrSr)
	if err != nil {
		logger.Error(err)
		return
	}
	attrOpt, err := attrMongo.GetAttrOptionByName(stndrdSizeVal)
	if err != nil {
		logger.Error(err)
		return
	}
	attr := make(AttributeMapSet, 0)
	attr[utils.ToString(attrMongo.SeqId)] = strconv.Itoa(attrOpt.SeqId)
	attrAdd, err := attr.ProcessAtributes(PRODUCT_TYPE_SIMPLE, true, dbAdapterName)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, val := range attrAdd {
		s.Global[utils.SnakeToCamel(val.Name)] = val
	}
	return
}

func (smp *ProductSimple) GetPriceDetails() PriceDetails {
	pd := PriceDetails{}

	if smp.Price == nil {
		logger.Error(fmt.Sprintf("GetPriceDetails(): PRICE NULL [%s]", smp.SKU))
		return pd
	}

	pd.Price = smp.Price

	dp := *smp.Price
	tmpSpFrom, tmpSpTo := SetSpecialPriceDates(smp)
	if (tmpSpFrom != nil && tmpSpTo != nil && smp.SpecialPrice != nil) &&
		time.Now().After(*tmpSpFrom) &&
		time.Now().Before(*tmpSpTo) {
		dp = *smp.SpecialPrice
	}

	if dp < *smp.Price {
		pd.SpecialPrice = &dp
		pd.SpecialFromDate = ToMySqlTimeNull(smp.SpecialFromDate)
		pd.SpecialToDate = ToMySqlTimeNull(smp.SpecialToDate)
		pd.DiscountedPrice = dp
		pd.MaxSavingPercentage = math.Floor((((*smp.Price - dp) * 100) / (*smp.Price)) + 0.5)
	}

	return pd
}

func (smp *ProductSimple) SetPriceDetailsFrmMysql(useMaster bool) error {
	var err error
	var driver sqldb.SqlDbInterface
	if useMaster {
		driver, err = factory.GetDefaultMysqlDriver()
	} else {
		driver, err = factory.GetDefaultMysqlDriverSlave()
	}
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return fmt.Errorf("Cannot initiate Mysql: %s", err.Error())
	}
	sql := `SELECT price, special_price, special_from_date, special_to_date
			FROM catalog_simple WHERE id_catalog_simple = ?`
	result, sqlerr := driver.Query(sql, smp.Id)
	if sqlerr != nil {
		logger.Error("(SetPriceDetailsFrmMysql)#Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return fmt.Errorf("(SetPriceDetailsFrmMysql)#Cannot initiate Mysql: %s", err.Error())
	}
	defer result.Close()
	for result.Next() {
		scerr := result.Scan(
			&smp.Price,
			&smp.SpecialPrice,
			&smp.SpecialFromDate,
			&smp.SpecialToDate,
		)
		if scerr != nil {
			return fmt.Errorf("(SetPriceDetailsFrmMysql)#Scan Error: %s", scerr.Error())
		}
	}
	if smp.SpecialToDate != nil {
		todate := smp.SpecialToDate.Add(time.Hour * 24)
		smp.SpecialToDate = &todate
	}
	return nil
}
