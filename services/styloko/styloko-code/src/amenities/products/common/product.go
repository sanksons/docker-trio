package common

import (
	"amenities/products/jbus"
	attributes "amenities/services/attributes"
	brandService "amenities/services/brands"
	factory "common/ResourceFactory"
	config "common/appconfig"
	"common/notification"
	"common/notification/datadog"
	"common/redis"
	"common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jabong/floRest/src/common/monitor"
	"github.com/jabong/floRest/src/common/utils/logger"
	"github.com/kennygrant/sanitize"
)

type Product struct {
	SeqId             int                   `bson:"seqId" json:"seqId"`
	SKU               string                `bson:"sku" json:"sku"`
	ProductSet        int                   `bson:"productSet" json:"productSet"`
	SupplierSKU       *string               `bson:"supplierSkuConfig,omitempty" json:"supplierSkuConfig"`
	Name              string                `bson:"name" json:"name"`
	Description       string                `bson:"description" json:"description"`
	UrlKey            string                `bson:"urlKey" json:"urlKey"`
	SellerId          int                   `bson:"sellerId" json:"sellerId"`
	DisplayStockedOut int                   `bson:"displayIfOutOfStock" json:"displayIfOutOfStock"`
	Status            string                `bson:"status" json:"status"`
	PetStatus         string                `bson:"petStatus" json:"petStatus"`
	PetApproved       int                   `bson:"petApproved" json:"petApproved"`
	CatalogImport     *int                  `bson:"catalogImport,omitempty" json:"catalogImport"`
	AttributeSet      ProAttributeSet       `bson:"attributeSet" json:"attributeSet"`
	BrandId           int                   `bson:"brandId" json:"brandId"`
	TY                int                   `bson:"ty" json:"ty"`
	SizeChart         ProSizeChart          `bson:"sizeChart" json:"sizeChart"`
	TaxClass          *int                  `bson:"taxClass" json:"taxClass"`
	Group             *ProductGroup         `bson:"group" json:"group"`
	ShipmentType      int                   `bson:"shipmentType" json:"shipmentType"`
	Categories        []int                 `bson:"categories" json:"categories"`
	PrimaryCategory   int                   `bson:"primaryCategory" json:"primaryCategory"`
	Leaf              []int                 `bson:"leaf" json:"leaf"`
	Images            []*ProductImage       `bson:"images,omitempty" json:"images"`
	Videos            []*ProductVideo       `bson:"videos,omitempty" json:"videos"`
	Global            map[string]*Attribute `bson:"global" json:"global"`
	Attributes        map[string]*Attribute `bson:"attributes" json:"attributes"`
	Simples           ProdSimples           `bson:"simples" json:"simples"`
	ApprovedAt        *time.Time            `bson:"approvedAt,omitempty" json:"approvedAt"`
	CreatedAt         *time.Time            `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt         *time.Time            `bson:"updatedAt,omitempty" json:"updatedAt"`
	ActivatedAt       *time.Time            `bson:"activatedAt,omitempty" json:"activatedAt"`
	Pricegetter       PriceGetter           `bson:"-" json:"-"`
}

//
// Load product based on the ProductSetId
//
func (p *Product) LoadByProductSet(id int) error {
	pro, err := GetCurrentAdapter().GetByProductSet(id)
	if err != nil {
		return err
	}
	*p = pro
	return nil
}

//
// Load Product based on SKU
//
func (p *Product) LoadBySku(sku string, adapter string) error {
	pro, err := GetAdapter(adapter).GetBySku(sku)
	if err != nil {
		return err
	}
	*p = pro
	return nil
}

//
// Load Product By SeqId
//
func (p *Product) LoadBySeqId(id int, adapter string) error {
	pro, err := GetAdapter(adapter).GetById(id)
	if err != nil {
		return err
	}
	*p = pro
	return nil
}

//
// Append Simple to the list of existing simples.
//
func (p *Product) AppendSimple(simple ProductSimple) {
	p.Simples = append(p.Simples, &simple)
	return
}

//
// Set Attributes for product
//
func (p *Product) SetAttributes(attrs AttributeMapSet,
	isAddDefault bool,
	adapter string,
	inclusionList bool,
) error {

	var (
		global map[string]*Attribute
		normal map[string]*Attribute
	)

	if inclusionList {
		//Product Update Request
		global = p.Global
		normal = p.Attributes
	} else {
		//Product Create Request
		global = make(map[string]*Attribute)
		normal = make(map[string]*Attribute)

		//Added only during product create request
		tmpAttr := p.SetIsReturnable(attrs)
		attrs = tmpAttr
	}

	//set default value attributes
	if isAddDefault {
		defattrs := GetMappedAttributesConfig(p.AttributeSet.Id, attrs, adapter)
		for k, s := range defattrs {
			if _, ok := attrs[k]; !ok {
				//key does not exist, add default val
				attrs[k] = s
			} else {
				//key Exists
				if (attrs[k] == nil || attrs[k] == "") && (utils.InArrayString(ServicibilityAttributes, k)) {
					attrs[k] = s
				}
			}
		}
	}
	attrSlice, err := attrs.ProcessAtributes(PRODUCT_TYPE_CONFIG, true, adapter)
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

	p.Global = global
	p.Attributes = normal
	return nil
}

func (p *Product) SetPetStatus(ps PetStatus) error {
	var st []string
	if ps.Created {
		st = append(st, "creation")
	}
	if ps.Edited {
		st = append(st, "edited")
	}
	if ps.Image {
		st = append(st, "images")
	}
	p.PetStatus = strings.Join(st, ",")
	return nil
}

func (p *Product) GetPetStatus() PetStatus {
	st := strings.Split(p.PetStatus, ",")
	ps := PetStatus{}
	for _, s := range st {
		if s == "creation" {
			ps.Created = true
		}
		if s == "edited" {
			ps.Edited = true
		}
		if s == "images" {
			ps.Image = true
		}
	}
	return ps
}

//
// Set Product Group
//
func (p *Product) SetProductGroup(name string) error {
	if name == "" {
		return nil
	}
	group, err := GetCurrentAdapter().GetProductGroupByName(name)
	if err != nil {
		return err
	}
	p.Group = &group
	return nil
}

//
// Set Url Key for product
//
func (p *Product) SetUrlKey() error {
	brand, _ := p.GetBrandInfo()
	p.UrlKey = fmt.Sprintf("%s-%s-%d",
		strings.Replace(brand.UrlKey, "/", "-", -1),
		sanitize.BaseName(strings.Replace(p.Name, "/", "-", -1)),
		p.SeqId)
	return nil
}

//
// Set Attributeset Data for product.
//
func (p *Product) SetAttributeSet(id int, adapter string) error {
	as, err := GetAdapter(adapter).GetProAttributeSetById(id)
	if err != nil {
		return errors.New("(p *Product)#SetAttributeSet" + err.Error())
	}
	p.AttributeSet = as
	return nil
}

//
// Get PrimaryCategoryId from supplied categories
//
func (p *Product) GetPrimaryCategory(adapter string) int {
	return p.PrimaryCategory
}

//
// Prepare Leaf categories
//
func (p *Product) PrepareLeafCategories(adapter string) error {

	if len(p.Categories) == 0 {
		msg := "Set categories first"
		return errors.New("(p *Product)#SetLeafCategories" + msg)
	}

	c, err := GetAdapter(adapter).GetCategoriesByIds(p.Categories)
	if err != nil {
		return errors.New("(p *Product)#SetLeafCategories" + err.Error())
	}
	var leafCats []int
	for _, v := range c {
		if (v.Left + 1) == v.Right {
			//its leaf
			leafCats = append(leafCats, v.Id)
		}
	}
	p.Leaf = leafCats
	return nil
}

//
// Prepare Price map
//
func (p *Product) PreparePriceMap() PriceMap {
	pricemap := PriceMap{}

	//check if we need to skip skuQuantitycheck
	var skipStockCheck bool = true
	for _, s := range p.Simples {
		if s.GetQuantity() > 0 {
			skipStockCheck = false
		}
	}

	var minInd int = -1
	var minPrice *float64

	for index, s := range p.Simples {
		//skip simple if quantity id 0
		if s.GetQuantity() <= 0 && !skipStockCheck {
			continue
		}
		// if price is less than 1 or simple is not active, re-loop
		if (s.Price == nil) || (*s.Price < 1) || (s.Status != STATUS_ACTIVE) {
			continue
		}
		// assign price to tmp var
		var price float64 = *s.Price
		//set max/ min price
		if price > pricemap.MaxPrice {
			pricemap.MaxPrice = price
		}

		//min, max original price
		if s.OriginalPrice != nil {
			if *s.OriginalPrice > pricemap.MaxOriginalPrice {
				pricemap.MaxOriginalPrice = *s.OriginalPrice
			}
		}

		dp := price
		tmpSpFrom, tmpSpTo := SetSpecialPriceDates(s)
		if (tmpSpFrom != nil && tmpSpTo != nil && s.SpecialPrice != nil) &&
			time.Now().After(*tmpSpFrom) &&
			time.Now().Before(*tmpSpTo) {
			dp = *s.SpecialPrice
		}

		if (minPrice == nil) || (dp < *minPrice) {
			minInd = index
			minPrice = &dp
		}
	}

	if minInd < 0 {
		return PriceMap{}
	}
	minSimple := p.Simples[minInd]
	pricemap.Price = *minSimple.Price
	if minSimple.OriginalPrice != nil {
		pricemap.OriginalPrice = *minSimple.OriginalPrice
	} else {
		pricemap.OriginalPrice = *minSimple.Price
	}
	if *minPrice < *minSimple.Price {
		pricemap.SpecialPrice = minPrice
		pricemap.SpecialPriceFrom = ToMySqlTimeNull(minSimple.SpecialFromDate)
		pricemap.SpecialPriceTo = ToMySqlTimeNull(minSimple.SpecialToDate)
		pricemap.DiscountedPrice = minPrice
		pricemap.MaxSavingPercentage = math.Floor((((*minSimple.Price - *minPrice) * 100) / *minSimple.Price) + 0.5)
	}
	return pricemap
}

//
// Generate new sku for product
//
func (p *Product) PrepareSKU(adapter string) error {
	//get SKu suffix
	suffix := PRODUCT_SUFFIX
	if p.BrandId == 0 {
		return errors.New("(p *Product)#PrepareSKU: Set Brand First")
	}
	var err error
	brnd, err := brandService.ById(p.BrandId)
	if err != nil {
		return err
	}
	illegalChars := regexp.MustCompile(`[^[a-zA-Z0-9]]*`)
	urlKey := illegalChars.ReplaceAllString(brnd.UrlKey, "")
	var sku string
	if len(urlKey) > 2 {
		urlKey = urlKey[0:2]
	}
	if len(urlKey) == 1 {
		urlKey = fmt.Sprintf("%s%s", "X", urlKey)
	}

	sku = strings.ToUpper(urlKey)
	sku += fmt.Sprintf("%03d", p.BrandId%1000)

	if p.AttributeSet.Id == 0 {
		return errors.New("(p *Product)#PrepareSKU: Set Attr Set First")
	}
	idAttrSet := p.AttributeSet.Id

	mgoSess := factory.GetMongoSession("Attributeset")
	defer mgoSess.Close()

	aSet := attributes.GetAttributeSetById(
		idAttrSet, mgoSess,
	)
	sku += strings.ToUpper(aSet.Identifier)

	part := func(id int) string {
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		chPart := ""
		value := id
		for i := 0; i < 3; i++ {
			chPart = string(chars[value%len(chars)]) + chPart
			value = int(math.Floor(float64(value / len(chars))))
		}
		num := 99 - (id % 100)
		numStr := strconv.Itoa(num)
		numPart := fmt.Sprintf("%02s", numStr)
		return numPart + chPart
	}(p.SeqId)
	sku += part
	sku += suffix
	p.SKU = sku
	return nil
}

//
// Get Brand Info for product
//
func (p *Product) GetBrandInfo() (Brand, error) {
	brand := Brand{}
	brandInfo, err := brandService.ById(p.BrandId)
	if err != nil {
		return brand, err
	}
	brand.SeqId = brandInfo.SeqId
	brand.Name = brandInfo.Name
	brand.Status = brandInfo.Status
	brand.IsExclusive = brandInfo.IsExclusive
	brand.ImageName = brandInfo.ImageName
	brand.BrandClass = brandInfo.BrandClass
	brand.UrlKey = brandInfo.UrlKey
	return brand, nil
}

//
// Get Seller Info for Product
//
func (p *Product) GetSellerInfo(org config.HttpAPIConfig, useCmsn bool) (Seller, error) {
	profiler := logger.NewProfiler()

	var (
		seller      Seller
		searchQuery string
	)
	if useCmsn {
		logger.StartProfile(profiler, FE_ORG_COM)
		defer logger.EndProfile(profiler, FE_ORG_COM)
		catArr := FetchCategoryTree(p.PrimaryCategory)
		catStr := ""
		for _, val := range catArr {
			if catStr == "" {
				catStr = fmt.Sprintf("%d", val)
			} else {
				catStr = fmt.Sprintf("%s,%d", catStr, val)
			}
		}
		searchQuery = fmt.Sprintf("/commissions/?q=id.eq~%d___categories.in~[%s]", p.SellerId, catStr)
	} else {
		logger.StartProfile(profiler, FE_ORG_API)
		defer logger.EndProfile(profiler, FE_ORG_API)
		searchQuery = fmt.Sprintf("/sellers/%d", p.SellerId)
	}
	url := fmt.Sprintf("%s%s%s", org.Host, org.Path, searchQuery)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return seller, errors.New(
			"(p *Product)#GetSellerInfo1(): Cannot getSeller Info" + err.Error())
	}
	client := HttpClient
	resp, err := client.Do(req)
	if err != nil {
		return seller, fmt.Errorf(
			"(p *Product)#GetSellerInfo2(): Cannot getSeller Info %s, %d",
			err.Error(), p.SellerId)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return seller, fmt.Errorf(
			"(p *Product)#GetSellerInfo3(): Cannot getSeller Info %d, Status code is %d",
			p.SellerId, resp.StatusCode,
		)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return seller, fmt.Errorf(
			"(p *Product)#GetSellerInfo4(): Cannot getSeller Info %s, %d. Body parsing failed",
			err.Error(), p.SellerId)
	}
	var data map[string]json.RawMessage
	err = json.Unmarshal(body, &data)
	if err != nil {
		return seller, fmt.Errorf(
			"(p *Product)#GetSellerInfo5(): Cannot getSeller Info %s, %d. Unmarshal failed",
			err.Error(), p.SellerId)
	}
	err = json.Unmarshal(data["data"], &seller)
	if err != nil {
		return seller, fmt.Errorf(
			"(p *Product)#GetSellerInfo6(): Cannot getSeller Info %s, %d. Unmarshal failed",
			err.Error(), p.SellerId)
	}
	if useCmsn && len(seller.UpdatedCommission) == 0 {
		return seller, fmt.Errorf(
			"(p *Product)#GetSellerInfo7(): Cannot getSeller commission %d.", p.SellerId)
	}
	return seller, nil
}

func (p *Product) GetSellerInfoWithRetries(org config.HttpAPIConfig,
	useCmsn bool) (Seller, error) {

	var (
		data Seller
		err  error
	)

	for try := 0; try < SELLER_RETRY_COUNT; try++ {
		data, err = p.GetSellerInfo(org, useCmsn)
		if err == nil {
			return data, err
		}
		time.Sleep(10 * time.Microsecond)
	}
	return data, err
}

//
// Get Size key based on product Atrribute set
//
func (p *Product) GetSizeKey() string {
	m := GetAttributeSet2VariationMapping()
	if v, ok := m[p.AttributeSet.Name]; ok {
		return v
	}
	return ""
}

//
// Get Total availbale stock for product
//
func (p *Product) GetTotalQuantity() int {
	var quantity int
	for _, simple := range p.Simples {
		q := simple.GetQuantity()
		quantity = quantity + q
	}
	return quantity
}

func (p *Product) CalculateScore() *Score {
	score := &Score{}
	sql := `SELECT
			  score_final,
			  score_final_wa,
			  score_mobile_a,
			  score_mobile_b,
			  score_app_a
			FROM product_solr_score
			WHERE sku = "` + p.SKU + `"`
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return score
	}
	result, sqlerr := driver.Query(sql)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return score
	}

	for result.Next() {
		result.Scan(
			&score.Final, &score.FinalWA, &score.MobileA,
			&score.MobileB, &score.AppA,
		)
	}
	result.Close()
	score.Final = Round(score.Final, .1, 6)
	score.FinalWA = Round(score.FinalWA, .1, 6)
	score.MobileA = Round(score.MobileA, .1, 6)
	score.MobileB = Round(score.MobileB, .1, 6)
	score.AppA = Round(score.AppA, .1, 6)

	score.Novelty = func(actTime *time.Time) float64 {
		var activatedAt time.Time
		if actTime != nil {
			activatedAt = *actTime
		}
		t := time.Since(activatedAt).Hours() / float64(24)
		return Round(math.Pow(2, -(math.Pow((t/float64(30)), 2))), .5, 4)
	}(p.ActivatedAt)

	score.Boost = score.GetBoost(p.SeqId)
	var availableSimplesCount int
	var totalSimplesCount int
	for _, s := range p.Simples {
		if s.GetQuantity() > 0 {
			availableSimplesCount += 1
		}
		totalSimplesCount += 1
	}
	score.SimpleAvailability, _ = score.GetSimpleAvailabilityScore(
		availableSimplesCount, totalSimplesCount,
	)
	score.Availability = float64(0)
	score.TopSeller = nil
	score.Random = score.GetRandomScore()
	score.New = score.IsNewProduct(p.ActivatedAt)
	return score
}

func (p *Product) GetWeightedAvailability() (float64, error) {
	if len(p.Categories) == 0 {
		return 0, nil
	}
	cat := p.Categories[len(p.Categories)-1]
	var sizes = make(map[string]int)
	var s []string
	for _, simple := range p.Simples {
		var size string
		attrName, err := p.AttributeSet.GetVariationAttributeName()
		if err != nil {
			logger.Error(err)
			continue
		}
		size = simple.GetSize(attrName)
		sizes[size] = simple.GetQuantity()
		s = append(s, size)
	}
	sql := `SELECT
			  size,
			  weight
			FROM category_size_weight
			WHERE fk_catalog_category = ?
			AND size IN ('` + strings.Join(s, "','") + `')`
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return 0, err
	}
	result, sqlerr := driver.Query(sql, cat)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return 0, err
	}
	defer result.Close()
	var mapping = make(map[string]float64)
	for result.Next() {
		var size string
		var weight float64
		err := result.Scan(&size, &weight)
		if err != nil {
			logger.Error("Error in size and weight row %s", err.Error())
			continue
		}
		mapping[size] = weight
	}
	//handle case for 'free size', 'one size', 'standard'
	exceptionSizes := []string{"free size", "one size", "standard"}
	if len(p.Simples) == 1 {
		for _, val := range exceptionSizes {
			for k, v := range sizes {
				if strings.ToLower(val) == strings.ToLower(k) {
					if v > 0 {
						return float64(1), nil
					}
				}
			}
		}
	}
	//todo: Confirm weighted availibility logic
	missCount := len(sizes) - len(mapping)
	var missing float64
	if missCount == len(p.Simples) {
		if len(p.Simples) > 4 {
			missing = float64(1) / float64(len(sizes))
		} else {
			missing = float64(1) / float64(4)
		}
	} else if missCount > 0 {
		missing = float64(1) / float64(len(sizes))
	}

	var weightedScore float64
	for key, val := range sizes {
		if val > 0 {
			if _, ok := mapping[key]; !ok || mapping[key] <= 0.00 {
				weightedScore = weightedScore + missing
			} else {
				weightedScore = weightedScore + mapping[key]
			}
		}
	}

	if weightedScore > 0.6 {
		return float64(1), nil
	}
	return strconv.ParseFloat(fmt.Sprintf("%.6f", weightedScore), 64)

	return 0, nil
}

//
// Insert Product
//
func (p *Product) InsertOrUpdate(adapter string) error {
	err := GetAdapter(adapter).SaveProduct(*p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) PushToMemcache(comment string) error {
	vc := VisibilityChecker{
		Product:        *p,
		VisibilityType: VISIBILITY_PDP,
	}
	isVisible := vc.IsVisible()

	var command string
	var alicedata []byte

	if isVisible {
		//prepare response.
		resp := ProductResponse{
			Product:        p,
			VisibilityType: VISIBILITY_PDP,
			Pricegetter:    p.Pricegetter,
		}
		data := resp.GetResponse("memcache")
		fresp := map[string]interface{}{
			utils.ToString(p.SeqId): data,
		}
		alicedata, _ = utils.JSONMarshal(fresp, true)
		command = "update"
	} else {
		alicedata, _ = utils.JSONMarshal([]string{p.SKU}, true)
		command = "delete"
	}
	stringTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	version := stringTime[0 : len(stringTime)-4]

	sqldriver, sqlerr := factory.GetDefaultMysqlDriver()
	if sqlerr != nil {
		return fmt.Errorf("(p *Product)#PushToMemcache(): %s", sqlerr.Error())
	}
	sql := `INSERT INTO alice_message
			 (timestamp, data, command, caller, comment, type)
			 VALUES (?,?,?,?,?,?)`

	_, serr := sqldriver.Execute(sql,
		version, alicedata, command, "styloko service", comment, "product",
	)
	if serr != nil {
		return fmt.Errorf("(p *Product)#PushToMemcache()2: %s", serr.DeveloperMessage)
	}
	return nil
}

func (p *Product) Publish(transactionId string, publishInVisible bool) {

	// Test Seller check
	sellersIgnore := AttributesInfo[SELLERS_IGNORE].([]interface{})
	slrStr := strconv.Itoa(p.SellerId)
	for _, slr := range sellersIgnore {
		if slrStr == slr.(string) {
			logger.Error(fmt.Sprintf("Product [%s]Publish Stopped, because this seller is not allowed", p.SKU))
			return
		}
	}

	if p.ActivatedAt == nil {
		//dont publish a product which dont have activatedAT
		logger.Error(fmt.Sprintf("Product [%s]Publish Stopped, because activated AT is NULL", p.SKU))
		return
	}
	vc := VisibilityChecker{
		Product:        *p,
		VisibilityType: VISIBILITY_DOOS,
	}
	isVisible := vc.IsVisible()

	//check if product is not visible
	//and we dont need to publish not visible products
	if (!publishInVisible) && (!isVisible) {
		//we dont need to do anything, return
		return
	}
	resp := ProductResponse{
		Product:        p,
		VisibilityType: VISIBILITY_DOOS,
		Pricegetter:    p.Pricegetter,
	}
	data := resp.GetResponse(EXPANSE_SOLR)

	// [NOTE]
	// log data that is being published to Bus.
	// this logging is only for testing purpose and should be
	// removed once the system becomes stable.
	//

	//var collectionName = "solrPushLogs_" + time.Now().Format(FORMAT_LOG_TIME)
	type SolrPublish struct {
		ConfigId    int       `bson:"configId"`
		Sku         string    `bson:"sku"`
		Data        string    `bson:"data"`
		Type        string    `bson:"type"`
		PublishedAt time.Time `bson:"publishedAt"`
	}
	var pdata SolrPublish
	pdata.ConfigId = p.SeqId
	pdata.Sku = p.SKU
	pdata.PublishedAt = time.Now()
	if isVisible {
		pdata.Type = "Update"
		monitor.GetInstance().Count("solrPublishLogs_update", 1, []string{"styloko"}, 0.7)
	} else {
		pdata.Type = "Delete"
		monitor.GetInstance().Count("solrPublishLogs_delete", 1, []string{"styloko"}, 0.7)
	}
	mData, _ := utils.JSONMarshal(data, true)
	pdata.Data = string(mData)
	//mSession := factory.GetMongoSession("logs")
	//defer mSession.Close()
	//mSession.SetCollection(collectionName).Insert(pdata)

	//
	// Logging code ends
	//

	message := jbus.GetNewProductMessage()
	if transactionId == "" {
		transactionId = jbus.GenerateTransactionId()
	}
	message.TransactionId = &transactionId
	message.Type = jbus.MESSAGE_TYPE
	message.TypeOfChange = jbus.TYPE_PRODUCT_DELETE
	if isVisible {
		message.TypeOfChange = jbus.TYPE_PRODUCT_UPDATE
	}
	m := M{}
	m["module"] = jbus.MODULE_NAME
	m["payload"] = data
	message.Data = m
	err := message.Publish()
	if err != nil {
		notification.SendNotification(
			"Product Publish Failed",
			fmt.Sprintf("Product:%d, Message:%s", p.SeqId, err.Error()),
			[]string{TAG_PRODUCT},
			datadog.ERROR,
		)
		logger.Error(err.Error())
	}
	return
}

func CreateNewProduct(adapter string) (Product, error) {
	p := Product{}
	seqId, err := GetAdapter(adapter).GenerateNextSequence(PRODUCT_COLLECTION)
	if seqId <= 0 {
		return p, errors.New("Unable to Generate Sequence" + err.Error())
	}
	p.SeqId = seqId
	t := time.Now()
	p.CreatedAt = &t
	p.UpdatedAt = &t
	return p, nil
}

//
// Contains Collection of Products
// Used when qurying for multiple products.
//
type ProductCollection struct {
	Products []Product
	Count    int
}

//Load products based on Sequence Id.
func (pc *ProductCollection) LoadByIds(ids []int) error {
	slice := []Product{}
	err := GetCurrentAdapter().GetByIds(ids, &slice)
	if err != nil {
		return err
	}
	pc.Products = slice
	pc.Count = len(slice)
	return nil
}

//Load products based on sku
func (pc *ProductCollection) LoadBySkus(skus []string) error {
	slice := []Product{}
	err := GetCurrentAdapter().GetBySkus(skus, &slice)
	if err != nil {
		return err
	}
	pc.Products = slice
	pc.Count = len(slice)
	return nil
}

func (p *Product) SetIsReturnable(attrs AttributeMapSet) AttributeMapSet {
	returnable, err := GetAdapter(DB_ADAPTER_MONGO).GetAtrributeByCriteria(AttrSearchCondition{
		Name:        IS_RETURNABLE,
		ProductType: PRODUCT_TYPE_CONFIG,
		IsGlobal:    true,
	})
	if err != nil {
		logger.Error(fmt.Errorf(
			"(p *Product) SetIsReturnable: Unable to get Returnable: %s", err.Error(),
		))
		return attrs
	}
	isReturnableIdStr := utils.ToString(returnable.SeqId)

	// If it is a market place seller product
	if !utils.InArrayInt(RetailPartners, p.SellerId) {
		isRtbMap := AttributesInfo[IS_RETURNABLE].(map[string]interface{})
		found := false
		for _, cat := range p.Leaf {
			if isReturnable, ok := isRtbMap[strconv.Itoa(cat)]; ok {
				attrs[isReturnableIdStr] = isReturnable.(string)
				found = true
				break
			}
		}
		if !found {
			//Setting 15 Days Free Return as default value
			logger.Warning("No is_returnable value for these leaf categories, setting default value")
			attrs[isReturnableIdStr] = DEF_IS_RETURNABLE
		}
	} else {
		// If it is a retail partner product
		_, ok := attrs[isReturnableIdStr]
		if !ok || attrs[isReturnableIdStr] == nil || attrs[isReturnableIdStr] == "" {
			//Setting 15 Days Free Return as default value
			logger.Warning("No is returnable value found in request, setting default value")
			attrs[isReturnableIdStr] = DEF_IS_RETURNABLE
		}
	}
	return attrs
}

func (p *Product) GetScProductValue() string {
	if v, ok := p.Global["scProduct"]; ok {
		val, err := v.GetValue("value")
		if err == nil {
			strVal, ok := val.(string)
			if ok {
				return strVal
			}
		}
	}
	return ""
}

func (p *Product) SetCatalogType(attrs AttributeMapSet, adapter string) (bool, error) {
	// Get id of attribute ty from DB.
	tyAttr, err := GetCatalogTyAttribute(adapter)
	if err != nil {
		return false, err
	}
	id := strconv.Itoa(tyAttr.SeqId)
	// check if ty attribute exist in SC-Request Attribute
	data, ok := attrs[id]
	if !ok {
		return false, nil
	}
	typeId, err := utils.GetInt(data)
	if err != nil {
		logger.Error(fmt.Errorf("(p *Product)#SetCatalogType: %s", err.Error()))
		return false, nil
	}
	//update only if typeId is greater than 0
	if typeId > 0 {
		p.TY = typeId
	}
	return true, nil
}

func (p *Product) GetPackQty() int {
	if qty, ok := p.Global["packQty"]; ok {
		pckQty, _ := utils.GetInt(qty.Value)
		return pckQty
	}
	return 1
}

func (p *Product) SortSimples() {
	varattrName, _ := p.AttributeSet.GetVariationAttributeName()
	for _, smp := range p.Simples {
		pos := smp.GetSizePosition(varattrName)
		posInt, err := strconv.Atoi(pos)
		if err != nil {
			smp.Position = 0
		} else {
			smp.Position = posInt
		}
	}
	sort.Sort(p.Simples)
}

func (p *Product) SetUpdateUID() (string, error) {
	id := utils.ToString(time.Now().UnixNano())
	driver, rerr := redis.GetDriver()
	if rerr != nil {
		return "", fmt.Errorf("(p *Product)#GetSetUpdateUUID()1: %s", rerr.Error())
	}
	key := fmt.Sprintf("styloko_productUpdateUUID_%d", p.SeqId)
	err := driver.Set(key, id)
	if err != nil {
		return "", fmt.Errorf("(p *Product)#GetSetUpdateUUID()2: %s", err.Error())
	}
	return id, nil
}

func (p *Product) GetUpdateUID() (string, error) {
	driver, rerr := redis.GetDriver()
	if rerr != nil {
		return "", fmt.Errorf("(p *Product)#GetUpdateUID()1: %s", rerr.Error())
	}
	key := fmt.Sprintf("styloko_productUpdateUUID_%d", p.SeqId)
	uid, err := driver.Get(key)
	if err != nil {
		return "", fmt.Errorf("(p *Product)#GetUpdateUID()2: %s", err.Error())
	}
	return uid, nil
}

func (p *Product) GetInactiveDeletedSimpleSizes() ([]string, error) {
	simpleSizesArray := []string{}
	p.SortSimples()
	prodSimpleArray := p.Simples
	sizeAttributeName, err := p.AttributeSet.GetVariationAttributeName()
	if err != nil {
		return simpleSizesArray, err
	}
	for _, simple := range prodSimpleArray {
		if simple.Status == "active" {
			continue
		}
		simpleSize := simple.GetSize(sizeAttributeName)
		simpleSizesArray = append(simpleSizesArray, simpleSize)
	}
	return simpleSizesArray, nil
}
