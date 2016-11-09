package productstest

func getPostBodyData(flag int) (result string) {
	switch flag {
	case 1:
		result = postDataValid
		break
		// case 2:
		// 	result = ""
		// 	break
		// case 3:
		// 	result = missingPostData
		// 	break

		// case 4:
		// 	result = invalidPostData
		// 	break
	}
	return result
}

func getPutBodyDataValid(flag int) (result string) {
	switch flag {
	case 1:
		result = putProductDataValid
		break
	case 2:
		result = putImageAddDataValid
		break
	case 3:
		result = putImageDelDataValid
		break
	case 4:
		result = putNodeDataValid
		break
	case 5:
		result = putPriceDataValid
		break
	case 6:
		result = putSellerDeactivateDataValid
		break
	case 7:
		result = putShipmentDataValid
		break
	case 8:
		result = putVideoStatusDataValid
		break
	case 9:
		result = putAttributeDataValid
		break
	case 10:
		result = putJabongDiscountDataValid
		break
	}
	return result
}

func getPutHeader(flag int) (result map[string]string) {
	putHeaders := make(map[string]string, 0)
	switch flag {
	case 1:
		putHeaders["Update-Type"] = "Product"
		break
	case 2:
		putHeaders["Update-Type"] = "ImageAdd"
		break
	case 3:
		putHeaders["Update-Type"] = "ImageDel"
		break
	case 4:
		putHeaders["Update-Type"] = "Node"
		break
	case 5:
		putHeaders["Update-Type"] = "Price"
		break
	case 6:
		putHeaders["Update-Type"] = "SellerDeactivate"
		break
	case 7:
		putHeaders["Update-Type"] = "Shipment"
		break
	case 8:
		putHeaders["Update-Type"] = "VideoStatus"
		break
	case 9:
		putHeaders["Update-Type"] = "Attribute"
		break
	case 10:
		putHeaders["Update-TYpe"] = "JabongDiscount"
		break
	}
	return putHeaders
}

var postDataValid = `[
			{
		        "approvalStatus": 1,
		        "attributeSet": 1,
		        "attributes": {
		            "4": "active",
		            "33": "1",
		            "39": "10",
		            "50": "1",
		            "58": "...",
		            "77": "Red",
		            "80": "Canon",
		            "140": "4",
		            "142": "0.00",
		            "143": "0.00",
		            "236": "active",
		            "248": "Canon",
		            "249": "Canon Canon"
		        },
		        "brand": 11,
		        "categories": [
		            1264,
		            1257,
		            1199,
		            490,
		            1
		        ],
		        "configId": null,
		        "description": "Canon Canon Canon Canon",
		        "platformIdentifier": "SellerCenter",
		        "price": 100,
		        "productIdentifier": "Canon",
		        "productSet": 3,
		        "seller": 624,
		        "sellerSku": "Canon",
		        "shipmentType": 2,
		        "specialFromDate": null,
		        "specialPrice": null,
		        "specialToDate": null,
		        "status": "active",
		        "stock": 5000,
		        "taxClass": 1,
		        "title": "Canon",
		        "idcatalogProduct": 174301
		    }
		]`

var putProductDataValid = `[
		    {
		        "productId": 6927966,
		        "sellerSku": "Canon12",
		        "sku": "AN011SH14YYZINDFAS",
		        "approvalStatus": 1,
		        "status": "active",
		        "title": "Canon12",
		        "brand": 11,
		        "attributeSet": 1,
		        "shipmentType": 2,
		        "description": "Canon12",
		        "price": 1012,
		        "specialFromDate": null,
		        "specialPrice": null,
		        "specialToDate": null,
		        "taxClass": 1,
		        "categories": [
		            1264,
		            1257,
		            1199,
		            490,
		            1
		        ],
		        "attributes": {
		            "4": "active",
		            "32": "8",
		            "33": "1",
		            "39": "10",
		            "50": "1",
		            "58": "...",
		            "77": "Red",
		            "80": "Canon",
		            "140": "4",
		            "142": "0.00",
		            "143": "0.00",
		            "236": "active",
		            "248": "Canon",
		            "249": "Canon Canon"
		        },
		        "productSet": 3,
		        "configId": 667185
		    }
		]`

var putImageAddDataValid = `[
		{
		  "productId":73,
		  "originalFilename":"test.jpg",
		  "isMain" :true,
		  "imageNo" : 1,
		  "orientation" : "portrait"
		}
	]`

var putImageDelDataValid = `{
	  "imageIds": [67654,25563272,111213]
	}`

var putNodeDataValid = `[
		{
		"type":"deleteNode",
		"nodeName":"SizeChart",
		"sku": "NI091SH99ZUE",
		"data": null
		}
	]`

var putPriceDataValid = `[
		 {
		  "simpleId":113878,
		  "configId":35000,
		  "price":300,
		  "specialFromDate":"2012-01-09 16:15:26",
		  "specialToDate":"2012-01-09 16:15:26",
		  "specialPrice":300
		 }
		]`

var putSellerDeactivateDataValid = `[{"sellerId":110}]`

var putShipmentDataValid = `[
			{
			"skuSimple":"NI091SH99ZUE-113878",
			"simpleId": 113878,
			"shipmentType": 2
			}
		]`

var putVideoStatusDataValid = `[
		    {
		        "videoId": 8087,
		        "status": "deleted"
		    },
		    {
		        "videoId": 8244,
		        "status": "active"
		    }
		]`

var putAttributeDataValid = `[
		    {
		        "attributeName": "cost",
		        "isGlobal": true,
		        "productSku": "NI091SH68ABF-196746",
		        "productType": "simple",
		        "action": 2,
		        "value": 2012.321,
		        "petApproved": 0
		    }
		] `

var putJabongDiscountDataValid = ` [
		    {
		        "productId": 73754,
		        "jabongDiscount": 191.88,
		        "jabongDiscountFromDate": "2013-12-29 18:30:00",
		        "jabongDiscountToDate": "2015-12-29 18:30:00"
		    }
		]`
