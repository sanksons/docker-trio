package common

import (
	categoryService "amenities/services/categories"
	factory "common/ResourceFactory"
	"common/notification"
	"common/notification/datadog"
	"common/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jabong/floRest/src/common/utils/logger"
)

type ProductMemcacheResponse struct {
	Attributes map[string]interface{} `json:"attributes"`
	Meta       map[string]interface{} `json:"meta"`
	Simples    map[string]interface{} `json:"simples"`
	Images     interface{}            `json:"images"`
}

func (memcache ProductMemcacheResponse) PrepareImages(p *Product) []map[string]string {
	response := make([]map[string]string, 0)
	var baseUrl string = "http://static.jabong.com"
	for _, im := range p.Images {
		resp := map[string]string{
			"image":             utils.ToString(im.ImageNo),
			"main":              utils.ToString(im.Main),
			"original_filename": im.OriginalFileName,
			"orientation":       im.Orientation,
			"path":              im.OriginalFileName,
			"name":              p.Name,
			"sku":               p.SKU,
			"url":               baseUrl + "/p/" + im.ImageName + utils.ToString(im.ImageNo),
			"sprite":            baseUrl + "/p/" + im.ImageName + "sprite.jpg",
		}
		response = append(response, resp)
	}
	return response
}

func (memcache ProductMemcacheResponse) PrepareSimples(p *Product) map[string]interface{} {

	response := make(map[string]interface{})
	for _, simple := range p.Simples {
		if simple.Status != "active" {
			continue
		}
		attribute, meta, _ := memcache.PrepareAttributes(simple.Global, simple.Attributes)
		attribute["barcode_ean"] = simple.BarcodeEan

		meta["sku"] = simple.SKU
		meta["seller_sku"] = simple.SellerSKU
		meta["ean_code"] = simple.EanCode
		meta["quantity"] = utils.ToString(simple.GetQuantity())
		meta["jabong_discount"] = utils.ToString(simple.JabongDiscount)
		meta["jabong_discount_from"] = ToMySqlTime(simple.JabongDiscountToDate)
		meta["jabong_discount_to"] = ToMySqlTime(simple.JabongDiscountFromDate)

		priceDetails := simple.GetPriceDetails()
		if priceDetails.Price != nil {
			meta["price"] = utils.ToString(*priceDetails.Price)
		}
		if priceDetails.SpecialPrice != nil {
			meta["special_price"] = utils.ToString(*priceDetails.SpecialPrice)
		}
		if priceDetails.MaxSavingPercentage > 0 {
			meta["max_saving_percentage"] = utils.ToString(priceDetails.MaxSavingPercentage)
		}
		resp := map[string]interface{}{
			"attributes":    attribute,
			"meta":          meta,
			"shipment_type": utils.ToString(p.ShipmentType),
		}
		response[simple.SKU] = resp
	}
	return response
}

func (memcache ProductMemcacheResponse) GetDispatchLocationId(
	global map[string]*Attribute,
) int {
	for k, v := range global {
		if k == "dispatchLocation" {
			val, err := v.GetValue("id")
			if err != nil {
				return 0
			}
			in, err := utils.GetInt(val)
			if err != nil {
				return 0
			}
			return in
		}
	}
	return 0
}

func (memcache ProductMemcacheResponse) PrepareGroupInfo(p *Product) string {
	if p.Group == nil {
		return ""
	}
	grp := p.Group
	result, err := GetAdapter(DB_READ_ADAPTER).GetProductsByGroupId(grp.Id)
	if err != nil {
		logger.Error(err)
		return ""
	}
	var skus []string
	for _, one := range result {
		skus = append(skus, one.SKU)
	}
	return strings.Join(skus, "|")
}

func (memcache ProductMemcacheResponse) PrepareAttributes(
	attributedata map[string]*Attribute,
	global map[string]*Attribute,
) (
	map[string]string,
	map[string]string,
	map[string]string,
) {

	attribute := make(map[string]string, 0)
	meta := make(map[string]string, 0)
	convey := make(map[string]string, 0)
	for _, data := range []map[string]*Attribute{attributedata, global} {
		for _, v := range data {
			val, err := v.GetValue("value")
			if err != nil {
				continue
			}
			switch v.AliceExport {
			case "meta":
				meta[v.Name] = utils.ToString(val)
			case "attribute":
				attribute[v.Name] = utils.ToString(val)
			case "no":
				convey[v.Name] = utils.ToString(val)
			}
		}
	}
	return attribute, meta, convey
}

func (cs *ProductResponse) PrepareMemcacheResponse() interface{} {
	response := ProductMemcacheResponse{}
	product := cs.Product

	if !cs.Pricegetter.DoNotPickPriceFrmMysql {
		for k, _ := range product.Simples {
			if product.Simples[k].Status != "active" {
				continue
			}
			setErr := product.Simples[k].SetPriceDetailsFrmMysql(cs.Pricegetter.UseMaster)
			if setErr != nil {
				//Remove special price
				product.Simples[k].SpecialPrice = nil
				//fire notification.
				notification.SendNotification(
					"Updating Price from Mysql Failed.",
					fmt.Sprintf(
						"Product:%s, Simple:%s, Message2:%s",
						product.SKU,
						product.Simples[k].SKU,
						setErr.Error(),
					),
					[]string{TAG_PRODUCT},
					datadog.ERROR,
				)
			}
		}
	}

	attribute, meta, _ := response.PrepareAttributes(
		product.Global, product.Attributes,
	)
	meta["sku"] = product.SKU

	meta["id_catalog_config"] = utils.ToString(product.SeqId)
	meta["attribute_set_id"] = utils.ToString(product.AttributeSet.Id)
	meta["attribute_set_label"] = product.AttributeSet.Label
	meta["name"] = product.Name
	meta["fk_catalog_attribute_option_global_dispatch_location"] = utils.ToString(response.GetDispatchLocationId(
		product.Global,
	))
	meta["status"] = product.Status
	meta["activated_at"] = ToMySqlTime(product.ActivatedAt)
	meta["fk_catalog_distinct_sizechart"] = utils.ToString(product.SizeChart.Id)
	meta["grouped_products"] = response.PrepareGroupInfo(product)

	//categories
	var strCats []string
	for _, cat := range product.Categories {
		strCats = append(strCats, utils.ToString(cat))
	}
	meta["categories"] = strings.Join(strCats, "|")
	brand, err := product.GetBrandInfo()
	if err != nil {
		cs.SetNoCache = true
	}
	meta["brand"] = brand.Name
	meta["ratings_total"] = "0"
	meta["ratings_single"] = ""
	meta["config_id"] = utils.ToString(product.SeqId)
	meta["shipment_type"] = utils.ToString(product.ShipmentType)

	//price
	priceMap := product.PreparePriceMap()
	meta["price"] = utils.ToString(priceMap.Price)
	meta["max_price"] = meta["price"]
	meta["max_original_price"] = meta["price"]
	meta["original_price"] = meta["price"]
	if priceMap.SpecialPrice != nil {
		meta["special_price"] = utils.ToString(*priceMap.SpecialPrice)
	}
	if priceMap.MaxSavingPercentage > 0 {
		meta["max_saving_percentage"] = utils.ToString(priceMap.MaxSavingPercentage)
	}

	attribute["description"] = product.Description
	if _, ok := attribute["short_description"]; !ok {
		attribute["short_description"] = attribute["description"]
	}
	response.Meta = make(map[string]interface{}, 0)
	response.Attributes = make(map[string]interface{}, 0)

	for k, v := range meta {
		response.Meta[k] = v
	}
	//add mandatory fields
	if _, ok := response.Meta["special_price"]; !ok {
		response.Meta["special_price"] = nil
	}
	if _, ok := response.Meta["max_saving_percentage"]; !ok {
		response.Meta["max_saving_percentage"] = nil
	}

	for k, v := range attribute {
		response.Attributes[k] = v
	}
	//prepare simples
	response.Simples = response.PrepareSimples(product)
	response.Images = response.PrepareImages(product)
	return response
}

