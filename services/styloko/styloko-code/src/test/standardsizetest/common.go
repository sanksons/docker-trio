package standardsizetest

import (
	"math/rand"
	"time"
)

const (
	brandLength = 4
	sizeLength  = 5
)

var startValidData = `{
    "attribute_set": "Men Apparel",
    "leaf_category": "3486",
    "brand": "",
    "brand_size": "`

var midValidData = `",
    "standard_size": "`

var endValidData = `"
}`

var invalidAttributeSet = `{
    "attribute_set": "Men Apparel1",
    "leaf_category": "3488",
    "brand": "Ocean",
    "brand_size": "30",
    "standard_size": "S"
}`

var invalidLeafCategory = `{
    "attribute_set": "Men Apparel1",
    "leaf_category": "abc",
    "brand": "Ocean",
    "brand_size": "30",
    "standard_size": "S"
}`

var invalidBrand = `{
    "attribute_set": "Men Apparel1",
    "leaf_category": "3488",
    "brand": "Ocean1",
    "brand_size": "30",
    "standard_size": "S"
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
	brandSize := randomStringGenerator(brandLength)
	standardSize := randomStringGenerator(sizeLength)

	switch flag {
	case 1:
		result = ""
		break
	case 2:
		result = invalidAttributeSet
		break
	case 3:
		result = invalidLeafCategory
		break
	case 4:
		result = invalidBrand
		break
	case 5:
		result = startValidData + midValidData + standardSize + endValidData
		break
	case 6:
		result = startValidData + brandSize + midValidData + standardSize + endValidData
		break
	}

	return result
}
