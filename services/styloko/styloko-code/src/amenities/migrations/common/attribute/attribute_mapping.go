package attribute

import (
	"common/ResourceFactory"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/jabong/floRest/src/common/utils/logger"
)

// StartMapping starts
func StartMapping() error {
	logger.Info("Started attribute mapping migration")
	csvdata, err := readCsv()
	if err != nil {
		return err
	}

	color_map := AttributeMap{
		From:    "color",
		To:      "color_family",
		Mapping: make(map[string]string),
	}
	frame_color_map := AttributeMap{
		From:    "frame_color",
		To:      "frame_colors",
		Mapping: make(map[string]string),
	}
	frame_material_detail_map := AttributeMap{
		From:    "frame_material_detail",
		To:      "frame_material",
		Mapping: make(map[string]string),
	}
	lens_color_map := AttributeMap{
		From:    "lens_color",
		To:      "lens_colors",
		Mapping: make(map[string]string),
	}
	strap_color_map := AttributeMap{
		From:    "strap_color",
		To:      "strap_colors",
		Mapping: make(map[string]string),
	}
	strap_material_detail_map := AttributeMap{
		From:    "strap_material_detail",
		To:      "strap_material",
		Mapping: make(map[string]string),
	}
	upper_material_details_map := AttributeMap{
		From:    "upper_material_details",
		To:      "upper_material",
		Mapping: make(map[string]string),
	}
	fabric_details_map := AttributeMap{
		From:    "fabric_details",
		To:      "fabric",
		Mapping: make(map[string]string),
	}

	color := getOptionsFromMongo("color")
	frame_color := getOptionsFromMongo("frame_color")
	frame_material_detail := getOptionsFromMongo("frame_material_detail")
	lens_color := getOptionsFromMongo("lens_color")
	strap_color := getOptionsFromMongo("strap_color")
	strap_material_detail := getOptionsFromMongo("strap_material_detail")
	upper_material_details := getOptionsFromMongo("upper_material_details")
	fabric_details := getOptionsFromMongo("fabric_details")

	color_family := getOptionsFromMongo("color_family")
	frame_colors := getOptionsFromMongo("frame_colors")
	frame_material := getOptionsFromMongo("frame_material")
	lens_colors := getOptionsFromMongo("lens_colors")
	strap_colors := getOptionsFromMongo("strap_colors")
	strap_material := getOptionsFromMongo("strap_material")
	upper_material := getOptionsFromMongo("upper_material")
	fabric := getOptionsFromMongo("fabric")

	for _, x := range csvdata {
		switch x[0] {
		case "color":
			attribute := findInArray(color, x[1])
			filter := findInArray(color_family, x[2])
			if attribute == "" || filter == "" {
				logger.Error(fmt.Sprintf("Missing Attribute: %s, Filter: %s", x[1], x[2]))
				break
			}
			color_map.Mapping[attribute] = filter
			break
		case "frame_color":
			attribute := findInArray(frame_color, x[1])
			filter := findInArray(frame_colors, x[2])
			if attribute == "" || filter == "" {
				logger.Error(fmt.Sprintf("Missing Attribute: %s, Filter: %s", x[1], x[2]))
				break
			}
			frame_color_map.Mapping[attribute] = filter
			break
		case "frame_material_detail":
			attribute := findInArray(frame_material_detail, x[1])
			filter := findInArray(frame_material, x[2])
			if attribute == "" || filter == "" {
				logger.Error(fmt.Sprintf("Missing Attribute: %s, Filter: %s", x[1], x[2]))
				break
			}
			frame_material_detail_map.Mapping[attribute] = filter
			break
		case "lens_color":
			attribute := findInArray(lens_color, x[1])
			filter := findInArray(lens_colors, x[2])
			if attribute == "" || filter == "" {
				logger.Error(fmt.Sprintf("Missing Attribute: %s, Filter: %s", x[1], x[2]))
				break
			}
			lens_color_map.Mapping[attribute] = filter
			break
		case "strap_color":
			attribute := findInArray(strap_color, x[1])
			filter := findInArray(strap_colors, x[2])
			if attribute == "" || filter == "" {
				logger.Error(fmt.Sprintf("Missing Attribute: %s, Filter: %s", x[1], x[2]))
				break
			}
			strap_color_map.Mapping[attribute] = filter
			break
		case "strap_material_detail":
			attribute := findInArray(strap_material_detail, x[1])
			filter := findInArray(strap_material, x[2])
			if attribute == "" || filter == "" {
				logger.Error(fmt.Sprintf("Missing Attribute: %s, Filter: %s", x[1], x[2]))
				break
			}
			strap_material_detail_map.Mapping[attribute] = filter
			break
		case "upper_material_details":
			attribute := findInArray(upper_material_details, x[1])
			filter := findInArray(upper_material, x[2])
			if attribute == "" || filter == "" {
				logger.Error(fmt.Sprintf("Missing Attribute: %s, Filter: %s", x[1], x[2]))
				break
			}
			upper_material_details_map.Mapping[attribute] = filter
			break
		case "fabric_details":
			attribute := findInArray(fabric_details, x[1])
			filter := findInArray(fabric, x[2])
			if attribute == "" || filter == "" {
				logger.Error(fmt.Sprintf("Missing Attribute: %s, Filter: %s", x[1], x[2]))
				break
			}
			fabric_details_map.Mapping[attribute] = filter
			break

		}
	}
	insertIntoMongo(color_map, frame_color_map, frame_material_detail_map,
		lens_color_map,
		strap_color_map,
		strap_material_detail_map,
		upper_material_details_map,
		fabric_details_map)
	logger.Info("Started attribute mapping migration")
	return nil
}

