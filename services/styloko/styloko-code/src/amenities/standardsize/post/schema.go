package post

type StandardSizeCreateInput struct {
	Brand          string `json:"brand"`
	AttributeSet   string `json:"attribute_set"`
	LeafCategory   string `json:"leaf_category"`
	BrandSize      string `json:"brand_size"`
	StandardSize   string `json:"standard_size"`
	Error          string `json:"error,omitempty"`
	AttributeSetID int    `json:"-"`
	LeafCategoryID int    `json:"-"`
	BrandID        int    `json:"-"`
}

type StandardSizeCreateReturn struct {
	SeqId int `json:"seqId"`
}
