package brandstest

// import (
// 	"strconv"
// )

var emptyPostData = `{
	"brandDataList": [
		{
        }
	]
}`

var validPostData = `{
	"brandDataList": [
		{
            "status": "active",
            "name": "Haha",
            "position":1,
            "urlKey": "haha",
            "imgName":"haha.jpg",
            "brndClass":"regular",
            "isExc":0,
            "brandInfo":"Haha is ainvi brand"
        }
	]
}`

var invalidPostData = `{
	"brandDataList": [
		{
            "status": "InvalidStatus",
            "name": "Haha",
            "position":1,
            "urlKey": "haha",
            "imgName":"haha.jpg",
            "brandInfo":"Haha is ainvi brand"
            "isExc":0			
        }
	]
}`

var missingPostData = `{
	"brandDataList": [
		{
            "status": "active",
            "name": "Haha",
            "position":1,
            "urlKey": "haha",
            "imgName":"haha.jpg",
            "brandInfo":"Haha is ainvi brand"			
        }
	]
}`

var validPutData = `{
    "brandDataList": [
        {
        	"seqId":2671,
            "status": "deleted"
        }
    ]
}`

var invalidPutData = `{
    "brandDataList": [
        {
            "status": "deleted"
        }
    ]
}`

func getPostBodyData(flag int) (result string) {
	switch flag {
	case 1:
		result = validPostData
		break
	case 2:
		result = ""
		break
	case 3:
		result = missingPostData
		break

	case 4:
		result = invalidPostData
		break
	}
	return result
}

func getPutBodyData(flag int) (result string) {
	switch flag {
	case 1:
		result = validPutData
		break
	case 2:
		result = invalidPutData
		break
	case 3:
		result = ""
		break
	}
	return result
}
