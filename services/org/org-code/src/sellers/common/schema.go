package common

import (
	"time"
)

type Schema struct {
	SeqId                int          `bson:"seqId,omitempty" json:"seqId,omitempty"`
	SellerId             string       `bson:"slrId,omitempty" json:"slrId"`
	VendorId             string       `bson:"-" json:"-"`
	OrgName              string       `bson:"orgName,omitempty" json:"orgName,omitempty"`
	SellerName           string       `bson:"slrName,omitempty" json:"slrName,omitempty"`
	Status               string       `bson:"status,omitempty" json:"status,omitempty"`
	OrderEmail           string       `bson:"ordrEml,omitempty" json:"ordrEml,omitempty"`
	Contact              string       `bson:"cntct,omitempty" json:"cntct,omitempty"`
	ContactName          string       `bson:"cntctName,omitempty" json:"cntctName,omitempty"`
	Phone                string       `bson:"phn,omitempty" json:"phn,omitempty"`
	Address1             string       `bson:"addr1,omitempty" json:"addr1,omitempty"`
	Address2             string       `bson:"addr2,omitempty" json:"addr2,omitempty"`
	City                 string       `bson:"city,omitempty" json:"city,omitempty"`
	CountryCode          string       `bson:"cntryCode,omitempty" json:"cntryCode,omitempty"`
	Postcode             int          `bson:"pstcode,omitempty" json:"pstcode,omitempty"`
	CustomercareEmail    string       `bson:"ccEmail,omitempty" json:"ccEmail,omitempty"`
	CustomercareContact  string       `bson:"ccContact,omitempty" json:"ccContact,omitempty"`
	CustomercarePhone    string       `bson:"ccPhone,omitempty" json:"ccPhone,omitempty"`
	CustomerCareAddress1 string       `bson:"ccAddr1,omitempty" json:"ccAddr1,omitempty"`
	CustomerCareAddress2 string       `bson:"ccAddr2,omitempty" json:"ccAddr2,omitempty"`
	CustomerCareCity     string       `bson:"ccCity,omitempty" json:"ccCity,omitempty"`
	CustomerCarePostCode int          `bson:"ccPstcode,omitempty" json:"ccPstcode,omitempty"`
	CustomerCareCountry  string       `bson:"ccCntry,omitempty" json:"ccCntry,omitempty"`
	BankAccName          string       `bson:"bankAccName,omitempty" json:"bankAccName,omitempty"`
	BankAccNo            string       `bson:"bankAccNo,omitempty" json:"bankAccNo,omitempty"`
	BankName             string       `bson:"bankName,omitempty" json:"bankName,omitempty"`
	BankCode             string       `bson:"bankAccCode,omitempty" json:"bankAccCode,omitempty"`
	BankAccountIban      string       `bson:"bankAccIban,omitempty" json:"bankAccIban,omitempty"`
	BankAccountSwift     string       `bson:"bankAccSwift,omitempty" json:"bankAccSwift,omitempty"`
	TermsCondition       string       `bson:"trmsCnd,omitempty" json:"trmsCnd,omitempty"`
	Tagline              string       `bson:"tagline,omitempty" json:"tagline,omitempty"`
	Description          string       `bson:"desc,omitempty" json:"desc,omitempty"`
	VatRegistered        bool         `bson:"vatReg" json:"vatReg"`
	VatNumber            string       `bson:"vatNum,omitempty" json:"vatNum,omitempty"`
	BuisnessRegNo        string       `bson:"buisRegNo,omitempty" json:"buisRegNo,omitempty"`
	FileUrl              string       `bson:"fileUrl,omitempty" json:"fileUrl,omitempty"`
	WarehouseName        string       `bson:"wrName,omitempty" json:"wrName,omitempty"`
	WarehouseAddress1    string       `bson:"wrAddr1,omitempty" json:"wrAddr1,omitempty"`
	WarehouseAddress2    string       `bson:"wrAddr2,omitempty" json:"wrAddr2,omitempty"`
	WarehouseEmail       string       `bson:"wrEml,omitempty" json:"wrEml,omitempty"`
	WarehousePhone       string       `bson:"wrPhn,omitempty" json:"wrPhn,omitempty"`
	WarehousePostcode    string       `bson:"wrPstcode,omitempty" json:"wrPstcode,omitempty"`
	WarehouseCity        string       `bson:"wrCity,omitempty" json:"wrCity,omitempty"`
	WarehouseCountry     string       `bson:"wrCntry,omitempty" json:"wrCntry,omitempty"`
	OrderLimitReached    int          `bson:"ordrLimRchd,omitempty" json:"ordrLimRchd,omitempty"`
	CreatedAt            *time.Time   `bson:"crtdAt,omitempty" json:"crtdAt,omitempty"`
	UpdatedAt            *time.Time   `bson:"updtdAt,omitempty" json:"updtdAt,omitempty"`
	Rating               string       `bson:"rating,omitempty" json:"rating,omitempty"`
	Sync                 bool         `bson:"sync" json:"sync"`
	UpdateCommission     []Commission `json:"updtComisn,omitempty"`
	SellerCustomInfo     interface{}  `bson:"slrCustInfo,omitempty" json:"slrCustInfo,omitempty"`
}

type Org struct {
	Orgdata []Schema `json:"sellerDataList"`
}

type ErpData struct {
	Method string        `json:"method"`
	Params ErpSellerData `json:"params"`
}

type ErpSellerData struct {
	SellerData interface{} `json:"sellerData"`
}

type Commission struct {
	CategoryId int     `json:"ctgrId,omitempty"`
	Percentage float64 `json:"percntge,omitempty"`
}

type ErpResponse struct {
	Status   interface{}              `json:"status"`
	Data     []map[string]interface{} `json:"data"`
	MetaData interface{}              `json:"_metaData"`
}

type SellerData struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type GetCommission struct {
	SeqId      int
	CategoryId int
	Percentage float64
}

type AliceStruct struct {
	SeqId    string `json:"id_catalog_supplier"`
	Status   string `json:"status"`
	SellerId string `json:"seller_id"`
}
