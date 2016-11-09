package product

import (
	"amenities/migrations/common/util"
	"common/ResourceFactory"
	"common/xorm/mysql"
	"errors"

	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type ProductGroup struct {
	Id   int    `bson:"seqId"`
	Name string `bson:"name"`
}

func MigrateProductGroup() {
	logger.Info("Migrating product groups")
	logger.Info("Starting to write Filter into Mongo")
	mgoSession := ResourceFactory.GetMongoSession("Products")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.ProductGroups)
	mongodb.DropCollection()
	sqlQ := `SELECT id_catalog_config_group, name FROM catalog_config_group;`
	response, err := mysql.GetInstance().Query(sqlQ, QUERY_MYSQL_MASTER)
	if err != nil {
		printErr(errors.New("MigrateProductGroup(): " + err.Error()))
	}
	rows, ok := response.(*core.Rows)
	if !ok {
		return
	}
	defer rows.Close()
	for rows.Next() {
		grp := ProductGroup{}
		err := rows.Scan(&grp.Id, &grp.Name)
		if err != nil {
			printErr(errors.New("MigrateProductGroup(): Unable to Scan Entry"))
		}
		mongodb.Insert(grp)
	}
	logger.Info("Product group migration done")
}
