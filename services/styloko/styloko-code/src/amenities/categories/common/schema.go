package common

// CategoryGet =>get all categories struct
type CategoryGet struct {
	CategoryId int    `bson:"seqId" json:"seqId"`
	Status     string `bson:"status" json:"status"`
	Left       int    `bson:"lft" json:"left"`
	Right      int    `bson:"rgt" json:"right"`
	Name       string `bson:"name" json:"name"`
	NameEn     string `json:"name_en"`
}

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

// CategoryCreate => Create new category struct. Temporary.
type CategoryCreate struct {
	Id                    int               `bson:"seqId"`
	Status                string            `bson:"status" json:"status,omitempty" required:"true" sql:"status,omitempty"`
	Name                  string            `bson:"name" json:"name,omitempty" required:"true" sql:"name,omitempty"`
	Left                  int               `bson:"lft" sql:"lft,omitempty"`
	Right                 int               `bson:"rgt" sql:"rgt,omitempty"`
	Parent                int               `bson:"parent,omitempty" json:"parent,omitempty" required:"true"`
	UrlKey                string            `bson:"urlKey,omitempty" json:"urlKey,omitempty" required:"true" sql:"url_key,omitempty"`
	SizeChartAcive        int               `bson:"szchrtActv,omitempty" json:"sizechartActive,omitempty" accepts:"[0,1]" sql:"sizechart_active"`
	PdfActive             int               `bson:"pdfActv,omitempty" json:"pdfActive,omitempty" accepts:"[0,1]" sql:"pdf_active"`
	DisplaySizeConversion string            `bson:"dispSzConv,omitempty" json:"displaySizeConversion,omitempty" sql:"display_size_conversion,omitempty"`
	GoogleTreeMapping     string            `bson:"gleTreeMapping,omitempty" json:"googleTreeMapping,omitempty" sql:"google_tree_mapping,omitempty"`
	SizeChartApplicable   int               `bson:"szchrtApp,omitempty" json:"sizeChartApplicable,omitempty" accepts:"[0,1]" sql:"sizechart_applicable,omitempty"`
	CategorySeg           []CategorySegment `bson:"segment,omitempty"`
	SegIds                []int             `json:"segIds,omitempty"`
}

// CategoryUpdate => Update category struct.
type CategoryUpdate struct {
	Status                string            `bson:"status,omitempty" json:"status,omitempty" sql:"status,omitempty"`
	Name                  string            `bson:"name,omitempty" json:"name,omitempty" sql:"name,omitempty"`
	Parent                int               `bson:"parent,omitempty"`
	UrlKey                string            `bson:"urlKey,omitempty" json:"urlKey,omitempty" sql:"url_key,omitempty"`
	SizeChartAcive        int               `bson:"szchrtActv" json:"sizechartActive,omitempty" accepts:"[0,1]" sql:"sizechart_active"`
	PdfActive             int               `bson:"pdfActv" json:"pdfActive,omitempty" accepts:"[0,1]" sql:"pdf_active"`
	DisplaySizeConversion string            `bson:"dispSzConv,omitempty" json:"displaySizeConversion,omitempty" sql:"display_size_conversion,omitempty"`
	GoogleTreeMapping     string            `bson:"gleTreeMapping,omitempty" json:"googleTreeMapping,omitempty" sql:"google_tree_mapping,omitempty"`
	SizeChartApplicable   int               `bson:"szchrtApp" json:"sizeChartApplicable,omitempty" accepts:"[0,1]" sql:"sizechart_applicable"`
	CategorySeg           []CategorySegment `bson:"segment,omitempty"`
	SegIds                []int             `json:"segIds,omitempty"`
}

// CategorySegment => Category Segment struct
type CategorySegment struct {
	IdCategorySegment int    `bson:"seqId" json:"seqId" required:"true"`
	Name              string `bson:"name,omitempty"`
	UrlKey            string `bson:"urlKey,omitempty"`
	Genders           string `bson:"genders,omitempty"`
}
