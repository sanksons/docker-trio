package search

import (
	proUtil "amenities/products/common"
	"common/appconfig"
	"encoding/json"
	_ "fmt"
)

var DbAdapterName string

var SellerSkuLimit int

var Pc proUtil.CacheManager

var Conf *appconfig.AppConfig

type M map[string]interface{}

type ProductData struct {
	Data       interface{}
	Cache      bool
	Visibility bool
	Identifier string
}

//
// This Query type will be used to load response for
// Exact type queries:
// -- get multi products by ID
// -- get multi products by SKU
// -- get single product by SKU
//
type ExactQuery struct {
	Id         []int    // List of product Ids
	Sku        []string // List of skus
	Expanse    string   // Expanse to be used
	Visibility string   // Visibility to be used
	IsSingle   bool     // Is single sku call
}

func (query ExactQuery) ToString() string {
	bytes, err := json.Marshal(query)
	if err != nil {
		return ""
	}
	return string(bytes)
}

type FilterQuery struct {
	Limit      int
	Offset     int
	Expanse    string
	Visibility string
}
