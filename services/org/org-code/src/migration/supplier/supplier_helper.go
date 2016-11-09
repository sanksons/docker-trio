package supplier

import (
	"common/ResourceFactory"
	"common/appconfig"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/cache"
	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"time"
)

var cacheObj cache.CacheInterface

//funtion takes input *sql.Rows, gets data from mysql and
//gives []OrgMongo as output
func processSupplierRows(rows *sql.Rows) ([]OrgMongo, error) {
	logger.Info("Processing Supplier Rows")
	orgInfo := []OrgMongo{}
	for rows.Next() {
		s := OrgMongo{}
		err := rows.Scan(
			&s.SeqId,
			&s.VendorId,
			&s.OrgName,
			&s.SellerName,
			&s.Status,
			&s.OrderEmail,
			&s.Contact,
			&s.Phone,
			&s.CustomercareEmail,
			&s.CustomercareContact,
			&s.CustomercarePhone,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.Street,
			&s.StreetNo,
			&s.City,
			&s.Postcode,
			&s.CountryCode)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in processSupplierRows() %v for id %d", err.Error(), s.SeqId))
			rows.Close()
			return nil, err
		}
		orgInfo = append(orgInfo, s)
	}
	rows.Close()
	return orgInfo, nil
}

//This function takes input []OrgMongo as input and
//places these into to final []OrgMongo with valid
//created and updated dates and returns it
func TransformSupplier(orgInfo []OrgMongo) []OrgMongo {
	logger.Info("Transforming Supplier")
	for k, v := range orgInfo {
		tmp := OrgMongo{}
		if v.CreatedAt == nil {
			time := time.Now()
			tmp.CreatedAt = &time
		} else {
			tmp.CreatedAt = v.CreatedAt
		}
		if v.UpdatedAt == nil {
			time := time.Now()
			tmp.UpdatedAt = &time
		} else {
			tmp.UpdatedAt = v.UpdatedAt
		}
		tmp.SeqId = v.SeqId
		tmp.VendorId = v.VendorId
		tmp.OrgName = v.OrgName
		tmp.SellerName = v.SellerName
		tmp.Status = v.Status
		tmp.OrderEmail = v.OrderEmail
		tmp.Contact = v.Contact
		tmp.Phone = v.Phone
		tmp.CustomercareContact = v.CustomercareContact
		tmp.CustomercareEmail = v.CustomercareEmail
		tmp.CustomercarePhone = v.CustomercarePhone
		tmp.StreetNo = v.SellerName
		tmp.Street = v.Street
		tmp.City = v.City
		tmp.Postcode = v.Postcode
		tmp.CountryCode = v.CountryCode
		orgInfo[k] = tmp
	}
	return orgInfo
}

//This function writes []OrgMongo into mongo
func writeSupplierToMongo(org []OrgMongo, chunkSize int) error {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, SUPPLIER_MIGRATION)
	defer func() {
		logger.EndProfile(profiler, SUPPLIER_MIGRATION)
	}()
	mgoSession := ResourceFactory.GetMongoSession(SUPPLIER_MIGRATION)
	mgoObj := mgoSession.SetCollection(SELLER_COLLECTION)
	defer mgoSession.Close()
	logger.Info(fmt.Sprintf("Starting to write Org Chunk into Mongo %v", chunkSize))
	for _, v := range org {
		logger.Debug(fmt.Sprintf("Inserting seqId %d in mongo", v.SeqId))
		var orgInfo OrgMongo
		upsertVal := true
		deleteVal := false
		returnNew := true
		updatedVal := bson.M{"$set": v}
		findCriteria := bson.M{"seqId": v.SeqId}
		change := mgo.Change{Update: bson.M(updatedVal), Upsert: upsertVal, Remove: deleteVal, ReturnNew: returnNew}
		_, err := mgoObj.Find(bson.M(findCriteria)).Apply(change, &orgInfo)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in Inserting Org into Mongo %v", err.Error()))
			return err
		}
	}
	return nil
}

//This breaks []Orgmongo into chunks of 1000
//,calls writeOrgToMongo from inside it
//also invalidates cache for all ids
func breakInChunks(org []OrgMongo) error {
	var ids []string
	limit := 1000
	for i := 0; i < len(org); {
		if limit > len(org) {
			limit = len(org)
		}
		var orgChunk []OrgMongo
		for j := i; j <= limit; j++ {
			if j >= len(org) {
				break
			}
			ids = append(ids, strconv.Itoa(org[j].SeqId))
			orgChunk = append(orgChunk, org[j])
		}
		logger.Info(fmt.Sprintf("Prepare chunk data for %d org", len(orgChunk)))
		//go func() error {
		err := writeSupplierToMongo(orgChunk, len(orgChunk))
		if err != nil {
			logger.Error(fmt.Sprintf("Error in calling writeOrgToMongo %v", err.Error()))
			return err
		}
		//return nil
		//}()
		go delete(ids)
		i = limit + 1
		limit = limit + 1000
	}
	return nil
}

