package servicetest

import (
	attr "amenities/attributes/get/search/attributes"
	attrSet "amenities/attributes/get/search/set"
	attrPost "amenities/attributes/post/create/attributes"
	attrSetPost "amenities/attributes/post/create/set"
	attrPut "amenities/attributes/put/update/attributes"
	attrSetPut "amenities/attributes/put/update/set"
	get "amenities/brands/get"
	post "amenities/brands/post"
	put "amenities/brands/put"
	catalogty "amenities/catalog_ty/get"
	categoryGet "amenities/categories/get"
	categoryPost "amenities/categories/post"
	categoryPut "amenities/categories/put"
	productDelete "amenities/products/delete"
	productSearch "amenities/products/get/search/api"
	productCreate "amenities/products/post"
	bootstrap "amenities/products/post/bootstrap"
	productUpdate "amenities/products/put/api"
	sizeChartPost "amenities/sizechart/post"
	standardSizeGet "amenities/standardsize/get"
	standardSizeCreate "amenities/standardsize/post"
	stock "amenities/stock"
	taxClassGet "amenities/taxclass/get"
	factory "common/ResourceFactory"
	"common/appconstant"
	"github.com/jabong/floRest/src/service"
)

func InitializeTestService() {
	service.RegisterHttpErrors(appconstant.AppErrorCodeToHttpCodeMap)
	initTestLogger()
	//register configuration.
	initTestConfig()
	//Initialize factories
	factory.InitializeFactories()
	//register n Init apis
	registernInitApis()

	service.InitVersionManager()
	initialiseTestWebServer()
	service.InitHealthCheck()
}

func PurgeTestService() {

}

func registernInitApis() {
	getBrand := new(get.BrandAPI)
	service.RegisterApi(getBrand)
	getBrand.Init()

	postBrand := new(post.BrandAPI)
	service.RegisterApi(postBrand)
	postBrand.Init()

	putBrand := new(put.BrandAPI)
	service.RegisterApi(putBrand)
	putBrand.Init()

	sizeCreate := new(standardSizeCreate.CreateStandardSizeApi)
	service.RegisterApi(sizeCreate)
	sizeCreate.Init()

	getCategory := new(categoryGet.CategoryAPI)
	service.RegisterApi(getCategory)
	getCategory.Init()

	postCategory := new(categoryPost.CategoryAPI)
	service.RegisterApi(postCategory)
	postCategory.Init()

	putcategory := new(categoryPut.CategoryAPI)
	service.RegisterApi(putcategory)
	putcategory.Init()

	getAttr := new(attr.GetAttributesApi)
	service.RegisterApi(getAttr)
	getAttr.Init()

	getAttrset := new(attrSet.GetAttributesSetsApi)
	service.RegisterApi(getAttrset)
	getAttrset.Init()

	setAttrset := new(attrSetPut.AttributeSetUpdateApi)
	service.RegisterApi(setAttrset)
	setAttrset.Init()

	createAttrset := new(attrSetPost.AttributeSetCreateApi)
	service.RegisterApi(createAttrset)
	createAttrset.Init()

	createAttr := new(attrPost.AttributesCreateApi)
	service.RegisterApi(createAttr)
	createAttr.Init()

	updateAtrr := new(attrPut.AttributesUpdateApi)
	service.RegisterApi(updateAtrr)
	updateAtrr.Init()

	proUpd := new(productUpdate.ProductAPI)
	service.RegisterApi(proUpd)
	proUpd.InitTest()

	getProduct := new(productSearch.ProductsApi)
	service.RegisterApi(getProduct)
	getProduct.Init()

	createProduct := new(productCreate.CreateProductsApi)
	service.RegisterApi(createProduct)
	createProduct.Init()

	deleteProduct := new(productDelete.DeleteProductsApi)
	service.RegisterApi(deleteProduct)
	deleteProduct.Init()

	//service.RegisterApi(new(productSearch.ProductsApi))
	// service.RegisterApi(new(productCreate.CreateProductsApi))
	// service.RegisterApi(new(productDelete.DeleteProductsApi))
	service.RegisterApi(new(standardSizeCreate.CreateStandardSizeApi))
	service.RegisterApi(new(standardSizeGet.GetStandardSizeApi))
	service.RegisterApi(new(taxClassGet.TaxClassApi))
	service.RegisterApi(new(bootstrap.BootstrapAPI))
	service.RegisterApi(new(sizeChartPost.SizeChartCreateApi))
	service.RegisterApi(new(catalogty.CatalogTyApi))
	service.RegisterApi(new(stock.StockAPI))
}
