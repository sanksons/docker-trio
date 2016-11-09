package common

import (
	"common/pool"
	"flag"
	"log"
	"migration/supplier"
	"strings"
)

// RunMigrationFromCli runs migrations with CLI args, and kills the server
func RunMigrationFromCli() {
	defer pool.RecoverHandler("RunMigrationFromCli")
	var flagvar string
	var id int
	flag.StringVar(&flagvar, "m", "", "Usage -m categories || -m=\"categories attributes products\"")
	flag.IntVar(&id, "i", 0, "Use in conjunction with -m")
	flag.Parse()
	if flagvar == "" {
		return
	}

	arrFlags := strings.Split(flagvar, " ")
	for _, x := range arrFlags {
		log.Printf("Running migration for: %s \n", x)
		switch strings.ToLower(x) {
		case "sellers":
			err := supplier.StartSupplierMigration()
			if err != nil {
				log.Println(err.Error())
			}
			break
		case "index":
			supplier.EnsureIndexInDb()
			break
		}
		log.Printf("Finished migration for: %s \n", x)
	}
}
