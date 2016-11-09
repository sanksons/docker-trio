package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func UpdateJsonPath(queries map[string]string, byt []byte, pathSep string) (newByt []byte, err error) {

	unMarshallObj := make(map[string]interface{})
	jerr := json.Unmarshal(byt, &unMarshallObj)
	if jerr != nil {
		return byt, jerr
	}

	for query, newNodeVal := range queries {
		path := strings.Split(query, pathSep)
		var v map[string]interface{}

		var jsPath map[string]interface{} = unMarshallObj
		for _, node := range path {

			nextJsPath, found := jsPath[node]
			if !found {
				return byt, errors.New(fmt.Sprintf("Not found node %s in path", node))
			} else {
				v = jsPath
				jsPath, _ = nextJsPath.(map[string]interface{})
			}

		}

		leafNode := path[len(path)-1]

		var newNodeValConv interface{}
		var convErr error

		switch v[leafNode].(type) {
		case float64:
			newNodeValConv, convErr = strconv.ParseFloat(newNodeVal, 64)
		case int64:
			newNodeValConv, convErr = strconv.ParseInt(newNodeVal, 10, 64)
		case uint64:
			newNodeValConv, convErr = strconv.ParseUint(newNodeVal, 10, 64)
		case string:
			newNodeValConv, convErr = newNodeVal, nil
		case bool:
			newNodeValConv, convErr = strconv.ParseBool(newNodeVal)
		default:
			newNodeValConv, convErr = nil, errors.New("Unsupported json value for json xpath update")
		}

		if convErr != nil {
			return byt, convErr
		}

		v[leafNode] = newNodeValConv
	}

	return json.Marshal(unMarshallObj)

}
