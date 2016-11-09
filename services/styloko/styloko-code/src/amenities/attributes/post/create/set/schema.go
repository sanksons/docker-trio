package set

import (
	"time"
)

type SetRequestJson struct {
	Name       string `json:"name" validate:"required"`
	Label      string `json:"label" validate:"required"`
	Identifier string `json:"identifier" validate:"required"`
}

type AttributeSet struct {
	SeqId      int       `bson:"seqId" json:"seqId"`
	Name       string    `bson:"name" json:"name" validate:"required"`
	Label      string    `bson:"label" json:"label" validate:"required"`
	Identifier string    `bson:"identifier" json:"identifier" validate:"required"`
	CreatedAt  time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt  time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
