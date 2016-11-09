package common

import (
	"common/appconfig"

	"github.com/jabong/floRest/src/common/config"
)

var mongoAdapter *MongoAdapter
var mysqlAdapter *MySqlAdapter

type ProductDBAdapter interface {
	//Load products via skus
	GetBySkus([]string, interface{}) error
	//Load Products via Ids
	GetByIds([]int, interface{}) error
	//Get product via Id
	GetById(int) (Product, error)
	//Get product via sku
	GetBySku(string) (Product, error)
	//Get Product via ProductSetId
	GetByProductSet(int) (Product, error)
	//Get Product By SimpleId
	GetProductBySimpleId(int) (Product, error)
	//Get Product via VideoId
	GetProductByVideoId(int) (Product, error)
	//Get Product via SKU
	GetProductBySkuAndType(string, string) (Product, error)

	//Get Products By GroupId
	GetProductsByGroupId(int) ([]Product, error)
	//Get Products by sellerId
	GetProductsBySellerId(int) ([]Product, error)
	//Get Products via BrandId
	GetProductsByBrandId(int) ([]Product, error)
	//Get Products via CategoryId
	GetProductsByCategoryId(int) ([]Product, error)
	//Get Products For Seller Center API
	GetProductsForSeller(sellers []int, limit int, offset int, lastSCId int) ([]Product, error)

	//Get SmallProduct via SellerId
	GetProductIdsBySellerId(int) ([]ProductSmall, error)
	//Get SmallProduct via BrandId
	GetProductIdsByBrandId(int) ([]ProductSmall, error)
	//Get SmallProduct via categoryId
	GetProductIdsByCategoryId(int) ([]ProductSmall, error)
	//Get Product Id by SimpleSku
	GetProductIdBySimpleSku(string) (ProductSmall, error)
	//Get Product Id by SimpleId
	GetProductIdBySimpleId(int) (ProductSmall, error)

	//Get GroupInfo by Name
	GetProductGroupByName(string) (ProductGroup, error)
	//Get AttributeMongo via search condition
	GetAtrributeByCriteria(AttrSearchCondition) (AttributeMongo, error)
	//Get AttributeMongo via Id
	GetAttributeMongoById(int) (AttributeMongo, error)

	GetAttributeMongoByName(string) (AttributeMongo, error)

	GetAttributeMapping(string) (AttrMapping, error)

	//Get AtrributeSet by ID
	GetProAttributeSetById(int) (ProAttributeSet, error)
	//Get category details via ids
	GetCategoriesByIds([]int) ([]Category, error)
	//Find Primary category from categoryIds
	FindPrimaryCategoryId(catIds []int) (int, error)

	//Insert or Update Full Product data
	// -> inserts if seqId not exists in db
	// -> updates if seqId already exists in db
	SaveProduct(Product) error

	// Generate next Sequence for given Type.
	// Used for insertionn of new entities.
	// Params:
	// -> Entity Type
	// Return
	// -> counter
	// -> error
	GenerateNextSequence(string) (int, error)

	//Add node to Product
	//Params:
	//1. productSKU
	//2. nodename
	//3. data
	AddNode(string, string, interface{}) error

	//Delete Node
	//1. productSKU
	//2. nodeName
	DeleteNode(string, string) error

	//Delete an Image based on the supplied ImageId
	//Return:
	// -> productId to which image belongs
	DeleteImage(int) (int, error)

	//Add image to the specified Product
	//Params:
	//1. ProductId to which image belongs
	//2. Image Data
	//Return
	//1. ImageId
	AddImage(int, ProductImage) (int, error)

	//Save video details(insert/update).
	//Params:
	//1. ConfigId
	//2. Video data
	//Return:
	// --> Id of the inserted video
	SaveVideo(int, ProductVideo) (int, error)

	//Update video status
	//Params:
	//1. Video Id
	//2. Status
	//Return:
	// --> Error in updating
	UpdateVideoStatus(int, string) error

	UpdateProductAttribute(PrdctAttrUpdateCndtn) error

	UpdatePrice(PriceUpdate) error

	// Update Shipment Type by product sku
	// Params:
	// 1. sku
	// 2. New shipment value
	UpdateShipmentBySKU(string, int) error

	// Update Product Status
	// Params:
	// 1. SeqId
	// 2. Status
	UpdateProductStatus(int, string) error

	// Update Product Simple Status
	// Params:
	// 1. Simple SeqId
	// 2. Status
	UpdateProductSimpleStatus(int, string) error

	// Update petApproved of supplied product
	// Params:
	// 1. ConfigId
	// 2. petApproved
	SetPetApproved(int, int) error

	UpdateJabongDiscount(JabongDiscount) error

	UpdateProductAttributeSystem(ProductAttrSystemUpdate) error

	UpdateProduct(configId int, criteria ProUpdateCriteria) error

	ResetSSRCounter() error

	GetProductBySellerIdSku(int, []string) ([]ProductSmallSimples, error)
}

func GetCurrentAdapter() ProductDBAdapter {
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	return GetAdapter(conf.DbAdapter)
}

func GetAdapter(adapterName string) ProductDBAdapter {

	switch adapterName {
	case "MySql":
		if mysqlAdapter != nil {
			return mysqlAdapter
		}
		mysql := &MySqlAdapter{}
		return mysql
	default:
		if mongoAdapter != nil {
			return mongoAdapter
		}
		mongo := &MongoAdapter{}
		return mongo
	}
	return nil
}
