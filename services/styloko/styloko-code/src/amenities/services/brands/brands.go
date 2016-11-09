package brands

import (
	"amenities/brands/common"
	"amenities/brands/get"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
)

func ById(brandId int) (common.Brand, error) {
	b := get.GetBrand{}
	brand, err := b.One(brandId)
	if err != nil {
		logger.Error(fmt.Sprintf("%s, %v", err.Error(), brandId))
		return brand, err
	}
	return brand, nil
}

func GetByName(name string) common.Brand {
	b := get.GetBrand{}
	brand, err := b.OneByName(name)
	if err != nil {
		logger.Error(fmt.Sprintf("%s, %s", err.Error(), name))
		return common.Brand{}
	}
	return brand
}
