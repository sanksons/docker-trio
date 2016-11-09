package common

import (
	"errors"
	"time"
)

type ProductVideo struct {
	Id        int       `bson:"seqId" json:"seqId"`
	FileName  string    `bson:"fileName" json:"fileName"`
	Thumbnail string    `bson:"thumbnail" json:"thumbnail"`
	Size      int       `bson:"size" json:"size"`
	Duration  int       `bson:"duration" json:"duration"`
	Status    string    `bson:"status" json:"status"`
	Hash      string    `bson:"hash" json:"-"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

//
// Link Video to Product.
// May be update or insert
//
func (pv *ProductVideo) Save(configId int, adapter string) (int, error) {
	_, err := GetAdapter(adapter).SaveVideo(configId, *pv)
	return pv.Id, err
}

//
// Create New Product Video.
//
func NewProductVideo(adapter string) (ProductVideo, error) {
	pv := ProductVideo{}
	seqId, _ := GetAdapter(adapter).GenerateNextSequence(PVIDEO_COLLECTION)
	if seqId <= 0 {
		return pv, errors.New("Unable to Generate Sequence")
	}
	pv.Id = seqId
	pv.CreatedAt = time.Now()
	pv.UpdatedAt = time.Now()
	return pv, nil
}
