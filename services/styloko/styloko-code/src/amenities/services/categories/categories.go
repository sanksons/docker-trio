package categories

import (
	"amenities/categories/common"
	"amenities/categories/get"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strconv"
)

func ByIds(categoryIds []int) []common.CategoryGetVerbose {
	c := get.CategoriesGetAll{}
	categories, e := c.FindByIds(categoryIds)
	if e != nil {
		logger.Error(fmt.Sprintf("%s, %v", e.Error(), categoryIds))
		return nil
	}
	return categories
}

func ById(categoryId int) common.CategoryGetVerbose {
	c := get.CategoriesGetID{}
	category, _, e := c.FindById(strconv.Itoa(categoryId))
	if e != nil {
		logger.Error(fmt.Sprintf("%s, %d", e.Error(), categoryId))
	}
	return category
}
