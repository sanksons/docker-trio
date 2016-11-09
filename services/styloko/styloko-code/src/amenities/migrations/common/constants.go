package common

const (
	Brands                = "brands"
	BrandsIndex           = "brandsindex"
	Filters               = "filters"
	Categories            = "categories"
	CategorySegments      = "categorysegments"
	AttributeSets         = "attributesets"
	Attributes            = "attributes"
	AttributesIndex       = "attributesindex"
	AttributesById        = "attributesbyid"
	Counter               = "counters"
	Products              = "products"
	ProductGroups         = "productgroups"
	ProductsIndex         = "productsindex"
	ProductsDrop          = "productsdrop"
	ProductsPartial       = "productspartial"
	ProductsActive        = "productsactive"
	ProductsDeleted       = "productsdeleted"
	ProductsInactive      = "productsinactive"
	ProductsBySeller      = "productsbyseller"
	ProductsByBrand       = "productsbybrand"
	ProductsByPromotion   = "productsbypromotion"
	ProductsById          = "productsbyid"
	DeleteProductById     = "deleteproductbyid"
	TaxClass              = "taxclass"
	SizeCharts            = "sizecharts"
	SizeChartsIndex       = "sizechartsindex"
	ProductsSizeChart     = "productsizechart"
	ProductsSizeChartById = "productsizechartbyid"
	SizeChartMapping      = "sizechartmapping"
	ResetCounter          = "resetcounter"
	judgedaemon           = "judgedaemon"
)

// Redis Flags
const (
	MIGRATION_DRIVER        = "PRODUCT_SOLR_BOOTSTRAP"
	MIGRATION_HASH_MAP      = "MIGRATION_HASH_MAP"
	PRODUCT_BOOTSTRAP_FLAGS = "PRODUCT_BOOTSTRAP_FLAGS"
)

// GetKeys returns a string with comma seperated possible keys
func GetKeys() string {
	keys := Brands + ", " + Filters + ", " + Categories + ", " + CategorySegments + ", " + AttributeSets + ", " + Attributes + ", " + Products + ", " + ProductGroups + ", " + ProductsIndex + ", " + ProductsDrop + ", " + ProductsActive + ", " + ProductsInactive + ", " + ProductsBySeller + ", " + ProductsById + ", " + TaxClass + ", " + SizeCharts + ", " + ProductsSizeChart

	return keys
}

// GetStatusKeys returns an array of valid status keys for locking flags
func GetStatusKeys() []string {
	return []string{Products, ProductGroups, ProductsIndex, ProductsDrop, ProductsActive, ProductsInactive, ProductsBySeller, ProductsById, ProductsSizeChart, ProductsBySeller, ProductsById}
}
