package rating

import (
	"sellers/common"
)

type Request struct {
	Data []common.Schema `json:"sellerRating"`
}
