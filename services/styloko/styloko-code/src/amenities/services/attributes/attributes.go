package attributes

import (
	attrCommon "amenities/attributes/common"
	attrSet "amenities/attributes/get/search/set"
	factory "common/ResourceFactory"
	mongodb "common/mongodb"
)

func GetAllAttributes() ([]attrCommon.AttributeMongo, error) {
	session := factory.GetMongoSession("attributes")
	defer session.Close()

	attr := []attrCommon.AttributeMongo{}
	obj := session.SetCollection("attributes")
	err := obj.Find(nil).All(&attr)
	if err != nil {
		return attr, err
	}
	return attr, nil
}

func GetAttributeSetById(attrSetId int, mgo *mongodb.MongoDriver) attrSet.AttributeSet {
	var attributeSet attrSet.AttributeSet
	type M map[string]interface{}
	attributeSetObj := mgo.SetCollection("attributeSets")
	attributeSetObj.Find(M{"seqId": attrSetId}).One(&attributeSet)
	return attributeSet
}

// Deprecated
func GetDefaultValueAttributes(attributeSetId int) map[string]interface{} {
	global := map[string]interface{}{
		"192": "1",
		"193": "1",
		"194": "0",
		"196": "0",
		"197": "0",
		"215": "0",
		"217": "N",
		"229": "0",
		"230": "1",
		"251": "0",
		"358": "0",
		"359": "0",
		"404": "0",
		"458": "landscape",
	}
	switch attributeSetId {
	//we dont have any such case till now
	default:
		return global
	}

}

func GetByName(name string, mgoSession *mongodb.MongoDriver) attrSet.AttributeSet {
	var attributeSet attrSet.AttributeSet
	type M map[string]interface{}
	attributeSetObj := mgoSession.SetCollection("attributeSets")
	attributeSetObj.Find(M{"label": name}).One(&attributeSet)
	return attributeSet
}