//
// Expanse-Type: Catalog
//
type ProductCatalogResponse struct {
	Id        int         `json:"id"`
	Sku       string      `json:"sku"`
	Name      string      `json:"name"`
	UrlKey    string      `json:"urlKey"`
	Brand     interface{} `json:"brand"`
	Price     interface{} `json:"price"`
	Group     interface{} `json:"group"`
	Image     interface{} `json:"image"`
	Meta      interface{} `json:"meta"`
	Simples   interface{} `json:"simples"`
	CreatedAt string      `json:"createdAt"`
}

func (catalog ProductCatalogResponse) PrepareGroupInfo(p *Product, visibilityType string) interface{} {
	type GroupResponse struct {
		Id        int         `json:"id"`
		Sku       string      `json:"sku"`
		Name      string      `json:"name"`
		UrlKey    string      `json:"urlKey"`
		Color     interface{} `json:"color"`
		Brand     interface{} `json:"brand"`
		PriceMap  interface{} `json:"priceMap"`
		Simples   interface{} `json:"simples"`
		Image     interface{} `json:"image"`
		CreatedAt string      `json:"createdAt"`
	}
	if p.Group == nil {
		return nil
	}
	grp := p.Group
	result, err := GetAdapter(DB_READ_ADAPTER).GetProductsByGroupId(grp.Id)
	if err != nil {
		logger.Error(err)
		return nil
	}
	response := make([]interface{}, 0)
	for _, one := range result {
		if one.SKU != p.SKU {
			//check if product is visible
			vc := VisibilityChecker{Product: one, VisibilityType: visibilityType}
			if !vc.IsVisible() {
				continue
			}
		}
		presp := ProductResponse{
			Product:        &one,
			VisibilityType: visibilityType,
		}
		m := GroupResponse{}
		m.Id = presp.Product.SeqId
		m.Sku = presp.Product.SKU
		m.UrlKey = presp.Product.UrlKey
		if _, ok := presp.Product.Global["colorFamily"]; ok {
			m.Color, _ = presp.Product.Global["colorFamily"].GetValue("value")
		}
		m.Brand, _ = presp.Product.GetBrandInfo()
		m.PriceMap = presp.Product.PreparePriceMap()

		for _, im := range presp.Product.Images {
			if im.Main == 1 {
				m.Image = im
			}
		}
		m.Simples = presp.PrepareCatalogSimples()
		m.CreatedAt = ToMySqlTime(presp.Product.CreatedAt)
		response = append(response, m)
	}
	if len(response) == 0 {
		return nil
	}
	return response
}

//
// Expanse-Type: XLarge
//
type ProductXLargeResponse struct {
	Id              int         `json:"id"`
	Sku             string      `json:"sku"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Status          string      `json:"status"`
	UrlKey          string      `json:"urlKey"`
	Ty              interface{} `json:"ty"`
	SizeKey         string      `json:"sizeKey"`
	Supplier        interface{} `json:"supplier"`
	Shipment        interface{} `json:"shipment"`
	AttributeSet    interface{} `json:"attributeSet"`
	Brand           interface{} `json:"brand"`
	Price           interface{} `json:"price"`
	PrimaryCategory int         `json:"primaryCategory"`
	Categories      interface{} `json:"categories"`
	Leaf            interface{} `json:"leaf"`
	Group           interface{} `json:"group"`
	SizeChart       interface{} `json:"sizeChart"`
	Rating          interface{} `json:"rating"`
	Images          interface{} `json:"images"`
	Videos          interface{} `json:"videos"`
	Meta            interface{} `json:"meta"`
	Attributes      interface{} `json:"attributes"`
	Convey          interface{} `json:"convey"`
	Simples         interface{} `json:"simples"`
}

//
// Expanse-Type: Large
//
type ProductLargeResponse struct {
	Id              int         `json:"id"`
	Sku             string      `json:"sku"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Status          string      `json:"status"`
	UrlKey          string      `json:"urlKey"`
	Ty              interface{} `json:"ty"`
	SizeKey         string      `json:"sizeKey"`
	Supplier        interface{} `json:"supplier"`
	Shipment        interface{} `json:"shipment"`
	AttributeSet    interface{} `json:"attributeSet"`
	Brand           interface{} `json:"brand"`
	Price           interface{} `json:"price"`
	PrimaryCategory int         `json:"primaryCategory"`
	Categories      interface{} `json:"categories"`
	Leaf            interface{} `json:"leaf"`
	Group           interface{} `json:"group"`
	SizeChart       interface{} `json:"sizeChart"`
	Rating          interface{} `json:"rating"`
	Images          interface{} `json:"images"`
	Videos          interface{} `json:"videos"`
	Meta            interface{} `json:"meta"`
	Attributes      interface{} `json:"attributes"`
	Simples         interface{} `json:"simples"`
}

type ProductSolrResponse struct {
	Id                   int         `json:"id"`
	Sku                  string      `json:"sku"`
	Name                 string      `json:"name"`
	Description          string      `json:"description"`
	Status               string      `json:"status"`
	UrlKey               string      `json:"urlKey"`
	Ty                   interface{} `json:"ty"`
	SizeKey              string      `json:"sizeKey"`
	Shipment             interface{} `json:"shipmentType"`
	Supplier             interface{} `json:"supplier"`
	AttributeSet         interface{} `json:"attributeSet"`
	Brand                interface{} `json:"brand"`
	Price                interface{} `json:"price"`
	Categories           interface{} `json:"categories"`
	Leaf                 interface{} `json:"leaf"`
	Group                interface{} `json:"group"`
	CuratedList          interface{} `json:"curatedList"`
	CTRN                 *float64    `json:"ctrN"`
	BundleInformation    interface{} `json:"bundleInfo"`
	Campaign             *string     `json:"campaign"`
	ShopLook             interface{} `json:"shopLook"`
	InventoryAvailable   int         `json:"inventoryAvailable"`
	IsStockAvailable     string      `json:"isStockAvailable"`
	WeightedAvailability float64     `json:"weightedAvailability"`
	CutSize              string      `json:"cutSize"`
	Score                *Score      `json:"score"`
	SizeChart            interface{} `json:"sizeChart"`
	Rating               interface{} `json:"rating"`
	Meta                 interface{} `json:"meta"`
	Attributes           interface{} `json:"attributes"`
	Convey               interface{} `json:"convey"`
	Simples              interface{} `json:"simples"`
	CreatedAt            string      `json:"createdAt"`
	ActivatedAt          string      `json:"activatedAt"`
	UpdatedAt            string      `json:"updatedAt"`
	DisplayStockedOut    int         `json:"displayFlag"`
}

