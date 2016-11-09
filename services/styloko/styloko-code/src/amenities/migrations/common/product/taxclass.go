package product

import (
	"amenities/migrations/common/util"
	"common/ResourceFactory"
	"common/xorm/mysql"
	"errors"

	"github.com/go-xorm/core"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type TaxClass struct {
	Id         int     `bson:"seqId"`
	Name       string  `bson:"name"`
	Position   int     `bson:"position"`
	IsDefault  int     `bson:"isDefault"`
	TaxPercent float64 `bson:"taxPercent"`
}

func MigrateTaxClass() {

	logger.Info("Migrating tax classes")
	mgoSession := ResourceFactory.GetMongoSession("Products")
	defer mgoSession.Close()
	mongodb := mgoSession.SetCollection(util.TaxClass)
	mongodb.DropCollection()
	sqlQ := `SELECT id_catalog_tax_class,
	name,
	position,
	is_default,
	tax_percent
	FROM catalog_tax_class;`
	response, err := mysql.GetInstance().Query(sqlQ, QUERY_MYSQL_MASTER)
	if err != nil {
		printErr(errors.New("MigrateTaxClasses(): " + err.Error()))
	}
	rows, _ := response.(*core.Rows)
	defer rows.Close()
	for rows.Next() {
		class := TaxClass{}
		err := rows.Scan(
			&class.Id,
			&class.Name,
			&class.Position,
			&class.IsDefault,
			&class.TaxPercent,
		)
		if err != nil {
			printErr(errors.New("MigrateProductGroup(): Unable to Scan Entry"))
		}
		mongodb.Insert(class)
	}
}
