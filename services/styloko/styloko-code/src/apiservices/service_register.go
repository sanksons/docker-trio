package apiservices

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
	categoryTree "amenities/category_tree/get"
	migrationsCli "amenities/migrations/common"
	migrations "amenities/migrations/get"
	productDelete "amenities/products/delete"
	productSearch "amenities/products/get/search/api"
	productCreate "amenities/products/post"
	bootstrap "amenities/products/post/bootstrap"
	productUpdate "amenities/products/put/api"
	sizeChartGet "amenities/sizechart/get"
	sizeChartPost "amenities/sizechart/post"
	standardSizeGet "amenities/standardsize/get"
	standardSizeCreate "amenities/standardsize/post"
	stock "amenities/stock"
	taxClassGet "amenities/taxclass/get"
	factory "common/ResourceFactory"
	"common/appconfig"
	"common/appconstant"
	mysqlMig "common/migration/mysql"
	"common/notification"
	simplifier "simplifier"

	"github.com/jabong/floRest/src/service"
)

// Register -> Registers everything.
func Register() {
	registerConfig()
	registerErrors()
	registerAllApis()
	registerCustomServices()
	overrideConfByEnvVariables()
}

func registerAllApis() {
	service.RegisterApi(new(get.BrandAPI))
	service.RegisterApi(new(post.BrandAPI))
	service.RegisterApi(new(put.BrandAPI))
	service.RegisterApi(new(standardSizeCreate.CreateStandardSizeApi))
	service.RegisterApi(new(categoryGet.CategoryAPI))
	service.RegisterApi(new(categoryPost.CategoryAPI))
	service.RegisterApi(new(categoryPut.CategoryAPI))
	service.RegisterApi(new(attr.GetAttributesApi))
	service.RegisterApi(new(attrSet.GetAttributesSetsApi))
	service.RegisterApi(new(attrSetPut.AttributeSetUpdateApi))
	service.RegisterApi(new(attrSetPost.AttributeSetCreateApi))
	service.RegisterApi(new(attrPost.AttributesCreateApi))
	service.RegisterApi(new(attrPut.AttributesUpdateApi))
	service.RegisterApi(new(productUpdate.ProductAPI))
	service.RegisterApi(new(productSearch.ProductsApi))
	service.RegisterApi(new(productCreate.CreateProductsApi))
	service.RegisterApi(new(productDelete.DeleteProductsApi))
	service.RegisterApi(new(standardSizeCreate.CreateStandardSizeApi))
	service.RegisterApi(new(standardSizeGet.GetStandardSizeApi))
	service.RegisterApi(new(taxClassGet.TaxClassApi))
	service.RegisterApi(new(bootstrap.BootstrapAPI))
	service.RegisterApi(new(sizeChartPost.SizeChartCreateApi))
	service.RegisterApi(new(migrations.MigrationAPI))
	service.RegisterApi(new(catalogty.CatalogTyApi))
	service.RegisterApi(new(stock.StockAPI))
	service.RegisterApi(new(sizeChartGet.SizeChartGetApi))
	service.RegisterApi(new(simplifier.ReturnError))
	service.RegisterApi(new(categoryTree.CategoryTreeApi))
}

func registerConfig() {
	service.RegisterConfig(new(appconfig.AppConfig))
}

func registerErrors() {
	service.RegisterHttpErrors(appconstant.AppErrorCodeToHttpCodeMap)
}

func registerCustomServices() {
	service.RegisterCustomApiInitFunc(func() {
		factory.InitializeFactories()
		mysqlMig.InitDotSql()
		migrationsCli.RunMigrationFromCli()
		notification.InitNotifpool()
		migrationsCli.ClearAllRedisFlags()
	})
}

func overrideConfByEnvVariables() {
	service.RegisterConfigEnvUpdateMap(appconfig.MapEnvVariables())
	service.RegisterGlobalEnvUpdateMap(appconfig.MapEnvGlobalVariables())
}
