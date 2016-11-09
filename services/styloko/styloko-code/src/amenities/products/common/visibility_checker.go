package common

import (
	brandService "amenities/services/brands"
	categoryService "amenities/services/categories"

	"github.com/jabong/floRest/src/common/utils/logger"
)

//
// Checks the visibility for the supplied product
// based on the supplied visibility condition.
//
type VisibilityChecker struct {
	Product        Product
	VisibilityType string
}

//
// Check is product satisfies visibility condition.
//
func (vc VisibilityChecker) IsVisible() bool {
	switch vc.VisibilityType {
	case VISIBILITY_PDP:
		return vc.PDP()
	case VISIBILITY_MSKU:
		return vc.MSKU()
	case VISIBILITY_DOOS:
		return vc.DOOS()
	case VISIBILITY_NONE:
		return true
	}
	return true
}

//
// Visibility condition for PDP
//
func (vc VisibilityChecker) PDP() bool {
	if (vc.isStatusActive() || vc.isStatusInActive()) &&
		vc.isSupplierActive() && vc.isPetApproved() &&
		vc.hasActiveImage() && vc.hasActiveBrand() &&
		vc.hasActiveCategory() {
		return true
	}
	return false
}

//
// Visibility condition for MSKU
//
func (vc VisibilityChecker) MSKU() bool {
	if vc.isStatusActive() &&
		vc.isSupplierActive() &&
		vc.isPetApproved() &&
		vc.hasActiveImage() &&
		vc.hasActiveBrand() &&
		vc.hasActiveCategory() &&
		vc.hasActiveSimple() {
		return true
	}
	return false
}

//
// Visibility condition for DOOS
//
func (vc VisibilityChecker) DOOS() bool {

	if vc.Product.ActivatedAt == nil {
		return false
	}
	var blockCatalog = "N"
	if bc, ok := vc.Product.Global["blockCatalog"]; ok {
		v := bc.Value
		blockCatalog, _ = v.(string)
	}
	var notBuyable = 0
	if nb, ok := vc.Product.Global["notBuyable"]; ok {
		v := nb.Value
		notBuyable, _ = v.(int)
	}
	// var doos = false
	// if vc.Product.DisplayStockedOut == 1 {
	// 	doos = true
	// }

	if vc.isStatusActive() &&
		vc.isSupplierActive() &&
		vc.isPetApproved() &&
		vc.hasActiveImage() &&
		vc.hasActiveBrand() &&
		vc.hasActiveCategory() &&
		vc.hasActiveSimple() &&
		blockCatalog == "N" &&
		notBuyable == 0 {
		return true
	}
	return false
}

//Checks if product status is active
func (vc VisibilityChecker) isStatusActive() bool {
	if vc.Product.Status == STATUS_ACTIVE {
		return true
	}
	return false
}

//Checks if product status is inactive
func (vc VisibilityChecker) isStatusInActive() bool {
	if vc.Product.Status == STATUS_INACTIVE {
		return true
	}
	return false
}

//Checks if petApproved set to 1
func (vc VisibilityChecker) isPetApproved() bool {
	if vc.Product.PetApproved == 1 {
		return true
	}
	return false
}

//Checks if selller is Active
func (vc VisibilityChecker) isSupplierActive() bool {
	conf := GetConfig()
	seller, err := vc.Product.GetSellerInfoWithRetries(conf.Org, false)
	if err != nil {
		logger.Error(err)
		return false
	}
	if seller.Status == STATUS_ACTIVE {
		return true
	}
	return false
}

//Checks if product has atleast one image
func (vc VisibilityChecker) hasActiveImage() bool {
	images := vc.Product.Images
	if len(images) > 0 {
		return true
	}
	return false
}

//Checks if product has atleast one active category
func (vc VisibilityChecker) hasActiveCategory() bool {
	cats := vc.Product.Categories
	verboseCats := categoryService.ByIds(cats)
	for _, cat := range verboseCats {
		if cat.Status == STATUS_ACTIVE {
			return true
		}
	}
	return false
}

//Checks if product brand is active
func (vc VisibilityChecker) hasActiveBrand() bool {
	brandId := vc.Product.BrandId
	brand, err := brandService.ById(brandId)
	if err != nil || brand.Status == STATUS_ACTIVE {
		return true
	}
	return false
}

//Checks if product has atleast one active simple
func (vc VisibilityChecker) hasActiveSimple() bool {
	simples := vc.Product.Simples
	pckQty := vc.Product.GetPackQty()
	for _, simple := range simples {
		if (simple.Status == STATUS_ACTIVE) && (*simple.Price > 0) {
			if simple.GetQuantityPck(pckQty) > 0 {
				return true
			}
			continue
		}
	}
	return false
}
