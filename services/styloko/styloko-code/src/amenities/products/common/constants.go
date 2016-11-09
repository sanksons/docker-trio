package common

const (
	//Adapter
	DB_ADAPTER_MONGO = "Mongo"
	DB_ADAPTER_MYSQL = "MySql"
	DB_READ_ADAPTER  = DB_ADAPTER_MONGO

	//Headers
	HEADER_UPDATE_TYPE      = "Update-Type"
	HEADER_REQUEST_PLATFORM = "Request-Platform"
	HEADER_VISIBILITY_TYPE  = "Visibility-Type"
	HEADER_EXPANSE          = "Expanse"
	HEADER_NOCACHE          = "No-Cache"
	HEADER_PUBLISH          = "Publish"

	//General
	DEBUG_KEY_NODE       = "node"
	PRODUCT_API          = "PRODUCTS"
	PRODUCT_UPDATE       = "PRODUCT_UPDATE"
	SELLER_CENTER        = "SC"
	REQUEST_STRUCT       = "REQUEST-STRUCT"
	IODATA               = "iodata"
	FORMAT_MYSQL_TIME    = "2006-01-02 15:04:05"
	FORMAT_LOG_TIME      = "20060102"
	DEFAULT_FROM_DATE    = "1980-01-01T00:00:00+00:00"
	DEFAULT_TO_DATE      = "2080-01-01T00:00:00+00:00"
	PRODUCT_SUFFIX       = "INDFAS"
	ATTRIBUTES_INFO_FILE = "attributes_info.json"
	ATTRIBUTES_IGNORE    = "attributes_ignore"
	ATTRIBUTES_ALLOWED   = "attributes_allowed"
	SELLERS_IGNORE       = "sellers_ignore"
	IS_RETURNABLE        = "is_returnable"
	SELLER_COMMISSION    = "Commissions"
	ACTION_REPLACE       = 1
	ACTION_ADD           = 2
	ACTION_REMOVE        = 3
	SSR_COUNTER_NAME     = "ssrcounter"
	MRK_AS_SC_PRODUCT    = "sellercenter"
	TO_DATE_DIFF         = 86399 //23:59:59 in seconds
	SELLER_RETRY_COUNT   = 3
	MIGRATE_PRODUCTS_KEY = "products-not-found"

	//Errors
	NO_PRODUCT_ATTRIBUTE = "This attribute is not present in this product"
	MANDATORY_ATTRIBUTE  = "Cannot remove a mandatory attribute"
	NO_ATTRIBUTE_OPTIONS = "No options found for this attribute value"

	//Addition and deletion of Nodes
	ADD_NODE    = "addNode"
	DELETE_NODE = "deleteNode"

	//Status
	STATUS_SUCCESS  = "success"
	STATUS_FAILURE  = "failure"
	STATUS_ACTIVE   = "active"
	STATUS_DELETE   = "delete"
	STATUS_DELETED  = "deleted"
	STATUS_INACTIVE = "inactive"
	STATUS_APPROVED = "approved"
	STATUS_PENDING  = "pending"

	//Option types
	OPTION_TYPE_MULTI  = "multi_option"
	OPTION_TYPE_SINGLE = "option"
	OPTION_TYPE_VALUE  = "value"
	OPTION_TYPE_CUSTOM = "custom"
	OPTION_TYPE_SYSTEM = "system"

	//Parameters
	PARAM_LIMIT   = "limit"
	PARAM_OFFSET  = "offset"
	PARAM_SELLERS = "sellers"
	PARAM_SKU     = "sku"
	PARAM_ID      = "id"

	//Product type simple or config
	PRODUCT_TYPE_SIMPLE = "simple"
	PRODUCT_TYPE_CONFIG = "config"

	//Notification tags
	TAG_PRODUCT        = "product"
	TAG_PRODUCT_UPDATE = "product-update"
	TAG_PRODUCT_CREATE = "product-create"
	TAG_PRODUCT_SYNC   = "product-sync"
)

