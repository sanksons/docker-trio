package supplier

import (
	"time"
)

type OrgMongo struct {
	SeqId               int        `bson:"seqId" mapstructure:"id_catalog_supplier" json:"seqId"`
	VendorId            *string    `bson:"vendorId" json:"vendorId"`
	OrgName             string     `bson:"orgName" mapstructure:"name" json:"orgName"`
	SellerName          *string    `bson:"slrName,omitempty" json:"sellerName"`
	Status              string     `xorm:"not null ENUM('inactive','active')" bson:"status" mapstructure:"status" json:"status"`
	OrderEmail          *string    `bson:"ordrEml,omitempty" mapstructure:"order_email,omitempty" json:"orderEmail"`
	Contact             *string    `bson:"cntctName,omitempty" mapstructure:"contact,omitempty" json:"contact"`
	Phone               *string    `bson:"phn,omitempty" mapstructure:"phone,omitempty" json:"phone"`
	CustomercareEmail   *string    `bson:"ccEmail,omitempty" mapstructure:"customercare_email,omitempty" json:"customercareEmail"`
	CustomercareContact *string    `bson:"ccName,omitempty" mapstructure:"customercare_contact,omitempty" json:"customercareContact"`
	CustomercarePhone   *string    `bson:"ccPhone,omitempty" mapstructure:"customercare_phone,omitempty" json:"customercarePhone"`
	Street              *string    `bson:"addr1,omitempty" json:"street"`
	StreetNo            *string    `bson:"addr2,omitempty" json:"streetNo"`
	City                *string    `bson:"city,omitempty" json:"city"`
	Postcode            *string    `bson:"pstcode,omitempty" json:"postcode"`
	CountryCode         *string    `bson:"cntryCode,omitempty" json:"countryCode"`
	Sync                bool       `bson:"sync" json:"sync"`
	CreatedAt           *time.Time `bson:"crtdAt,omitempty" mapstructure:"created_at" json:"createdAt"`
	UpdatedAt           *time.Time `bson:"updtdAt,omitempty" mapstructure:"updated_at" json:"updatedAt"`
}

type CounterInfo struct {
	Id    string `bson:"_id"`
	SeqId int    `bson:"seqId"`
}
