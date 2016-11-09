package apiservice

import (
	migrationsCli "common"
	factory "common/ResourceFactory"
	"common/appconfig"
	"common/appconstant"
	"common/notification"
	erp "erp/post"
	"github.com/jabong/floRest/src/service"
	"migration"
	"sellers/commissions"
	"sellers/get"
	"sellers/post"
	"sellers/put"
	"sellers/rating"
	"sellers/sync"
)

func Register() {
	registerConfig()
	registerErrors()
	registerAllApis()
	overrideConfByEnvVariables()
	registerAllFactories()
}

func registerAllApis() {
	service.RegisterApi(new(get.GetSellerApi))
	service.RegisterApi(new(post.CreateSellerApi))
	service.RegisterApi(new(put.UpdateSellerApi))
	service.RegisterApi(new(rating.UploadRatingApi))
	service.RegisterApi(new(erp.UpdateErpApi))
	service.RegisterApi(new(migration.MigrationApi))
	service.RegisterApi(new(commissions.CommissionsApi))
	service.RegisterApi(new(sync.Sync))
}

func registerConfig() {
	service.RegisterConfig(new(appconfig.AppConfig))
}

func registerErrors() {
	service.RegisterHttpErrors(appconstant.AppErrorCodeToHttpCodeMap)
}

func registerAllFactories() {
	service.RegisterCustomApiInitFunc(func() {
		factory.InitializeFactories()
		migrationsCli.RunMigrationFromCli()
		notification.InitNotifpool()
	})
}

func overrideConfByEnvVariables() {
	service.RegisterConfigEnvUpdateMap(appconfig.MapEnvVariables())
}
