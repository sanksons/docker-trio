package products

import (
	"amenities/products/common"
	updaters "amenities/products/put/updaters"
	"errors"
	"strconv"
	"strings"
)

// UpdateCacheForCategories -> Invalidates cache.
func PurgeCacheByCategories(categoryIds []int) bool {
	if len(categoryIds) == 0 {
		return true
	}
	for _, id := range categoryIds {
		c := updaters.CacheInvalidate{
			Id:   strconv.Itoa(id),
			Type: updaters.CACHE_INV_CATEGORY,
		}
		c.InvalidateCache()
		c.Publish()
	}
	return true
}

func PurgeCacheByBrands(brandIds []int) bool {
	if len(brandIds) == 0 {
		return true
	}
	for _, id := range brandIds {
		c := updaters.CacheInvalidate{
			Id:   strconv.Itoa(id),
			Type: updaters.CACHE_INV_BRAND,
		}
		c.InvalidateCache()
		c.Publish()
	}
	return true
}

func AddNode(name string, sku string, data interface{}) error {
	n := updaters.Node{
		NodeName: name,
		SKU:      sku,
		NodeData: data,
		Type:     common.ADD_NODE,
	}
	err := UpdateProduct(n)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNode(name string, sku string) error {
	n := updaters.Node{
		NodeName: name,
		SKU:      sku,
		NodeData: nil,
		Type:     common.DELETE_NODE,
	}
	err := UpdateProduct(n)
	if err != nil {
		return err
	}
	return nil
}

func BySku(sku string, adapter string) (*common.Product, error) {
	p := common.Product{}
	err := p.LoadBySku(sku, adapter)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdateProduct(node updaters.Node) error {
	errs := node.Validate()
	if errs != nil {
		return errors.New(strings.Join(errs, ";"))
	}
	_, err := node.Update()
	if err != nil {
		return err
	}
	err = node.InvalidateCache()
	if err != nil {
		return err
	}
	return nil
}

func GetSimpleSizesForSKU(sku string, adapter string) ([]string, error) {
	simpleSizesArray := []string{}
	prd, err := BySku(sku, adapter)
	prd.SortSimples()
	if err != nil {
		return simpleSizesArray, err
	}
	prodSimpleArray := prd.Simples
	sizeAttributeName, err := prd.AttributeSet.GetVariationAttributeName()

	if err != nil {
		return simpleSizesArray, err
	}
	for _, simple := range prodSimpleArray {
		simpleSize := simple.GetSize(sizeAttributeName)
		simpleSizesArray = append(simpleSizesArray, simpleSize)
	}
	return simpleSizesArray, nil
}

func GetInactiveDeletedSimpleSizes(sku string, adapter string) ([]string, error) {
	simpleSizesArray := []string{}
	prd, err := BySku(sku, adapter)
	if err != nil {
		return simpleSizesArray, err
	}
	prd.SortSimples()

	prodSimpleArray := prd.Simples
	sizeAttributeName, err := prd.AttributeSet.GetVariationAttributeName()

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
