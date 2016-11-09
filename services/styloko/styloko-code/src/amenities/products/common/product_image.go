package common

import (
	brands "amenities/services/brands"
	"errors"
	"strconv"
	"strings"
	"time"
)

type ProductImage struct {
	SeqId            int        `bson:"seqId" json:"seqId"`
	ImageNo          int        `bson:"imageNo" json:"imageNo"`
	Main             int        `bson:"main" json:"main"`
	OriginalFileName string     `bson:"originalFilename" json:"originalFilename"`
	ImageName        string     `bson:"imageName" json:"imageName"`
	Orientation      string     `bson:"orientation" json:"orientation"`
	UpdatedAt        *time.Time `bson:"updatedAt" json:"updatedAt"`
	IfUpdate         bool       `bson:"-" json:"ifUpdate"`
}

//
// Prepare Image name from supplies product data
//
func (pi *ProductImage) PrepareImageName(p Product) string {
	brand, _ := brands.ById(p.BrandId)
	configIdStr := strconv.Itoa(p.SeqId)
	var name string = configIdStr
	if brand.Name != "" {
		name = brand.Name
	}
	name += "-"
	name += p.Name
	name += "-"
	unixTimeStr := strconv.Itoa(int(pi.UpdatedAt.Unix()))
	runes := []rune(unixTimeStr)
	lenStr := len(runes)
	name += string(runes[lenStr-1]) +
		string(runes[lenStr-2]) +
		string(runes[lenStr-3]) +
		string(runes[lenStr-4])
	name += "-"
	if p.SeqId < 100 {
		name += RightPad2Len(Reverse(configIdStr), "0", 3)
	} else {
		name += Reverse(configIdStr)
	}
	name += "-"
	name = strings.Replace(name, " ", "-", -1)
	name = SanitizeImageName(name)
	return name
}

//
// Link Image to Product.
//
func (pi *ProductImage) Add(configId int, adapter string) (int, error) {
	_, err := GetAdapter(adapter).AddImage(configId, *pi)
	return pi.SeqId, err
}

//
// Create New Product Image.
//
func NewProductImage(adapter string) (ProductImage, error) {
	pi := ProductImage{}
	seqId, _ := GetAdapter(adapter).GenerateNextSequence(PIMAGE_COLLECTION)
	if seqId <= 0 {
		return pi, errors.New("Unable to Generate Sequence")
	}
	pi.SeqId = seqId
	now := time.Now()
	pi.UpdatedAt = &now
	return pi, nil
}
