package post

import (
	sizeUtils "amenities/sizechart/common"
	"common/appconstant"
	_ "common/notification"
	utils "common/utils"
	"fmt"
	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
	"strings"
)

type MapSizeChartProduct struct {
	id string
}

func (n *MapSizeChartProduct) SetID(id string) {
	n.id = id
}

func (n MapSizeChartProduct) GetID() (id string, err error) {
	return n.id, nil
}

func (a MapSizeChartProduct) Name() string {
	return "MapSizeChartProduct"
}

func (a MapSizeChartProduct) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, sizeUtils.SizeChartMapping)
	defer func() {
		logger.EndProfile(profiler, sizeUtils.SizeChartMapping)
	}()
	// Get the sizecharts which were created in DB
	collection, _ := io.IOData.Get(sizeUtils.SavedSizeChCollec)
	ty, _ := utils.GetRequestHeader(io, sizeUtils.SizeChartHeader)
	// upload sizchart for SKU level sizechart
	if ty == sizeUtils.SKU {
		// get the skus for which which sizchart were successfully uploaded
		var errMsg string
		uploadedSkus, _ := io.IOData.Get(sizeUtils.SuccessSkus)
		// Get skus for which sizechart was invalid
		failedValidationSkus, err := io.IOData.Get(sizeUtils.FailedSkus)
		if err != nil {
			errMsg = "Could not get validation-failed skus."
		} else if len(failedValidationSkus.([]string)) != 0 {
			errMsg = "Validation Failed for skus: " + strings.Join(failedValidationSkus.([]string), ",") + "."
		} else {
			errMsg = ""
		}
		failedToUploadScSku := uploadProdforSkuLevelSc(uploadedSkus, collection)

		if len(failedToUploadScSku) == 0 && errMsg == "" {
			return io, nil
		} else if len(failedToUploadScSku) == 0 && errMsg != "" {
			return io, &florest_constants.AppError{
				Code:             appconstant.FailedToCreateErrorCode,
				Message:          errMsg + "Rest are uploaded.",
				DeveloperMessage: sizeUtils.SkuNotUpdated,
			}
		} else if len(failedToUploadScSku) != 0 && errMsg == "" {
			return io, &florest_constants.AppError{
				Code: appconstant.FailedToCreateErrorCode,
				Message: "Sizechart is created but not added to skus " + strings.Join(failedToUploadScSku, ", ") +
					".For rest sizechart is created and added to skus as well.",
				DeveloperMessage: sizeUtils.SkuNotUpdated,
			}
		} else {
			return io, &florest_constants.AppError{
				Code: appconstant.FailedToCreateErrorCode,
				Message: errMsg + ".Sizechart not added to skus :" + strings.Join(failedToUploadScSku, ", ") +
					".Rest uploaded succesfully.",
				DeveloperMessage: sizeUtils.SkuNotUpdated,
			}
		}
	}
	// update the product for brand-brick level sizechart
	//for each sizechart doc update the skus associated with them
	flag := false
	notUpdated := []string{}
	updatedSkus := []string{}
	mysqlSizChArr := []sizeUtils.SizeChartForMysql{}

	for _, sizechartDoc := range collection.([]sizeUtils.SizeChartMongo) {
		total := 0
		res := getSkusForSizeChart(sizechartDoc)
		if len(res) == 0 {
			logger.Error(fmt.Sprintf("No skus found for sizechart with Id: %d", sizechartDoc.IdCatalogSizeChart))
			continue
		}

		for sku, preTy := range res {
			resp := UpdateProduct(sku, preTy, sizechartDoc)
			if !resp {
				notUpdated = append(notUpdated, sku)
				flag = true
			} else {
				updatedSkus = append(updatedSkus, sku)
				total = total + 1
			}
		}
		/*
			// Send notification to Datadog about successfully mapped skus.
			title := fmt.Sprintf("Brand-Brick type sizechart updated to Skus.")
			text := fmt.Sprintf("Skus updated for sizechart with ID (%d), Name: %s and sizechart-type: (%d) are: %s .", sizechartDoc.IdCatalogSizeChart,
				sizechartDoc.SizeChartName, sizechartDoc.SizeChartType, strings.Join(updatedSkus, ","))
			tags := []string{"sizechart", "brand-brick"}
			notification.SendNotification(title, text, tags, "info")
			// Send notification about skus failed to be mapped.
			if len(notUpdated) != 0 {
				titleF := fmt.Sprintf("Brand-Brick type sizechart not updated to Skus.")
				textF := fmt.Sprintf("Failed skus mapping for sizechart with ID (%d), Name: %s and sizechart-type: (%d) are: %s .", sizechartDoc.IdCatalogSizeChart,
					sizechartDoc.SizeChartName, sizechartDoc.SizeChartType, strings.Join(notUpdated, ","))
				tagsF := []string{"sizechart", "brand-brick"}
				notification.SendNotification(titleF, textF, tagsF, "info")
			}
		*/

		// create a mysqlsizech struct to be dumped into mysql
		sizChstrct := sizeUtils.SizeChartForMysql{
			updatedSkus,
			sizechartDoc.IdCatalogSizeChart,
			sizechartDoc.SizeChartType,
			sizechartDoc.FkCatalogCategory,
			sizechartDoc.CreatedAt,
		}
		mysqlSizChArr = append(mysqlSizChArr, sizChstrct)
	}
	// push the sizechart data to the worker
	sizechartProdPool.StartJob(mysqlSizChArr)
	if flag == true {
		return io, &florest_constants.AppError{
			Code:             appconstant.FailedToCreateErrorCode,
			Message:          "Few skus were not updated with new sizecharts: " + strings.Join(notUpdated, ", "),
			DeveloperMessage: sizeUtils.SkuNotUpdated,
		}
	}

	return io, nil
}

//This function upload the sku level sizechart to
//concerned skus
func uploadProdforSkuLevelSc(uploadedSkus interface{}, collection interface{}) []string {
	sizeCollec := collection.([]sizeUtils.SizeChartMongo)
	uploadedSkuArr := uploadedSkus.([]string)
	failedSku := []string{}
	mysqlSizChStrctArr := []sizeUtils.SizeChartForMysql{}

	for k, doc := range sizeCollec {

		// SKus have highest level priority
		resp := UpdateProduct(uploadedSkuArr[k], "", doc)
		if !resp {
			failedSku = append(failedSku, uploadedSkuArr[k])
		} else {
			// Store sku and sizechart mapping in collection
			err := storeSizechartSkuMapping(uploadedSkuArr[k], doc.IdCatalogSizeChart)
			if err != nil {
				logger.Error("#uploadProdforSkuLevelSc(): Unable to store sizechart mapping for sku %s : %s", uploadedSkuArr[k], err.Error())
			}
			// create sizechart struct to be dumped to mysql database
			sizChStrct := sizeUtils.SizeChartForMysql{
				[]string{uploadedSkuArr[k]},
				doc.IdCatalogSizeChart,
				doc.SizeChartType,
				doc.FkCatalogCategory,
				doc.CreatedAt,
			}
			mysqlSizChStrctArr = append(mysqlSizChStrctArr, sizChStrct)
		}
	}
	// Push the sizechart product data to worker
	sizechartProdPool.StartJob(mysqlSizChStrctArr)
	return failedSku
}