//This function creates indexes for the collection to be created if the collection does not exist
func checkAndEnsureIndex() error {
	flag := false
	mgoSession := ResourceFactory.GetMongoSession(SUPPLIER_MIGRATION)
	defer mgoSession.Close()
	logger.Info("Checking if collection already exists")
	colNames, err := mgoSession.CollectionExists()
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting collection names from mongo :%s", err.Error()))
		return err
	}
	for _, v := range colNames {
		if v == SELLER_COLLECTION {
			flag = true
		}
	}
	if flag == true {
		logger.Info("Collection already exists so skipping creating indexes")
		return nil
	}
	EnsureIndexInDb()
	return nil
}

//This function creates indexes after dropping initial indexes(if any already existed in db)
func EnsureIndexInDb() {
	logger.Info("Creating Indexes for new collection to be created")
	mgoSession := ResourceFactory.GetMongoSession(SUPPLIER_MIGRATION)
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection(SELLER_COLLECTION)
	var normalIndexes = []string{
		"sync",
		"ordrEml",
		"slrId",
	}
	var uniqueIndexes = []string{
		//"slrId",
		"seqId",
	}
	for _, v := range normalIndexes {
		err := mgoObj.DropIndex(v)
		if err != nil {
			log.Println(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
		}
		err = mgoObj.EnsureIndex(mgo.Index{
			Key:    []string{v},
			Unique: false,
			Sparse: false,
		})
		if err != nil {
			log.Println(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
		}
	}
	for _, v := range uniqueIndexes {
		err := mgoObj.DropIndex(v)
		if err != nil {
			log.Println(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while dropping existing indexes :%s", err.Error()))
		}
		err = mgoObj.EnsureIndex(mgo.Index{
			Key:    []string{v},
			Unique: true,
			Sparse: true,
		})
		if err != nil {
			log.Println(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
			logger.Error(fmt.Sprintf("Error while ensuring indexes :%s", err.Error()))
		}
	}
	logger.Info("New indexes created")
}

//This queries mysql to get the max id and returns it
func getCounter() (int, error) {
	logger.Info("Getting Counter for org")
	sql := getMaxIdForSupplier()
	logger.Debug(fmt.Sprintf("Sql being fired to get max id is : %s", sql))
	driver, derr := ResourceFactory.GetMySqlDriver(SUPPLIER_MIGRATION)
	if derr != nil {
		logger.Error(fmt.Sprintf("Couldn't acquire mysql resource. Error: %s", derr.Error()))
		return 0, derr
	}
	var counter int
	rows, err := driver.Query(sql)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while execting query for getting max id :%s", err.DeveloperMessage))
		return 0, errors.New(fmt.Sprintf("Error while executing query for getting max id :%s", err.DeveloperMessage))
	}
	for rows.Next() {
		err := rows.Scan(&counter)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while scanning max id as counter :%s", err.Error()))
			return 0, err
		}
	}
	return counter, nil
}

//This function updates counter in mongo
func updateCounter(counter int) error {
	logger.Info("Updating Counter for org")
	mgoSession := ResourceFactory.GetMongoSession(SUPPLIER_MIGRATION)
	mgoObj := mgoSession.SetCollection(COUNTER_COLLECTION)
	defer mgoSession.Close()
	c := new(CounterInfo)
	change := mgo.Change{Update: bson.M{"_id": SELLER_COLLECTION, "seqId": counter}, Upsert: true}
	_, err := mgoObj.Find(bson.M{"_id": SELLER_COLLECTION}).Apply(change, &c)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while updating counters in mongo :%s", err.Error()))
		return err
	}
	return nil
}

//Deleting key from cache
func delete(keys []string) {
	err := cacheObj.DeleteBatch(keys)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while deleting key from cache ;%s", err.Error()))
	}
}

//initializes cache object
func initializeCacheObj() error {
	config := config.ApplicationConfig.(*appconfig.AppConfig)
	var err error
	cacheObj, err = cache.Get(config.Cache)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while getting cache object: %s", err.Error()))
		return err
	}
	return nil
}