//
// Expanse-Type: Promotion
//
type ProductPromotionResponse struct {
	Id              int         `json:"id"`
	Sku             string      `json:"sku"`
	Name            string      `json:"name"`
	Status          string      `json:"status"`
	SizeKey         string      `json:"sizeKey"`
	Shipment        interface{} `json:"shipment"`
	AttributeSet    interface{} `json:"attributeSet"`
	Brand           interface{} `json:"brand"`
	Price           interface{} `json:"price"`
	PrimaryCategory int         `json:"primaryCategory"`
	Categories      interface{} `json:"categories"`
	Leaf            interface{} `json:"leaf"`
	Meta            interface{} `json:"meta"`
	Attributes      interface{} `json:"attributes"`
	Simples         interface{} `json:"simples"`
}

//
// Expanse-Type: Medium
//
type ProductMediumResponse struct {
	Id     int    `json:"id"`
	Sku    string `json:"sku"`
	Name   string `json:"name"`
	Status string `json:"status"`
	UrlKey string `json:"urlKey"`

	Supplier        interface{} `json:"supplier"`
	AttributeSet    interface{} `json:"attributeSet"`
	Brand           interface{} `json:"brand"`
	Price           interface{} `json:"price"`
	PrimaryCategory interface{} `json:"primaryCategory"`
	Categories      interface{} `json:"categories"`
	Leaf            interface{} `json:"leaf"`
	Images          interface{} `json:"images"`
	Videos          interface{} `json:"videos"`
	Global          interface{} `json:"global"`
	Simples         interface{} `json:"simples"`
}

//
// Expanse-Type: Small
//
type ProductSmallResponse struct {
	Id           int         `json:"id"`
	Sku          string      `json:"sku"`
	Name         string      `json:"name"`
	AttributeSet interface{} `json:"attributeSet"`
	Brand        interface{} `json:"brand"`
	Price        interface{} `json:"price"`
	Simples      interface{} `json:"simples"`
}

//
// Expanse-Type: XSmall
//
type ProductXSmallResponse struct {
	Id   int    `json:"id"`
	Sku  string `json:"sku"`
	Name string `json:"name"`
}

//
// This struct is used to prepare response of particular type
// based on the supplied expanse
//
type ProductResponse struct {
	Product        *Product
	VisibilityType string
	CacheTTL       int
	Pricegetter    PriceGetter
	SetNoCache     bool
}

//
// Main method which will be used to get response, based on expanse
//
func (cs *ProductResponse) GetResponse(expanse string) interface{} {
	switch expanse {
	case EXPANSE_CATALOG:
		return cs.PrepareCatalogResponse()
	case EXPANSE_XLARGE:
		return cs.PrepareXLargeResponse()
	case EXPANSE_LARGE:
		return cs.PrepareLargeResponse()
	case EXPANSE_XSMALL:
		return cs.PrepareXSmallResponse()
	case EXPANSE_SMALL:
		return cs.PrepareSmallResponse()
	case EXPANSE_MEDIUM:
		return cs.PrepareMediumResponse()
	case EXPANSE_SOLR:
		return cs.PrepareSolrResponse()
	case EXPANSE_MEMCACHE:
		return cs.PrepareMemcacheResponse()
	case EXPANSE_PROMOTION:
		return cs.PreparePromotionResponse()
	default:
		return cs.PrepareXLargeResponse()
	}
	return nil
}

//
// Preape response for Expanse XSmall
//
func (cs *ProductResponse) PrepareXSmallResponse() interface{} {
	response := ProductXSmallResponse{}
	response.Id = cs.Product.SeqId
	response.Sku = cs.Product.SKU
	response.Name = cs.Product.Name
	return response
}

//
// Preape response for Expanse Small
//
func (cs *ProductResponse) PrepareSmallResponse() interface{} {
	response := ProductSmallResponse{}
	response.Id = cs.Product.SeqId
	response.Sku = cs.Product.SKU
	response.Name = cs.Product.Name
	response.Simples = cs.PrepareSmallSimples()
	response.AttributeSet = cs.Product.AttributeSet
	response.Price = cs.Product.PreparePriceMap()
	var err error
	response.Brand, err = cs.Product.GetBrandInfo()
	if err != nil {
		cs.SetNoCache = true
	}
	return response
}

//
// Prepare response for Expanse Medium
//
func (cs *ProductResponse) PrepareMediumResponse() interface{} {
	response := ProductMediumResponse{}
	response.Id = cs.Product.SeqId
	response.Status = cs.Product.Status
	response.AttributeSet = cs.Product.AttributeSet
	var err error
	response.Brand, err = cs.Product.GetBrandInfo()
	if err != nil {
		cs.SetNoCache = true
	}
	response.Categories = cs.PrepareCategories()
	response.Global = cs.PrepareGlobalAttributes(cs.Product.Global)
	response.Images = cs.PrepareImages()
	response.PrimaryCategory = cs.Product.PrimaryCategory
	response.Leaf = cs.PrepareLeaf()
	response.Name = cs.Product.Name
	pMap := cs.Product.PreparePriceMap()
	response.Price = pMap
	cs.CacheTTL = cs.GetTTL(pMap)
	response.Sku = cs.Product.SKU
	response.Videos = cs.PrepareVideos()
	response.UrlKey = cs.Product.UrlKey
	seller, err := cs.Product.GetSellerInfoWithRetries(GetConfig().Org, true)
	if err != nil {
		logger.Error(err)
		cs.SetNoCache = true
	}
	seller.Id = cs.Product.SellerId
	response.Supplier = seller
	response.Simples = cs.PrepareMediumSimples(seller.UpdatedCommission)
	return response
}

