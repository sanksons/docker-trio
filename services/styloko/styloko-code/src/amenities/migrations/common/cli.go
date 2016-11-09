package common

import (
	"amenities/migrations/common/attribute"
	"amenities/migrations/common/brand"
	"amenities/migrations/common/category"
	"amenities/migrations/common/filter"
	"amenities/migrations/common/product"
	"amenities/migrations/common/sizechart"
	"amenities/migrations/common/util"
	proUtil "amenities/products/common"
	"common/ResourceFactory"
	"flag"
	"fmt"
	"os"
	"simplifier"
	"strings"
)

// RunMigrationFromCli runs migrations with CLI args, and kills the server
func RunMigrationFromCli() {
	defer RecoverHandler("RunMigrationFromCli")
	var flagvar string
	var flagDaemon string
	var id int
	var minMax string
	flag.StringVar(&flagvar, "m", "", "Usage -m categories || -m=\"categories attributes products\"")
	flag.IntVar(&id, "i", 0, "Use in conjunction with -m")
	flag.StringVar(&minMax, "partial", "", "Use in conjunction with -m")
	flag.StringVar(&flagDaemon, "d", "", "Usage -d judgedaemon || -d=\"judgedaemon\"")
	flag.Parse()
	if flagvar == "" && flagDaemon == "" {
		return
	}
	if flagvar == "" && flagDaemon == judgedaemon {
		simplifier.StartJudgeDaemon()
		return
	}

	arrFlags := strings.Split(flagvar, " ")

	for _, x := range arrFlags {
		fmt.Printf("Running migration for: %s \n", x)
		switch strings.ToLower(x) {
		case AttributeSets:
			attribute.StartAttributeSetMigration()
			err := setSequence(util.AttributeSets)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error())
			}
			break

		case Attributes:
			attribute.StartAttributeMigration()
			err := setSequence(util.Attributes)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error())
			}
			fmt.Println("Started attribute mapping")
			attribute.StartMapping()
			break
		case AttributesIndex:
			attribute.EnsureIndexInDb()
			break
		case AttributesById:
			attribute.MigrateSingleAttribute(id)
			break
		case Brands:
			brand.StartBrandMigration()
			err := setSequence(util.Brands)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error())
			}
			break
		case BrandsIndex:
			brand.EnsureIndexInDb()
			break

		case Filters:
			filter.StartFilterMigration()
			err := setSequence(util.Filters)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error())
			}
			break

		case Categories:
			category.StartCategoryMigration()
			err := setSequence(util.Categories)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error())
			}
			break

		case SizeCharts:
			err := sizechart.StartSizeChartMigrationPartial()
			if err != nil {
				fmt.Println(err.Error())
			}
			err = setSequence(util.SizeCharts)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error())
			}
			break
		case TaxClass:
			product.MigrateTaxClass()
			err := setSequence(util.TaxClass)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error())
			}
			break

		// Product and related migration cases below
		case Products:
			product.StartActiveMigration(false)
			product.StartInActiveMigration()
			err := setSequence(util.Products)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error())
			}
			setSequence(util.Simples)
			setSequence(util.ProductImages)
			setSequence(util.ProductVideos)
			break

		case ProductsActive:
			product.StartActiveMigration(false)
			break

		case ProductsPartial:
			product.StartPartialMigration(minMax)
			break

		case ProductsDrop:
			product.StartActiveMigration(true)
			product.StartInActiveMigration()
			err := setSequence(util.Products)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error())
			}
			setSequence(util.Simples)
			setSequence(util.ProductImages)
			setSequence(util.ProductVideos)
			break

		case ProductsInactive:
			product.StartInActiveMigration()
			break

		case ProductGroups:
			product.MigrateProductGroup()
			err := setSequence(util.ProductGroups)
			if err != nil && err.Error() != "" {
				fmt.Println(err.Error)
			}
			break

		case ProductsIndex:
			product.ReCreateIndexes()
			break

		case ProductsById:
			product.MigrateSingleProduct(id)
			break

		case DeleteProductById:
			err := product.DeleteProductById(id)
			if err != nil {
				fmt.Println(err.Error())
			}
			break

		case ProductsBySeller:
			product.StartSellerMigration(id)
			break

		case ProductsByBrand:
			product.StartMigrationByBrand(id)
			break

		case ProductsByPromotion:
			product.StartMigrationByPromotion(id)
			break

		case ProductsSizeChart:
			err := sizechart.WriteSizeChartToProduct()
			if err != nil {
				fmt.Println(err.Error())
			}
			break

		case ProductsSizeChartById:
			err := sizechart.SingleProductSizeChartUpdate(id)
			if err != nil {
				fmt.Println(err.Error())
			}
			break

		case ResetCounter:
			err := setSequence(
				proUtil.DUMMY_IMAGES,
				util.Products,
				util.Simples,
				util.ProductImages,
				util.ProductVideos,
				util.AttributeSets,
				util.Attributes,
				util.Brands,
				util.Categories,
				util.Filters,
				util.ProductGroups,
				util.TaxClass,
				util.SizeCharts,
				util.PrePack,
			)
			if err != nil {
				fmt.Println(err.Error())
			}
			break
		}
		fmt.Printf("Finished migration for: %s \n", x)
	}
	os.Exit(0)
}

func setSequence(args ...string) error {
	mgoSession := ResourceFactory.GetMongoSession("Sequencing")
	defer mgoSession.Close()
	fmt.Println(args)

	for _, v := range args {
		fmt.Printf("Starting sequencing for %s\n", v)
		var counter int
		var err error
		switch v {

		case
			proUtil.DUMMY_IMAGES,
			util.SizeCharts,
			util.ProductGroups,
			util.Brands,
			util.Products,
			util.Simples,
			util.ProductImages,
			util.ProductVideos:
			counter, err = GetSeqCounter(v)
			if err != nil {
				fmt.Println(err)
				return err
			}
		case util.PrePack:
			counter, err = GetSeqCounterMysql(v)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		err = mgoSession.SetCollectionInCounter(v, counter)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Finished sequencing for %s\n", v)
	}
	return nil
}