func insertIntoMongo(attrMap ...AttributeMap) {
	mgoSession := ResourceFactory.GetMongoSession("AttributeMappings")
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection("attributeMapping")
	_ = mgoObj.DropCollection()
	for _, x := range attrMap {
		logger.Info(fmt.Sprintf("Inserting %s map into Mongo\n", x.From))
		mgoObj.Insert(x)
	}
}

func getOptionsFromMongo(attr string) []Option {
	mgoSession := ResourceFactory.GetMongoSession("AttributeMappings")
	defer mgoSession.Close()
	mgoObj := mgoSession.SetCollection("attributes")
	opt := new(AttributeRow)
	q := bson.M{"name": attr}
	mgoObj.Find(q).One(&opt)
	return opt.Options
}

func readCsv() ([]([]string), error) {
	data, err := ioutil.ReadFile("attribute_mappings.csv")
	var finalArr []([]string)
	if err != nil {
		logger.Error("Error occured while reading mappings csv", err.Error())
		return finalArr, err
	}
	d := string(data)
	dataArr := strings.Split(d, "\n")
	dataArr = dataArr[1 : len(dataArr)-1]
	for _, x := range dataArr {
		tmp := strings.Split(x, ";")
		tmp[0] = hyphenate(tmp[0])
		tmp[1] = regexpReplace(tmp[1])
		tmp[2] = regexpReplace(tmp[2])
		finalArr = append(finalArr, tmp)
	}
	return finalArr, nil
}

func hyphenate(data string) string {
	data = strings.ToLower(data)
	data = strings.Replace(data, " ", "_", -1)
	return data
}

func regexpReplace(data string) string {
	pattern := `[!@#$%^&*()-,.<>:;{}_=?/\"\ \+]+`
	reg, _ := regexp.Compile(pattern)
	safe := reg.ReplaceAllString(data, "")
	safe = strings.ToLower(safe)
	return safe
}

func findInArray(opt []Option, data string) string {
	for _, x := range opt {
		if data == regexpReplace(*x.Value) {
			return *x.Value
		}
	}
	return ""
}