func (cs *ProductResponse) PrepareSolrResponse() interface{} {
	response := ProductSolrResponse{}
	p := cs.Product

	if !cs.Pricegetter.DoNotPickPriceFrmMysql {
		for k, _ := range p.Simples {
			setErr := p.Simples[k].SetPriceDetailsFrmMysql(cs.Pricegetter.UseMaster)
			if setErr != nil {
				//Remove special price
				p.Simples[k].SpecialPrice = nil
				//fire notification.
				notification.SendNotification(
					"Updating Price from Mysql Failed.",
					fmt.Sprintf(
						"Product:%s, Simple:%s, Message2:%s",
						p.SKU,
						p.Simples[k].SKU,
						setErr.Error(),
					),
					[]string{TAG_PRODUCT},
					datadog.ERROR,
				)
			}
		}
	}

	response.Id = p.SeqId
	response.Sku = p.SKU
	response.Name = p.Name
	response.Description = p.Description
	response.Status = cs.Product.Status
	response.UrlKey = p.UrlKey
	response.Ty = cs.PrepareTy()
	response.Rating = cs.PrepareRating()
	response.SizeKey = p.GetSizeKey()
	var err error
	response.Brand, err = p.GetBrandInfo()
	if err != nil {
		cs.SetNoCache = true
	}
	response.Shipment = cs.Product.ShipmentType
	response.CuratedList = cs.PrepareCuratedList()
	response.BundleInformation = cs.GetBundleInformation()
	response.Campaign = cs.GetSpecialBucketCampaign()
	response.ShopLook = cs.GetLookDetail()
	response.DisplayStockedOut = p.DisplayStockedOut
	seller, err := p.GetSellerInfoWithRetries(GetConfig().Org, false)
	if err != nil {
		logger.Error(err)
		cs.SetNoCache = true
	}
	response.Supplier = seller
	response.AttributeSet = p.AttributeSet
	response.Price = p.PreparePriceMap()
	attribute, meta, convey := cs.PrepareAttributes(p.Attributes, p.Global)
	response.Attributes = attribute
	response.Meta = meta
	response.Convey = convey
	response.Group = cs.PrepareGroupInfo()
	response.Categories = cs.PrepareCategories()
	response.Leaf = cs.PrepareLeaf()
	response.SizeChart = cs.PrepareSizeChart()
	response.CreatedAt = ToMySqlTime(cs.Product.CreatedAt)
	response.ActivatedAt = ToMySqlTime(cs.Product.ActivatedAt)
	response.UpdatedAt = ToMySqlTime(cs.Product.UpdatedAt)
	response.Simples = cs.PrepareSimples(true, true)
	response.CTRN = cs.GetSolrCTR()
	stock := p.GetTotalQuantity()
	response.InventoryAvailable = stock
	response.WeightedAvailability, _ = p.GetWeightedAvailability()
	if stock > 0 {
		response.IsStockAvailable = "Y"
	} else {
		response.IsStockAvailable = "N"
	}
	var cutSize string
	var availableSimples int
	var simpleCount int
	for _, s := range p.Simples {
		if s.GetQuantity() > 0 {
			availableSimples += 1
		}
		simpleCount += 1
	}
	if simpleCount >= 4 && availableSimples == 1 {
		cutSize = "Y"
	} else {
		cutSize = "N"
	}
	response.CutSize = cutSize
	response.Score = p.CalculateScore()
	return response
}

//
// Prepare response for Expanse catalog
//
func (cs *ProductResponse) PrepareCatalogResponse() interface{} {
	response := ProductCatalogResponse{}
	p := cs.Product
	response.Id = p.SeqId
	response.Sku = p.SKU
	response.Name = p.Name
	response.UrlKey = p.UrlKey
	var err error
	response.Brand, err = p.GetBrandInfo()
	if err != nil {
		cs.SetNoCache = true
	}
	pMap := p.PreparePriceMap()
	response.Price = pMap
	//set cache TTL
	cs.CacheTTL = cs.GetTTL(pMap)

	for _, im := range cs.Product.Images {
		response.Image = im
	}
	_, meta, _ := cs.PrepareAttributes(p.Attributes, p.Global)
	response.Meta = meta
	response.Group = response.PrepareGroupInfo(p, cs.VisibilityType)
	response.Simples = cs.PrepareCatalogSimples()
	response.CreatedAt = ToMySqlTime(cs.Product.CreatedAt)
	return response
}

//
// Prepare response for Expanse Large
//
func (cs *ProductResponse) PrepareXLargeResponse() interface{} {
	response := ProductXLargeResponse{}
	p := cs.Product
	response.Id = p.SeqId
	response.Sku = p.SKU
	response.Name = p.Name
	response.Description = p.Description
	response.Status = cs.Product.Status
	response.UrlKey = p.UrlKey
	response.Ty = cs.PrepareTy()
	response.Rating = cs.PrepareRating()
	response.SizeKey = p.GetSizeKey()
	var err error
	response.Brand, err = p.GetBrandInfo()
	if err != nil {
		cs.SetNoCache = true
	}
	response.Shipment = cs.PrepareShipmentInfo()
	response.AttributeSet = p.AttributeSet
	pMap := p.PreparePriceMap()
	response.Price = pMap
	//set cache TTL
	cs.CacheTTL = cs.GetTTL(pMap)

	response.Images = cs.PrepareImages()
	response.Videos = cs.PrepareVideos()
	attribute, meta, convey := cs.PrepareAttributes(p.Attributes, p.Global)
	response.Attributes = attribute
	response.Meta = meta
	response.Convey = convey
	response.Group = cs.PrepareGroupInfo()
	response.Categories = cs.PrepareCategories()
	response.Leaf = cs.PrepareLeaf()
	response.PrimaryCategory = p.PrimaryCategory

	response.SizeChart = cs.PrepareSizeChart()
	seller, err := p.GetSellerInfoWithRetries(GetConfig().Org, false)
	if err != nil {
		logger.Error(err)
		cs.SetNoCache = true
	}
	response.Supplier = seller
	response.Simples = cs.PrepareSimples(false, true)
	return response
}

//
// Prepare response for Expanse Large
//
func (cs *ProductResponse) PrepareLargeResponse() interface{} {
	response := ProductLargeResponse{}
	p := cs.Product
	response.Id = p.SeqId
	response.Sku = p.SKU
	response.Name = p.Name
	response.Description = p.Description
	response.Status = cs.Product.Status
	response.UrlKey = p.UrlKey
	response.Ty = cs.PrepareTy()
	response.Rating = cs.PrepareRating()
	response.SizeKey = p.GetSizeKey()
	var err error
	response.Brand, err = p.GetBrandInfo()
	if err != nil {
		cs.SetNoCache = true
	}
	response.Shipment = cs.PrepareShipmentInfo()
	response.AttributeSet = p.AttributeSet
	pMap := p.PreparePriceMap()
	response.Price = pMap
	//set cache TTL
	cs.CacheTTL = cs.GetTTL(pMap)

	response.Images = cs.PrepareImages()
	response.Videos = cs.PrepareVideos()
	attribute, meta := response.PrepareAttributes(p.Attributes, p.Global)
	response.Attributes = attribute
	response.Meta = meta
	response.Group = cs.PrepareGroupInfo()
	response.Categories = cs.PrepareCategories()
	response.Leaf = cs.PrepareLeaf()
	response.PrimaryCategory = p.PrimaryCategory

	response.SizeChart = cs.PrepareSizeChart()
	seller, err := p.GetSellerInfoWithRetries(GetConfig().Org, false)
	if err != nil {
		logger.Error(err)
		cs.SetNoCache = true
	}
	response.Supplier = seller
	response.Simples = response.PrepareSimples(cs.Product, cs.VisibilityType)
	return response
}

