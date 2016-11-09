package common

import (
	"errors"
)

type AttributeSet struct {
	SeqId int    `bson:"seqId"`
	Name  string `bson:"name"`
}

type ProAttributeSet struct {
	Id    int    `bson:"seqId" json:"seqId"`
	Name  string `bson:"name" json:"name"`
	Label string `bson:"label" json:"label"`
}

func (set ProAttributeSet) GetVariationAttributeName() (string, error) {
	mapping := GetAttributeSet2VariationMapping()
	if val, ok := mapping[set.Name]; ok {
		return val, nil
	}
	return "", errors.New("(set ProAttributeSet)#getVariationAttributeName(): Invalid Index")
}

// Get size key for attributeset
//Mapping for attributeset and size keys
func GetAttributeSet2VariationMapping() map[string]string {
	return map[string]string{
		"bags":             "variations",
		"beauty":           "variations",
		"fragrances":       "variations",
		"home":             "variations",
		"sports_equipment": "variations",
		"toys":             "variations",
		"shoes":            "sh_size",
		"app_men":          "apm_size",
		"app_women":        "apw_size",
		"app_kids":         "apk_size",
		"jewellery":        "variation",
	}
}
