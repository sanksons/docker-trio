package set

import (
	"time"
)

type AttributeSet struct {
	Name       string    `bson:"name" json:"name" validate:"required"`
	Label      string    `bson:"label" json:"label" validate:"omitempty"`
	Identifier string    `bson:"identifier" json:"identifier" validate:"omitempty"`
	UpdatedAt  time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