//
// Prepare response for Expanse Promotion
//
func (cs *ProductResponse) PreparePromotionResponse() interface{} {
	response := ProductPromotionResponse{}
	p := cs.Product
	response.Id = p.SeqId
	response.Sku = p.SKU
	response.Name = p.Name
	response.Status = cs.Product.Status
	response.SizeKey = p.GetSizeKey()
	var err error
	response.Brand, err = p.GetBrandInfo()
	if err != nil {
		cs.SetNoCache = true
	}
	response.Shipment = cs.PrepareShipmentInfo()
	response.AttributeSet = p.AttributeSet
	pMap := p.PreparePriceMap()
	response.Price = pMap
	//set cache TTL
	cs.CacheTTL = cs.GetTTL(pMap)
	attribute, meta := response.PrepareAttributesForPromotion(p.Attributes, p.Global)
	response.Attributes = attribute
	response.Meta = meta
	response.Categories = cs.PrepareCategories()
	response.Leaf = cs.PrepareLeaf()
	response.PrimaryCategory = p.PrimaryCategory
	response.Simples = ProductLargeResponse{}.PrepareSimples(cs.Product, cs.VisibilityType)
	return response
}

//prepare data for global attributes
func (cs *ProductResponse) PrepareGlobalAttributes(
	global map[string]*Attribute,
) interface{} {
	type AttributeRes struct {
		Label string      `json:"label"`
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}
	response := []AttributeRes{}
	for _, v := range global {
		var attr AttributeRes
		attr.Label = v.Label
		attr.Name = v.Name
		attr.Value, _ = v.GetValue("value")
		response = append(response, attr)
	}
	return response
}

//Prepare data for Attributes and global attributes combined
func (cs *ProductResponse) PrepareAttributes(
	attributedata map[string]*Attribute,
	global map[string]*Attribute,
) (
	interface{},
	interface{},
	interface{},
) {
	type AttributeRes struct {
		Label           string      `json:"label"`
		Name            string      `json:"name"`
		OptionType      string      `json:"attributeType"`
		Value           interface{} `json:"value"`
		SolrSearchable  int         `json:"solrSearchable"`
		SolrFilter      int         `json:"solrFilter"`
		SolrSuggestions int         `json:"solrSuggestions"`
	}
	attribute := []AttributeRes{}
	meta := []AttributeRes{}
	convey := []AttributeRes{}
	for _, data := range []map[string]*Attribute{
		attributedata,
		global,
	} {
		attributei, metai, conveyi := func(data map[string]*Attribute) (
			[]AttributeRes,
			[]AttributeRes,
			[]AttributeRes,
		) {
			attribute := []AttributeRes{}
			meta := []AttributeRes{}
			convey := []AttributeRes{}
			for _, v := range data {
				a := AttributeRes{}
				a.Label = v.Label
				a.Name = v.Name
				val, err := v.GetValue("value")
				if err != nil {
					continue
				}
				a.Value = val
				a.SolrSearchable = v.SolrSearchable
				a.SolrSuggestions = v.SolrSuggestions
				a.SolrFilter = v.SolrFilter
				a.OptionType = v.OptionType

				switch v.AliceExport {
				case "meta":
					meta = append(meta, a)
				case "attribute":
					attribute = append(attribute, a)
				case "no":
					convey = append(convey, a)
				}
			}
			return attribute, meta, convey
		}(data)
		attribute = append(attribute, attributei...)
		meta = append(meta, metai...)
		convey = append(convey, conveyi...)
	}
	return attribute, meta, convey
}

//Prepare data for Attributes and global attributes combined
func (ppr ProductPromotionResponse) PrepareAttributesForPromotion(
	attributedata map[string]*Attribute,
	global map[string]*Attribute,
) (
	interface{},
	interface{},
) {
	type AttributeRes struct {
		Id    int         `json:"seqId"`
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}
	attribute := []AttributeRes{}
	meta := []AttributeRes{}
	for _, data := range []map[string]*Attribute{
		attributedata,
		global,
	} {
		attributei, metai := func(data map[string]*Attribute) (
			[]AttributeRes,
			[]AttributeRes,
		) {
			attribute := []AttributeRes{}
			meta := []AttributeRes{}
			for _, v := range data {
				a := AttributeRes{}
				a.Id = v.Id
				a.Name = v.Name
				val, err := v.GetValueForPromotion()
				if err != nil {
					continue
				}
				a.Value = val

				switch v.AliceExport {
				case "meta":
					meta = append(meta, a)
				case "attribute":
					attribute = append(attribute, a)
				}
			}
			return attribute, meta
		}(data)
		attribute = append(attribute, attributei...)
		meta = append(meta, metai...)
	}
	return attribute, meta
}

//Prepare data for Attributes and global attributes combined
func (plr ProductLargeResponse) PrepareAttributes(
	attributedata map[string]*Attribute,
	global map[string]*Attribute,
) (
	interface{},
	interface{},
) {
	type AttributeRes struct {
		Label      string      `json:"label"`
		Name       string      `json:"name"`
		OptionType string      `json:"attributeType"`
		Value      interface{} `json:"value"`
	}
	attribute := []AttributeRes{}
	meta := []AttributeRes{}
	for _, data := range []map[string]*Attribute{
		attributedata,
		global,
	} {
		attributei, metai := func(data map[string]*Attribute) (
			[]AttributeRes,
			[]AttributeRes,
		) {
			attribute := []AttributeRes{}
			meta := []AttributeRes{}
			for _, v := range data {
				a := AttributeRes{}
				a.Label = v.Label
				a.Name = v.Name
				val, err := v.GetValue("value")
				if err != nil {
					continue
				}
				a.Value = val
				a.OptionType = v.OptionType

				switch v.AliceExport {
				case "meta":
					meta = append(meta, a)
				case "attribute":
					attribute = append(attribute, a)
				}
			}
			return attribute, meta
		}(data)
		attribute = append(attribute, attributei...)
		meta = append(meta, metai...)
	}
	return attribute, meta
}

//Preape videos data
func (cs *ProductResponse) PrepareVideos() interface{} {
	p := cs.Product
	response := make([]interface{}, 0)
	for _, v := range p.Videos {
		m := M{}
		m["id"] = v.Id
		m["fileName"] = v.FileName
		m["thumbNail"] = v.Thumbnail
		response = append(response, m)
	}
	return response
}

//Prepare GroupInfo data
func (cs *ProductResponse) PrepareGroupInfo() interface{} {
	type GroupResponse struct {
		Id          int         `json:"id"`
		Sku         string      `json:"sku"`
		Status      string      `json:"status"`
		PetApproved int         `json:"petApproved"`
		UrlKey      string      `json:"urlKey"`
		Color       interface{} `json:"color"`
		Quantity    int         `json:"quantity"`
	}
	p := cs.Product
	visibilityType := cs.VisibilityType
	if p.Group == nil {
		return nil
	}
	grp := p.Group
	result, err := GetAdapter(DB_READ_ADAPTER).GetProductsByGroupId(grp.Id)
	if err != nil {
		logger.Error(err)
		return nil
	}
	response := make([]interface{}, 0)
	for _, one := range result {
		if one.SKU != p.SKU {
			//check if product is visible
			vc := VisibilityChecker{Product: one, VisibilityType: visibilityType}
			if !vc.IsVisible() {
				continue
			}
		}
		m := GroupResponse{}
		m.Id = one.SeqId
		m.Sku = one.SKU
		m.Status = one.Status
		m.PetApproved = one.PetApproved
		m.UrlKey = one.UrlKey
		if _, ok := one.Global["color"]; ok {
			m.Color, _ = one.Global["color"].GetValue("value")
		}
		m.Quantity = one.GetTotalQuantity()
		response = append(response, m)
	}
	if len(response) == 0 {
		return nil
	}
	return response
}