//
// Supported Product Updates
//
const (
	UPDATE_TYPE_PRICE             = "price"
	UPDATE_TYPE_SHIPMENT          = "shipment"
	UPDATE_TYPE_CACHE             = "cache"
	UPDATE_TYPE_NODE              = "node"
	UPDATE_TYPE_PRODUCT           = "product"
	UPDATE_TYPE_IMAGEADD          = "imageadd"
	UPDATE_TYPE_IMAGEDEL          = "imagedel"
	UPDATE_TYPE_VIDEO             = "video"
	UPDATE_TYPE_MYSQL             = "mysql"
	UPDATE_TYPE_VIDEO_STATUS      = "videostatus"
	UPDATE_TYPE_PRODUCT_ATTRIBUTE = "attribute"
	UPDATE_TYPE_PRODUCT_STATUS    = "productstatus"
	UPDATE_TYPE_JABONG_DISCOUNT   = "jabongdiscount"
)

//
// Profiler constants
//
const (
	PUT_INVALIDATE_CACHE_NODE = "product_invalidate_cache"
	PUT_RESPONSE_NODE         = "product_put_response_node"
	PUT_UPDATE_NODE           = "product_put_update_node"
	PUT_VALIDATE_NODE         = "product_put_validate_node"
	GET_SC_VALIDATE_NODE      = "product_get_sc_validate_node"
	GET_SC_FETCH_NODE         = "product_get_sc_fetch_node"
	GET_SC_RESPONSE_NODE      = "product_get_sc_resp_node"
	GET_CACHE_GET_NODE        = "product_cache_get_node"
	GET_PRO_RESPONSE_NODE     = "product_get_response_node"
	GET_PREPARE_QUERY_NODE    = "product_get_prepare_query_node"
	GET_LOAD_DATA_NODE        = "product_get_load_data_node"
	GET_VISIBILITY_NODE       = "product_get_visibility_node"
	POST_INSERT_NODE          = "product_post_insert_node"
	POST_RESPONSE_NODE        = "product_post_response_node"
	POST_VALIDATE_NODE        = "product_post_validate_node"
	DELETE_CACHE_NODE         = "product_delete_cache_node"
	PRODUCT_CACHE_SET         = "product_cache_set"
	PRODUCT_CACHE_GET         = "product_cache_get"
	BOOTSTRAP_NODE            = "product_bootstrap_node"
	FE_ORG_API                = "fe_org_api"
	FE_ORG_COM                = "fe_org_commision"
	STOCK_CALL                = "fe_stock_call"
)

//
//Mongo Collection names
//
const (
	PRODUCT_COLLECTION       = "products"
	SIMPLE_COLLECTION        = "simples"
	PIMAGE_COLLECTION        = "productImages"
	PVIDEO_COLLECTION        = "productVideos"
	PGROUP_COLLECTION        = "productGroups"
	TAXCLASS_COLLECTION      = "taxClass"
	CATEGORY_COLLECTION      = "categories"
	BRAND_COLLECTION         = "brand"
	ATTRIBUTESET_COLLECTION  = "attributeSets"
	ATTRIBUTE_COLLECTION     = "attributes"
	ATTRIBUTE_MAP_COLLECTION = "attributeMapping"
	COUNTER_COLLECTION       = "counters"
	PREPACK_COUNTER          = "prepack"
	DUMMY_IMAGES             = "dummyImages"
)

//
// Expanse
//
const (
	EXPANSE_CATALOG   = "Catalog"
	EXPANSE_XLARGE    = "XLarge"
	EXPANSE_LARGE     = "Large"
	EXPANSE_SMALL     = "Small"
	EXPANSE_XSMALL    = "XSmall"
	EXPANSE_MEDIUM    = "Medium"
	EXPANSE_SOLR      = "Solr"
	EXPANSE_PROMOTION = "Promotion"
	EXPANSE_MEMCACHE  = "memcache"
)

//
// Visbility conditions
//
const (
	VISIBILITY_PDP  = "PDP"
	VISIBILITY_MSKU = "MULTI-SKU"
	VISIBILITY_DOOS = "DOOS"
	VISIBILITY_NONE = "NONE"
)

// Bootstrap Flags

const (
	BOOTSTRAP_DRIVER   = "PRODUCT_SOLR_BOOTSTRAP"
	BOOTSTRAP_HASH_MAP = "PRODUCT_BOOTSTRAP_FLAGS"
)

