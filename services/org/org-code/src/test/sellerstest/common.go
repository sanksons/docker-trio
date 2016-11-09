package sellerstest

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	emailLength = 10
)

var startPostData = `{
    "sellerDataList": [
        {
            "slrId": "`

var midPostDataValid = `",
            "orgName": "ghfv",
            "slrName": "Sundaram Jewels",
            "status": "active",
            "addr1": "Sundaram Jewels",
            "city": "delhi",
            "pstcode": 122011,
            "cntryCode": "IN",
            "ordrEml": "`

var midPostDataInvalid = `",
            "orgName": "",
            "slrName": "Sundaram Jewels",
            "status": "active",
            "addr2": "Sundaram Jewels",
            "city": "delhi",
            "pstcode": 122011,
            "cntryCode": "IN",
            "ordrEml": "`

var endPostData = `",
            "cntctName": "sundaram",
            "phn": "00000"
        }
    ]
}`

var validPutData = `{
    "sellerDataList": [
        {
            "seqId": 2,
            "orgName": "ghfv",
            "slrName": "Sundaram Jewels",
            "status": "active",
            "addr2": "Sundaram Jewels",
            "city": "delhi",
            "pstcode": 122011,
            "cntryCode": "IN",
            "ordrEml": "s@abc.com",
            "cntctName": "sundaram",
            "phn": "00000"
        },
        {
            "seqId": 3,
            "orgName": "Sunwels",
            "slrName": "Sunls",
            "status": "active",
            "addr2": "Sundaram Jewels",
            "city": "delhi",
            "pstcode": 122011,
            "cntryCode": "IN",
            "ordrEml": "s@abc.com",
            "cntctName": "sundaram",
            "phn": "00000"
        }
    ]
}`

var invalidPutData = `{
    "sellerDataList": [
        {
            "seqId": 0,
            "orgName": "ghfv",
            "slrName": "Sundaram Jewels",
            "status": "active",
            "addr2": "Sundaram Jewels",
            "city": "delhi",
            "pstcode": 122011,
            "cntryCode": "IN",
            "ordrEml": "s@abc.com",
            "cntctName": "sundaram",
            "phn": "00000"
        },
        {
           "seqId": 3,
            "orgName": "Sunwels",
            "slrName": "Sunls",
            "status": "active",
            "addr2": "Sundaram Jewels",
            "city": "delhi",
            "pstcode": 122011,
            "cntryCode": "IN",
            "ordrEml": "s@abc.com",
            "cntctName": "sundaram",
            "phn": "00000"
        }
    ]
}`

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomStringGenerator(n int) string {
	source := rand.NewSource(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		seed := rand.New(source)
		randNum := seed.Int63()
		b[i] = letterRunes[randNum%int64(len(letterRunes))]
	}
	return string(b)
}

func getPostBodyData(flag int) (result string) {
	source := rand.NewSource(time.Now().UnixNano())
	seed := rand.New(source)
	sellerId := seed.Int63()
	sellerIdString := strconv.FormatInt(sellerId, 10)

	emailId := randomStringGenerator(emailLength)

	switch flag {
	case 1:
		result = startPostData + sellerIdString + midPostDataValid + emailId + endPostData
		break
	case 2:
		result = startPostData + sellerIdString + midPostDataInvalid + emailId + endPostData
		break
	case 3:
		result = startPostData + midPostDataValid + emailId + endPostData
		break
	case 4:
		result = ""
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