//Prepare Images data
func (cs *ProductResponse) PrepareImages() interface{} {
	p := cs.Product
	var response []map[string]interface{}
	portraitArr := []*ProductImage{}
	landArr := []*ProductImage{}
	for _, im := range p.Images {
		if im.Orientation == "portrait" {
			portraitArr = append(portraitArr, im)
		} else {
			landArr = append(landArr, im)
		}
	}
	if len(portraitArr) > 0 {
		m := make(map[string]interface{}, 0)
		m["orientation"] = "portrait"
		m["imageList"] = portraitArr
		response = append(response, m)
	}
	if len(landArr) > 0 {
		m := make(map[string]interface{}, 0)
		m["orientation"] = "landscape"
		m["imageList"] = landArr
		response = append(response, m)
	}
	return response
}

//prepare sizechart data
func (cs *ProductResponse) PrepareSizeChart() interface{} {
	p := cs.Product
	catData := categoryService.ById(p.PrimaryCategory)
	if catData.SizeChartAcive == 1 {
		return p.SizeChart.Data
	}
	return nil
}

//Prepare Categories data
func (cs *ProductResponse) PrepareCategories() interface{} {
	type CategoryResponse struct {
		Id            int    `json:"id"`
		Status        string `json:"status"`
		Name          string `json:"name"`
		UrlKey        string `json:"urlKey"`
		Segment       string `json:"segment"`
		SegmentUrlKey string `json:"segmentUrlKey"`
	}
	p := cs.Product
	catlist := categoryService.ByIds(p.Categories)
	response := make([]CategoryResponse, 0)
	if catlist == nil {
		logger.Error("Category returned NIL")
		cs.SetNoCache = true
		return response
	}
	// get category segment
	var index int
	if len(catlist) == 3 {
		index = 2
	} else {
		index = 1
	}
	sgmntName := ""
	sgmntUrlKey := ""
	if len(catlist[len(catlist)-index].CategorySeg) > 0 {
		sgmntName = catlist[len(catlist)-index].CategorySeg[0].Name
		sgmntUrlKey = catlist[len(catlist)-index].CategorySeg[0].UrlKey
	}
	for _, cat := range catlist {
		if cat.CategoryId == 1 {
			//skip root category from response.
			continue
		}
		data := CategoryResponse{}
		data.Id = cat.CategoryId
		data.Name = cat.Name
		data.Status = cat.Status
		data.UrlKey = cat.UrlKey
		data.Segment = sgmntName
		data.SegmentUrlKey = sgmntUrlKey
		response = append(response, data)
	}
	//sort categories
	newResponse := make([]CategoryResponse, 0)
	tree := FetchCategoryTree(cs.Product.PrimaryCategory)
	cc := len(tree)
	revTree := make([]int, cc)
	for k, v := range tree {
		revTree[cc-(k+1)] = v
	}
	for _, v1 := range revTree {
		for _, resp := range response {
			if v1 == resp.Id {
				newResponse = append(newResponse, resp)
			}
		}

	}
	for _, v1 := range response {
		var found bool
		for _, v2 := range newResponse {
			if v1.Id == v2.Id {
				found = true
			}
		}
		if !found {
			newResponse = append(newResponse, v1)
		}
	}

	return newResponse
}

//Prepare TY data
func (cs *ProductResponse) PrepareTy() interface{} {
	type TyResponse struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		UrlKey string `json:"urlKey"`
	}
	ty := cs.Product.TY
	if ty <= 0 {
		return nil
	}
	sql := `SELECT
			id_catalog_ty,
			name,
			url_key
		FROM catalog_ty
		WHERE id_catalog_ty = ` + strconv.Itoa(ty)
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return nil
	}
	result, sqlerr := driver.Query(sql)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return nil
	}
	defer result.Close()
	var tyresponse TyResponse
	for result.Next() {
		result.Scan(&tyresponse.Id, &tyresponse.Name, &tyresponse.UrlKey)
	}
	return tyresponse
}

//Prepare Rating data
func (cs *ProductResponse) PrepareRating() interface{} {
	type RatingResponse struct {
		Total  float64 `json:"total"`
		Single string  `json:"single"`
	}
	pId := cs.Product.SeqId

	sql := `SELECT
			total,
			single
		FROM rating_aggregated
		WHERE fk_catalog_config = ` + strconv.Itoa(pId)
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return nil
	}
	result, sqlerr := driver.Query(sql)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return nil
	}
	defer result.Close()
	var ratresponse RatingResponse
	var count int
	for result.Next() {
		result.Scan(&ratresponse.Total, &ratresponse.Single)
		count = count + 1
	}
	if count > 0 {
		return ratresponse
	}
	return nil
}

func (cs *ProductResponse) PrepareShipmentInfo() interface{} {
	type ShipmentResponse struct {
		Id   int    `json:"id"`
		Type string `json:"type"`
	}
	ship := GetShipmentById(cs.Product.ShipmentType)
	data := ShipmentResponse{}
	data.Id = cs.Product.ShipmentType
	data.Type = ship
	return data
}

//Preape leaf data
func (cs *ProductResponse) PrepareLeaf() interface{} {
	type LeafResponse struct {
		Id                    int    `json:"id"`
		Status                string `json:"status"`
		Lft                   int    `json:"lft"`
		Rgt                   int    `json:"rgt"`
		Name                  string `json:"name"`
		UrlKey                string `json:"urlKey"`
		SizeChartActive       int    `json:"sizeChartActive"`
		PdfActive             int    `json:"pdfActive"`
		DisplaySizeConversion string `json:"displaySizeConversion"`
		GoogleTreeMapping     string `json:"googleTreeMapping"`
		SizeChartApplicable   int    `json:"sizeChartApplicable"`
	}
	p := cs.Product
	catlist := categoryService.ByIds(p.Leaf)
	response := make([]interface{}, 0)
	for _, cat := range catlist {
		data := LeafResponse{}
		data.Id = cat.CategoryId
		data.Status = cat.Status
		data.Lft = cat.Left
		data.Rgt = cat.Right
		data.Name = cat.Name
		data.UrlKey = cat.UrlKey
		data.SizeChartActive = cat.SizeChartAcive
		data.PdfActive = cat.PdfActive
		data.DisplaySizeConversion = cat.DisplaySizeConversion
		data.GoogleTreeMapping = cat.GoogleTreeMapping
		data.SizeChartApplicable = cat.SizeChartApplicable
		response = append(response, data)
	}
	return response
}

