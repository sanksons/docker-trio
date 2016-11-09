package sellercenter

import (
	proUtil "amenities/products/common"
	utils "common/utils"
	"fmt"
	"time"

	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type ResponseNodeSC struct {
	id string
}

func (cs *ResponseNodeSC) SetID(id string) {
	cs.id = id
}

func (cs ResponseNodeSC) GetID() (id string, err error) {
	return cs.id, nil
}

func (cs ResponseNodeSC) Name() string {
	return "ResponseNodeSC"
}

func (cs ResponseNodeSC) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {

	logger.Debug("Enter Response node SC")
	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.GET_SC_RESPONSE_NODE)
	defer logger.EndProfile(profiler, proUtil.GET_SC_RESPONSE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Response node sc")

	ssr, err := FetchDataNodeSC{}.GetSSR(io)
	if err != nil {
		return io, &constants.AppError{
			Code:    constants.ResourceErrorCode,
			Message: err.Error(),
		}
	}

	if ssr.ResetCounter {
		io.IOData.Set(constants.RESULT, "Counter reset")
		return io, nil
	}

	d, _ := io.IOData.Get(proUtil.IODATA)
	products, ok := d.([]proUtil.Product)
	if !ok {
		return io, &constants.AppError{
			Code:    constants.ResourceErrorCode,
			Message: "(cs ResponseNodeSC) Execute(): Invalid IO Data"}
	}
	io.IOData.Set(constants.RESULT, cs.PrepareResponse(products))
	logger.Debug("Exit Response node SC")
	return io, nil
}

//
// Prepare response for the user.
//
func (cs ResponseNodeSC) PrepareResponse(
	products []proUtil.Product,
) []*SellerCenterResponse {

	var response []*SellerCenterResponse
	for _, product := range products {
		scr := &SellerCenterResponse{}
		scr.Prepare(product)
		response = append(response, scr)
	}
	return response
}

//
// Define and setup response pattern for seller center
//
type SellerCenterResponse struct {
	SellerId                int                      `json:"sellerId"`
	ProductSetId            int                      `json:"productSetId"`
	Status                  string                   `json:"status"`
	Name                    string                   `json:"name"`
	Description             string                   `json:"description"`
	PrimaryCategoryId       int                      `json:"primaryCategoryId"`
	AdditionalCategoriesIds []int                    `json:"additionalCategoriesIds"`
	BrandId                 int                      `json:"brandId"`
	AttributeSetId          int                      `json:"attributeSetId"`
	Attributes              []*SellerCenterAttribute `json:"attributes"`
	TaxClassId              int                      `json:"taxClassId"`
	ShipmentType            string                   `json:"shipmentType"`
	ApprovalStatus          string                   `json:"approvalStatus"`
	UpdatedAt               string                   `json:"updatedAt"`
	Images                  []*SellerCenterImage     `json:"images"`
	Simples                 []*SellerCenterSimple    `json:"variations"`
}

type SellerCenterAttribute struct {
	AttributeId int         `json:"attributeId"`
	Value       interface{} `json:"value"`
}

type SellerCenterImage struct {
	ImageId    int    `json:"imageId"`
	DisplayUrl string `json:"displayUrl"`
	UpdatedAt  string `json:"updatedAt"`
	Position   int    `json:"position"`
	IsMain     bool   `json:"isMain"`
}

type SellerCenterSimple struct {
	Id                   int         `json:"productId"`
	SKU                  string      `json:"simpleSku"`
	SellerSku            string      `json:"sellerSku"`
	ProductIdentifier    string      `json:"productIdentifier"`
	Variation            interface{} `json:"variation"`
	Price                float64     `json:"price"`
	SpecialPrice         float64     `json:"specialPrice,omitempty"`
	SpecialPriceFromDate string      `json:"specialPriceFromDate,omitempty"`
	SpecialPriceToDate   string      `json:"specialPriceToDate,omitempty"`
	Stock                int         `json:"stock"`
}

func (scr *SellerCenterResponse) SetAttributes(product proUtil.Product) error {
	var scrAttributes []*SellerCenterAttribute
	for _, attribute := range product.Attributes {
		scrAttr := &SellerCenterAttribute{}
		var err error
		scrAttr.AttributeId = attribute.Id
		scrAttr.Value, err = attribute.GetValue("value")
		if err != nil {
			logger.Error("(scr *SellerCenterResponse)#SetAttributes:" + err.Error())
			continue
		}
		if attribute.OptionType == proUtil.OPTION_TYPE_VALUE {
			scrAttr.Value = utils.ToString(scrAttr.Value)
		}
		scrAttributes = append(scrAttributes, scrAttr)
	}
	for _, attribute := range product.Global {
		scrAttr := &SellerCenterAttribute{}
		var err error
		scrAttr.AttributeId = attribute.Id
		scrAttr.Value, err = attribute.GetValue("value")
		if err != nil {
			logger.Error("(scr *SellerCenterResponse)#SetAttributes:" + err.Error())
			continue
		}
		if attribute.OptionType == proUtil.OPTION_TYPE_VALUE {
			scrAttr.Value = utils.ToString(scrAttr.Value)
		}
		scrAttributes = append(scrAttributes, scrAttr)
	}
	scr.Attributes = scrAttributes
	return nil
}

func (scr *SellerCenterResponse) SetImages(product proUtil.Product) error {
	var scrImages []*SellerCenterImage
	var flag = true
	for _, im := range product.Images {
		scrIm := &SellerCenterImage{}
		scrIm.ImageId = im.SeqId
		var isMain = false
		if im.Main == 1 && flag {
			isMain = true
			flag = false
		}
		scrIm.IsMain = isMain
		scrIm.Position = im.ImageNo
		scrIm.DisplayUrl = fmt.Sprintf("/p/%s%d.jpg", im.ImageName, im.ImageNo)
		scrIm.UpdatedAt = proUtil.ToMySqlTime(im.UpdatedAt)
		scrImages = append(scrImages, scrIm)
	}
	scr.Images = scrImages
	return nil
}

func (scr *SellerCenterResponse) SetCategories(product proUtil.Product) error {

	scr.PrimaryCategoryId = product.PrimaryCategory
	var secondaryCats []int
	for _, v := range product.Categories {
		if v != scr.PrimaryCategoryId {
			secondaryCats = append(secondaryCats, v)
		}
	}
	scr.AdditionalCategoriesIds = secondaryCats
	return nil
}

func (scr *SellerCenterResponse) SetShipmentType(product proUtil.Product) error {
	shipment := proUtil.GetShipmentById(product.ShipmentType)
	scr.ShipmentType = shipment
	return nil
}

func (scr *SellerCenterResponse) SetSimples(product proUtil.Product) error {
	var scrSimples []*SellerCenterSimple
	varOptionName, _ := product.AttributeSet.GetVariationAttributeName()
	for _, simple := range product.Simples {
		scrSimple := &SellerCenterSimple{}
		scrSimple.Id = simple.Id
		scrSimple.SKU = simple.SKU
		if simple.SellerSKU != "" {
			scrSimple.SellerSku = simple.SellerSKU
		} else {
			scrSimple.SellerSku = simple.SKU
		}

		scrSimple.ProductIdentifier = simple.EanCode
		if val, ok := simple.Attributes[utils.SnakeToCamel(varOptionName)]; ok {
			var err error
			scrSimple.Variation, err = val.GetValue("value")
			if err != nil {
				logger.Error(err)
			}
		}
		scrSimple.Price = *simple.Price
		var specialPrice float64
		if (simple.SpecialPrice != nil) &&
			(*simple.SpecialPrice < *simple.Price) &&
			((simple.SpecialToDate != nil) || (simple.SpecialFromDate != nil)) {
			specialPrice = *simple.SpecialPrice
		}
		scrSimple.SpecialPrice = specialPrice

		//case where special todate is less than from date
		if (simple.SpecialFromDate != nil) &&
			(simple.SpecialToDate != nil) {
			frmdate := *simple.SpecialFromDate
			todate := *simple.SpecialToDate
			if frmdate.After(todate) {
				simple.SpecialFromDate = simple.SpecialToDate
			}
		}
		scrSimple.SpecialPriceFromDate = proUtil.ToMySqlTime(
			simple.SpecialFromDate)
		scrSimple.SpecialPriceToDate = proUtil.ToMySqlTime(simple.SpecialToDate)

		//if FromDate equals ToDate, increment to date via a day
		if (scrSimple.SpecialPriceFromDate != "") &&
			(scrSimple.SpecialPriceToDate != "") &&
			(scrSimple.SpecialPriceFromDate == scrSimple.SpecialPriceToDate) {
			fromTime, err := proUtil.FromMysqlTime(scrSimple.SpecialPriceToDate, true)
			if err != nil {
				logger.Error(err.Error())
			}
			toTime := fromTime.Add(24 * 58 * time.Minute)
			scrSimple.SpecialPriceToDate = proUtil.ToMySqlTime(&toTime)
		}
		scrSimple.Stock = simple.GetQuantity()
		scrSimples = append(scrSimples, scrSimple)
	}
	scr.Simples = scrSimples
	return nil
}

func (scr *SellerCenterResponse) Prepare(product proUtil.Product) error {

	if product.ProductSet > 0 {
		scr.ProductSetId = product.ProductSet
	} else {
		scr.ProductSetId = product.SeqId
	}
	scr.SellerId = product.SellerId
	scr.Name = product.Name
	scr.Description = product.Description
	scr.BrandId = product.BrandId
	scr.ApprovalStatus = proUtil.STATUS_PENDING
	if product.PetApproved == 1 {
		scr.ApprovalStatus = proUtil.STATUS_APPROVED
	}

	if product.Name == "" {
		scr.Name = product.SKU
	}

	scr.PrimaryCategoryId = product.GetPrimaryCategory(proUtil.DB_ADAPTER_MONGO)
	var addCatIds []int
	for _, c := range product.Categories {
		if c == scr.PrimaryCategoryId {
			continue
		}
		addCatIds = append(addCatIds, c)
	}

	scr.AdditionalCategoriesIds = addCatIds
	scr.UpdatedAt = proUtil.ToMySqlTime(product.UpdatedAt)
	scr.Status = product.Status
	scr.Description = product.Description
	scr.AttributeSetId = product.AttributeSet.Id
	//As discussed, with apporva taxClassId will always be sent as 1.
	scr.TaxClassId = 1
	//prepare attributes
	scr.SetAttributes(product)
	//prepare Images
	scr.SetImages(product)

	// PATCH:
	// Add some default image incase images is blank
	// This patch needs to be removed once done with recreate on live.
	var currTime = time.Now()
	if len(scr.Images) == 0 {
		sequence, _ := proUtil.GetAdapter(proUtil.DB_READ_ADAPTER).GenerateNextSequence(proUtil.DUMMY_IMAGES)
		if sequence > 0 {
			scr.Images = []*SellerCenterImage{
				&SellerCenterImage{
					ImageId:    sequence,
					DisplayUrl: fmt.Sprintf("/p/%s%d.jpg", "dummyImage", 1),
					UpdatedAt:  proUtil.ToMySqlTime(&currTime),
					Position:   1,
					IsMain:     true,
				},
			}
		}
	}

	//prepare shipment type
	scr.SetShipmentType(product)
	// prepare primary category and secondary categories.
	scr.SetCategories(product)
	//prepare simples
	scr.SetSimples(product)
	return nil
}
