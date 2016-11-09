package erp

import ()

type CsvParser struct {
	CsvFile         string
	CsvSeparator    rune
	SkipFirstLine   bool
	SkipEmptyValues bool
}

type Parser interface {
	Parse(resultType interface{})
}

type csvSchema struct {
	SellerId                      string      `bson:"slrId,omitempty" json:"slrId,omitempty"`
	OrderEmail                    string      `bson:"ordrEml,omitempty" json:"ordrEml,omitempty"`
	ContactName                   string      `bson:"cntctName,omitempty" json:"cntctName,omitempty"`
	OrgName                       string      `bson:"orgName,omitempty" json:"orgName,omitempty"`
	SellerName                    string      `bson:"slrName,omitempty" json:"slrName,omitempty"`
	Phone                         string      `bson:"phn,omitempty" json:"phn,omitempty"`
	Address1                      string      `bson:"addr1,omitempty" json:"addr1,omitempty"`
	Address2                      string      `bson:"addr2,omitempty" json:"addr2,omitempty"`
	City                          string      `bson:"city,omitempty" json:"city,omitempty"`
	Postcode                      int         `bson:"pstcode,omitempty" json:"pstcode,omitempty"`
	Contact                       string      `bson:"cntct,omitempty" json:"cntct,omitempty"`
	BuisnessRegNo                 string      `bson:"buisRegNo,omitempty" json:"buisRegNo,omitempty"`
	CountryCode                   string      `bson:"cntryCode,omitempty" json:"cntryCode,omitempty"`
	BankAccName                   string      `bson:"bankAccName,omitempty" json:"bankAccName,omitempty"`
	BankAccNo                     string      `bson:"bankAccNo,omitempty" json:"bankAccNo,omitempty"`
	BankName                      string      `bson:"bankName,omitempty" json:"bankName,omitempty"`
	BankCode                      string      `bson:"bankAccCode,omitempty" json:"bankAccCode,omitempty"`
	BankAccountIban               string      `bson:"bankAccIban,omitempty" json:"bankAccIban,omitempty"`
	BankAccountSwift              string      `bson:"bankAccSwift,omitempty" json:"bankAccSwift,omitempty"`
	CustomercareEmail             string      `bson:"ccEmail,omitempty" json:"ccEmail,omitempty"`
	CustomercareContact           string      `bson:"ccContact,omitempty" json:"ccContact,omitempty"`
	CustomercarePhone             string      `bson:"ccPhone,omitempty" json:"ccPhone,omitempty"`
	CustomerCareAddress1          string      `bson:"ccAddr1,omitempty" json:"ccAddr1,omitempty"`
	CustomerCareAddress2          string      `bson:"ccAddr2,omitempty" json:"ccAddr2,omitempty"`
	CustomerCareCity              string      `bson:"ccCity,omitempty" json:"ccCity,omitempty"`
	CustomerCarePostCode          int         `bson:"ccPstcode,omitempty" json:"ccPstcode,omitempty"`
	CustomerCareCountry           string      `bson:"ccCntry,omitempty" json:"ccCntry,omitempty"`
	TermsCondition                string      `bson:"trmsCnd,omitempty" json:"trmsCnd,omitempty"`
	Tagline                       string      `bson:"tagline,omitempty" json:"tagline,omitempty"`
	Description                   string      `bson:"desc,omitempty" json:"desc,omitempty"`
	VatRegistered                 bool        `bson:"vatReg,omitempty" json:"vatReg,omitempty"`
	VatNumber                     string      `bson:"vatNum,omitempty" json:"vatNum,omitempty"`
	FileUrl                       string      `bson:"fileUrl,omitempty" json:"fileUrl,omitempty"`
	WarehouseName                 string      `bson:"wrName,omitempty" json:"wrName,omitempty"`
	WarehouseAddress1             string      `bson:"wrAddr1,omitempty" json:"wrAddr1,omitempty"`
	WarehouseAddress2             string      `bson:"wrAddr2,omitempty" json:"wrAddr2,omitempty"`
	WarehouseEmail                string      `bson:"wrEml,omitempty" json:"wrEml,omitempty"`
	WarehousePhone                string      `bson:"wrPhn,omitempty" json:"wrPhn,omitempty"`
	WarehousePostcode             string      `bson:"wrPstcode,omitempty" json:"wrPstcode,omitempty"`
	WarehouseCity                 string      `bson:"wrCity,omitempty" json:"wrCity,omitempty"`
	WarehouseCountry              string      `bson:"wrCntry,omitempty" json:"wrCntry,omitempty"`
	OrderLimitReached             int         `bson:"ordrLimRchd,omitempty" json:"ordrLimRchd,omitempty"`
	Status                        string      `bson:"status,omitempty" json:"status,omitempty"`
	PaymentTermsCode              string      `json:"pymntcode,omitempty"`
	TinNo                         string      `json:"tinNo,omitempty"`
	PanNo                         string      `json:"panNo,omitempty"`
	SerTaxRegNo                   string      `json:"serTaxRegNo,omitempty"`
	CntrctExpDate                 string      `json:"cntrctExpDate,omitempty"`
	StateCode                     string      `json:"stCode,omitempty"`
	IfscCode                      string      `json:"bankIFSCCode,omitempty"`
	CstNo                         string      `json:"cstNo,omitempty"`
	CinNo                         string      `json:"cinNo,omitempty"`
	NatureOfEntity                string      `json:"natOfEntity,omitempty"`
	NatureOfBuisness              string      `json:"natOfBuis,omitempty"`
	ProcessingTime                string      `json:"proctime,omitempty"`
	OneshipCentreCode             string      `json:"oneshipCntrCode,omitempty"`
	OneshipAddress                string      `json:"oneshipAddr,omitempty"`
	OneshipCity                   string      `json:"oneshipCity,omitempty"`
	OneshipZipcode                string      `json:"oneshipZipcode,omitempty"`
	OneshipState                  string      `json:"oneshipState,omitempty"`
	ReturnProcessingCentreCode    string      `json:"retrnProcCntrCode,omitempty"`
	ReturnProcessingCentreAddress string      `json:"retrnProcCntrAddr,omitempty"`
	ReturnProcessingCentreCity    string      `json:"retrnProcCntrCity,omitempty"`
	ReturnProcessingCentreZipcode string      `json:"retrnProcCntrZipcode,omitempty"`
	ReturnProcessingCentreState   string      `json:"retrnProcCntrState,omitempty"`
	DelistingReason               string      `json:"delistReasn,omitempty"`
	PickupTime                    string      `json:"pickupTime,omitempty"`
	ReturnPolicy                  string      `json:"rtrnPolicy,omitempty"`
	PenaltyClause                 string      `json:"penltyClause,omitempty"`
	PenaltyPercentage             string      `json:"penaltyPerc,omitempty"`
	PickUpPartnerCode             string      `json:"pickupPrtnrCode,omitempty"`
	ReversePickUpOneship          string      `json:"reversepickUpCodeOSS,omitempty"`
	ReversePickUpRPC              string      `json:"reversepickUpCodeRPC,omitempty"`
	DispatchLocation              string      `json:"dispatchLoc,omitempty"`
	SellerCustomInfo              interface{} `bson:"slrCustInfo,omitempty" json:"slrCustInfo,omitempty"`
}

type Response struct {
	SellerId string `json:"sellerId"`
	Error    string `json:"error"`
}