// Validation Types
const (
	VALIDATION_DECIMAL = "decimal"
	VALIDATION_INTEGER = "integer"
)

const (
	SYNC_ATTRIBUTE_SYSTEM  = "AttributeSystem"
	SYNC_ATTRIBUTE_GENERAL = "AttributeGeneral"
)

const (
	SYSTEM_TY         = "ty"
	SYSTEM_PET_STATUS = "pet_status"
)

const (
	ATTR_SC_PRODUCT         = "sc_product"
	ATTR_IS_RETURNABLE      = "is_returnable"
	ATTR_IS_CANCELABLE      = "is_cancelable"
	ATTR_IS_COD             = "is_cod"
	ATTR_IS_FRAGILE         = "is_fragile"
	ATTR_IS_SURFACE         = "is_surface"
	ATTR_PROCESSING_TIME    = "processing_time"
	ATTR_SHIPPING_AMOUNT    = "shipping_amount"
	ATTR_BLOCK_CATALOG      = "block_catalog"
	ATTR_PACK_ID            = "pack_id"
	ATTR_PACK_QTY           = "pack_qty"
	ATTR_VAT_CMPNY_CONTRI   = "vat_company_contribution"
	ATTR_VAT_CUST_CONTRI    = "vat_customer_contribution"
	ATTR_CUSTOMIZATION_TIME = "customization_time"
	ATTR_IMAGE_ORIENTATION  = "image_orientation"
	ATTR_BG_COLOR           = "bgcolor"
	ATTR_COLOR              = "color"
	ATTR_UPPR_MTRL_DTLS     = "upper_material_details"
	ATTR_FABRIC_DETAILS     = "fabric_details"
	ATTR_FRAME_MTRL_DTLS    = "frame_material_detail"
	ATTR_STRAP_COLOR        = "strap_color"
	ATTR_STRAP_MTRL_DTL     = "strap_material_detail"
	ATTR_FRAME_COLOR        = "frame_color"
	ATTR_DIAL_COLOR         = "dial_color"
	ATTR_STANDARD_SIZE      = "standard_size"
	ATTR_VARIATION          = "variation"
	ATTR_VARIATIONS         = "variations"
	ATTR_FIT                = "fit"
	ATTR_FITS               = "fits"
	ATTR_QUALITIES          = "qualities"
	ATTR_PACKQTY            = "pack_qty"
	ATTR_PACKID             = "pack_id"
	ATTR_QUALITY            = "quality"
	ATTR_MTRLS_CODE         = "materials_code"
	ATTR_MTRL_CODE          = "material_code"
	ATTR_PRDCT_WRNTY        = "product_warranty"
	ATTR_PRDCTS_WRNTY       = "products_warranty"
	ATTR_SCNDRY_CLRS        = "secondary_colors"
	ATTR_SCNDRY_CLR         = "secondary_color"
	ATTR_LENS_TYP           = "lens_type"
	ATTR_LENS_TYPS          = "lens_types"
	ATTR_JEANS_WE           = "jeans_wash_effect"
	ATTR_JEANS_WES          = "jeans_wash_effects"
	ATTR_POCKETS            = "pockets"
	ATTR_POCKET             = "pocket"

	//@todo: below attributes needs to be loaded from seller
	ATTR_DISPATCH_LOCATION = "dispatch_location"

	DEF_IS_CANCELABLE      = "1"
	DEF_IS_COD             = "1"
	DEF_IS_FRAGILE         = "0"
	DEF_IS_PRECIOUS        = "0"
	DEF_IS_SURFACE         = "0"
	DEF_PROCESSING_TIME    = "0"
	DEF_SHIPPING_AMOUNT    = "1"
	DEF_BLOCK_CATALOG      = "N"
	DEF_PACK_ID            = "0"
	DEF_PACK_QTY           = "1"
	DEF_VAT_CMPNY_CONTRI   = "0"
	DEF_VAT_CUST_CONTRI    = "0"
	DEF_CUSTOMIZATION_TIME = "0"
	DEF_IMAGE_ORIENTATION  = "landscape"
	DEF_IS_RETURNABLE      = "1"
	DEF_DISPATCH_LOCATION  = "3"
)
