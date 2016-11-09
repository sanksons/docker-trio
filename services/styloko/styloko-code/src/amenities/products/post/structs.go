package post

import (
	proUtil "amenities/products/common"
	_ "fmt"
	"github.com/jabong/floRest/src/common/constants"
)

type ProIOData struct {
	Product     *proUtil.Product
	ReqData     *ProductCreateRequestData
	Error       *constants.AppError
	Status      string
	IsDuplicate bool
}

func (data *ProIOData) setSuccess(p *proUtil.Product) {
	data.Status = proUtil.STATUS_SUCCESS
	data.Product = p
	data.Error = nil
	return
}

func (data *ProIOData) setFailure(err constants.AppError) {
	data.Status = proUtil.STATUS_FAILURE
	data.Error = &err
	return
}

type ProductCreateRequestData struct {
	ApprovalStatus         *int                    `json:"approvalStatus" validate:"-"`
	AttributeSet           int                     `json:"attributeSet" validate:"required"`
	Attributes             proUtil.AttributeMapSet `json:"attributes" validate:"required"`
	Brand                  int                     `json:"brand" validate:"required"`
	Categories             []int                   `json:"categories" validate:"required"`
	PrimaryCategory        int                     `json:"primaryCategory" validate:"-"`
	ConfigId               int                     `json:"configId" validate:"omitempty"`
	Description            string                  `json:"description" validate:"required"`
	PlatformIdentifier     string                  `json:"platformIdentifier" validate:"required"`
	Price                  float64                 `json:"price" validate:"required"`
	EanCode                string                  `json:"productIdentifier"`
	GroupName              string                  `json:"groupName"`
	ProductSet             int                     `json:"productSet" validate:"required"`
	Seller                 int                     `json:"seller" validate:"required"`
	SellerSKU              string                  `json:"sellerSku" validate:"required"`
	ShipmentType           int                     `json:"shipmentType" validate:"required"`
	SpecialPrice           *float64                `json:"specialPrice" validate:"omitempty,ltfield=Price"`
	SpecialFromDate        *string                 `json:"specialFromDate"`
	SpecialToDate          *string                 `json:"specialToDate"`
	Status                 string                  `json:"status" validate:"required"`
	TaxClass               int                     `json:"taxClass" validate:"required"`
	Title                  string                  `json:"title" validate:"required"`
	IdCatalogProduct       int                     `json:"idCatalogProduct"`
	ShipmentMatrixTemplate interface{}             `json:"shipmentMatrixTemplate"`
}
