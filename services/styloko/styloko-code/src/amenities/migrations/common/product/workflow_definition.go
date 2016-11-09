package product

import (
	"amenities/migrations/common/util"
	productsApi "amenities/products/common"
	"common/ResourceFactory"
	"common/appconfig"
	"common/notification"
	"common/notification/datadog"
	"common/utils"
	"common/xorm/mysql"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"strconv"

	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/config"
	mgo "gopkg.in/mgo.v2"
	"strings"
	"time"
)

var ErrExhausted = errors.New("Limit exhausted")
var ErrEmpty = errors.New("No value found")

var AttributesCache map[string][]*CatalogAttribute

const (
	PRO_ERR_COLL            = "productErrors"
	CRITERIA_TYPE_SELLER    = "seller"
	CRITERIA_TYPE_STATUS    = "status"
	CRITERIA_TYPE_PARTIAL   = "partial"
	CRITERIA_TYPE_BRAND     = "brand"
	CRITERIA_TYPE_PROMOTION = "promotion"
)

type ProductErr struct {
	Id  int    `bson:"id"`
	Msg string `bson:"msg"`
}

func StartPartialMigration(minMax string) {
	minMaxArr := strings.Split(minMax, ",")

	var criteria ProductFetchCriteria

	criteria.MinId, _ = utils.GetInt(minMaxArr[0])
	criteria.MaxId, _ = utils.GetInt(minMaxArr[1])
	criteria.Type = CRITERIA_TYPE_PARTIAL
	startMongoMigration(criteria)
}

func StartActiveMigration(drpCollection bool) {
	if drpCollection {
		DropCollection()
		ReCreateIndexes()
	}
	var criteria ProductFetchCriteria
	criteria.Status = "active"
	criteria.Type = CRITERIA_TYPE_STATUS
	startMongoMigration(criteria)
}

func StartSellerMigration(sellerId int) {
	if sellerId <= 0 {
		return
	}
	var criteria ProductFetchCriteria
	criteria.SellerId = sellerId
	criteria.Type = CRITERIA_TYPE_SELLER
	startMongoMigration(criteria)
}

func StartMigrationByBrand(brandId int) {
	if brandId <= 0 {
		return
	}
	var criteria ProductFetchCriteria
	criteria.BrandId = brandId
	criteria.Type = CRITERIA_TYPE_BRAND
	startMongoMigration(criteria)
}

func StartMigrationByPromotion(promotionId int) {
	if promotionId <= 0 {
		return
	}
	var criteria ProductFetchCriteria
	criteria.PromotionId = promotionId
	criteria.Type = CRITERIA_TYPE_PROMOTION
	startMongoMigration(criteria)
}

func StartInActiveMigration() {
	var criteria ProductFetchCriteria
	criteria.Status = "inactive"
	criteria.Type = CRITERIA_TYPE_STATUS
	startMongoMigration(criteria)
}

func startMongoMigration(criteria ProductFetchCriteria) {
	StartWorkflow(criteria)
	FinishWorkflow()
}

func DropCollection() {
	mgoSession := ResourceFactory.GetMongoSession(util.Products)
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(PRO_ERR_COLL)
	mongodb.DropCollection()

	sess := mgoSession.SetCollection(util.Products)
	sess.DropCollection()
}

func ReCreateIndexes() {
	mgoSession := ResourceFactory.GetMongoSession(util.Products)
	defer mgoSession.Close()
	sess := mgoSession.SetCollection(util.Products)

	var normalIndexes = []string{
		"sellerId",
		"-createdAt",
		"-updatedAt",
		"group.seqId",
		"brandId",
		"categories",
		"leaf",
		"petApproved",
		"status",
		"productSet",
	}
	var uniqueIndexes = []string{
		"seqId",
		"sku",
		"simples.seqId",
		"simples.sku",
	}
	var uniqueSparseIndexes = []string{
		"videos.seqId",
		"images.seqId",
	}
	for _, v := range normalIndexes {
		sess.DropIndex(v)
		sess.EnsureIndex(mgo.Index{
			Key:    []string{v},
			Unique: false,
			Sparse: false,
		})
	}
	for _, v := range uniqueIndexes {
		sess.DropIndex(v)
		sess.EnsureIndex(mgo.Index{
			Key:    []string{v},
			Unique: true,
			Sparse: false,
		})
	}
	for _, v := range uniqueSparseIndexes {
		sess.DropIndex(v)
		sess.EnsureIndex(mgo.Index{
			Key:    []string{v},
			Unique: true,
			Sparse: true,
		})
	}
}

func FinishWorkflow() {
	printInfo("Started: Finish Workflow")
	printInfo("Finished: Finish Workflow")
}

