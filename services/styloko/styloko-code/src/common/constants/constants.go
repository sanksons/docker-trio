package constants

// Category Constants
const (
	CATEGORY_API        = "CATEGORIES"
	CATEGORY_SEARCH     = "CategorySearch"
	CATEGORY_CREATE     = "CategoryCreate"
	CATEGORY_UPDATE     = "CategoryUpdate"
	CATEGORY_COLLECTION = "categories"
)

//constants pertaining to BrandAPI
const (
	BRAND_API        = "BRANDS"
	BRAND_DATA       = "brandData"
	BRAND_UPDATE     = "brandUpdate"
	BRAND_COLLECTION = "brands"
	BRAND_CREATE     = "brandCreate"
)

// constants for category_ty API
const (
	CATALOG_TY_API = "CATALOGTY"
)

const (
	ATTRIBUTEGLOBALAPI          = "ATTRIBUTEGLOBAL"
	ATTRIBUTESETAPI             = "ATTRIBUTESETS"
	ATTRIBUTEAPI                = "ATTRIBUTES"
	ATTRIBUTEMAPPING_COLLECTION = "attributeMapping"
	ATTRIBUTESETS_COLLECTION    = "attributeSets"
	ATTRIBUTES_COLLECTION       = "attributes"
	Attributes                  = "attributes"
)

// Migration Constants
const (
	MIGRATION_API = "MIGRATIONS"
)
const (
	PRODUCT_RESOURCE_NAME = "Product"
)

// Attribute Cache Constants
const (
	ATTR_CACHE_KEY_FORMAT_ID       = "attributes-%d"
	ATTR_CACHE_KEY_FORMAT_CRITERIA = "attributes-%s-%s-%d"
	ATTR_CACHE_EXPIRY              = 1800
	ATTR_CACHE_DELETE_RETRY_COUNT  = 3
)
