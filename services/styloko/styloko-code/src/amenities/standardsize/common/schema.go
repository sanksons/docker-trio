package common

import (
	"time"
)

type BrandStandardSize struct {
	BrandSize    string `bson:"brndSz"`
	StandardSize string `bson:"stdrdSz"`
}

type StandardSizeStore struct {
	SeqId          int                 `bson:"seqId"`
	BrandId        int                 `bson:"brndId"`
	AttributeSetId int                 `bson:"attrbtStId"`
	LeafCategoryId int                 `bson:"lfCtgryId"`
	Size           []BrandStandardSize `bson:"size"`
}

type StandardSizeResult struct {
	SeqId        int
	Brand        string
	AttributeSet string
	LeafCategory string
	Size         []BrandStandardSize
}

type FinalResult struct {
	Count        int
	SizesMapping []StandardSizeResult
}

type StandardSizeError struct {
	SeqId          int       `bson:"seqId"`
	AttributeSetId int       `bson:"attrbtStId"`
	LeafCategoryId int       `bson:"lfCtgryId"`
	BrandId        int       `bson:"brndId"`
	BrandSize      string    `bson:"brndSz"`
	CreatedAt      time.Time `bson:"crtdAt"`
	Fixed          bool      `bson:"fixed"`
}
