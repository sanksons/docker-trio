package json

import (
	"encoding/json"
)

//Converts the array of interface to array of string
func JsonInterfaceArrToStringArr(v []interface{}) []string {
	var varr []string = make([]string, len(v))
	for i := range v {
		s, _ := v[i].(string)
		varr[i] = s
	}
	return varr
}

//Gets the Map key-value from the input json
func GetMapFromJson(js string) (map[string]interface{}, error) {
	if len(js) == 0 {
		return nil, nil
	}
	var res map[string]interface{}
	err := json.Unmarshal([]byte(js), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetStringArrFromJson(js string) ([]string, error) {
	if len(js) == 0 {
		return nil, nil
	}
	var res []string
	err := json.Unmarshal([]byte(js), &res)
	if err != nil {
		return nil, err
	}

	return res, nil

}