//Prepare simple data for XLarge, Solr response
func (cs *ProductResponse) PrepareSimples(isSolr bool,
	includeConvey bool) interface{} {

	type SimpleResponse struct {
		Id                  int         `json:"id"`
		Sku                 string      `json:"sku"`
		BarcodeEan          string      `json:"barcodeEan"`
		Meta                interface{} `json:"meta"`
		Convey              interface{} `json:"convey"`
		Attribute           interface{} `json:"attribute"`
		Price               *float64    `json:"price"`
		Stock               *int        `json:"stock,omitempty"`
		Size                string      `json:"size"`
		Position            string      `json:"position"`
		SpecialPrice        *float64    `json:"specialPrice"`
		SpecialFromDate     *string     `json:"specialFromDate"`
		SpecialToDate       *string     `json:"specialToDate"`
		DiscountedPrice     float64     `json:"discountedPrice"`
		MaxSavingPercentage float64     `json:"maxSavingPercentage"`
	}
	p := cs.Product

	// Sort Simples based on position
	p.SortSimples()

	response := make([]interface{}, 0)
	varattrName, _ := p.AttributeSet.GetVariationAttributeName()
	for _, simple := range p.Simples {
		if cs.VisibilityType != VISIBILITY_NONE && simple.Status != STATUS_ACTIVE {
			continue
		}
		m := SimpleResponse{}
		m.Id = simple.Id
		m.Sku = simple.SKU
		m.BarcodeEan = simple.SellerSKU
		m.Size = simple.GetSize(varattrName)
		m.Position = strconv.Itoa(simple.Position)
		if isSolr {
			var stock = simple.GetQuantity()
			m.Stock = &stock
		}
		attribute, meta, convey := cs.PrepareAttributes(
			simple.Attributes, simple.Global,
		)
		m.Meta = meta
		if includeConvey {
			m.Convey = convey
		}
		m.Attribute = attribute

		pd := simple.GetPriceDetails()
		m.Price = pd.Price
		m.SpecialPrice = pd.SpecialPrice
		m.SpecialFromDate = pd.SpecialFromDate
		m.SpecialToDate = pd.SpecialToDate
		m.DiscountedPrice = pd.DiscountedPrice
		m.MaxSavingPercentage = pd.MaxSavingPercentage

		response = append(response, m)
	}
	if len(response) == 0 {
		response = nil
	}
	return response
}

//Prepare simple data for large response
func (plr ProductLargeResponse) PrepareSimples(pro *Product,
	visibilityType string) interface{} {

	type SimpleResponse struct {
		Id                  int         `json:"id"`
		Sku                 string      `json:"sku"`
		BarcodeEan          string      `json:"barcodeEan"`
		Meta                interface{} `json:"meta"`
		Attribute           interface{} `json:"attribute"`
		Price               *float64    `json:"price"`
		Size                string      `json:"size"`
		Position            string      `json:"position"`
		SpecialPrice        *float64    `json:"specialPrice"`
		SpecialFromDate     *string     `json:"specialFromDate"`
		SpecialToDate       *string     `json:"specialToDate"`
		DiscountedPrice     float64     `json:"discountedPrice"`
		MaxSavingPercentage float64     `json:"maxSavingPercentage"`
	}
	p := pro

	// Sort Simples based on position
	p.SortSimples()

	response := make([]interface{}, 0)
	varattrName, _ := p.AttributeSet.GetVariationAttributeName()
	for _, simple := range p.Simples {
		if visibilityType != VISIBILITY_NONE && simple.Status != STATUS_ACTIVE {
			continue
		}
		m := SimpleResponse{}
		m.Id = simple.Id
		m.Sku = simple.SKU
		m.BarcodeEan = simple.SellerSKU
		m.Size = simple.GetSize(varattrName)
		m.Position = strconv.Itoa(simple.Position)

		attribute, meta := plr.PrepareAttributes(
			simple.Attributes, simple.Global,
		)
		m.Meta = meta
		m.Attribute = attribute

		pd := simple.GetPriceDetails()
		m.Price = pd.Price
		m.SpecialPrice = pd.SpecialPrice
		m.SpecialFromDate = pd.SpecialFromDate
		m.SpecialToDate = pd.SpecialToDate
		m.DiscountedPrice = pd.DiscountedPrice
		m.MaxSavingPercentage = pd.MaxSavingPercentage

		response = append(response, m)
	}
	if len(response) == 0 {
		response = nil
	}
	return response
}

func (cs *ProductResponse) PrepareCatalogSimples() interface{} {
	type SimpleResponse struct {
		Id   int    `json:"id"`
		Sku  string `json:"sku"`
		Size string `json:"size"`
	}
	p := cs.Product

	// Sort simples based on Position
	cs.Product.SortSimples()

	response := make([]interface{}, 0)
	varattrName, _ := p.AttributeSet.GetVariationAttributeName()
	for _, simple := range p.Simples {
		if cs.VisibilityType != VISIBILITY_NONE && simple.Status != STATUS_ACTIVE {
			continue
		}
		m := SimpleResponse{}
		m.Id = simple.Id
		m.Sku = simple.SKU
		m.Size = simple.GetSize(varattrName)
		response = append(response, m)
	}
	if len(response) == 0 {
		response = nil
	}
	return response
}

//Prepare simple data for medium expanse
func (cs *ProductResponse) PrepareMediumSimples(cmsnArr []SellerCommission) interface{} {

	type SimpleResponse struct {
		Id                 int         `json:"id"`
		Sku                string      `json:"sku"`
		BarcodeEan         string      `json:"barcodeEan"`
		SellerSku          string      `json:"sellerSku"`
		Size               string      `json:"size"`
		Global             interface{} `json:"global"`
		FixedCommission    float64     `json:"fixedCommission"`
		VariableCommission float64     `json:"variableCommission"`
	}
	p := cs.Product

	// Sort simples based on Position
	p.SortSimples()

	response := make([]interface{}, 0)
	varattrName, _ := p.AttributeSet.GetVariationAttributeName()
	for _, simple := range p.Simples {
		if cs.VisibilityType != VISIBILITY_NONE && simple.Status != STATUS_ACTIVE {
			continue
		}
		m := SimpleResponse{}
		m.Id = simple.Id
		m.Sku = simple.SKU
		m.BarcodeEan = simple.BarcodeEan
		m.SellerSku = simple.SellerSKU
		m.Size = simple.GetSize(varattrName)
		m.Global = cs.PrepareGlobalAttributes(simple.Global)
		if len(cmsnArr) != 0 && cmsnArr[0].CategoryId != 0 {
			m.FixedCommission = cmsnArr[0].Percentage
		}
		if simple.JabongDiscountFromDate != nil && simple.JabongDiscountToDate != nil {
			if time.Now().After(*simple.JabongDiscountFromDate) && time.Now().Before(*simple.JabongDiscountToDate) {
				m.VariableCommission = simple.JabongDiscount
			}
		}
		response = append(response, m)
	}
	if len(response) == 0 {
		response = nil
	}
	return response
}