func StartWorkflow(criteria ProductFetchCriteria) bool {
	printInfo("Started: Workflow Import")
	AttributesCache = make(map[string][]*CatalogAttribute)
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	var offset = 0
	var limit = conf.ProductConLimit

	var count = 1
	for count != 0 {
		var offsetStr = strconv.Itoa(offset)
		var limitStr = strconv.Itoa(limit)

		//print logs
		switch criteria.GetCriteriaType() {
		case CRITERIA_TYPE_PARTIAL:
			printInfo(
				"partial:" + strconv.Itoa(criteria.MinId) + "," + strconv.Itoa(criteria.MaxId) +
					"| Fetching products " + offsetStr + "," + limitStr)
		case CRITERIA_TYPE_SELLER:
			printInfo("Seller:" + strconv.Itoa(criteria.SellerId) + "| Fetching products " + offsetStr + "," + limitStr)
		case CRITERIA_TYPE_STATUS:
			printInfo("Status:" + criteria.Status + "| Fetching products " + offsetStr + "," + limitStr)
		case CRITERIA_TYPE_BRAND:
			printInfo("Brand:" + strconv.Itoa(criteria.BrandId) + "| Fetching products " + offsetStr + "," + limitStr)
		case CRITERIA_TYPE_PROMOTION:
			printInfo("Promotion:" + strconv.Itoa(criteria.PromotionId) + "| Fetching products " + offsetStr + "," + limitStr)
		}

		response, err := GetProductsFromMysql(offsetStr, limitStr, criteria)
		if err != nil {
			printErr(err)
			break
		}
		rows, ok := response.(*core.Rows)
		if !ok {
			printErr(errors.New("importWorkflow(): error parsing result set"))
			break
		}
		count = ProcessProductRows(rows)
		rows.Close()
		offset = offset + limit
	}

	printInfo("Finished: Workflow Import")

	return true
}

//
// Migrate Single Product from MySQl -> Mongo
//
func MigrateSingleProduct(pid int) {

	printInfo("Started: Migrating product(" + strconv.Itoa(pid) + ")")
	AttributesCache = make(map[string][]*CatalogAttribute)
	var sqlQ string = `SELECT * FROM catalog_config WHERE id_catalog_config=` + strconv.Itoa(pid)
	response, err := mysql.GetInstance().Query(sqlQ, true)
	if err != nil {
		notification.SendNotification(
			"Product Migration Failed[Mysql -> Mongo]",
			fmt.Sprintf("Product:%d, Message:%s", pid, err.Error()),
			[]string{"MongoMigration", "product"},
			datadog.ERROR,
		)
		printErr(errors.New("MigrateSingleProduct(): " + err.Error()))
		return
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		notification.SendNotification(
			"Product Migration Failed[Mysql -> Mongo]",
			fmt.Sprintf("Product:%d, Message:%s", pid, "(*core.Rows) Assertion failed"),
			[]string{"MongoMigration", "product"},
			datadog.ERROR,
		)
		printErr(errors.New("MigrateSingleProduct():(*core.Rows) Assertion failed"))
		return
	}
	defer rows.Close()
	var productWriteSuccess bool
	for rows.Next() {
		product := make(generalMap)
		err := rows.ScanMap(&product)
		if err != nil {
			printErr(errors.New("MigrateSingleProduct(): Unable to Scan Product"))
			return
		}
		s := ProcessProduct(product)
		if !s.State {
			notification.SendNotification(
				"Product Mongo Write Failed[Mysql -> Mongo]",
				fmt.Sprintf("Product:%d, Message:%s", pid, s.Msg),
				[]string{"MongoMigration", "product"},
				datadog.ERROR,
			)
			printErr(errors.New("MigrateSingleProduct(): " + s.Msg))
		} else {
			productWriteSuccess = true
		}
	}
	rows.Close()
	printInfo("Finished: Migrating product(" + strconv.Itoa(pid) + ")")

	printInfo("Start: Cache Invalidation(" + strconv.Itoa(pid) + ")")
	var cacheMngr productsApi.CacheManager
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	cacheMngr.CacheObj, err = cache.Get(conf.Cache)
	if err != nil {
		notification.SendNotification(
			"Product Cache Invalidation Failed[Mysql -> Mongo]",
			fmt.Sprintf("Product:%d, Message:%s", pid, err.Error()),
			[]string{"MongoMigration", "product"},
			datadog.ERROR,
		)
		printErr(fmt.Errorf("unable to clear cache: %s", err.Error()))
	} else {
		cacheMngr.DeleteById([]int{pid}, true)
	}
	printInfo("End: Cache Invalidation(" + strconv.Itoa(pid) + ")")
	printInfo("Start: Solr Publish(" + strconv.Itoa(pid) + ")")

	if productWriteSuccess {
		go func() {
			defer utils.RecoverHandler("Mysql -> Mongo")
			time.Sleep(1 * time.Millisecond)
			pro, err := productsApi.GetAdapter(productsApi.DB_ADAPTER_MONGO).GetById(pid)
			if err != nil {
				notification.SendNotification(
					"Product Publish Failed[Mysql -> Mongo]",
					fmt.Sprintf("Product:%d, Message:%s", pid, err.Error()),
					[]string{"MongoMigration", "product"},
					datadog.ERROR,
				)
				printErr(fmt.Errorf("unable to publish: %s", err.Error()))
			} else {
				pro.Publish("", true)
			}
		}()
	}
	printInfo("End: Solr Publish(" + strconv.Itoa(pid) + ")")
	return
}
