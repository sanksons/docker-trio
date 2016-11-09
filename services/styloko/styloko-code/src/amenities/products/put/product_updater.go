//
// This file contains the "ProductUpdater" interface to be implemented via all the
// form of product updates.
// Currently Supported Product update types are (Update-Type):
// A:) Price
// B:) Shipment
// C:) SellerDeactivate
// D:) Node
// E:) Product
// F:) ImageAdd
// G:) ImageDel

package put

import (
	proUtil "amenities/products/common"
)

//
// Any request to Product update must implement this interface.
//
type ProductUpdater interface {

	//Defines Validation to be applied, returns list of errors
	Validate() []string

	//Contains logic for updating, returns the updated product
	Update() (proUtil.Product, error)

	//Invalidates the product cache based on the updation done.
	InvalidateCache() error

	//Publish product data to Jabong Bus.
	Publish() error

	//Prepares response to be shown to the User.
	Response(*proUtil.Product) interface{}

	//Acquire Lock
	Lock() bool

	//Release Lock
	UnLock() bool
}