//Prepare simple data for small expanse
func (cs *ProductResponse) PrepareSmallSimples() interface{} {
	// Sort simples based on Position
	cs.Product.SortSimples()

	var response []string
	for _, simple := range cs.Product.Simples {
		if cs.VisibilityType != VISIBILITY_NONE && simple.Status != STATUS_ACTIVE {
			continue
		}
		response = append(response, simple.SKU)
	}
	return response
}

func (cs *ProductResponse) PrepareCuratedList() interface{} {
	type Curated struct {
		Id   *int    `json:"id"`
		Name *string `json:"name"`
	}
	sql := `SELECT
			  curated_list.id_curated_list,
			  catalog_attribute_option_global_curated_list.name
			FROM curated_list
			INNER JOIN catalog_attribute_option_global_curated_list
			  ON catalog_attribute_option_global_curated_list.
			  id_catalog_attribute_option_global_curated_list = curated_list.
			  fk_catalog_attribute_option_global_curated_list
			WHERE fk_catalog_config = ` + strconv.Itoa(cs.Product.SeqId)
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return nil
	}
	result, sqlerr := driver.Query(sql)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return nil
	}
	defer result.Close()
	var curatedList []Curated
	for result.Next() {
		var tmp Curated
		err := result.Scan(&tmp.Id, &tmp.Name)
		if err != nil {
			logger.Error(err)
			continue
		}
		curatedList = append(curatedList, tmp)
	}
	result.Close()
	return curatedList
}

func (cs *ProductResponse) GetSolrCTR() *float64 {
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return nil
	}
	sql := `SELECT ctr_n FROM product_solr_ctr WHERE sku = '` + cs.Product.SKU + `'`
	result, sqlerr := driver.Query(sql)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return nil
	}
	defer result.Close()
	var ctrN *float64
	for result.Next() {
		result.Scan(&ctrN)
	}
	result.Close()
	return ctrN
}

func (cs *ProductResponse) GetTTL(pm PriceMap) int {
	return 0
	/*
		var defaultCacheTTL int = 3600 //1 hour in seconds
		var maxTTL int = 0
		var (
			intVal    int
			laterDate *time.Time
			fromDate  *time.Time
			toDate    *time.Time
			err       error
		)

		if pm.SpecialPrice == nil {
			return maxTTL
		}

		if pm.SpecialPriceFrom != nil {
			fromDate, err = FromMysqlTime(*pm.SpecialPriceFrom, true)
			if err != nil {
				logger.Error(fmt.Sprintf("(cs *ProductResponse)#GetTTL(pm PriceMap)1:%s",
					err.Error()))
			}
		}
		if pm.SpecialPriceTo != nil {
			toDate, err = FromMysqlTime(*pm.SpecialPriceTo, true)
			if err != nil {
				logger.Error(fmt.Sprintf("(cs *ProductResponse)#GetTTL(pm PriceMap)2:%s",
					err.Error()))
			}
		}

		if fromDate != nil && time.Now().Before(*fromDate) {
			laterDate = fromDate
		} else if toDate != nil && time.Now().Before(*toDate) {
			laterDate = toDate
		} else {
			return maxTTL
		}

		timeSpan := laterDate.Sub(time.Now())
		intVal, err = utils.GetInt(timeSpan.Seconds())
		if err != nil {
			logger.Error(fmt.Sprintf("(cs *ProductResponse)#GetTTL(pm PriceMap)3:%s", err.Error()))
			return defaultCacheTTL
		}
			// if intVal > defaultCacheTTL {
			// 	intVal = defaultCacheTTL
			// }
		return intVal
	*/
}

// get activated bundle id's against product
func (cs *ProductResponse) GetBundleInformation() []string {
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"(cs *ProductResponse)#GetBundleInformation() : Cannot get mysql driver: %s",
			err.Error(),
		))
		return nil
	}
	sql := `SELECT
			  catalog_attribute_option_global_bundle.name
			FROM catalog_attribute_option_global_bundle
			INNER JOIN catalog_attribute_link_global_bundle
			  ON id_catalog_attribute_option_global_bundle = catalog_attribute_link_global_bundle.fk_catalog_attribute_option_global_bundle
			INNER JOIN sku_bundle
			  ON sku_bundle.name = catalog_attribute_option_global_bundle.name
			WHERE catalog_attribute_link_global_bundle.fk_catalog_config = ?
			AND is_active = 1
			AND from_date <= ?
			AND to_date >= ?`
	result, sqlerr := driver.Query(sql,
		cs.Product.SeqId,
		time.Now().Format(FORMAT_MYSQL_TIME),
		time.Now().Format(FORMAT_MYSQL_TIME),
	)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return nil
	}
	defer result.Close()
	var bundles []string
	for result.Next() {
		var bundle string
		err = result.Scan(&bundle)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in scanning row %s", err.Error()))
			continue
		}
		bundles = append(bundles, bundle)
	}
	result.Close()
	return bundles
}

// get special bucket campaign
func (cs *ProductResponse) GetSpecialBucketCampaign() *string {
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"(cs *ProductResponse)#GetSpecialBucketCampaign() : Cannot get mysql driver",
			err.Error(),
		))
		return nil
	}
	sql := `SELECT
				GROUP_CONCAT(CONCAT_WS("_",NULLIF(CAST(campaign_code AS CHAR),""),
				NULLIF(CAST(special_bucket_url_parameter AS CHAR),"")) SEPARATOR "##|##")
			FROM banner_campaign
			INNER JOIN campaign_special_bucket_products
				ON id_banner_campaign=fk_banner_campaign
			INNER JOIN campaign_type
				ON fk_campaign_type=id_campaign_type
			WHERE (campaign_special_bucket_products.sku = ?)
			AND (campaign_special_bucket_products.is_primary = (1))
			AND (banner_campaign.is_active = 1 )
			AND (campaign_type.status = "active")
			AND (campaign_type.id_campaign_type <> 1 )
			AND ((? BETWEEN banner_campaign.from_date AND banner_campaign.to_date)
			OR (banner_campaign.from_date IS NULL AND banner_campaign.to_date IS NULL )
			OR (banner_campaign.from_date >= ? AND banner_campaign.from_date <= ?))`
	result, sqlerr := driver.Query(sql,
		cs.Product.SKU,
		time.Now().Format(FORMAT_MYSQL_TIME),
		time.Now().Format(FORMAT_MYSQL_TIME),
		time.Now().Add(time.Hour*24).Format(FORMAT_MYSQL_TIME),
	)
	if sqlerr != nil {
		logger.Error("(cs *ProductResponse)#GetSpecialBucketCampaign()" + sqlerr.DeveloperMessage)
		return nil
	}
	defer result.Close()
	var campaign *string
	for result.Next() {
		err = result.Scan(&campaign)
		if err != nil {
			logger.Error(fmt.Sprintf(
				"(cs *ProductResponse)#GetSpecialBucketCampaign(): Error in scanning row %s",
				err.Error(),
			))
		}
	}
	return campaign
}

func (cs *ProductResponse) GetLookDetail() interface{} {
	resp, err := GetLookDetail(cs.Product)
	if err != nil {
		logger.Error(
			fmt.Sprintf("(cs *ProductResponse)#GetLookDetail() %s:", err.Error()))
	}
	return resp
}
