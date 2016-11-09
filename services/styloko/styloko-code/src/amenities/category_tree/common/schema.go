package common

// CategoryGetVerbose => Returns a verbose category by ID.
type CategoryGetVerbose struct {
	CategoryId            int               `bson:"seqId" json:"seqId"`
	Status                string            `bson:"status" json:"status"`
	Left                  int               `bson:"lft" json:"left"`
	Right                 int               `bson:"rgt" json:"right"`
	Name                  string            `bson:"name" json:"name"`
	Parent                int               `bson:"parent,omitempty" json:"parent"`
	UrlKey                string            `bson:"urlKey,omitempty" json:"urlKey"`
	SizeChartAcive        int               `bson:"szchrtActv,omitempty" json:"sizechartActive" accepts:"[0,1]"`
	PdfActive             int               `bson:"pdfActv,omitempty" json:"pdfActive" accepts:"[0,1]"`
	DisplaySizeConversion string            `bson:"dispSzConv,omitempty" json:"displaySizeConversion"`
	GoogleTreeMapping     string            `bson:"gleTreeMapping,omitempty" json:"googleTreeMapping"`
	SizeChartApplicable   int               `bson:"szchrtApp,omitempty" json:"sizeChartApplicable"`
	CategorySeg           []CategorySegment `bson:"segment,omitempty" json:"segment"`
}

// CategorySegment => Category Segment struct
type CategorySegment struct {
	IdCategorySegment int    `bson:"seqId" json:"seqId" required:"true"`
	Name              string `bson:"name,omitempty"`
	UrlKey            string `bson:"urlKey,omitempty"`
	Genders           string `bson:"genders,omitempty"`
}
